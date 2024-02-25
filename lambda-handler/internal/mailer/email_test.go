package mailer

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"transaction-mailer/internal/s3aws"
	"transaction-mailer/internal/transactions"
)

type mockS3Client struct {
	mock.Mock
}

func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

type mockSESClient struct {
	mock.Mock
}

func (m *mockSESClient) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*ses.SendEmailOutput), args.Error(1)
}

type testCase struct {
	name                string
	emailData           EmailData
	recipient           string
	expectedErr         error
	s3Client            s3aws.S3API
	sesClient           SESAPI
	mockGetObjectOutput *s3.GetObjectOutput
	mockSendEmailOutput *ses.SendEmailOutput
	mockSendEmailError  error
	mockGetObjectError  error
	expectTemplateError bool
}

func TestSentTransactionsEmail(t *testing.T) {
	testCases := []testCase{
		{
			name:                "Success",
			s3Client:            new(mockS3Client),
			sesClient:           new(mockSESClient),
			mockGetObjectOutput: &s3.GetObjectOutput{Body: io.NopCloser(bytes.NewReader([]byte("<html></html>")))},
			mockSendEmailOutput: &ses.SendEmailOutput{},
			emailData:           EmailData{Transactions: &transactions.TransactionsData{}, Username: "testuser"},
			recipient:           "joserivera@example.com",
		},
		{
			name:                "Error retrieving template from S3",
			emailData:           EmailData{},
			mockGetObjectOutput: &s3.GetObjectOutput{},
			s3Client:            new(mockS3Client),
			sesClient:           new(mockSESClient),
			mockSendEmailOutput: &ses.SendEmailOutput{},
			mockGetObjectError:  errors.New("S3 error"),
			recipient:           "joserivera@example.com",
			expectedErr:         errors.New("S3 error"),
			mockSendEmailError:  nil,
		},
		{
			name:                "Error parsing template",
			emailData:           EmailData{},
			mockGetObjectOutput: &s3.GetObjectOutput{Body: io.NopCloser(bytes.NewReader([]byte("Invalid template")))},
			mockSendEmailOutput: &ses.SendEmailOutput{},
			s3Client:            new(mockS3Client),
			sesClient:           new(mockSESClient),
			recipient:           "joserivera@example.com",
		},
		{
			name:                "Error executing template",
			s3Client:            new(mockS3Client),
			sesClient:           new(mockSESClient),
			mockGetObjectOutput: &s3.GetObjectOutput{Body: io.NopCloser(bytes.NewReader([]byte("{{.InvalidField}}")))},
			mockSendEmailOutput: &ses.SendEmailOutput{},
			emailData:           EmailData{},
			recipient:           "joserivera@example.com",
			expectTemplateError: true,
		},
		{
			name:                "Empty email data",
			s3Client:            new(mockS3Client),
			sesClient:           new(mockSESClient),
			mockGetObjectOutput: &s3.GetObjectOutput{Body: io.NopCloser(bytes.NewReader([]byte("<html></html>")))},
			mockSendEmailOutput: &ses.SendEmailOutput{},
			emailData:           EmailData{},
			recipient:           "joserivera@example.com",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockS3Client, ok := tc.s3Client.(*mockS3Client)
			if ok {
				mockS3Client.On("GetObject", mock.Anything).Return(tc.mockGetObjectOutput, tc.mockGetObjectError)
			}

			mockSESClient, ok := tc.sesClient.(*mockSESClient)
			if ok {
				mockSESClient.On("SendEmail", mock.Anything).Return(tc.mockSendEmailOutput, tc.mockSendEmailError)
			}

			err := SentTransactionsEmail(tc.s3Client, tc.sesClient, tc.emailData, tc.recipient)

			if tc.mockGetObjectError != nil || tc.expectTemplateError == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSendEmail(t *testing.T) {
	testCases := []struct {
		name                string
		recipient           string
		emailBody           string
		expectedErr         bool
		mockSendEmailError  error
		sesClient           SESAPI
		mockSendEmailOutput *ses.SendEmailOutput
	}{
		{
			name:        "Success",
			recipient:   "joserivera@example.com",
			emailBody:   "<html></html>",
			sesClient:   new(mockSESClient),
			expectedErr: false,
		},
		{
			name:               "SES Error",
			recipient:          "test@example.com",
			emailBody:          "<html></html>",
			sesClient:          new(mockSESClient),
			mockSendEmailError: errors.New("SES error"),
			expectedErr:        true,
		},
		{
			name:               "Empty Body",
			recipient:          "test@example.com",
			emailBody:          "",
			sesClient:          new(mockSESClient),
			expectedErr:        true,
			mockSendEmailError: errors.New("Empty body"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockSESClient, ok := testCase.sesClient.(*mockSESClient)
			if ok {
				mockSESClient.On("SendEmail", mock.Anything).Return(testCase.mockSendEmailOutput, testCase.mockSendEmailError)
			}

			err := SendEmail(mockSESClient, testCase.recipient, testCase.emailBody)
			if testCase.expectedErr {
				assert.Error(t, err)
			}

		})
	}
}
