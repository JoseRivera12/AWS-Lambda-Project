package fileutils

import (
	"errors"
	"strings"

	"transaction-mailer/internal/s3aws"
)

type CustomerData struct {
	Name  string
	Email string
}

func NewCustomerData(s3Client s3aws.S3API, data []byte) (*CustomerData, error) {
	customerData := &CustomerData{}
	err := customerData.ReadDataFromFile(s3Client, data)
	return customerData, err
}

func (customer *CustomerData) ReadDataFromFile(s3Client s3aws.S3API, data []byte) error {
	/*
		File format
		customer name
		email
	*/
	lines := string(data)
	if lines == "" {
		if s3Client == nil {
			return errors.New("invalid AWS session provided")
		}
		txtFile, err := s3aws.RetrieveFileFromS3(s3Client, "stori-challenge-jose-rivera", "user.txt")
		if err != nil {
			return err
		}
		lines = string(txtFile)
	}
	userData := strings.Split(lines, "\n")
	if len(userData) != 2 {
		return errors.New("customers data invalid file format")
	}
	customer.Name = userData[0]
	customer.Email = userData[1]

	return nil
}
