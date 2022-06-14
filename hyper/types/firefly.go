package types

import "time"

type ServerInfo struct {
	LastActivity time.Time   `json:"last_activity"`
	Name         string      `json:"name"`
	Pending      interface{} `json:"pending"`
	ProgressURL  string      `json:"progress_url"`
	Ready        bool        `json:"ready"`
	Started      time.Time   `json:"started"`
	State        struct {
		PodName string `json:"pod_name"`
	} `json:"state"`
	URL         string   `json:"url"`
	UserOptions struct{} `json:"user_options"`
}

type ListServersResponse struct {
	Admin        bool                  `json:"admin"`
	AuthState    interface{}           `json:"auth_state"`
	Created      time.Time             `json:"created"`
	Groups       []interface{}         `json:"groups"`
	Kind         string                `json:"kind"`
	LastActivity time.Time             `json:"last_activity"`
	Name         string                `json:"name"`
	Pending      interface{}           `json:"pending"`
	Server       interface{}           `json:"server"`
	Servers      map[string]ServerInfo `json:"servers"`
}

type CreateServerOptions struct {
	Profile string `json:"profile"`
}

type UploadType string
type UploadFormat string

type UploadDataBody struct {
	Content  string       `json:"content"`
	Format   UploadFormat `json:"format"`
	FileType UploadType   `json:"type"`
}

type TrainingStatus string

type DownloadFileResponse struct {
	Content string `json:"content"`
}
