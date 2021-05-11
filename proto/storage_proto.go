package proto

type MkdirREQ struct {
	Pwd  string `json:"pwd,omitempty"`
	Name string `json:"name,omitempty"`
}

type ListRES struct {
	Pwd           string `json:"pwd,omitempty"`
	Name          string `json:"name,omitempty"`
	Ext           string `json:"ext,omitempty"`
	DownloadTimes int64  `json:"download_times,omitempty"`
	Type          int64  `json:"type,omitempty"`
}

type TryUpload struct {
	Filehash string `json:"file_hash,omitempty"`
}

type ConfirmUpload struct {
	UploadToken string `json:"upload_token,omitempty"`
	Pwd         string `json:"pwd,omitempty"`
	FileName    string `json:"file_name,omitempty"`
}
