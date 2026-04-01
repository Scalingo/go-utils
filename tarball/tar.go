package tarball

import (
	"archive/tar"
	"context"
	"io"
	"time"

	"github.com/Scalingo/go-utils/errors/v3"
)

type TarFileReader struct {
	Size int64
	io.Reader
}

// Tar is a methods to write a tarball out of a map of TarFileReader, data can
// then come from disk or memory.
func Tar(ctx context.Context, writer io.Writer, files map[string]TarFileReader) error {
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	for path, reader := range files {
		h := &tar.Header{
			Name:     path,
			ModTime:  time.Now(),
			Mode:     0600 | c_ISREG,
			Typeflag: TypeReg,
			Size:     reader.Size,
		}
		err := tarWriter.WriteHeader(h)
		if err != nil {
			return errors.Wrapf(ctx, err, "add %v to archive", path)
		}

		_, err = io.Copy(tarWriter, reader)
		if err != nil {
			return errors.Wrapf(ctx, err, "copy file content of %v to tar archive", path)
		}
	}
	return nil
}
