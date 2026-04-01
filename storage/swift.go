package storage

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"

	"github.com/ncw/swift/v2"

	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/storage/types"
)

const contentType = "application/octet-stream"

type SwiftConfig struct {
	Prefix    string
	Container string
	ChunkSize int64
}

type Swift struct {
	cfg  SwiftConfig
	conn *swift.Connection
}

// NewSwift instantiate a new connection to a Swift object storage. The
// configuration is taken from the environment. Refer to the
// github.com/ncw/swift documentation for more information.
func NewSwift(ctx context.Context, cfg SwiftConfig) (*Swift, error) {
	conn := new(swift.Connection)
	err := conn.ApplyEnvironment()
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get Swift configuration from the environment")
	}

	err = conn.Authenticate(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "authenticate to Swift")
	}

	return &Swift{cfg: cfg, conn: conn}, nil
}

func (s *Swift) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	path = s.fullPath(path)
	log := logger.Get(ctx)
	log.WithField("path", path).Info("Get object")

	object, _, err := s.conn.ObjectOpen(ctx, s.cfg.Container, path, false, swift.Headers{})
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get object %v", path)
	}
	return object, nil
}

func (s *Swift) Upload(ctx context.Context, reader io.Reader, path string) error {
	path = s.fullPath(path)
	segmentPath, err := s.segmentPath(ctx, path)
	if err != nil {
		return errors.Wrap(ctx, err, "generate segment path")
	}
	object, err := s.conn.DynamicLargeObjectCreateFile(ctx, &swift.LargeObjectOpts{
		ObjectName:       path,
		ContentType:      contentType,
		Container:        s.cfg.Container,
		SegmentContainer: s.cfg.Container,
		SegmentPrefix:    segmentPath,
		ChunkSize:        s.cfg.ChunkSize,
	})
	if err != nil {
		return errors.Wrapf(ctx, err, "create dynamic large object %v", path)
	}
	defer object.Close()

	_, err = io.Copy(object, reader)
	if err != nil {
		return errors.Wrapf(ctx, err, "upload content of object %v", path)
	}

	err = object.Flush(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "flush object %v", path)
	}

	return nil
}

// Size returns the size of the content of the object. A retry mechanism is
// implemented because of the eventual consistency of Swift backends NotFound
// error are sometimes returned when the object was just uploaded.
func (s *Swift) Size(ctx context.Context, path string) (int64, error) {
	path = s.fullPath(path)
	info, _, err := s.conn.Object(ctx, s.cfg.Container, path)
	if err != nil {
		return -1, errors.Wrapf(ctx, err, "get object info %v", path)
	}
	return info.Bytes, nil
}

func (s *Swift) Delete(ctx context.Context, path string) error {
	path = s.fullPath(path)
	err := s.conn.DynamicLargeObjectDelete(ctx, s.cfg.Container, path)
	if err != nil {
		if err.Error() == swift.ObjectNotFound.Error() {
			return ObjectNotFound{Path: path}
		}
		return errors.Wrapf(ctx, err, "delete object %v", path)
	}
	return nil
}

func (s *Swift) Info(ctx context.Context, path string) (types.Info, error) {
	path = s.fullPath(path)
	info, _, err := s.conn.Object(ctx, s.cfg.Container, path)
	if err != nil {
		if err.Error() == swift.ObjectNotFound.Error() {
			return types.Info{}, ObjectNotFound{Path: path}
		}
		return types.Info{}, errors.Wrapf(ctx, err, "get object info %v", path)
	}

	return types.Info{
		ContentLength: info.Bytes,
		Checksum:      info.Hash,
	}, nil
}

func (s *Swift) Move(ctx context.Context, srcPath, dstPath string) error {
	err := s.conn.ObjectMove(ctx, s.cfg.Container, srcPath, s.cfg.Container, dstPath)
	if err != nil {
		return errors.Wrapf(ctx, err, "move Swift object %v", srcPath)
	}

	return nil
}

func (s *Swift) List(ctx context.Context, prefix string, opts types.ListOpts) ([]string, error) {
	objects, err := s.conn.ObjectNames(ctx, s.cfg.Container, &swift.ObjectsOpts{
		Prefix: prefix,
		Limit:  int(opts.MaxKeys),
	})
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "list objects in %v", prefix)
	}

	return objects, nil
}

func (s *Swift) segmentPath(ctx context.Context, path string) (string, error) {
	checksum := sha1.New()
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return "", errors.Wrap(ctx, err, "read random bytes for Swift segment path")
	}
	path = hex.EncodeToString(checksum.Sum(append([]byte(path), random...)))
	return strings.TrimLeft(strings.TrimRight(s.cfg.Prefix+"/segments/"+path[0:3]+"/"+path[3:], "/"), "/"), nil
}

func (s *Swift) fullPath(path string) string {
	return strings.TrimLeft(s.cfg.Prefix+"/"+fullPath(path), "/")
}
