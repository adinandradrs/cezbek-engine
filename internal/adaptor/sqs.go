package adaptor

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQS struct {
	SQS *sqs.SQS
}

type SQSAdapter interface {
	GetMessages(q string) *sqs.Message
	DeleteMessages(q string, h string) *sqs.DeleteMessageOutput
	SendMessage(q string, h string) error
}

func NewSQS(c SQS) SQSAdapter {
	return &c
}

func (s *SQS) GetMessages(t string) *sqs.Message {
	msgs, _ := s.SQS.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &t,
		MaxNumberOfMessages: aws.Int64(1),
	})
	if len(msgs.Messages) > 0 {
		return msgs.Messages[0]
	}

	return nil
}

func (s *SQS) DeleteMessages(q string, h string) *sqs.DeleteMessageOutput {
	msgs, _ := s.SQS.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &q,
		ReceiptHandle: &h,
	})
	return msgs
}

func (s *SQS) SendMessage(q string, h string) error {
	_, err := s.SQS.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &q,
		MessageBody: &h,
	})
	if err != nil {
		return err
	}
	return nil
}
