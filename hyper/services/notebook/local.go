package notebook

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gohypergiant/hyperdrive/hyper/types"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/gohypergiant/hyperdrive/hyper/client/cli"
	"github.com/gohypergiant/hyperdrive/hyper/client/manifest"
	"github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/pkg/browser"
)

var (
	id         string
	publicPort uint16
)

type LocalNotebookService struct {
	ManifestPath  string
	S3Credentials types.S3Credentials
}

func (s LocalNotebookService) Start(jupyterOptions types.JupyterLaunchOptions, _ types.EC2StartOptions, _ types.WorkspaceSyncOptions) {

	dockerClient := cli.NewDockerClient()
	cwdPath, _ := os.Getwd()
	name := GetNotebookName(s.ManifestPath)
	hostIP := "0.0.0.0"
	execute := false
	projectName := manifest.GetProjectName(s.ManifestPath)

	imageOptions := GetNotebookImageOptions("local")
	clientImages, _ := dockerClient.ListImages()
	inImageCache := false
	awsAccessKeyId := ""
	awsSecretAccessKey := ""
	awsSessionToken := ""
	region := ""
	if jupyterOptions.S3AwsProfile != "" {
		fmt.Printf("Using AWS named profile '%s' to retrieve AWS creds\n", jupyterOptions.S3AwsProfile)
		namedProfileConfig := config.GetNamedProfileConfig(jupyterOptions.S3AwsProfile)
		awsAccessKeyId = namedProfileConfig.AccessKey
		awsSecretAccessKey = namedProfileConfig.Secret
		awsSessionToken = namedProfileConfig.Token
		region = namedProfileConfig.Region
	} else {
		awsAccessKeyId = s.S3Credentials.AccessKey
		awsSecretAccessKey = s.S3Credentials.AccessSecret
		region = s.S3Credentials.Region
	}
	env := []string{"JUPYTER_TOKEN=firefly",
		fmt.Sprintf("NB_TOKEN=%s", jupyterOptions.APIKey),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", awsAccessKeyId),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", awsSecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%s", awsSessionToken),
		fmt.Sprintf("AWS_DEFAULT_REGION=%s", region),
		fmt.Sprintf("HYPER_PROJECT_NAME=%s", projectName),
	}
	for _, clientImage := range clientImages {
		for _, tag := range clientImage.RepoTags {
			if tag == imageOptions.Image {
				inImageCache = true
				break
			}
		}
	}
	if !inImageCache {
		jupyterOptions.PullImage = true
	}

	runningContainers, _ := dockerClient.ListContainers(name)

	imageName := ""
	if jupyterOptions.Requirements {
		imageName = fmt.Sprintf("hyperdrive-jupyter-reqs:%s", name)
	} else {
		imageName = imageOptions.Image
	}

	contConfig := &container.Config{
		Hostname: name,
		Image:    imageName,
		Tty:      true,
		Env:      env,
	}

	restartPolicy := container.RestartPolicy{
		Name: "unless-stopped",
	}
	if jupyterOptions.RestartAlways {
		restartPolicy = container.RestartPolicy{
			Name: "always",
		}
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8888/tcp": []nat.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: strconv.Itoa(jupyterOptions.HostPort),
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
		RestartPolicy: restartPolicy,
	}

	if jupyterOptions.Requirements {
		if len(runningContainers) != 0 {
			err := dockerClient.RemoveContainer(name)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		dockerClient.CreateDockerFile("", "Dockerfile.reqs", true, types.WorkspaceSyncOptions{})
		dockerClient.BuildImage("Dockerfile.reqs", []string{imageName})

		createdIdReqs, errReqs := dockerClient.CreateContainer(imageName, name, contConfig, hostConfig, false)
		id = createdIdReqs
		if errReqs != nil {
			fmt.Println(errReqs)
			os.Exit(1)
		}
		execute = true
	} else if len(runningContainers) == 0 {
		createdId, err := dockerClient.CreateContainer(imageOptions.Image, name, contConfig, hostConfig, jupyterOptions.PullImage)
		id = createdId
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		execute = true
	} else {
		id = runningContainers[0].ID
	}

	if execute {
		dockerClient.ExecuteContainer(id, false)
	}
	time.Sleep(1 * time.Second)

	nowRunningContainers, _ := dockerClient.ListContainers(name)

	for _, runningContainer := range nowRunningContainers {
		publicPort = runningContainer.Ports[0].PublicPort
		fmt.Println("Jupyter Lab Now Running via Docker Container", runningContainer.ID[:10], "on port", publicPort)

	}

	if jupyterOptions.LaunchBrowser {
		url := fmt.Sprintf("http://%s:%d/lab?token=firefly", hostIP, publicPort)
		fmt.Println("Launching Jupyter Lab")
		fmt.Println("    Mount Point:", cwdPath)
		fmt.Println("    Opening:", url)
		if execute {
			time.Sleep(2 * time.Second)
		}
		err := browser.OpenURL(url)
		if err != nil {
			fmt.Println("failed to open browser")
			os.Exit(0) // Probably fine to exit 0 if it's just the browser that didn't open
		}
	}
}
func (s LocalNotebookService) List() {

	dockerClient := cli.NewDockerClient()

	runningContainers, _ := dockerClient.ListAllRunningContainers()

	for _, runningContainer := range runningContainers {
		for _, name := range runningContainer.Names {
			fmt.Println(name)
			if len(name) > 16 && name[:16] == "/firefly-jupyter" {
				image := strings.Replace(runningContainer.Image, "docker.io/", "", -1)
				image = strings.Replace(image, ":latest", "", -1)
				volMount := runningContainer.Mounts[0]
				publicPort := runningContainer.Ports[0].PublicPort
				fmt.Println(
					"Mount:", volMount.Source[9:], "\n",
					"Image:", image, "\n",
					"Container Id:", runningContainer.ID[:10], "\n",
					"Url:", fmt.Sprintf("http://127.0.0.1:%d/lab?token=firefly", publicPort),
				)
			}
		}
	}

}
func (s LocalNotebookService) Stop(mountPoint string) {
	dockerClient := cli.NewDockerClient()

	runningContainers, _ := dockerClient.ListAllRunningContainers()
	containerId := ""
	name := GetNotebookName(s.ManifestPath)

	for _, runningContainer := range runningContainers {
		if strings.ToLower(runningContainer.Names[0][1:]) == strings.ToLower(name) {
			containerId = runningContainer.ID
			break
		}
	}

	if containerId == "" {
		fmt.Println("No container found for mount point", mountPoint)
	} else {
		err := dockerClient.RemoveContainer(containerId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}
func (s LocalNotebookService) UploadTrainingJobData() {

	manifestConfig := manifest.GetManifest(s.ManifestPath)
	//Right now, these two are the same, but in the future I'm sure that will change
	//studyName := manifestConfig.StudyName
	//notebookName := GetNotebookName(s.ManifestPath)
	jobsPath := s.GetJobsPath()

	fmt.Println("Uploading features data")
	//upload data
	featuresDataFilePath := strings.TrimLeft(manifestConfig.Training.Data.Features.Source, "./")
	s.CopyFile(featuresDataFilePath, fmt.Sprintf("%s/%s", jobsPath, featuresDataFilePath))
	fmt.Println("Uploading target data")
	targetDataFilePath := strings.TrimLeft(manifestConfig.Training.Data.Target.Source, "./")
	s.CopyFile(targetDataFilePath, fmt.Sprintf("%s/%s", jobsPath, targetDataFilePath))
	fmt.Println("Uploading Study Manifest")
	s.CopyFile(s.ManifestPath, fmt.Sprintf("%s/_study.yaml", jobsPath))

	fmt.Println("Upload complete")
}
func (s LocalNotebookService) CopyFile(srcPath string, dstPath string) {

	err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	in, err := os.Open(srcPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		closeErr := in.Close()
		if closeErr != nil {
			fmt.Println(closeErr)
		}
	}()

	out, err := os.Create(dstPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		closeErr := out.Close()
		if closeErr != nil {
			fmt.Println(closeErr)
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		closeErr := out.Close()
		if closeErr != nil {
			fmt.Println(closeErr)
		}
	}()
}
func (s LocalNotebookService) GetStudyRoot() string {

	manifestConfig := manifest.GetManifest(s.ManifestPath)
	//Right now, these two are the same, but in the future I'm sure that will change
	studyName := manifestConfig.StudyName
	return fmt.Sprintf("/%s/%s", jobsDir, studyName)
}
func (s LocalNotebookService) WaitForTrainingToComplete(timeout int) {

	jobsPath := s.GetJobsPath()

	fmt.Print("Waiting for training to complete")
	fmt.Println()
	for i := 0; i <= timeout; i++ {
		if i%3 == 0 || i == timeout {
			status := s.GetTrainingStatus(fmt.Sprintf("%s/", jobsPath))
			if status == TrainingComplete {
				fmt.Println()
				fmt.Println("Training completed")
				return
			}
			fmt.Printf("\nTraining status: %s.\nWaiting.", status)
		} else {
			fmt.Print(".")
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
	fmt.Println("Timed out waiting for study to complete")
	os.Exit(1)
}

type TrainingStatus string

const (
	TrainingPending  TrainingStatus = "pending"
	TrainingStarted  TrainingStatus = "started"
	TrainingComplete TrainingStatus = "completed"
)

func (s LocalNotebookService) GetTrainingStatus(studyDir string) TrainingStatus {
	startedPath := fmt.Sprintf("%s/STARTED", studyDir)
	completedPath := fmt.Sprintf("%s/COMPLETED", studyDir)
	if s.FileExists(startedPath) {
		return TrainingStarted
	} else if s.FileExists(completedPath) {
		return TrainingComplete
	}
	return TrainingPending
}
func (s LocalNotebookService) FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !errors.Is(err, os.ErrNotExist)
}
func (s LocalNotebookService) GetServerPath(rootPath string) string {
	dockerClient := cli.NewDockerClient()

	runningContainers, _ := dockerClient.ListAllRunningContainers()
	containerMount := ""

	for _, runningContainer := range runningContainers {
		c := dockerClient.InspectContainer(runningContainer.ID)
		if strings.HasPrefix(c.Mounts[0].Source, rootPath) {
			containerMount = c.Mounts[0].Source
			break
		}
	}

	if containerMount == "" {
		fmt.Println("No container found running with root path", rootPath)
		os.Exit(1)
	}

	return containerMount

}

func (s LocalNotebookService) GetJobsPath() string {
	studyName := manifest.GetName(s.ManifestPath)
	return fmt.Sprintf("_jobs/%s", studyName)
}

func (s LocalNotebookService) GetHyperpackArtifactPath() string {
	studyName := manifest.GetName(s.ManifestPath)
	jobsPath := s.GetJobsPath()
	return fmt.Sprintf("%s/%s.hyperpack.zip", jobsPath, studyName)
}
func (s LocalNotebookService) GetHyperpackSavePath() string {
	studyName := manifest.GetName(s.ManifestPath)
	return fmt.Sprintf("%s.hyperpack.zip", studyName)
}
func (s LocalNotebookService) DownloadHyperpack() {

	hyperpackPath := s.GetHyperpackArtifactPath()
	savePath := s.GetHyperpackSavePath()
	fmt.Printf("Saving to %s \n", savePath)

	s.CopyFile(hyperpackPath, savePath)

	fmt.Println("Done")
}
