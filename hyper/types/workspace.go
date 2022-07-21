package types

type IWorkspaceService interface {
	Sync(localPath string, watch bool, studyName string)
	Pull(localPath string, studyName string)
	Pack(studyName string, packPath string)
}

type WorkspaceSyncOptions struct {
	StudyName string
	S3Config  S3WorkspacePersistenceRemoteConfiguration
}
