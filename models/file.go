package models

type File struct {
	Type         string `json:"type"`
	S3_File_URL  string
	TempFilePath string
	From         string
	To           string
	Timestamp    int64 `json:"timestamp"`
}
