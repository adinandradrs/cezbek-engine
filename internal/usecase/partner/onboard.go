package partner

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"go.uber.org/zap"
	"time"
)

type Onboard struct {
	Dao           repository.PartnerPersister
	CiamWatcher   adaptor.CiamWatcher
	Cacher        storage.Cacher
	Logger        *zap.Logger
	ClientAuthTTL time.Duration
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
	json, err := json.Marshal(resp)
	if err != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	o.Cacher.Set("CLIENTSESSION", p.Code.String, json, o.ClientAuthTTL)

	return &resp, nil
}

func (o Onboard) AuthenticateOfficer(inp *model.OfficerAuthenticationRequest) (*model.OfficerAuthenticationResponse, *model.BusinessError) {
	//TODO implement me
	panic("implement me")
}

func (o Onboard) ValidateAuthOfficer(inp *model.OfficerValidationRequest) (*model.OfficerValidationResponse, *model.BusinessError) {
	//TODO implement me
	panic("implement me")
}
