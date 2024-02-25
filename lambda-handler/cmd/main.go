package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"transaction-mailer/internal/fileutils"
	"transaction-mailer/internal/mailer"
	"transaction-mailer/internal/s3aws"
	"transaction-mailer/internal/store"
	"transaction-mailer/internal/transactions"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func handle(ctx context.Context, event events.S3Event) error {
	dbContext := context.Background()
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	s3Client := s3.New(session)
	sesClient := ses.New(session)
	user_data_txt := "user.txt"
	connectionString := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=require",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
	)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}
	store.InitDB(db)
	for _, record := range event.Records {
		var txtData []byte
		// Only triggers for CSV files
		if filepath.Ext(record.S3.Object.Key) != ".csv" {
			return nil
		}

		// Read csv file
		csvData, err := s3aws.RetrieveFileFromS3(s3Client, record.S3.Bucket.Name, record.S3.Object.Key)
		if err != nil {
			return err
		}

		// Read customer information
		txtData, err = s3aws.RetrieveFileFromS3(s3Client, record.S3.Bucket.Name, user_data_txt)
		if err != nil {
			log.Println("No user file found")
		}

		userData, err := fileutils.NewCustomerData(s3Client, txtData)
		if err != nil {
			return err
		}

		// Parse CSV file
		records, err := fileutils.GetCSVRows(csvData)
		if err != nil {
			return err
		}
		transaction := transactions.NewTransactionsData(records)
		// Send email with transactions information
		err = mailer.SentTransactionsEmail(
			s3Client,
			sesClient,
			mailer.EmailData{
				Transactions: transaction,
				Username:     userData.Name,
			},
			userData.Email,
		)
		if err != nil {
			return err
		}

		// Save user transactions in the DB.
		user := &store.User{Name: userData.Name, Email: userData.Email}
		err = store.CreateTransactions(db, dbContext, records, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(handle)
}
