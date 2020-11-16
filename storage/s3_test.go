package storage

import (
	"context"
	"testing"
	"time"

	"github.com/Scalingo/go-utils/storage/storagemock"
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
		expectMock func(t *testing.T, m *storagemock.MockS3Client)
		err        string
	}{
		"it should make a HEAD request on the object": {
			expectMock: func(t *testing.T, m *storagemock.MockS3Client) {
				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(&s3.HeadObjectOutput{ContentLength: aws.Int64(10)}, nil)
			},
		},
		"it should retry if the first HEAD request return 404": {
			expectMock: func(t *testing.T, m *storagemock.MockS3Client) {
				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(nil, KeyNotFoundErr{})

				m.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
					Bucket: aws.String("bucket"), Key: aws.String("key"),
				}).Return(&s3.HeadObjectOutput{ContentLength: aws.Int64(10)}, nil)
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
