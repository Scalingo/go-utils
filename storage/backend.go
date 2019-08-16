package storage

import (
	"context"
	"io"
)

// BackendMethod represents the name of a Method included in the Backend interface
type BackendMethod string

const (
	GetMethod    BackendMethod = "Get"
	UploadMethod BackendMethod = "Upload"
	SizeMethod   BackendMethod = "Size"
	DeleteMethod BackendMethod = "Delete"
)

// Backend represents something which is able to store files on an object
// storage service
type Backend interface {
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	Upload(ctx context.Context, file io.Reader, path string) error
	Size(ctx context.Context, path string) (int64, error)
	Delete(ctx context.Context, path string) error
}

var _ Backend = &S3{}
var _ Backend = &Swift{}
