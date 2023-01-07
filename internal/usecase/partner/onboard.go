package partner

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Onboard struct {
	Dao                       repository.PartnerPersister
	CiamWatcher               adaptor.CiamWatcher
	SqsAdapter                adaptor.SQSAdapter
	Cacher                    storage.Cacher
	Logger                    *zap.Logger
	AuthTTL                   time.Duration
	OtpTTL                    time.Duration
	QueueNotificationEmailOtp *string
	CDN                       *string
}

const kotp = "OTPB2B:"
const otpB2B = "OTPB2B"

type OnboardProvider interface {
	Authenticate(inp *model.OfficerAuthenticationRequest) (*model.OfficerAuthenticationResponse, *model.BusinessError)
	Validate(inp *model.OfficerValidationRequest) (*model.OfficerValidationResponse, *model.BusinessError)
}

func NewOnboard(o Onboard) OnboardProvider {
	return &o
}

func (o *Onboard) authenticate(p model.Partner) (*model.CiamAuthenticationResponse, *model.BusinessError) {
	secret, bx := o.decryptedSecret(p.Secret, p.Salt.String)
	if bx != nil {
		return nil, bx
	}

	auth, ex := o.CiamWatcher.Authenticate(model.CiamAuthenticationRequest{
		Username: p.Code.String,
		Secret:   *secret,
	})
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeUnauthorized,
			ErrorMessage: apps.ErrMsgUnauthorized,
		}
	}
	return auth, nil
}

func (o *Onboard) decryptedSecret(secret []byte, salt string) (*string, *model.BusinessError) {
	d, ex := apps.Decrypt(secret, salt, o.Logger)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeUnauthorized,
			ErrorMessage: apps.ErrMsgUnauthorized,
		}
	}
	return &d, nil
}

func (o *Onboard) Authenticate(inp *model.OfficerAuthenticationRequest) (*model.OfficerAuthenticationResponse, *model.BusinessError) {
	p, ex := o.Dao.FindActiveByEmail(inp.Email)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeUnauthorized,
			ErrorMessage: apps.ErrMsgUnauthorized,
		}
	}
	otp, err := apps.RandomOtp(6)
	if err != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}

	trx := apps.TransactionId(apps.DefaultTrxId)
	ttl, _ := o.Cacher.Ttl(otpB2B, p.Email.String)
	if ttl.Seconds() > 0 {
		trx, _ = o.Cacher.Get(otpB2B, p.Email.String)
		return &model.OfficerAuthenticationResponse{
			RemainingSeconds: ttl.Seconds(),
			TransactionResponse: model.TransactionResponse{
				TransactionId:        strings.Split(trx, "#")[1],
				TransactionTimestamp: time.Now().Unix(),
			},
		}, nil
	}
	cache, err := json.Marshal(p)
	if err != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	o.Cacher.Set(otpB2B, p.Email.String, otp+"#"+trx, o.OtpTTL)
	o.Cacher.Set(kotp+trx, otp, cache, o.OtpTTL)
	bx := o.queueEmailOtp(otp, p)
	if bx != nil {
		return nil, bx
	}
	return &model.OfficerAuthenticationResponse{
		RemainingSeconds: o.OtpTTL.Seconds(),
		TransactionResponse: model.TransactionResponse{
			TransactionId:        trx,
			TransactionTimestamp: time.Now().Unix(),
		},
	}, nil
}

func (o *Onboard) queueEmailOtp(otp string, p *model.Partner) *model.BusinessError {
	sbj, _ := o.Cacher.Hget("EMAIL_SUBJECT", "OTP")
	tmpl, _ := o.Cacher.Hget("EMAIL_TEMPLATE", "OTP")
	tmpl = strings.ReplaceAll(tmpl, "${otp}", otp)
	tmpl = strings.ReplaceAll(tmpl, "${partner}", p.Partner.String)
	tmpl = strings.ReplaceAll(tmpl, "\n", "")
	tmpl = strings.ReplaceAll(tmpl, "\t", "")
	msg, err := json.Marshal(model.SendEmailRequest{
		Content:     tmpl,
		Subject:     sbj,
		Destination: p.Email.String,
	})
	if err != nil {
		return &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	err = o.SqsAdapter.SendMessage(*o.QueueNotificationEmailOtp, string(msg))
	if err != nil {
		return &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	return nil
}

func (o *Onboard) Validate(inp *model.OfficerValidationRequest) (*model.OfficerValidationResponse, *model.BusinessError) {
	cp, ex := o.Cacher.Get(kotp+inp.TransactionId, inp.Otp)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussPartnerOTPInvalid,
			ErrorMessage: apps.ErrMsgBussPartnerOTPInvalid,
		}
	}
	var p model.Partner
	err := json.Unmarshal([]byte(cp), &p)
	if err != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}

	auth, bx := o.authenticate(p)
	if bx != nil {
		return nil, bx
	}
	resp := model.OfficerValidationResponse{
		Id:      p.Id,
		UrlLogo: *o.CDN + p.Logo.String,
		Company: p.Partner.String,
		Email:   p.Email.String,
		Msisdn:  p.Msisdn.String,
		Code:    p.Code.String,
		SessionResponse: model.SessionResponse{
			RefreshToken: auth.RefreshToken,
			Token:        auth.Token,
			AccessToken:  auth.AccessToken,
			Expired:      &auth.ExpiresIn,
		},
	}

	cache, err := json.Marshal(resp)
	if err != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	_ = o.Cacher.Delete(kotp+inp.TransactionId, inp.Otp)
	_ = o.Cacher.Delete(otpB2B, p.Email.String)
	o.Cacher.Set("B2BSESSION", p.Email.String, cache, o.AuthTTL)
	resp.Id = 0
	return &resp, nil
}
