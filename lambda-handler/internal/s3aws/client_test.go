package s3aws

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testCase struct {
	name                string
	bucketName          string
	key                 string
	mockGetObjectOutput *s3.GetObjectOutput
	mockGetObjectError  error
	s3Client            S3API
	expectedData        []byte
	expectedError       error
}

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func TestRetrieveFileFromS3(t *testing.T) {
	testCases := []testCase{
		{
			name:                "Success",
			bucketName:          "test-bucket",
			key:                 "test-key",
			mockGetObjectOutput: &s3.GetObjectOutput{Body: io.NopCloser(bytes.NewReader([]byte("test-data")))},
			mockGetObjectError:  nil,
			s3Client:            new(MockS3Client),
			expectedData:        []byte("test-data"),
			expectedError:       nil,
		},
		{
			name:               "Invalid AWS session",
			s3Client:           nil,
			mockGetObjectError: errors.New("invalid AWS session provided"),
			expectedError:      errors.New("invalid AWS session provided"),
		},
		{
			name:                "S3 GetObject error",
			bucketName:          "test-bucket",
			key:                 "test-key",
			mockGetObjectOutput: nil,
			s3Client:            new(MockS3Client),
			mockGetObjectError:  errors.New("S3 error"),
			expectedData:        nil,
			expectedError:       errors.New("S3 error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockS3Client, ok := tc.s3Client.(*MockS3Client)
			if ok {
				mockS3Client.On("GetObject", mock.Anything).Return(tc.mockGetObjectOutput, tc.mockGetObjectError)
			}
			data, err := RetrieveFileFromS3(tc.s3Client, tc.bucketName, tc.key)

			assert.Equal(t, tc.expectedData, data)
			assert.Equal(t, tc.expectedError, err)

		})
	}
}
