package storage

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Scalingo/go-utils/storage/s3mock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock 404 NotFound error from AWS API
type KeyNotFoundErr struct{}

func (err KeyNotFoundErr) Code() string {
	return NotFoundErrCode
}

func (err KeyNotFoundErr) Error() string {
	return err.Code()
}

func (err KeyNotFoundErr) Message() string {
	return err.Error()
}

func (err KeyNotFoundErr) OrigErr() error {
	return err
}

func TestS3_Size(t *testing.T) {
	cases := map[string]struct {
		expectMock func(t *testing.T, m *s3mock.MockS3Client)
		err        string
	}{
		"it should make a HEAD request on the object": {
			expectMock: func(t *testing.T, m *s3mock.MockS3Client) {
				m.EXPECT().HeadObjectRequest(&s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(s3.HeadObjectRequest{Request: &aws.Request{
					// Mandatory to create an empty request, otherwise it panics
					HTTPRequest: new(http.Request),
					Data:        &s3.HeadObjectOutput{ContentLength: aws.Int64(10)},
				}})
			},
		},
		"it should retry if the first HEAD request return 404": {
			expectMock: func(t *testing.T, m *s3mock.MockS3Client) {
				m.EXPECT().HeadObjectRequest(&s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(s3.HeadObjectRequest{Request: &aws.Request{
					HTTPRequest: new(http.Request),
					Error:       KeyNotFoundErr{},
				}})

				m.EXPECT().HeadObjectRequest(&s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(s3.HeadObjectRequest{Request: &aws.Request{
					// Mandatory to create an empty request, otherwise it panics
					HTTPRequest: new(http.Request),
					Data:        &s3.HeadObjectOutput{ContentLength: aws.Int64(10)},
				}})
			},
		},
		"it should fail if the max amount of retried is passed": {
			expectMock: func(t *testing.T, m *s3mock.MockS3Client) {
				m.EXPECT().HeadObjectRequest(&s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(s3.HeadObjectRequest{Request: &aws.Request{
					HTTPRequest: new(http.Request),
					Error:       KeyNotFoundErr{},
				}}).Times(3)
			},
			err: "NotFound",
		},
	}
	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := s3mock.NewMockS3Client(ctrl)
			storage := &S3{
				cfg:      S3Config{Bucket: "bucket"},
				s3client: mock,
				retryPolicy: RetryPolicy{
					Attempts: 3, WaitDuration: 50 * time.Millisecond,
					MethodHandlers: map[BackendMethod][]string{SizeMethod: []string{NotFoundErrCode}},
				},
			}

			c.expectMock(t, mock)
			_, err := storage.Size(context.Background(), "/key")
			if c.err != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), c.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
