package types

type IWorkspaceService interface {
	Sync(localPath string, watch bool, studyName string)
}

type WorkspaceSyncOptions struct {
	StudyName string
	S3Config  S3WorkspacePersistenceRemoteConfiguration
}
