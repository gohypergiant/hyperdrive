package types

type JupyterLaunchOptions struct {
	Flavor        string
	APIKey        string
	Password      string
	HostPort      string
	PullImage     bool
	LaunchBrowser bool
	Requirements  bool
	RestartAlways bool
}
type INotebookService interface {
	Start(jupyterOptions JupyterLaunchOptions, ec2Options EC2StartOptions)
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