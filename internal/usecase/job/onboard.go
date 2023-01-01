package job

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"go.uber.org/zap"
)

type Onboard struct {
	SesAdapter                adaptor.SESAdapter
	SqsAdapter                adaptor.SQSAdapter
	Logger                    *zap.Logger
	QueueNotificationEmailOtp *string
}

type OnboardManager interface {
	SendOtpEmail() *model.BusinessError
}

func NewOnboard(o Onboard) OnboardManager {
	return &o
}

func (o Onboard) SendOtpEmail() *model.BusinessError {
	msg := o.SqsAdapter.GetMessages(*o.QueueNotificationEmailOtp)
	if msg != nil {
		inp := model.SendEmailRequest{}
		_ = json.Unmarshal([]byte(*msg.Body), &inp)
		o.Logger.Info("send email check payload", zap.Any("input", inp))
		tx, ex := o.SesAdapter.SendEmail(inp)
		if ex != nil {
			o.Logger.Error("SES send OTP to email failed", zap.String("email", inp.Destination), zap.Any("ex", ex))
			return &model.BusinessError{
				ErrorCode:    apps.ErrCodeSomethingWrong,
				ErrorMessage: apps.ErrMsgSomethingWrong,
			}
		}
		o.SqsAdapter.DeleteMessages(*o.QueueNotificationEmailOtp, *msg.ReceiptHandle)
		o.Logger.Info("SES send OTP to email success", zap.String("email", inp.Destination), zap.Any("tx", tx))
	}
	return nil
}
