package firefly

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func ListServers(configuration types.FireflyComputeRemoteConfiguration) types.ListServersResponse {
	rootUrl := GetHubAPIRoot(configuration)
	endpoint := fmt.Sprintf("%s/users/%s", rootUrl, configuration.Username)
	token := fmt.Sprintf("token %s", configuration.HubToken)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var listServerResponse types.ListServersResponse
	err = json.Unmarshal(body, &listServerResponse)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return listServerResponse
}

func StartServer(configuration types.FireflyComputeRemoteConfiguration, name string, profile string) {

	rootUrl := GetHubAPIRoot(configuration)
	endpoint := fmt.Sprintf("%s/users/%s/servers/%s", rootUrl, configuration.Username, name)
	token := fmt.Sprintf("token %s", configuration.HubToken)
	postBody, err := json.Marshal(types.CreateServerOptions{
		Profile: profile,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	notebookUrl := fmt.Sprintf("%s/user/%s/%s", rootUrl, configuration.Username, name)
	fmt.Println(fmt.Sprintf("Your notebook should be available at %s shortly", notebookUrl))
}
func StopServer(configuration types.FireflyComputeRemoteConfiguration, name string) {

	rootUrl := GetHubAPIRoot(configuration)
	endpoint := fmt.Sprintf("%s/users/%s/servers/%s", rootUrl, configuration.Username, name)
	token := fmt.Sprintf("token %s", configuration.HubToken)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	//fmt.Println(string(body))
}

const (
	NotebookUploadType  types.UploadType = "notebook"
	FileUploadType                       = "file"
	DirectoryUploadType                  = "directory"
)
const (
	JsonUploadFormat   types.UploadFormat = "json"
	TextUploadFormat                      = "text"
	Base64UploadFormat                    = "base64"
)

func GetEncodedFile(path string) string {

	f, _ := os.Open(path)
	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	// Print encoded data to console.
	// ... The base64 image can be used as a data URI in a browser.
	//fmt.Println("ENCODED: " + encoded)
	return encoded
}
func GetHubAPIRoot(configuration types.FireflyComputeRemoteConfiguration) string {
	return fmt.Sprintf("%s/hub/api", configuration.Url)
}
func GetNotebookAPIRoot(configuration types.FireflyComputeRemoteConfiguration, notebookName string) string {
	return fmt.Sprintf("%s/user/%s/%s/api", configuration.Url, configuration.Username, notebookName)
}
func MkDir(configuration types.FireflyComputeRemoteConfiguration, notebookName string, remotePath string) {

	// Recursively create parents directories first
	splitPath := strings.Split(remotePath, "/")
	if len(splitPath) > 2 { //Greater than 2 since the leading / adds an element
		MkDir(configuration, notebookName, strings.Join(splitPath[:len(splitPath)-1], "/"))
	}

	rootUrl := GetNotebookAPIRoot(configuration, notebookName)
	endpoint := fmt.Sprintf("%s/contents%s", rootUrl, remotePath)
	reqBody, _ := json.Marshal(types.UploadDataBody{
		Content:  "",
		Format:   "",
		FileType: DirectoryUploadType,
	})
	token := fmt.Sprintf("token %s", configuration.HubToken)
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	//fmt.Println(string(body))
}

func UploadData(configuration types.FireflyComputeRemoteConfiguration, notebookName string, localPath string, remotePath string) {

	//Create parent directory
	splitPath := strings.Split(remotePath, "/")
	parentDir := strings.Join(splitPath[:len(splitPath)-1], "/")
	if len(splitPath) > 2 { //Greater than 2 since the leading / adds an element
		MkDir(configuration, notebookName, parentDir)
	}

	rootUrl := GetNotebookAPIRoot(configuration, notebookName)
	endpoint := fmt.Sprintf("%s/contents%s", rootUrl, remotePath)
	encodedFile := GetEncodedFile(localPath)
	reqBody, _ := json.Marshal(types.UploadDataBody{
		Content:  encodedFile,
		Format:   Base64UploadFormat,
		FileType: FileUploadType,
	})
	token := fmt.Sprintf("token %s", configuration.HubToken)
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	//fmt.Println(string(body))
}

const (
	TrainingPending  types.TrainingStatus = "pending"
	TrainingStarted  types.TrainingStatus = "started"
	TrainingComplete types.TrainingStatus = "completed"
)

func GetTrainingStatus(configuration types.FireflyComputeRemoteConfiguration, notebookName string, studyDir string) types.TrainingStatus {
	startedPath := fmt.Sprintf("%s/STARTED", studyDir)
	completedPath := fmt.Sprintf("%s/COMPLETED", studyDir)
	if FileExists(configuration, notebookName, startedPath) {
		return TrainingStarted
	} else if FileExists(configuration, notebookName, completedPath) {
		return TrainingComplete
	}
	return TrainingPending
}

func FileExists(configuration types.FireflyComputeRemoteConfiguration, notebookName string, filepath string) bool {
	rootUrl := GetNotebookAPIRoot(configuration, notebookName)
	endpoint := fmt.Sprintf("%s/contents%s?content=0", rootUrl, filepath)
	token := fmt.Sprintf("token %s", configuration.HubToken)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func DownloadFile(configuration types.FireflyComputeRemoteConfiguration, notebookName string, filepath string) string {
	rootUrl := GetNotebookAPIRoot(configuration, notebookName)
	endpoint := fmt.Sprintf("%s/contents%s?content=1&format=base64", rootUrl, filepath)
	token := fmt.Sprintf("token %s", configuration.HubToken)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Could not download hyperpack")
		os.Exit(1)
	}
	defer resp.Body.Close()
	var responseBody types.DownloadFileResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	json.Unmarshal(body, &responseBody)
	return responseBody.Content
}
