package mailer

import (
	"bytes"
	"errors"
	"html/template"
	"os"

	"transaction-mailer/internal/s3aws"
	"transaction-mailer/internal/transactions"

	"github.com/aws/aws-sdk-go/service/ses"
)

type SESAPI interface {
	SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error)
}

type EmailData struct {
	Transactions *transactions.TransactionsData
	Username     string
}

func SentTransactionsEmail(s3Client s3aws.S3API, sesClient SESAPI, emailData EmailData, recipient string) error {
	// templatePath := filepath.Join("internal", "email", "template.html")
	// emailTemplateHTML, err := os.ReadFile(templatePath)
	if s3Client == nil {
		return errors.New("non S3 session provided")
	}
	if sesClient == nil {
		return errors.New("non SES session provided")
	}
	templateData, err := s3aws.RetrieveFileFromS3(s3Client, "stori-challenge-jose-rivera", "template.html")
	if err != nil {
		return err
	}
	template.ParseFiles()
	tmpl, err := template.New("emailTemplate").Funcs(template.FuncMap{
		"MonthMapper":    monthMapper,
		"MoneyFormatter": formatMoney,
		"GetYear":        getYear,
	}).Parse(string(templateData))
	if err != nil {
		return err
	}
	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, emailData); err != nil {
		return err
	}

	return SendEmail(sesClient, recipient, emailBody.String())
}

func SendEmail(sesClient SESAPI, recipient, emailBody string) error {
	sender := os.Getenv("SES_EMAIL")
	subject := "Transaction Summary Report"
	bodyData := &ses.Content{Data: &emailBody}
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{&recipient},
		},
		Message: &ses.Message{
			Body:    &ses.Body{Html: bodyData},
			Subject: &ses.Content{Data: &subject},
		},
		Source: &sender,
	}

	_, err := sesClient.SendEmail(input)
	if err != nil {
		return err
	}
	return nil
}
