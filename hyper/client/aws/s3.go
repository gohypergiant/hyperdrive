package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	config2 "github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"github.com/seqsense/s3sync"
	"os"
)

func SyncDirectory(s3Config types.S3WorkspacePersistenceRemoteConfiguration, srcPath string, destPath string) {
	syncManager := GetSyncManger(s3Config)
	err := syncManager.Sync(srcPath, destPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func GetSyncManger(s3Config types.S3WorkspacePersistenceRemoteConfiguration) *s3sync.Manager {
	var sess *session.Session
	sess = getSession(s3Config, sess)
	return s3sync.New(sess)
}

func getSession(s3Config types.S3WorkspacePersistenceRemoteConfiguration, sess *session.Session) *session.Session {
	awsConfig := aws.Config{Region: &s3Config.Region}
	accessKey := s3Config.AccessKey
	secret := s3Config.Secret
	token := s3Config.Token
	if s3Config.Profile != "" {

		namedProfileConfig := config2.GetNamedProfileConfig(s3Config.Profile)
		accessKey = namedProfileConfig.AccessKey
		secret = namedProfileConfig.Secret
		token = namedProfileConfig.Token
	}
	creds := credentials.NewStaticCredentials(accessKey, secret, token)
	sess, err := session.NewSession(awsConfig.WithCredentials(creds))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return sess

}
