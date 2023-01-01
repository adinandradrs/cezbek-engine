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
	ClientAuthTTL             time.Duration
	OtpTTL                    time.Duration
	QueueNotificationEmailOtp *string
}

type OnboardManager interface {
	AuthenticateClient(inp *model.ClientAuthenticationRequest) (*model.ClientAuthenticationResponse, *model.BusinessError)
	AuthenticateOfficer(inp *model.OfficerAuthenticationRequest) (*model.OfficerAuthenticationResponse, *model.BusinessError)
	ValidateAuthOfficer(inp *model.OfficerValidationRequest) (*model.OfficerValidationResponse, *model.BusinessError)
}

func NewOnboard(o Onboard) OnboardManager {
	return &o
}

func (o Onboard) ciamAuthenticate(p model.Partner) (*model.CiamAuthenticationResponse, *model.BusinessError) {
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

func (o Onboard) decryptedSecret(secret []byte, salt string) (*string, *model.BusinessError) {
	d, ex := apps.Decrypt(secret, salt, o.Logger)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeUnauthorized,
			ErrorMessage: apps.ErrMsgUnauthorized,
		}
	}
	return &d, nil
}

func (o Onboard) AuthenticateClient(inp *model.ClientAuthenticationRequest) (*model.ClientAuthenticationResponse, *model.BusinessError) {
	p, ex := o.Dao.FindActiveByCodeAndApiKey(inp.Code, inp.ApiKey)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeUnauthorized,
			ErrorMessage: apps.ErrMsgUnauthorized,
		}
	}

	auth, bx := o.ciamAuthenticate(*p)
	if bx != nil {
		return nil, bx
	}

	resp := model.ClientAuthenticationResponse{
		Code:    p.Code.String,
		Company: p.Partner.String,
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
	o.Cacher.Set("CLIENTSESSION", p.Code.String, cache, o.ClientAuthTTL)

	return &resp, nil
}

func (o Onboard) AuthenticateOfficer(inp *model.OfficerAuthenticationRequest) (*model.OfficerAuthenticationResponse, *model.BusinessError) {
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
	ttl, _ := o.Cacher.Ttl("OTPB2B", p.Email.String)
	if ttl.Seconds() > 0 {
		trx, _ = o.Cacher.Get("OTPB2B", p.Email.String)
		return &model.OfficerAuthenticationResponse{
			RemainingSeconds: ttl.Seconds(),
			TransactionResponse: model.TransactionResponse{
				TransactionId:        strings.Split(trx, "#")[1],
				TransactionTimestamp: time.Now().Unix(),
			},
		}, nil
	}
	o.Cacher.Set("OTPB2B", p.Email.String, otp+"#"+trx, o.OtpTTL)
	cache, err := json.Marshal(p)
	if err != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	o.Cacher.Set("OTPB2B:"+trx, otp, cache, o.OtpTTL)
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

func (o Onboard) queueEmailOtp(otp string, p *model.Partner) *model.BusinessError {
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

func (o Onboard) ValidateAuthOfficer(inp *model.OfficerValidationRequest) (*model.OfficerValidationResponse, *model.BusinessError) {
	//TODO implement me
	panic("implement me")
}
