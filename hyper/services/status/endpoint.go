package status

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
  port      string
  filePath  string
)

/**
  A public function that updates the status file upon `hyper remoteStatus update "<message>"`
  The `<message>` is supplied through `args[0]`, trimmed of parentheses. 
*/
func UpdateStatus(args []string, statusFilePath string) {
  updateMessage :=  strings.Trim(args[0], "\"")
  if(updateMessage == "") {
    fmt.Println("[remoteStatus] No message provided, --help for more.")
    os.Exit(2)
  }

  err := writeStatusFile(updateMessage, generateStatusFilePath(statusFilePath))

  if err != nil {
    fmt.Println("[remoteStatus] Could not update")
    panic(err)
  }

  fmt.Println("[remoteStatus] updated")
}

/**
  Public function that starts the http server upon `hyper remoteStatus` invokation
*/
func StartEndpoint(statusEndpointPort string, statusFilePath string) {
  if(statusEndpointPort == "" || statusFilePath == "") {
    panic("[remoteStatus] Unable to start server, port & filepath not specified")
  }
  port = statusEndpointPort
  filePath = generateStatusFilePath(statusFilePath)
  
  startHttpServer()
}

/**
  Starts the http server for status discovery.
 */
func startHttpServer() {

  // ensure the desired port is available
  if(!portAvailable()) {
    fmt.Println("[remoteStatus] Port "+port+" already in use")
    os.Exit(2);
  }

  http.HandleFunc("/status", statusPage)
  http.HandleFunc("/", statusPage)

  fmt.Println("[remoteStatus] Endpoint available at http://localhost:"+port+"")
  log.Fatal(http.ListenAndServe(":"+port, nil))
}


/**
  Supports route: `/status`
*/
func statusPage(w http.ResponseWriter, r *http.Request){
  jsonMap, jsonBytes, err :=  readStatusFile(filePath)
  w.Header().Set("Content-Type", "application/json")

  if (err != nil) {
    w.Write(routeError("Internal server error", w, http.StatusInternalServerError))
    fmt.Println("[remoteStatus] Could not retrieve statusFile: ", err)
  }
  
  if (jsonMap.Message == "" && err == nil) {
    w.Write(routeError("No status set", w, http.StatusNoContent))
    fmt.Println("[remoteStatus] No status set")
  }
  
  if (jsonMap.Message != "" && err == nil){
    w.WriteHeader(http.StatusOK)
    w.Write(jsonBytes)
  }

}
  
/**
  Utility function to generate a json message on errors.
*/
func routeError(message string, w http.ResponseWriter, statusCode int)  []byte {
  w.WriteHeader(statusCode)
  jsonBytes, err := createJsonMessage(message)
  if err != nil {
    fmt.Println("[remoteStatus] Cannot format json")
    panic(err)
  }

  return jsonBytes;
}



/**
  Utility function to determine if a port is available, or already in-use.
*/
func portAvailable() bool{
  connection, err := net.Listen("tcp", ":"+port)

  if err != nil {
    return false
  }

  connection.Close()
  return true
}
