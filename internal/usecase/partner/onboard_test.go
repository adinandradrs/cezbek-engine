package partner

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func TestOnboard_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)

	authTTL, _ := time.ParseDuration("1s")
	otpTTL, _ := time.ParseDuration("1s")
	dao, ciamWatcher, sqsAdapter, cacher, q, cdn := repository.NewMockPartnerPersister(ctrl),
		adaptor.NewMockCiamWatcher(ctrl), adaptor.NewMockSQSAdapter(ctrl), storage.NewMockCacher(ctrl),
		"mock-queue", "https://cdn-mock.id"
	svc := NewOnboard(Onboard{
		Logger:                    logger,
		Cacher:                    cacher,
		SqsAdapter:                sqsAdapter,
		QueueNotificationEmailOtp: &q,
		CDN:                       &cdn,
		Dao:                       dao,
		CiamWatcher:               ciamWatcher,
		OtpTTL:                    otpTTL,
		AuthTTL:                   authTTL,
	})
	inp := &model.OfficerAuthenticationRequest{
		Email: "someone@email.net",
	}
	gpass, _ := apps.RandomPassword(12, 5, 3, logger)
	salt := apps.Hash("DUMMY-CODE" + ":" + uuid.NewString())
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
	t.Run("should success with new generated OTP", func(t *testing.T) {
		ttl, _ := time.ParseDuration("-5s")
		dao.EXPECT().FindActiveByEmail(inp.Email).Return(p, nil)
		cacher.EXPECT().Ttl("OTPB2B", p.Email.String).Return(ttl, nil)
		cacher.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
		cacher.EXPECT().Hget(gomock.Any(), gomock.Any()).Times(2).Return("subject", nil)
		sqsAdapter.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil)
		v, ex := svc.Authenticate(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})
	t.Run("should success with existing OTP", func(t *testing.T) {
		ttl, _ := time.ParseDuration("5s")
		dao.EXPECT().FindActiveByEmail(inp.Email).Return(p, nil)
		cacher.EXPECT().Ttl("OTPB2B", p.Email.String).Return(ttl, nil)
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return("1#1", nil)
		v, ex := svc.Authenticate(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})
	t.Run("should return exception on send message", func(t *testing.T) {
		ttl, _ := time.ParseDuration("-5s")
		dao.EXPECT().FindActiveByEmail(inp.Email).Return(p, nil)
		cacher.EXPECT().Ttl("OTPB2B", p.Email.String).Return(ttl, nil)
		cacher.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
		cacher.EXPECT().Hget(gomock.Any(), gomock.Any()).Times(2).Return("subject", nil)
		sqsAdapter.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(fmt.Errorf("something went wrong"))
		v, ex := svc.Authenticate(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}

func TestOnboard_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)

	authTTL, _ := time.ParseDuration("1s")
	otpTTL, _ := time.ParseDuration("1s")
	dao, ciamWatcher, sqsAdapter, cacher, q, cdn := repository.NewMockPartnerPersister(ctrl),
		adaptor.NewMockCiamWatcher(ctrl), adaptor.NewMockSQSAdapter(ctrl), storage.NewMockCacher(ctrl),
		"mock-queue", "https://cdn-mock.id"
	svc := NewOnboard(Onboard{
		Logger:                    logger,
		Cacher:                    cacher,
		SqsAdapter:                sqsAdapter,
		QueueNotificationEmailOtp: &q,
		CDN:                       &cdn,
		Dao:                       dao,
		CiamWatcher:               ciamWatcher,
		OtpTTL:                    otpTTL,
		AuthTTL:                   authTTL,
	})
	inp := &model.OfficerValidationRequest{
		TransactionId: "TX-001",
		Otp:           "123456",
	}
	gpass, _ := apps.RandomPassword(12, 5, 3, logger)
	salt := apps.Hash("CODE_A:" + uuid.NewString())
	secret, _ := apps.Encrypt(gpass, salt, logger)
	t.Run("should success", func(t *testing.T) {
		p := model.Partner{
			Id:      int64(1),
			Partner: sql.NullString{String: "Partner A", Valid: true},
			Code:    sql.NullString{String: "CODE_A", Valid: true},
			Msisdn:  sql.NullString{String: "628123456789", Valid: true},
			Email:   sql.NullString{String: "someone@email.net", Valid: true},
			Secret:  secret,
			Salt:    sql.NullString{String: salt, Valid: true},
		}
		cache, _ := json.Marshal(p)
		cacher.EXPECT().Get("OTPB2B:"+inp.TransactionId, inp.Otp).
			Return(string(cache), nil)
		cacher.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(2)
		cacher.EXPECT().Set("B2BSESSION", p.Email.String, gomock.Any(), authTTL)
		ciamWatcher.EXPECT().Authenticate(gomock.Any()).Return(&model.CiamAuthenticationResponse{
			Token:        "token-abc",
			ExpiresIn:    int64(1),
			AccessToken:  "access-token-abc",
			RefreshToken: "ref-token-abc",
		}, nil)
		v, ex := svc.Validate(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})
	t.Run("should return exception on wrong OTP", func(t *testing.T) {
		cacher.EXPECT().Get("OTPB2B:"+inp.TransactionId, inp.Otp).
			Return("", &model.TechnicalError{
				Exception: "something went wrong",
				Occurred:  time.Now().Unix(),
				Ticket:    "ERR-001",
			})
		v, ex := svc.Validate(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
	t.Run("should return exception on CIAM failure", func(t *testing.T) {
		p := model.Partner{
			Id:      int64(1),
			Partner: sql.NullString{String: "Partner A", Valid: true},
			Code:    sql.NullString{String: "CODE_A", Valid: true},
			Msisdn:  sql.NullString{String: "628123456789", Valid: true},
			Email:   sql.NullString{String: "someone@email.net", Valid: true},
			Secret:  secret,
			Salt:    sql.NullString{String: salt, Valid: true},
		}
		cache, _ := json.Marshal(p)
		cacher.EXPECT().Get("OTPB2B:"+inp.TransactionId, inp.Otp).
			Return(string(cache), nil)
		ciamWatcher.EXPECT().Authenticate(gomock.Any()).Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		v, ex := svc.Validate(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}
