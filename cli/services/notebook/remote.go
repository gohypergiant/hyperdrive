package notebook

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gohypergiant/hyperdrive/cli/client/firefly"
	"github.com/gohypergiant/hyperdrive/cli/client/manifest"
	"github.com/gohypergiant/hyperdrive/cli/services/config"
)

const jobsDir string = "_jobs"

type RemoteNotebookService struct {
	RemoteConfiguration config.RemoteConfiguration
	ManifestPath        string
}

func (s RemoteNotebookService) Start(flavor string, pullImage bool, jupyterBrowser bool) {

	imageOptions := GetNotebookImageOptions(flavor)
	name := GetNotebookName(s.ManifestPath)
	fmt.Println("Starting remote notebook instance")
	firefly.StartServer(s.RemoteConfiguration, name, imageOptions.Profile)
}
func (s RemoteNotebookService) List() {

	resp := firefly.ListServers(s.RemoteConfiguration)

	for name, info := range resp.Servers {
		fmt.Println(fmt.Sprintf("%s:", name))
		fmt.Println("URL: ", info.URL)
	}
}
func (s RemoteNotebookService) Stop(identifier string) {
	name := GetNotebookName(s.ManifestPath)
	firefly.StopServer(s.RemoteConfiguration, name)
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
	featuresDataFilePath := strings.TrimLeft(manifestConfig.Training.Data.Features.Source, "./")
	firefly.UploadData(s.RemoteConfiguration, notebookName, featuresDataFilePath, fmt.Sprintf("%s/%s", studyRoot, featuresDataFilePath))
	fmt.Println("Uploading target data")
	targetDataFilePath := strings.TrimLeft(manifestConfig.Training.Data.Target.Source, "./")
	firefly.UploadData(s.RemoteConfiguration, notebookName, targetDataFilePath, fmt.Sprintf("%s/%s", studyRoot, targetDataFilePath))
	fmt.Println("Uploading Study Manifest")
	firefly.UploadData(s.RemoteConfiguration, notebookName, s.ManifestPath, fmt.Sprintf("%s/_study.yaml", studyRoot))

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
			status := firefly.GetTrainingStatus(s.RemoteConfiguration, notebookName, studyRoot)
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
	return fmt.Sprintf("%s/%s.hyperpack.zip", studyRoot, studyName)
}
func (s RemoteNotebookService) GetHyperpackSavePath() string {

	studyName := manifest.GetName(s.ManifestPath)
	return fmt.Sprintf("./%s.hyperpack.zip", studyName)
}
func (s RemoteNotebookService) DownloadHyperpack() {

	hyperpackPath := s.GetRemoteHyperpackPath()
	notebookName := GetNotebookName(s.ManifestPath)
	savePath := s.GetHyperpackSavePath()
	fmt.Println("Downloading hyperpack from remote")
	base64File := firefly.DownloadFile(s.RemoteConfiguration, notebookName, hyperpackPath)
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
