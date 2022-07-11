package types

type IWorkspaceService interface {
	Sync(localPath string, watch bool, studyName string)
	Pull(localPath string, studyName string)
}

type WorkspaceSyncOptions struct {
	StudyName string
	S3Config  S3WorkspacePersistenceRemoteConfiguration
}
