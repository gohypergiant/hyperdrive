package status

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gohypergiant/hyperdrive/hyper/types"
)

func readStatusFile(filePath string) (types.RemoteStatus, []byte, error) {
  if !fileExists(filePath) {
    return types.RemoteStatus{Message: ""}, nil, nil
  }

  jsonBytes, readErr := ioutil.ReadFile(filePath)
  if(readErr != nil) {
    return types.RemoteStatus{Message: ""}, nil, readErr
  }

  var jsonMap types.RemoteStatus
  err := json.Unmarshal(jsonBytes, &jsonMap)
  if(err != nil) {
    return types.RemoteStatus{Message: ""}, nil, err
  }

  return jsonMap, jsonBytes, nil;
}

func writeStatusFile(status string, filepath string) error {
  jsonString, jsonErr := createJsonMessage(status)

  if jsonErr != nil {
    return jsonErr
  }

  err := ioutil.WriteFile(filepath, jsonString, os.ModePerm) 
  if err != nil {
    return err
  }

  return nil;
}




func generateStatusFilePath(path string) string {
  workingdir, err := os.Getwd()
  if err != nil {
    fmt.Println(err)
  }
  return workingdir + path
}

func fileExists(filename string) bool {
  info, err := os.Stat(filename)

  if os.IsNotExist(err) {
    return false
  }
  
  return !info.IsDir()
}

func createJsonMessage(message string) ([]byte, error) {
  status := types.RemoteStatus{
    Message: message,
  }
  
  jsonBytes, jsonErr := json.Marshal(status)
  if jsonErr != nil {
    return nil, jsonErr
  }

  return jsonBytes, nil

}