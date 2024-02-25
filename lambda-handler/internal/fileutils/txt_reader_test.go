package fileutils

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"transaction-mailer/internal/s3aws"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func TestNewCustomerData(t *testing.T) {
	type testData struct {
		name               string
		data               []byte
		s3Client           s3aws.S3API
		expectedErr        error
		mockGetObjectError error
		expectedName       string
		expectedEmail      string
	}

	tests := []testData{
		{
			name:          "Success",
			data:          []byte("Jose Rivera\njoserivera@example.com"),
			expectedErr:   nil,
			s3Client:      nil,
			expectedName:  "Jose Rivera",
			expectedEmail: "joserivera@example.com",
		},
		{
			name:          "Success retrieved default from s3",
			data:          []byte(""),
			s3Client:      new(MockS3Client),
			expectedErr:   nil,
			expectedName:  "Jose Rivera",
			expectedEmail: "joserivera@example.com",
		},
		{
			name:        "Invalid AWS session",
			data:        []byte(""),
			expectedErr: errors.New("invalid AWS session provided"),
		},
		{
			name:               "S3 error",
			data:               []byte(""),
			s3Client:           new(MockS3Client),
			mockGetObjectError: errors.New("S3 error"),
			expectedErr:        errors.New("S3 error"),
		},
		{
			name:        "Invalid file format",
			data:        []byte("invalid format"),
			expectedErr: errors.New("customers data invalid file format"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockS3Client, ok := test.s3Client.(*MockS3Client)
			if ok {
				mockS3Client.On("GetObject", mock.Anything).Return(&s3.GetObjectOutput{Body: io.NopCloser(bytes.NewReader([]byte("Jose Rivera\njoserivera@example.com")))}, test.mockGetObjectError)
			}
			customerData, err := NewCustomerData(test.s3Client, test.data)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedName, customerData.Name)
			assert.Equal(t, test.expectedEmail, customerData.Email)

		})
	}
}
