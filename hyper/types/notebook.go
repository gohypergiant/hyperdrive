package types

type JupyterLaunchOptions struct {
	Flavor        string
	APIKey        string
	HostPort      int
	PullImage     bool
	LaunchBrowser bool
	Requirements  bool
	RestartAlways bool
	S3AwsProfile  string
}
type INotebookService interface {
	Start(jupyterOptions JupyterLaunchOptions, ec2Options EC2StartOptions, syncOptions WorkspaceSyncOptions)
	List()
	Stop(mountPointOrIdentifier string)
	UploadTrainingJobData()
	WaitForTrainingToComplete(timeout int)
	DownloadHyperpack()
}
type S3Credentials struct {
	AccessKey    string
	AccessSecret string
	Region       string
}

type ImageOptions struct {
	Image   string
	RepoTag string
	Profile string
}
