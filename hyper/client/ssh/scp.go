package ssh

import (
	"context"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func CopyToRemote(username string, privateKeyPath string, remoteServerIP string, filePath string, saveFolderPath string) error {
	publickKey := GetPublicKeyFromPrivateKey(privateKeyPath)
	clientConfig, _ := auth.PrivateKey(username, privateKeyPath, ssh.FixedHostKey(publickKey))

	client := scp.NewClient(remoteServerIP, &clientConfig)
	err := client.Connect()
	if err != nil {
		return err
	}

	f, _ := os.Open(filePath)

	defer client.Close()
	defer f.Close()
	err = client.CopyFile(context.TODO(), f, saveFolderPath, "0655")

	if err != nil {
		return err
	}
	return nil
}
