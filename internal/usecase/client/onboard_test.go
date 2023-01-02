package client

import (
	"database/sql"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/adaptor"
	"github.com/adinandradrs/cezbek-engine/mock/repository"
	"github.com/adinandradrs/cezbek-engine/mock/storage"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOnboard_AuthenticateClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)

	clientAuthTTL, _ := time.ParseDuration("1s")
	dao, ciamWatcher, cacher := repository.NewMockPartnerPersister(ctrl),
		adaptor.NewMockCiamWatcher(ctrl), storage.NewMockCacher(ctrl)
	manager := NewOnboard(Onboard{
		Logger:      logger,
		Cacher:      cacher,
		Dao:         dao,
		CiamWatcher: ciamWatcher,
		AuthTTL:     clientAuthTTL,
	})
	inp := &model.ClientAuthenticationRequest{
		Code:   "DUMMY-CODE",
		ApiKey: "api-123-456",
	}

	gpass, _ := apps.RandomPassword(12, 5, 3, logger)
	salt := apps.Hash(inp.Code + ":" + uuid.NewString())
	secret, _ := apps.Encrypt(gpass, salt, logger)
	p := &model.Partner{
		Id:      int64(1),
		Partner: sql.NullString{String: "PT. Partner A", Valid: true},
		Code:    sql.NullString{String: "DUMMY-CODE", Valid: true},
		ApiKey:  sql.NullString{String: "api-123-456", Valid: true},
		Salt:    sql.NullString{String: salt, Valid: true},
		Secret:  secret,
		Email:   sql.NullString{String: "someone@email.id", Valid: true},
		Msisdn:  sql.NullString{String: "628123456789", Valid: true},
	}
	t.Run("should success", func(t *testing.T) {
		dao.EXPECT().FindActiveByCodeAndApiKey(inp.Code, inp.ApiKey).
			Return(p, nil)
		ciamWatcher.EXPECT().Authenticate(gomock.Any()).Return(&model.CiamAuthenticationResponse{
			Token:        "token-abc",
			ExpiresIn:    int64(1),
			AccessToken:  "access-token-abc",
			RefreshToken: "ref-token-abc",
		}, nil)
		cacher.EXPECT().Set("CLIENTSESSION", p.Code.String, gomock.Any(), clientAuthTTL)
		v, ex := manager.Authenticate(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})
	t.Run("should return exception on failed to decrypt", func(t *testing.T) {
		invPartner := *p
		invPartner.Salt = sql.NullString{String: "invalid-salt", Valid: true}
		dao.EXPECT().FindActiveByCodeAndApiKey(inp.Code, inp.ApiKey).
			Return(&invPartner, nil)
		v, ex := manager.Authenticate(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
	t.Run("should return exception on ciam to sign-in", func(t *testing.T) {
		dao.EXPECT().FindActiveByCodeAndApiKey(inp.Code, inp.ApiKey).
			Return(p, nil)
		ciamWatcher.EXPECT().Authenticate(gomock.Any()).Return(nil,
			&model.TechnicalError{
				Exception: "something went wrong",
				Occurred:  time.Now().Unix(),
				Ticket:    "ERR-001",
			})
		v, ex := manager.Authenticate(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
	t.Run("should return exception on dao to find data", func(t *testing.T) {
		dao.EXPECT().FindActiveByCodeAndApiKey(inp.Code, inp.ApiKey).
			Return(nil, &model.TechnicalError{
				Exception: "something went wrong",
				Occurred:  time.Now().Unix(),
				Ticket:    "ERR-001",
			})
		v, ex := manager.Authenticate(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}
