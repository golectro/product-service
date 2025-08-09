package model

import "io"

type UploadFileInput struct {
	Bucket      string
	ObjectKey   string
	Content     []byte
	ContentType string
}

type PresignedURLInput struct {
	Bucket    string
	ObjectKey string
	Expiry    int64
}

type UploadFileRequest struct {
	Bucket string `form:"bucket" validate:"required"`
}

type DeleteFileRequest struct {
	Bucket    string `json:"bucket" validate:"required"`
	ObjectKey string `json:"object_key" validate:"required"`
}

type MinioObjectResponse struct {
	Bucket    string
	ObjectKey string
	Object    io.ReadCloser
}

// // Read implements io.Reader.
// func (m *MinioObjectResponse) Read(p []byte) (n int, err error) {
// 	panic("unimplemented")
// }

// func (m *MinioObjectResponse) Close() {
// 	panic("unimplemented")
// }
