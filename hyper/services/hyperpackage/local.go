package hyperpackage

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/gohypergiant/hyperdrive/hyper/client/cli"
	"github.com/gohypergiant/hyperdrive/hyper/client/manifest"
	"github.com/gohypergiant/hyperdrive/hyper/services/notebook"
	"github.com/gohypergiant/hyperdrive/hyper/types"
)

type LocalHyperpackageService struct {
	HyperpackagePath string
	ManifestPath     string
}

func (s LocalHyperpackageService) BuildAndRun(dockerfileSavePath string, imageTags []string, jupyterOptions types.JupyterLaunchOptions, ec2Options types.EC2StartOptions, syncOptions types.WorkspaceSyncOptions, dockerOptions types.DockerOptions) {

	studyName := manifest.GetName(s.ManifestPath)
	if len(imageTags) == 0 {
		imageTags = []string{fmt.Sprintf("%s:latest", studyName)}
	}
	runTag := imageTags[0]

	s.Build(dockerfileSavePath, imageTags, syncOptions)
	s.Run(runTag, dockerOptions)
}
func (s LocalHyperpackageService) Build(dockerfileSavePath string, imageTags []string, syncOptions types.WorkspaceSyncOptions) {
	dockerClient := cli.NewDockerClient()
	dockerClient.CreateDockerFile(s.HyperpackagePath, dockerfileSavePath, false, syncOptions)
	dockerClient.BuildImage(strings.TrimLeft(dockerfileSavePath, "./"), imageTags)
}
func (s LocalHyperpackageService) Run(imageTag string, dockerOptions types.DockerOptions) {
	var hostIP, hostPort string
	dockerClient := cli.NewDockerClient()
	studyName := fmt.Sprintf("%s_%s", HYPERPACK_CONTAINER_PREFIX, manifest.GetName(s.ManifestPath))

	if dockerOptions.LocalOnly {
		hostIP = "127.0.0.1"
	} else {
		hostIP = "0.0.0.0"
	}

	if dockerOptions.HostPort <= -1 {
		hostPort = ""
	} else {
		hostPort = strconv.Itoa(dockerOptions.HostPort)
	}

	contConfig := &container.Config{
		Hostname: studyName,
		Image:    imageTag,
		Tty:      true,
		Env:      []string{"JUPYTER_TOKEN=firefly"},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8001/tcp": []nat.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: hostPort,
				},
			},
		},
		Mounts: []mount.Mount{},
	}
	createdId, err := dockerClient.CreateContainer(imageTag, studyName, contConfig, hostConfig, false)
	id := createdId
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dockerClient.ExecuteContainer(id, false)

	nowRunningContainers, _ := dockerClient.ListAllRunningContainers()

	for _, runningContainer := range nowRunningContainers {
		if runningContainer.ID == id {

			publicPort := runningContainer.Ports[0].PublicPort
			fmt.Println("Hyperpackage now running via Docker Container", runningContainer.ID[:10], "on port", publicPort)
		}
	}
}
func (s LocalHyperpackageService) Import(importModelFileName string, modelFlavor string, trainShape string) {

	if importModelFileName == "" {
		fmt.Println("Error: Must specify filename of trained model to be imported with the --filename flag.")
		os.Exit(1)
	}

	if trainShape == "" {
		fmt.Println("Error: Must specify the number of columns in the training data. Use the --shape flag.")
		os.Exit(1)
	}

	mod := fmt.Sprintf("You are importing a %s model", modelFlavor)
	fmt.Println(mod)

	dockerClient := cli.NewDockerClient()
	cwdPath, _ := os.Getwd()
	name := fmt.Sprintf("imported_%s", modelFlavor)
	hostIP := "127.0.0.1"
	imageOptions := notebook.GetNotebookImageOptions("local")
	imageName := imageOptions.Image
	env := []string{"JUPYTER_TOKEN=firefly"}
	inImageCache := false
	pullImage := false

	clientImages, _ := dockerClient.ListImages()
	for _, clientImage := range clientImages {
		for _, tag := range clientImage.RepoTags {
			if tag == imageOptions.Image {
				inImageCache = true
				break
			}
		}
	}
	if !inImageCache {
		pullImage = true
	}

	contConfig := &container.Config{
		Hostname: name,
		Image:    imageName,
		Tty:      true,
		Env:      env,
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8888/tcp": []nat.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: "",
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cwdPath,
				Target: "/home/jovyan",
			},
		},
	}

	runningContainers, _ := dockerClient.ListContainers(name)
	if len(runningContainers) != 0 {
		dockerClient.RemoveContainer(name)
	}

	createdId, err := dockerClient.CreateContainer(imageOptions.Image, name, contConfig, hostConfig, pullImage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dockerClient.ExecuteContainer(createdId, false)

	nowRunningContainers, _ := dockerClient.ListAllRunningContainers()

	for _, runningContainer := range nowRunningContainers {
		if runningContainer.ID == createdId {
			publicPort := runningContainer.Ports[0].PublicPort
			fmt.Println("Importing process now running via Docker Container", runningContainer.ID[:10], "on port", publicPort)
		}
	}

	notebookOutPath := "/home/jovyan/import_outs.ipynb"

	_, errExec := exec.Command("docker", "exec", name, "papermill",
		"/home/jovyan/.executor/notebooks/importer.ipynb", notebookOutPath, "-p", "filename", importModelFileName, "-p", "flavor", modelFlavor, "-p", "shape", trainShape).Output()
	if errExec != nil {
		fmt.Println("Error with importer notebook execution in the docker container: ", err)
		os.Exit(1)
	}
}
func (s LocalHyperpackageService) List() {

	fmt.Println("Currently running hyperpackages:")

	dockerClient := cli.NewDockerClient()
	formattedPrefix := fmt.Sprintf("/%s_", HYPERPACK_CONTAINER_PREFIX)
	prefixLength := len(formattedPrefix)

	runningContainers, _ := dockerClient.ListAllRunningContainers()

	for _, runningContainer := range runningContainers {
		for _, name := range runningContainer.Names {
			if len(name) > prefixLength && name[:prefixLength] == formattedPrefix {
				image := strings.Replace(runningContainer.Image, "docker.io/", "", -1)
				image = strings.Replace(image, ":latest", "", -1)
				publicPort := runningContainer.Ports[0].PublicPort
				fmt.Println(
					"Name:", name[prefixLength:], "\n",
					"Image:", image, "\n",
					"Container Id:", runningContainer.ID[:10], "\n",
					"Url:", fmt.Sprintf("http://127.0.0.1:%d/info", publicPort),
				)
			}
		}
	}
}
func (s LocalHyperpackageService) Stop(name string) {
	dockerClient := cli.NewDockerClient()

	runningContainers, _ := dockerClient.ListAllRunningContainers()
	containerId := ""
	formattedName := fmt.Sprintf("/%s_%s", HYPERPACK_CONTAINER_PREFIX, name)

	for _, runningContainer := range runningContainers {
		if runningContainer.Names[0] == formattedName {
			containerId = runningContainer.ID
			break
		}
	}

	if containerId == "" {
		fmt.Println("No container found for ", name)
	} else {
		err := dockerClient.RemoveContainer(containerId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
