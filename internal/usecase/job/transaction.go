package job

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"go.uber.org/zap"
)

type Transaction struct {
	SesAdapter                adaptor.SESAdapter
	SqsAdapter                adaptor.SQSAdapter
	Logger                    *zap.Logger
	QueueNotificationEmailTrx *string
}

type TransactionWatcher interface {
	SendInvoiceEmail() *model.BusinessError
}

func NewTransaction(t Transaction) TransactionWatcher {
	return &t
}

func (t *Transaction) SendInvoiceEmail() *model.BusinessError {
	msg := t.SqsAdapter.GetMessages(*t.QueueNotificationEmailTrx)
	if msg != nil {
		inp := model.SendEmailRequest{}
		_ = json.Unmarshal([]byte(*msg.Body), &inp)
		t.Logger.Info("send email check payload", zap.Any("input", inp))
		tx, ex := t.SesAdapter.SendEmail(inp)
		if ex != nil {
			t.Logger.Error("SES send invoice to email failed", zap.String("email", inp.Destination), zap.Any("ex", ex))
			return &model.BusinessError{
				ErrorCode:    apps.ErrCodeSomethingWrong,
				ErrorMessage: apps.ErrMsgSomethingWrong,
			}
		}
		t.SqsAdapter.DeleteMessages(*t.QueueNotificationEmailTrx, *msg.ReceiptHandle)
		t.Logger.Info("SES send invoice to email success", zap.String("email", inp.Destination), zap.Any("tx", tx))
	}
	return nil
}
