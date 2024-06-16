package models

type File struct {
	S3_File_URL  string
	TempFilePath string
	From         string
	To           string
	Timestamp    int64 `json:"timestamp"`
}
