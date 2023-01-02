package partner

import (
	"database/sql"
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

func TestOnboard_AuthenticateOfficer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)

	authTTL, _ := time.ParseDuration("1s")
	otpTTL, _ := time.ParseDuration("1s")
	dao, ciamWatcher, sqsAdapter, cacher, q, cdn := repository.NewMockPartnerPersister(ctrl),
		adaptor.NewMockCiamWatcher(ctrl), adaptor.NewMockSQSAdapter(ctrl), storage.NewMockCacher(ctrl),
		"mock-queue", "https://cdn-mock.id"
	manager := NewOnboard(Onboard{
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
		v, ex := manager.Authenticate(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})
	t.Run("should success with existing OTP", func(t *testing.T) {
		ttl, _ := time.ParseDuration("5s")
		dao.EXPECT().FindActiveByEmail(inp.Email).Return(p, nil)
		cacher.EXPECT().Ttl("OTPB2B", p.Email.String).Return(ttl, nil)
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return("1#1", nil)
		v, ex := manager.Authenticate(inp)
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
		v, ex := manager.Authenticate(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}
