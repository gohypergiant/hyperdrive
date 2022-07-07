package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/fs"
	"io/ioutil"
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
func MarshalPublicKey(publicKey interface{}) []byte {
	publicRsaKey, err := ssh.NewPublicKey(publicKey)
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
	publicKeyBytes := GetPublicKeyBytes(&privateKey.PublicKey)

	return privateKeyBytes, publicKeyBytes
}
func ParsePrivateKey(keyName string) *rsa.PrivateKey {

	privateKeyBytes, err := ioutil.ReadFile(keyName)
	if err != nil {
		panic("error reading private key file " + err.Error())
	}

	pemBlock, _ := pem.Decode(privateKeyBytes)
	if pemBlock == nil {
		panic("error decoding private key " + err.Error())
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		panic("error parsing private key " + err.Error())
	}

	return privateKey
}

func GetPublicKeyBytes(privateKeyName string) []byte {
	privateKey := ParsePrivateKey(privateKeyName)

	publicKeyBytes := MarshalPublicKey(&privateKey.PublicKey)

	return publicKeyBytes
}
func GetPublicKeyFromPrivateKey(privateKeyName string) ssh.PublicKey {
	privateKey := ParsePrivateKey(privateKeyName)

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic("error getting public key " + err.Error())
	}
	return publicKey
}
