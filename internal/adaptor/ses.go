package adaptor

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"go.uber.org/zap"
)

type SES struct {
	Sender *string
	SES    *sesv2.SESV2
	Logger *zap.Logger
}

type SESAdapter interface {
	SendEmail(m model.SendEmailRequest) (*model.TransactionResponse, *model.TechnicalError)
}

func NewSES(s SES) SESAdapter {
	return &s
}

func (s *SES) SendEmail(m model.SendEmailRequest) (*model.TransactionResponse, *model.TechnicalError) {
	utf8 := "utf-8"
	out, err := s.SES.SendEmail(&sesv2.SendEmailInput{
		Destination: &sesv2.Destination{
			ToAddresses: []*string{
				aws.String(m.Destination),
			},
		},
		Content: &sesv2.EmailContent{
			Simple: &sesv2.Message{
				Body: &sesv2.Body{
					Html: &sesv2.Content{
						Charset: aws.String(utf8),
						Data:    aws.String(m.Content),
					},
				},
				Subject: &sesv2.Content{
					Charset: aws.String(utf8),
					Data:    aws.String(m.Subject),
				},
			},
		},
		FromEmailAddress: s.Sender,
	})
	if err != nil {
		return nil, apps.Exception("failed to send email",
			err, zap.String("destination", m.Destination), s.Logger)
	}
	trx := apps.Transaction(*out.MessageId)
	return &trx, nil
}
