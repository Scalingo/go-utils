package storage

import (
	"context"
	stderrors "errors"
	"fmt"
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
	"github.com/Scalingo/go-utils/storage/types"
)

const (
	NotFoundErrCode = "NotFound"
	// DefaultPartSize 16 MB part size define the size in bytes of the parts
	// uploaded in a multipart upload
	DefaultPartSize = int64(16777216)
	// DefaultUploadConcurrency defines that multipart upload will be done in
	// parallel in 2 routines
	DefaultUploadConcurrency = int(2)
	// S3ListMaxKeys is defined here https://github.com/aws/aws-sdk-go-v2/blob/v1.17.1/service/s3/api_op_ListObjectsV2.go#L107-L109
	S3ListMaxKeys = 1000
)

type S3Client interface {
	GetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	HeadObject(ctx context.Context, input *s3.HeadObjectInput, opts ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	CopyObject(ctx context.Context, params *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
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
		return nil, errors.Wrapf(err, "fail to get S3 object %v", path)
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
		return errors.Wrapf(err, "fail to upload file to S3 %v", path)
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
		res = aws.ToInt64(stat.ContentLength)
		return nil
	})

	if err != nil {
		return -1, errors.Wrapf(err, "fail to HEAD S3 object '%v'", path)
	}
	return res, nil
}

func (s *S3) Delete(ctx context.Context, path string) error {
	path = fullPath(path)
	input := &s3.DeleteObjectInput{Bucket: &s.cfg.Bucket, Key: &path}
	_, err := s.s3client.DeleteObject(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if stderrors.As(err, &apiErr) {
			if apiErr.ErrorCode() == NotFoundErrCode {
				return ObjectNotFound{Path: path}
			}
		}
		return errors.Wrapf(err, "fail to delete S3 object %v", path)
	}

	return nil
}

// Info returns several information contained in the header.
// It returns ObjectNotFound custom error if the object does not exists.
func (s *S3) Info(ctx context.Context, path string) (types.Info, error) {
	path = fullPath(path)
	var res *s3.HeadObjectOutput
	err := s.retryWrapper(ctx, InfoMethod, func(ctx context.Context) error {
		input := &s3.HeadObjectInput{Bucket: &s.cfg.Bucket, Key: &path}
		var err error
		res, err = s.s3client.HeadObject(ctx, input)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		var apiErr smithy.APIError
		if stderrors.As(err, &apiErr) {
			if apiErr.ErrorCode() == NotFoundErrCode {
				return types.Info{}, ObjectNotFound{Path: path}
			}
		}
		return types.Info{}, errors.Wrapf(err, "fail to HEAD S3 object '%v'", path)
	}

	info := types.Info{
		ContentLength: aws.ToInt64(res.ContentLength),
	}

	if res.ContentType != nil {
		info.ContentType = *res.ContentType
	}

	if res.ETag != nil {
		info.Checksum = *res.ETag
	}

	return info, nil
}

func (s *S3) Move(ctx context.Context, srcPath, dstPath string) error {
	dstPath = fullPath(dstPath)
	srcPathWithBucket := fmt.Sprintf("%s/%s", s.cfg.Bucket, fullPath(srcPath))
	_, err := s.s3client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     &s.cfg.Bucket,
		Key:        &dstPath,
		CopySource: &srcPathWithBucket,
	})
	if err != nil {
		return errors.Wrapf(err, "fail to copy S3 object '%v' to '%v'", srcPathWithBucket, dstPath)
	}

	err = s.Delete(ctx, srcPath)
	if err != nil {
		return errors.Wrapf(err, "fail to delete the old object '%v' on S3", srcPath)
	}
	return nil
}

// List function lists object contained in bucket up to 1,000 objects.
// If maxKeys > 1,000, S3 will set maxKeys to 1,000. Source: https://github.com/aws/aws-sdk-go-v2/blob/v1.17.1/service/s3/api_op_ListObjectsV2.go#L16
func (s *S3) List(ctx context.Context, prefix string, opts types.ListOpts) ([]string, error) {
	objects, err := s.s3client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  &s.cfg.Bucket,
		Prefix:  &prefix,
		MaxKeys: &opts.MaxKeys,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "fail to list S3 objects")
	}

	strObjects := make([]string, aws.ToInt32(objects.KeyCount))
	for i, o := range objects.Contents {
		strObjects[i] = *o.Key
	}
	return strObjects, nil
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
		config.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
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
