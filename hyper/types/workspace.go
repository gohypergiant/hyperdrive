package types

type IWorkspaceService interface {
	Sync(localPath string)
}
