package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/fs"
	"os"

	"golang.org/x/crypto/ssh"
)

const PRIVATE_KEY_FILE_MODE uint32 = 0600
const SSH_FOLDER_FILE_MODE uint32 = 0700
const DEFAULT_KEY string = "id_rsa"

func WriteKey(fileName string, keyBytes []byte, permissions fs.FileMode) error {

	err := os.WriteFile(fileName, keyBytes, permissions)
	return err
}
func GetPublicKeyBytes(privateKey *rsa.PrivateKey) []byte {
	publicRsaKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic("error getting public key " + err.Error())
	}

	return ssh.MarshalAuthorizedKey(publicRsaKey)
}
func GetPrivateKeyBytes(key *rsa.PrivateKey) []byte {
	privateKeyDer := x509.MarshalPKCS1PrivateKey(key)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}

	return pem.EncodeToMemory(&privateKeyBlock)
}
func CreateRSAKeyPair(keyName string) ([]byte, []byte) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		panic("error generating private key " + err.Error())
	}

	privateKeyBytes := GetPrivateKeyBytes(privateKey)
	publicKeyBytes := GetPublicKeyBytes(privateKey)

	return privateKeyBytes, publicKeyBytes
}
