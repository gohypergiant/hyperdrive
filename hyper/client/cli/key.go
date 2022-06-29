package cli

import (
    "fmt"
    "os"
    
    "github.com/sethvargo/go-password/password"
)

func generateFastAppAPIKey() string {
    fastAppApiKey, err := password.Generate(64, 10, 0, false, true)
    if err != nil {
        fmt.Println(err)
		os.Exit(1)
    }
    
    fmt.Println("Here is the generated FastApp API key:", fastAppApiKey)

	return fastAppApiKey
}
