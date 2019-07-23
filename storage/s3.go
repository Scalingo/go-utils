package storage

import (
	"context"
	"io"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"github.com/pkg/errors"
)

const (
	NotFoundErrCode = "NotFound"
)

type S3Client interface {
	GetObjectRequest(input *s3.GetObjectInput) s3.GetObjectRequest
	HeadObjectRequest(input *s3.HeadObjectInput) s3.HeadObjectRequest
	DeleteObjectRequest(input *s3.DeleteObjectInput) s3.DeleteObjectRequest
}

type S3Config struct {
	AK       string
	SK       string
	Region   string
	Endpoint string
	Bucket   string
}

type S3 struct {
	cfg               S3Config
	s3client          S3Client
	s3uploader        *s3manager.Uploader
	retryWaitDuration time.Duration
	retryAttempts     int
}

func NewS3(cfg S3Config, bucket string) *S3 {
	s3config := s3Config(cfg)
	return &S3{
		cfg: cfg, s3client: s3.New(s3config), s3uploader: s3manager.NewUploader(s3config),
		retryWaitDuration: time.Second, retryAttempts: 3,
	}
}

func (s *S3) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	log := logger.Get(ctx)
	log.WithField("path", path).Info("Get object")

	input := &s3.GetObjectInput{
		Bucket: &s.cfg.Bucket,
		Key:    &path,
	}
	out, err := s.s3client.GetObjectRequest(input).Send(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to get object %v", path)
	}
	return out.Body, nil
}

func (s *S3) Upload(ctx context.Context, file io.Reader, path string) error {
	input := &s3manager.UploadInput{
		Body:   file,
		Bucket: &s.cfg.Bucket,
		Key:    &path,
	}
	_, err := s.s3uploader.UploadWithContext(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "fail to save file to %v", path)
	}

	return nil
}

// Size returns the size of the content of the object. A retry mecanism is
// implemented because of the eventual consistency of S3 backends NotFound
// error are sometimes returned when the object was just uploaded.
func (s *S3) Size(ctx context.Context, path string) (int64, error) {
	var (
		err   error
		stat  *s3.HeadObjectResponse
		log   = logger.Get(ctx)
		input = &s3.HeadObjectInput{Bucket: &s.cfg.Bucket, Key: &path}
	)
	for i := 0; i < s.retryAttempts; i++ {
		stat, err = s.s3client.HeadObjectRequest(input).Send(ctx)
		if err == nil {
			return *stat.ContentLength, nil
		}

		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == NotFoundErrCode {
			log.WithField("key", path).Info("[s3] retry HEAD")
			time.Sleep(s.retryWaitDuration)
			continue
		}
		return -1, errors.Wrapf(err, "fail to HEAD object '%v'", path)
	}

	return -1, errors.Wrapf(err, "fail to HEAD object '%v' after %v retries", path, s.retryAttempts)
}

func (s *S3) Delete(ctx context.Context, path string) error {
	input := &s3.DeleteObjectInput{Bucket: &s.cfg.Bucket, Key: &path}
	req := s.s3client.DeleteObjectRequest(input)
	_, err := req.Send(ctx)
	if err != nil {
		return errors.Wrapf(err, "fail to delete object %v", path)
	}

	return nil
}

func s3Config(cfg S3Config) aws.Config {
	credentials := aws.NewStaticCredentialsProvider(cfg.AK, cfg.SK, "")
	config := aws.Config{
		Region:      cfg.Region,
		Handlers:    defaults.Handlers(),
		HTTPClient:  defaults.HTTPClient(),
		Credentials: credentials,
		EndpointResolver: aws.ResolveWithEndpoint(aws.Endpoint{
			URL:           "https://" + cfg.Endpoint,
			SigningRegion: cfg.Endpoint,
		}),
	}
	if cfg.Endpoint == "" {
		config.EndpointResolver = endpoints.NewDefaultResolver()
	}
	return config
}
