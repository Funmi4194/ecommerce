package types

type Object struct {
	Name          string `json:"name"`
	OriginalName  string `json:"original_name"`
	RemoteAddress string `json:"remote_address"`
	Size          int64  `json:"size"`
	FileFormat    string `json:"file_format"`
	Error         error  `json:"error"`
}
