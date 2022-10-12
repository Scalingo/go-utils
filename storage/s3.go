package storage

import (
	"context"
	stderrors "errors"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/pkg/errors"

	"github.com/Scalingo/go-utils/logger"
)

const (
	NotFoundErrCode = "NotFound"
	// DefaultPartSize 16 MB part size define the size in bytes of the parts
	// uploaded in a multipart upload
	DefaultPartSize = int64(16777216)
	// DefaultUploadConcurrency defines that multipart upload will be done in
	// parallel in 2 routines
	DefaultUploadConcurrency = int(2)
)

type S3Client interface {
	GetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	HeadObject(ctx context.Context, input *s3.HeadObjectInput, opts ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

type S3Config struct {
	AK       string
	SK       string
	Region   string
	Endpoint string
	Bucket   string
}

type RetryPolicy struct {
	WaitDuration   time.Duration
	Attempts       int
	MethodHandlers map[BackendMethod][]string
}

type S3 struct {
	cfg               S3Config
	s3client          S3Client
	s3uploader        *manager.Uploader
	retryPolicy       RetryPolicy
	uploadConcurrency int
	partSize          int64
}

type s3Opt func(s3 *S3)

// WithRetryPolicy is an option to constructor NewS3 to add a Retry Policy
// impacting GET operations
func WithRetryPolicy(policy RetryPolicy) s3Opt {
	return s3Opt(func(s3 *S3) {
		s3.retryPolicy = policy
	})
}

func WithPartSize(size int64) s3Opt {
	return s3Opt(func(s3 *S3) {
		s3.partSize = size
	})
}

func WithUploadConcurrency(concurrency int) s3Opt {
	return s3Opt(func(s3 *S3) {
		s3.uploadConcurrency = concurrency
	})
}

func NewS3(cfg S3Config, opts ...s3Opt) *S3 {
	s3config := s3Config(cfg)
	s3client := s3.NewFromConfig(s3config)
	s3 := &S3{
		cfg:      cfg,
		s3client: s3client,
		retryPolicy: RetryPolicy{
			WaitDuration: time.Second,
			Attempts:     3,
			MethodHandlers: map[BackendMethod][]string{
				SizeMethod: {NotFoundErrCode},
			},
		},
	}
	for _, opt := range opts {
		opt(s3)
	}

	partSize := DefaultPartSize
	if s3.partSize != 0 {
		partSize = s3.partSize
	}
	concurrency := DefaultUploadConcurrency
	if s3.uploadConcurrency != 0 {
		concurrency = s3.uploadConcurrency
	}
	s3.s3uploader = manager.NewUploader(s3client, func(u *manager.Uploader) {
		u.PartSize = partSize
		u.Concurrency = concurrency
	})
	return s3
}

func (s *S3) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	path = fullPath(path)
	log := logger.Get(ctx)
	log.WithField("path", path).Info("Get object")

	input := &s3.GetObjectInput{
		Bucket: &s.cfg.Bucket,
		Key:    &path,
	}
	out, err := s.s3client.GetObject(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to get object %v", path)
	}
	return out.Body, nil
}

func (s *S3) Upload(ctx context.Context, file io.Reader, path string) error {
	path = fullPath(path)
	input := &s3.PutObjectInput{
		Body:   file,
		Bucket: &s.cfg.Bucket,
		Key:    &path,
	}
	_, err := s.s3uploader.Upload(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "fail to save file to %v", path)
	}

	return nil
}

// Size returns the size of the content of the object. A retry mechanism is
// implemented because of the eventual consistency of S3 backends NotFound
// error are sometimes returned when the object was just uploaded.
func (s *S3) Size(ctx context.Context, path string) (int64, error) {
	path = fullPath(path)
	var res int64
	err := s.retryWrapper(ctx, SizeMethod, func(ctx context.Context) error {
		log := logger.Get(ctx).WithField("key", path)
		log.Infof("[s3] Size()")

		input := &s3.HeadObjectInput{Bucket: &s.cfg.Bucket, Key: &path}
		stat, err := s.s3client.HeadObject(ctx, input)
		if err != nil {
			return err
		}
		res = stat.ContentLength
		return nil
	})

	if err != nil {
		return -1, errors.Wrapf(err, "fail to HEAD object '%v'", path)
	}
	return res, nil
}

func (s *S3) Delete(ctx context.Context, path string) error {
	path = fullPath(path)
	input := &s3.DeleteObjectInput{Bucket: &s.cfg.Bucket, Key: &path}
	_, err := s.s3client.DeleteObject(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "fail to delete object %v", path)
	}

	return nil
}

func (s *S3) retryWrapper(ctx context.Context, method BackendMethod, fun func(ctx context.Context) error) error {
	var err error

	errorCodes := s.retryPolicy.MethodHandlers[method]
	// no-op is no retry policy on the method
	if errorCodes == nil {
		return fun(ctx)
	}
	for i := 0; i < s.retryPolicy.Attempts; i++ {
		log := logger.Get(ctx).WithField("attempt", i+1)
		ctx := logger.ToCtx(ctx, log)
		err = fun(ctx)
		if err == nil {
			return nil
		}
		var apiErr smithy.APIError
		if stderrors.As(err, &apiErr) {
			for _, code := range errorCodes {
				if apiErr.ErrorCode() == code {
					time.Sleep(s.retryPolicy.WaitDuration)
				}
			}
		}
	}
	return err
}

func s3Config(cfg S3Config) aws.Config {
	credentials := credentials.NewStaticCredentialsProvider(cfg.AK, cfg.SK, "")
	config := aws.Config{
		Region:      cfg.Region,
		Credentials: credentials,
	}
	if cfg.Endpoint != "" {
		config.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           "https://" + cfg.Endpoint,
				SigningRegion: cfg.Region,
			}, nil
		})
	}

	return config
}

func fullPath(path string) string {
	return strings.TrimLeft("/"+path, "/")
}
