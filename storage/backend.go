package storage

import (
	"context"
	"io"

	"github.com/Scalingo/go-utils/storage/types"
)

// BackendMethod represents the name of a Method included in the Backend interface
type BackendMethod string

// ObjectNotFound is a custom error in case the object is not found
type ObjectNotFound struct{}

func (o ObjectNotFound) Error() string {
	return "Object not found"
}

const (
	GetMethod    BackendMethod = "Get"
	UploadMethod BackendMethod = "Upload"
	SizeMethod   BackendMethod = "Size"
	InfoMethod   BackendMethod = "Info"
	DeleteMethod BackendMethod = "Delete"
	MoveMethod   BackendMethod = "Move"
	ListMethod   BackendMethod = "List"
)

// Backend represents something which is able to store files on an object
// storage service
type Backend interface {
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	Upload(ctx context.Context, file io.Reader, path string) error
	Size(ctx context.Context, path string) (int64, error)
	Delete(ctx context.Context, path string) error
	Info(ctx context.Context, path string) (types.Info, error)
	Move(ctx context.Context, srcPath, dstPath string) error
	List(ctx context.Context, prefix string, opts types.ListOpts) ([]string, error)
}

var _ Backend = &S3{}
var _ Backend = &Swift{}
