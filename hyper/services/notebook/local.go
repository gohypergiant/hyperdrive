package notebook

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/gohypergiant/hyperdrive/hyper/client/cli"
	"github.com/gohypergiant/hyperdrive/hyper/client/manifest"
	"github.com/pkg/browser"
)

var (
	id             string
	image          string
	jupyterBrowser bool
	mountPoint     string
	pullImage      bool
	repoTag        string
	requirements   bool
	publicPort     uint16
)

type LocalNotebookService struct {
	ManifestPath  string
	S3Credentials S3Credentials
}

func (s LocalNotebookService) Start(flavor string, pullImage bool,
	jupyterBrowser bool, requirements bool, ec2Options EC2StartOptions,
	hostPort string, restartAlways bool, s3AwsProfile string) {

	awsConfigFilePath := config.DefaultSharedConfigFilename()
	if _, errFile := os.Stat(awsConfigFilePath); errFile == nil {
		// AWS config file exists at $HOME/.aws/config. We're good.
	} else if errors.Is(errFile, os.ErrNotExist) {
		fmt.Println("Error:", awsConfigFilePath, "does not exist. Please create one.")
		os.Exit(1)
	}

	ctx := context.TODO()
	cfg, errConfig := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(s3AwsProfile))
	if errConfig != nil {
		fmt.Println("Error:", errConfig)
		os.Exit(1)
	}
	creds, errCreds := cfg.Credentials.Retrieve(ctx)
	if errCreds != nil {
		fmt.Println(errCreds)
		os.Exit(1)
	}
	fmt.Println("creds access key id:", creds.AccessKeyID)
	fmt.Println("creds secret access key:", creds.SecretAccessKey)
	fmt.Println("creds session token:", creds.SessionToken)
	os.Exit(1)

	dockerClient := cli.NewDockerClient()
	cwdPath, _ := os.Getwd()
	name := GetNotebookName(s.ManifestPath)
	hostIP := "0.0.0.0"
	execute := false
	projectName := manifest.GetProjectName(s.ManifestPath)

	imageOptions := GetNotebookImageOptions("dev") // change to "local" later
	clientImages, _ := dockerClient.ListImages()
	inImageCache := false
	env := []string{"JUPYTER_TOKEN=firefly",
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", s.S3Credentials.AccessKey),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", s.S3Credentials.AccessSecret),
		fmt.Sprintf("AWS_DEFAULT_REGION=%s", s.S3Credentials.Region),
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
		pullImage = true
	}

	runningContainers, _ := dockerClient.ListContainers(name)

	imageName := ""
	if requirements {
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
		Name:              "unless-stopped",
	}
	if restartAlways {
		restartPolicy = container.RestartPolicy{
			Name: "always",
		}
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8888/tcp": []nat.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: hostPort,
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

	if requirements {
		if len(runningContainers) != 0 {
			dockerClient.RemoveContainer(name)
		}
		dockerClient.CreateDockerFile("", "Dockerfile.reqs", true)
		dockerClient.BuildImage("Dockerfile.reqs", []string{imageName})

		createdIdReqs, errReqs := dockerClient.CreateContainer(imageName, name, contConfig, hostConfig, false)
		id = createdIdReqs
		if errReqs != nil {
			fmt.Println(errReqs)
			os.Exit(1)
		}
		execute = true
	} else if len(runningContainers) == 0 {
		createdId, err := dockerClient.CreateContainer(imageOptions.Image, name, contConfig, hostConfig, pullImage)
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

	nowRunningContainers, _ := dockerClient.ListContainers(name)

	for _, runningContainer := range nowRunningContainers {
		publicPort = runningContainer.Ports[0].PublicPort
		fmt.Println("Jupyter Lab Now Running via Docker Container", runningContainer.ID[:10], "on port", publicPort)

	}

	if jupyterBrowser {
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

	os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)

	in, err := os.Open(srcPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer in.Close()

	out, err := os.Create(dstPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	out.Close()
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
		container := dockerClient.InspectContainer(runningContainer.ID)
		if strings.HasPrefix(container.Mounts[0].Source, rootPath) {
			containerMount = container.Mounts[0].Source
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
