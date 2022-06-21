package types

type ComputeRemoteType string
type WorkspacePersistenceRemoteType string

const (
	S3 WorkspacePersistenceRemoteType = "s3"
)

var ValidWorkspacePersistenceRemoteTypes = []WorkspacePersistenceRemoteType{
	S3,
}

const (
	EC2     ComputeRemoteType = "ec2"
	Firefly ComputeRemoteType = "firefly"
)

var ValidRemoteTypes = []ComputeRemoteType{
	Firefly,
	EC2,
}

type ComputeRemoteConfiguration struct {
	Type                 ComputeRemoteType                 `mapstructure:"type" json:"type"`
	FireflyConfiguration FireflyComputeRemoteConfiguration `mapstructure:"firefly" json:"firefly"`
	EC2Configuration     EC2ComputeRemoteConfiguration     `mapstructure:"ec2" json:"ec2"`
	JupyterAPIKey        string                            `mapstructure:"jupyter_api_key" json:"jupyter_api_key"`
}
type WorkspacePersistenceRemoteConfiguration struct {
	Type      WorkspacePersistenceRemoteType `mapstructure:"type" json:"type"`
	IsDefault bool                           `mapstructure:"default" json:"default"`
}
type FireflyComputeRemoteConfiguration struct {
	Url      string `mapstructure:"url" json:"url"`
	HubToken string `mapstructure:"hub_token" json:"hub_token"`
	Username string `mapstructure:"username" json:"username"`
}
type EC2ComputeRemoteConfiguration struct {
	Profile   string `mapstructure:"profile" json:"profile"`
	AccessKey string `mapstructure:"access_key" json:"access_key"`
	Secret    string `mapstructure:"secret" json:"secret"`
	Region    string `mapstructure:"region" json:"region"`
	Token     string
}
type Configuration struct {
	SchemaVersion               string                                             `mapstructure:"schema_version" json:"schema_version"`
	ComputeRemotes              map[string]ComputeRemoteConfiguration              `mapstructure:"compute_remotes" json:"compute_remotes"`
	WorkspacePersistenceRemotes map[string]WorkspacePersistenceRemoteConfiguration `mapstructure:"workspace_persistence_remotes" json:"workspace_persistence_remotes"`
}
type NamedProfileConfiguration struct {
	AccessKey string
	Secret    string
	Token     string
	Region    string
}
