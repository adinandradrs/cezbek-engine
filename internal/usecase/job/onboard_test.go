package job

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/adaptor"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOnboard_SendOtpEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	sesAdapter, sqsAdapter, q := adaptor.NewMockSESAdapter(ctrl),
		adaptor.NewMockSQSAdapter(ctrl), "mock-queue"
	manager := NewOnboard(Onboard{
		Logger:                    logger,
		QueueNotificationEmailOtp: &q,
		SesAdapter:                sesAdapter,
		SqsAdapter:                sqsAdapter,
	})
	t.Run("should success", func(t *testing.T) {
		msg := sqs.Message{
			ReceiptHandle: aws.String("mock-handler-msg"),
			Body:          aws.String("mock-email-body"),
		}
		inp := model.SendEmailRequest{}
		_ = json.Unmarshal([]byte(*msg.Body), &inp)
		sqsAdapter.EXPECT().GetMessages(q).Return(&msg)
		sesAdapter.EXPECT().SendEmail(inp).Return(
			&model.TransactionResponse{
				TransactionId:        "trx-001",
				TransactionTimestamp: time.Now().Unix(),
			}, nil)
		sqsAdapter.EXPECT().DeleteMessages(q, gomock.Any())
		ex := manager.SendOtpEmail()
		assert.Nil(t, ex)
	})
	t.Run("should return exception on send email", func(t *testing.T) {
		msg := sqs.Message{
			ReceiptHandle: aws.String("mock-handler-msg"),
			Body:          aws.String("mock-email-body"),
		}
		inp := model.SendEmailRequest{}
		_ = json.Unmarshal([]byte(*msg.Body), &inp)
		sqsAdapter.EXPECT().GetMessages(q).Return(&msg)
		sesAdapter.EXPECT().SendEmail(inp).Return(
			nil, &model.TechnicalError{
				Exception: "something went wrong",
				Occurred:  time.Now().Unix(),
				Ticket:    "err-001",
			})
		ex := manager.SendOtpEmail()
		assert.NotNil(t, ex)
	})
}
