package job

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/adaptor"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransaction_SendInvoiceEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	sesAdapter, sqsAdapter := adaptor.NewMockSESAdapter(ctrl),
		adaptor.NewMockSQSAdapter(ctrl)
	q := "mock-queue"
	svc := NewTransaction(Transaction{
		Logger:                    logger,
		SqsAdapter:                sqsAdapter,
		SesAdapter:                sesAdapter,
		QueueNotificationEmailTrx: &q,
	})

	t.Run("should success", func(t *testing.T) {
		b := "html content - bla bla"
		h := "q-handler"
		sqsAdapter.EXPECT().GetMessages(q).Return(&sqs.Message{
			Body:          &b,
			ReceiptHandle: &h,
		})
		sesAdapter.EXPECT().SendEmail(gomock.Any()).Return(nil, nil)
		sqsAdapter.EXPECT().DeleteMessages(q, h)
		ex := svc.SendInvoiceEmail()
		assert.Nil(t, ex)
	})

	t.Run("should skip when no message in queue", func(t *testing.T) {
		sqsAdapter.EXPECT().GetMessages(q).Return(nil)
		ex := svc.SendInvoiceEmail()
		assert.Nil(t, ex)
	})

	t.Run("should return exception when failed to send email", func(t *testing.T) {
		b := "html content - bla bla"
		h := "q-handler"
		sqsAdapter.EXPECT().GetMessages(q).Return(&sqs.Message{
			Body:          &b,
			ReceiptHandle: &h,
		})
		sesAdapter.EXPECT().SendEmail(gomock.Any()).Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		ex := svc.SendInvoiceEmail()
		assert.NotNil(t, ex)
	})
}
