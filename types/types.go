package types

type FileType string

var (
	TypeUrl        FileType = "url"
	TypeUrlFile    FileType = "urlfile"
	TypeFile       FileType = "file"
	TypeRemoteFile FileType = "remoteFile"
)

type Request struct {
	Path    string
	Type    FileType
	Options map[string]string
}

func NewRequest(path string, fileType FileType, opts map[string]string) Request {
	return Request{path, fileType, opts}
}

type ResponseItem struct {
	URL     string `json:"url"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
