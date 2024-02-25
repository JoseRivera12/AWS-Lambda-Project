package s3aws

import (
	"bytes"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3API interface {
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

func RetrieveFileFromS3(s3Client S3API, bucketName, key string) ([]byte, error) {
	if s3Client == nil {
		return nil, errors.New("invalid AWS session provided")
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	result, err := s3Client.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	buffer := bytes.NewBuffer(nil)
	_, err = io.Copy(buffer, result.Body)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
