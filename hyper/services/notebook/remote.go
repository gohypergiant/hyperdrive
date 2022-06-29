package notebook

import (
	"encoding/base64"
	"fmt"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"os"
	"path"
	"time"

	"github.com/gohypergiant/hyperdrive/hyper/client/aws"
	"github.com/gohypergiant/hyperdrive/hyper/client/firefly"
	"github.com/gohypergiant/hyperdrive/hyper/client/manifest"
	"github.com/gohypergiant/hyperdrive/hyper/services/config"
)

const jobsDir string = "_jobs"

type RemoteNotebookService struct {
	RemoteConfiguration types.ComputeRemoteConfiguration
	ManifestPath        string
}

func (s RemoteNotebookService) Start(jupyterOptions types.JupyterLaunchOptions, ec2Options types.EC2StartOptions, syncOptions types.WorkspaceSyncOptions) {

	imageOptions := GetNotebookImageOptions(jupyterOptions.Flavor)
	name := GetNotebookName(s.ManifestPath)
	fmt.Println("Starting remote notebook instance")
	jupyterOptions.APIKey = s.RemoteConfiguration.JupyterAPIKey
	if s.RemoteConfiguration.Type == types.Firefly {
		firefly.StartServer(s.RemoteConfiguration.FireflyConfiguration, name, imageOptions.Profile)
	} else if s.RemoteConfiguration.Type == types.EC2 {
		if jupyterOptions.S3AwsProfile != "" {
			fmt.Printf("Using AWS named profile '%s' to retrieve AWS creds\n", jupyterOptions.S3AwsProfile)
			namedProfileConfig := config.GetNamedProfileConfig(jupyterOptions.S3AwsProfile)
			s.RemoteConfiguration.EC2Configuration.AccessKey = namedProfileConfig.AccessKey
			s.RemoteConfiguration.EC2Configuration.Secret = namedProfileConfig.Secret
			s.RemoteConfiguration.EC2Configuration.Region = namedProfileConfig.Region
			s.RemoteConfiguration.EC2Configuration.Token = namedProfileConfig.Token
		}
		aws.StartJupyterEC2(s.ManifestPath, s.RemoteConfiguration.EC2Configuration, ec2Options.InstanceType, ec2Options.AmiId, jupyterOptions, syncOptions)
	} else {
		fmt.Println("Not Implemented")
	}
}
func (s RemoteNotebookService) List() {

	if s.RemoteConfiguration.Type == types.Firefly {
		resp := firefly.ListServers(s.RemoteConfiguration.FireflyConfiguration)

		for name, info := range resp.Servers {
			fmt.Println(fmt.Sprintf("%s:", name))
			fmt.Println("URL: ", info.URL)
		}
	} else if s.RemoteConfiguration.Type == types.EC2 {
		aws.ListServers(s.RemoteConfiguration.EC2Configuration)
	} else {
		fmt.Println("Not Implemented")
	}
}
func (s RemoteNotebookService) Stop(identifier string) {
	name := GetNotebookName(s.ManifestPath)
	if s.RemoteConfiguration.Type == types.Firefly {
		firefly.StopServer(s.RemoteConfiguration.FireflyConfiguration, name)
	} else if s.RemoteConfiguration.Type == types.EC2 {
		aws.StopServer(s.ManifestPath, s.RemoteConfiguration.EC2Configuration)
	} else {
		fmt.Println("Not Implemented")
	}
}
func (s RemoteNotebookService) UploadTrainingJobData() {

	manifestConfig := manifest.GetManifest(s.ManifestPath)
	//Right now, these two are the same, but in the future I'm sure that will change
	//studyName := manifestConfig.StudyName
	notebookName := GetNotebookName(s.ManifestPath)
	studyRoot := s.GetStudyRoot()
	//upload study

	fmt.Println("Uploading features data")
	//upload data
	featuresDataFilePath := path.Clean(manifestConfig.Training.Data.Features.Source)
	firefly.UploadData(s.RemoteConfiguration.FireflyConfiguration, notebookName, featuresDataFilePath, fmt.Sprintf("%s/%s", studyRoot, featuresDataFilePath))
	fmt.Println("Uploading target data")
	targetDataFilePath := path.Clean(manifestConfig.Training.Data.Target.Source)
	firefly.UploadData(s.RemoteConfiguration.FireflyConfiguration, notebookName, targetDataFilePath, fmt.Sprintf("%s/%s", studyRoot, targetDataFilePath))
	fmt.Println("Uploading Study Manifest")
	firefly.UploadData(s.RemoteConfiguration.FireflyConfiguration, notebookName, s.ManifestPath, fmt.Sprintf("%s/_study.yaml", studyRoot))

	fmt.Println("Upload complete")
}
func (s RemoteNotebookService) GetStudyRoot() string {

	manifestConfig := manifest.GetManifest(s.ManifestPath)
	//Right now, these two are the same, but in the future I'm sure that will change
	studyName := manifestConfig.StudyName
	return fmt.Sprintf("/%s/%s", jobsDir, studyName)
}
func (s RemoteNotebookService) WaitForTrainingToComplete(timeout int) {

	notebookName := GetNotebookName(s.ManifestPath)
	studyRoot := s.GetStudyRoot()
	fmt.Print("Waiting for training to complete")
	fmt.Println()
	for i := 0; i <= timeout; i++ {
		if i%3 == 0 || i == timeout {
			status := firefly.GetTrainingStatus(s.RemoteConfiguration.FireflyConfiguration, notebookName, studyRoot)
			if status == firefly.TrainingComplete {
				fmt.Println()
				fmt.Println("Training completed")
				return
			}
			fmt.Print(fmt.Sprintf("\nTraining status: %s.\nWaiting.", status))
		} else {
			fmt.Print(".")
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
	fmt.Println("Timed out waiting for study to complete")
	os.Exit(1)
}
func (s RemoteNotebookService) GetRemoteHyperpackPath() string {

	studyRoot := s.GetStudyRoot()
	studyName := manifest.GetName(s.ManifestPath)
	return path.Join(studyRoot, fmt.Sprintf("%s.hyperpack.zip", studyName))
}
func (s RemoteNotebookService) GetHyperpackSavePath() string {

	studyName := manifest.GetName(s.ManifestPath)
	return path.Join(".", fmt.Sprintf("%s.hyperpack.zip", studyName))
}
func (s RemoteNotebookService) DownloadHyperpack() {

	hyperpackPath := s.GetRemoteHyperpackPath()
	notebookName := GetNotebookName(s.ManifestPath)
	savePath := s.GetHyperpackSavePath()
	fmt.Println("Downloading hyperpack from remote")
	base64File := firefly.DownloadFile(s.RemoteConfiguration.FireflyConfiguration, notebookName, hyperpackPath)
	decodedFile, err := base64.StdEncoding.DecodeString(base64File)

	fmt.Println(fmt.Sprintf("Saving to %s", savePath))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file, err := os.Create(savePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	if _, err := file.Write(decodedFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := file.Sync(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done")
}
