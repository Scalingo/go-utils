package tarball

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	fspkg "github.com/Scalingo/go-utils/fs"
	iopkg "github.com/Scalingo/go-utils/io"
	"github.com/Scalingo/go-utils/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// CreateOpts gathers all possible options to create a tarball archive
type CreateOpts struct {
	Fs          fspkg.Fs
	IncludeRoot bool

	// CopyBufferSize is the size of the memory buffer used to perform the copy
	// from the file to the tarball writer
	CopyBufferSize int64
}

func Create(ctx context.Context, src string, dst io.Writer, opts CreateOpts) error {
	log := logger.Get(ctx).WithField("src", src)
	log.Debug("Create archive")

	if src == "" {
		return errors.New("src is empty")
	}
	if dst == nil {
		return errors.New("dst writer is nil")
	}
	if opts.CopyBufferSize == 0 {
		// 512kB of read buffer when reading file (instead of 32kB default)
		opts.CopyBufferSize = 512 * 1024
	}

	copierOpts := []iopkg.CopierOpt{
		iopkg.WithBufferSize(opts.CopyBufferSize),
		iopkg.WithNoDiskCacheRead,
	}
	copier := iopkg.NewCopier(copierOpts...)

	fs := opts.Fs
	if fs == nil {
		fs = fspkg.NewOsFs()
	}

	gzWriter := gzip.NewWriter(dst)
	defer gzWriter.Close()
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	err := afero.Walk(fs, src, func(fullpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path := strings.TrimPrefix(fullpath, src)
		if opts.IncludeRoot {
			path = filepath.Join(filepath.Base(src), path)
		}

		header, err := FileInfoHeader(info, path, fullpath)
		if err != nil {
			return err
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			file, err := fs.Open(fullpath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = copier.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return errors.Wrapf(err, "fail to create tar archive of %v", src)
	}
	return nil
}

type ExtractOpts struct {
	Fs             fspkg.Fs
	User           string
	UID            int
	GID            int
	CopyBufferSize int64
}

func Extract(ctx context.Context, dst string, reader io.Reader, opts *ExtractOpts) error {
	log := logger.Get(ctx).WithField("dst", dst)
	log.Debug("Extract archive")

	fs := opts.Fs
	if fs == nil {
		fs = fspkg.NewOsFs()
	}

	if opts.CopyBufferSize == 0 {
		// 512kB of read buffer when reading file (instead of 32kB default)
		opts.CopyBufferSize = 512 * 1024
	}
	copierOpts := []iopkg.CopierOpt{
		iopkg.WithBufferSize(opts.CopyBufferSize),
		iopkg.WithNoDiskCacheWrite,
	}
	copier := iopkg.NewCopier(copierOpts...)

	var (
		archiveUserUID int
		archiveUserGID int
	)

	if opts.User != "" {
		u, err := user.Lookup(opts.User)
		if err != nil {
			return errors.Wrapf(err, "fail to get user %v", opts.User)
		}
		archiveUserUID, _ = strconv.Atoi(u.Uid)
		archiveUserGID, _ = strconv.Atoi(u.Gid)
	}
	if opts.UID != 0 {
		archiveUserUID = opts.UID
	}
	if opts.GID != 0 {
		archiveUserGID = opts.GID
	}

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return errors.Wrapf(err, "invalid gzip")
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrapf(err, "invalid tar")
		}

		path := dst + "/" + header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			fs.MkdirAll(path, header.FileInfo().Mode())
			if archiveUserUID != 0 && archiveUserGID != 0 {
				fs.Chown(path, archiveUserUID, archiveUserGID)
			} else {
				fs.Chown(path, header.Uid, header.Gid)
			}
		case tar.TypeReg:
			fd, err := fs.OpenFile(path, os.O_CREATE|os.O_WRONLY, header.FileInfo().Mode())
			if err != nil {
				return errors.Wrapf(err, "invalid err file")
			}
			_, err = copier.Copy(fd, tarReader)
			if err != nil {
				return errors.Wrapf(err, "invalid err file during copy")
			}
			fd.Close()
			if archiveUserUID != 0 && archiveUserGID != 0 {
				fs.Chown(path, archiveUserUID, archiveUserGID)
			} else {
				fs.Chown(path, header.Uid, header.Gid)
			}
		case tar.TypeSymlink:
			err := fs.Symlink(header.Linkname, path)
			if err != nil {
				return errors.Wrapf(err, "fail to create symlink from cache %v -> %v", path, header.Linkname)
			}
		case tar.TypeLink:
			fs.Link(path, dst+"/"+header.Linkname)
		default:
		}
	}
	return nil
}
