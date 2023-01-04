package client

import (
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"time"
)

type Onboard struct {
	Dao         repository.PartnerPersister
	CiamWatcher adaptor.CiamWatcher
	Cacher      storage.Cacher
	Logger      *zap.Logger
	AuthTTL     time.Duration
}

type OnboardProvider interface {
	Authenticate(inp *model.ClientAuthenticationRequest) (*model.ClientAuthenticationResponse, *model.BusinessError)
}

func NewOnboard(o Onboard) OnboardProvider {
	return &o
}

func (o *Onboard) Authenticate(inp *model.ClientAuthenticationRequest) (*model.ClientAuthenticationResponse, *model.BusinessError) {
	p, ex := o.Dao.FindActiveByCodeAndApiKey(inp.Code, inp.ApiKey)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeUnauthorized,
			ErrorMessage: apps.ErrMsgUnauthorized,
		}
	}

	d, ex := apps.Decrypt(p.Secret, p.Salt.String, o.Logger)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeUnauthorized,
			ErrorMessage: apps.ErrMsgUnauthorized,
		}
	}

	auth, ex := o.CiamWatcher.Authenticate(model.CiamAuthenticationRequest{
		Username: p.Code.String,
		Secret:   d,
	})
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeUnauthorized,
			ErrorMessage: apps.ErrMsgUnauthorized,
		}
	}

	resp := model.ClientAuthenticationResponse{
		Id:      &p.Id,
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
	o.Cacher.Set("CLIENTSESSION", p.Code.String, cache, o.AuthTTL)
	resp.Id = nil
	resp.AccessToken = ""
	return &resp, nil
}
