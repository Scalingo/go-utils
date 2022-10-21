package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-utils/storage/storagemock"
	storageTypes "github.com/Scalingo/go-utils/storage/types"
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
		expectMock func(t *testing.T, m *storagemock.MockS3Client)
		err        string
	}{
		"it should make a HEAD request on the object": {
			expectMock: func(t *testing.T, m *storagemock.MockS3Client) {
				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(&s3.HeadObjectOutput{ContentLength: int64(10)}, nil)
			},
		},
		"it should retry if the first HEAD request return 404": {
			expectMock: func(t *testing.T, m *storagemock.MockS3Client) {
				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(nil, KeyNotFoundErr{})

				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(&s3.HeadObjectOutput{ContentLength: int64(10)}, nil)
			},
		},
		"it should fail if the max amount of retried is passed": {
			expectMock: func(t *testing.T, m *storagemock.MockS3Client) {
				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(nil, KeyNotFoundErr{}).Times(3)
			},
			err: "NotFound",
		},
	}
	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := storagemock.NewMockS3Client(ctrl)
			storage := &S3{
				cfg:      S3Config{Bucket: "bucket"},
				s3client: mock,
				retryPolicy: RetryPolicy{
					Attempts: 3, WaitDuration: 50 * time.Millisecond,
					MethodHandlers: map[BackendMethod][]string{SizeMethod: {NotFoundErrCode}},
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

func TestS3_Info(t *testing.T) {
	cases := map[string]struct {
		expectMock   func(t *testing.T, m *storagemock.MockS3Client, key string)
		key          string
		err          string
		expectedInfo storageTypes.Info
	}{
		"it should return ObjectNotFound error if the object does not exists": {
			expectMock: func(t *testing.T, m *storagemock.MockS3Client, key string) {
				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String(key),
				}).Return(nil, &smithy.GenericAPIError{
					Code:    "NotFound",
					Message: "Not Found",
				})
			},
			key:          "unknown_key",
			err:          (&ObjectNotFound{}).Error(),
			expectedInfo: storageTypes.Info{},
		},
		"it should return information about the object if the object exists": {
			expectMock: func(t *testing.T, m *storagemock.MockS3Client, key string) {
				contentType := "test"
				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String(key),
				}).Return(&s3.HeadObjectOutput{
					ContentType:   &contentType,
					ContentLength: 4,
					ETag:          aws.String("checksum"),
				}, nil)
			},
			err: "",
			expectedInfo: storageTypes.Info{
				ContentLength: 4,
				ContentType:   "test",
				Checksum:      "checksum",
			},
		},
		"it should fail if s3 does not respond": {
			expectMock: func(t *testing.T, m *storagemock.MockS3Client, key string) {
				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String(key),
				}).Return(nil, errors.New("there is an issue"))
			},
			err:          "there is an issue",
			expectedInfo: storageTypes.Info{},
		},
		"it should not fail if s3 return no Content-Type": {
			expectMock: func(t *testing.T, m *storagemock.MockS3Client, key string) {
				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String(key),
				}).Return(&s3.HeadObjectOutput{
					ContentType:   nil,
					ContentLength: 4,
					ETag:          aws.String("checksum"),
				}, nil)
			},
			expectedInfo: storageTypes.Info{
				ContentType:   "",
				ContentLength: 4,
				Checksum:      "checksum",
			},
		},
		"it should not fail if s3 return no ETag": {
			expectMock: func(t *testing.T, m *storagemock.MockS3Client, key string) {
				contentType := "test"
				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String(key),
				}).Return(&s3.HeadObjectOutput{
					ContentType:   &contentType,
					ContentLength: 4,
					ETag:          nil,
				}, nil)
			},
			expectedInfo: storageTypes.Info{
				ContentType:   "test",
				ContentLength: 4,
				Checksum:      "",
			},
		},
	}
	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := storagemock.NewMockS3Client(ctrl)
			storage := &S3{
				cfg:      S3Config{Bucket: "bucket"},
				s3client: mock,
				retryPolicy: RetryPolicy{
					Attempts:     3,
					WaitDuration: 50 * time.Millisecond,
				},
			}

			if c.key == "" {
				c.key = "key"
			}

			c.expectMock(t, mock, c.key)
			info, err := storage.Info(context.Background(), c.key)
			if c.err != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), c.err)
				return
			}
			require.Equal(t, c.expectedInfo, info)
		})
	}
}
