package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	config2 "github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"github.com/seqsense/s3sync"
)

var sess *session.Session
var syncManager *s3sync.Manager

func SyncDirectory(s3Config types.S3WorkspacePersistenceRemoteConfiguration, srcPath string, destPath string) {
	syncManager := GetSyncManger(s3Config)
	err := syncManager.Sync(srcPath, destPath)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
}
func GetSyncManger(s3Config types.S3WorkspacePersistenceRemoteConfiguration) *s3sync.Manager {
	if syncManager == nil {
		sess = getSession(s3Config, sess)
		syncManager = s3sync.New(sess)
	}
	return syncManager
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
		os.Exit(1)
	}
	return sess

}
func DownloadObject(s3Config types.S3WorkspacePersistenceRemoteConfiguration, filename string, key string) error {
	sess = getSession(s3Config, sess)
	downloader := s3manager.NewDownloader(sess)

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", filename, err)
	}

	fmt.Println("Downloading " + key + " from bucket " + s3Config.BucketName)
	_, err = downloader.Download(f,
		&s3.GetObjectInput{
			Bucket: aws.String(s3Config.BucketName),
			Key:    aws.String(key),
		})

	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}
	return nil
}
