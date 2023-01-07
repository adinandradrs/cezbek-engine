package handler

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/handler/middleware"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/adaptor"
	"github.com/adinandradrs/cezbek-engine/mock/storage"
	"github.com/adinandradrs/cezbek-engine/mock/usecase/partner"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestPartnerTransactionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	ciamPartner := adaptor.NewMockCiamWatcher(ctrl)
	cacher := storage.NewMockCacher(ctrl)
	transactionProvider := partner.NewMockTransactionProvider(ctrl)
	jwtAuthenticator := middleware.NewJwtAuthenticator(&middleware.JwtAuthenticator{
		Logger:      logger,
		CiamPartner: ciamPartner,
		Cacher:      cacher,
	})
	jwtAuthPartnerFilter := jwtAuthenticator.PartnerFilter()

	api := fiber.New()
	partnerTransactions := api.Group("/api/partner/v1/transactions")
	PartnerTransactionHandler(partnerTransactions, PartnerTransaction{
		TransactionProvider: transactionProvider,
		PartnerFilter:       jwtAuthPartnerFilter,
	})
	jwtInfo := map[string]interface{}{
		"email":            "someone@email.net",
		"cognito:username": "someone",
	}

	id := int64(1)
	cauth := model.ClientAuthenticationResponse{
		Id:      &id,
		Code:    "CORP_A",
		Company: "Company A",
	}
	c, _ := json.Marshal(cauth)

	t.Run("should return 200 success to search transaction", func(t *testing.T) {
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(c), nil)
		ciamPartner.EXPECT().JwtInfo(gomock.Any()).Return(jwtInfo, nil)
		transactionProvider.EXPECT().Search(gomock.Any()).Return(&model.PartnerTransactionSearchResponse{
			Transactions: []model.PartnerTransactionProjection{
				{
					Transaction: decimal.NewFromInt(385000),
					WalletCode:  "WCODE_A",
					Msisdn:      "628118770510",
					Id:          1,
					Reward:      decimal.NewFromInt(5000),
					Email:       "someone@email.net",
					Cashback:    decimal.NewFromInt(12500),
					Qty:         2,
				},
			},
		}, nil)
		req := httptest.NewRequest(fiber.MethodGet, "/api/partner/v1/transactions?limit=5&start=0", nil)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(fiber.HeaderAuthorization, "Bearer *secret*")
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelEBizKezbek)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.NotNil(t, m.Data)
	})

	t.Run("should return 200 failed to search transaction", func(t *testing.T) {
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(c), nil)
		ciamPartner.EXPECT().JwtInfo(gomock.Any()).Return(jwtInfo, nil)
		transactionProvider.EXPECT().Search(gomock.Any()).Return(nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeNotFound,
			ErrorMessage: apps.ErrMsgNotFound,
		})
		req := httptest.NewRequest(fiber.MethodGet, "/api/partner/v1/transactions?limit=5&start=0", nil)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(fiber.HeaderAuthorization, "Bearer *secret*")
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelEBizKezbek)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.Nil(t, m.Data)
		assert.Equal(t, apps.ErrCodeNotFound, m.Meta.Code)
	})

	t.Run("should return 200 success to view detail transaction", func(t *testing.T) {
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(c), nil)
		ciamPartner.EXPECT().JwtInfo(gomock.Any()).Return(jwtInfo, nil)
		transactionProvider.EXPECT().Detail(gomock.Any()).Return(&model.PartnerTransactionProjection{
			Transaction: decimal.NewFromInt(385000),
			WalletCode:  "WCODE_A",
			Msisdn:      "628118770510",
			Id:          1,
			Reward:      decimal.NewFromInt(5000),
			Email:       "someone@email.net",
			Cashback:    decimal.NewFromInt(12500),
			Qty:         2,
		}, nil)
		req := httptest.NewRequest(fiber.MethodGet, "/api/partner/v1/transactions/1", nil)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(fiber.HeaderAuthorization, "Bearer *secret*")
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelEBizKezbek)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.NotNil(t, m.Data)
	})

	t.Run("should return 200 failed to view detail transaction", func(t *testing.T) {
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(c), nil)
		ciamPartner.EXPECT().JwtInfo(gomock.Any()).Return(jwtInfo, nil)
		transactionProvider.EXPECT().Detail(gomock.Any()).Return(nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeNotFound,
			ErrorMessage: apps.ErrMsgNotFound,
		})
		req := httptest.NewRequest(fiber.MethodGet, "/api/partner/v1/transactions/1", nil)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(fiber.HeaderAuthorization, "Bearer *secret*")
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelEBizKezbek)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.Nil(t, m.Data)
		assert.Equal(t, apps.ErrCodeNotFound, m.Meta.Code)
	})
}
