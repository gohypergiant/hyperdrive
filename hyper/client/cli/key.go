package cli

import (
    "fmt"
    "os"
    
    "github.com/sethvargo/go-password/password"
)

func generateFastAppAPIKey() string {
    fastAppApiKey, err := password.Generate(64, 10, 0, true, true)
    if err != nil {
        fmt.Println(err)
		os.Exit(1)
    }
    
    log.Printf("A FastApp Token of %s has been generated.", fastAppApiKey)

	return fastAppApiKey
}
