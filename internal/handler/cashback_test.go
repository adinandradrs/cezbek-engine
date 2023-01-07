package handler

import (
	"bytes"
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/handler/middleware"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/adaptor"
	"github.com/adinandradrs/cezbek-engine/mock/storage"
	"github.com/adinandradrs/cezbek-engine/mock/usecase/client"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCashbackHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	transactionProvider, cacher, ciamPartner := client.NewMockTransactionProvider(ctrl), storage.NewMockCacher(ctrl),
		adaptor.NewMockCiamWatcher(ctrl)
	jwtAuthenticator := middleware.NewJwtAuthenticator(&middleware.JwtAuthenticator{
		Logger:      logger,
		CiamPartner: ciamPartner,
		Cacher:      cacher,
	})
	jwtAuthClientFilter := jwtAuthenticator.ClientFilter()

	api := fiber.New()

	cashbacks := api.Group("/api/v1/cashbacks")
	CashbackHandler(cashbacks, Cashback{
		TransactionProvider: transactionProvider,
		ClientFilter:        jwtAuthClientFilter,
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
	inp := model.TransactionRequest{
		Msisdn:               "628118770510",
		Email:                "someone@email.net",
		Qty:                  1,
		Amount:               decimal.NewFromInt(50000),
		MerchantCode:         "LSAJA",
		TransactionReference: "CORPA/123/456",
		SessionRequest: model.SessionRequest{
			Id:       1,
			Username: "CORP_A",
			Fullname: "Company A",
		},
	}
	b, _ := json.Marshal(inp)
	t.Run("should return 200 success to apply cashback", func(t *testing.T) {
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(c), nil)
		ciamPartner.EXPECT().JwtInfo(gomock.Any()).Return(jwtInfo, nil)
		transactionProvider.EXPECT().Add(gomock.Any()).Return(&model.TransactionResponse{
			TransactionId:        "TRX-001",
			TransactionTimestamp: time.Now().Unix(),
		}, nil)
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/cashbacks", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(fiber.HeaderAuthorization, "Bearer *secret*")
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelB2BClient)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.NotNil(t, m.Data)
	})

	t.Run("should return 400 failed to apply cashback due invalid merchant code", func(t *testing.T) {
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(c), nil)
		ciamPartner.EXPECT().JwtInfo(gomock.Any()).Return(jwtInfo, nil)
		transactionProvider.EXPECT().Add(gomock.Any()).Return(nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussMerchantCodeInvalid,
			ErrorMessage: apps.ErrMsgBussMerchantCodeInvalid,
		})
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/cashbacks", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(fiber.HeaderAuthorization, "Bearer *secret*")
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelB2BClient)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Equal(t, apps.ErrCodeBussMerchantCodeInvalid, m.Meta.Code)
	})

	t.Run("should return 400 failed to apply cashback due no cashback", func(t *testing.T) {
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(c), nil)
		ciamPartner.EXPECT().JwtInfo(gomock.Any()).Return(jwtInfo, nil)
		transactionProvider.EXPECT().Add(gomock.Any()).Return(nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussNoCashback,
			ErrorMessage: apps.ErrMsgBussNoCashback,
		})
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/cashbacks", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(fiber.HeaderAuthorization, "Bearer *secret*")
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelB2BClient)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Equal(t, apps.ErrCodeBussNoCashback, m.Meta.Code)
	})

	t.Run("should return 500 failed to apply cashback due error on wallet services", func(t *testing.T) {
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(c), nil)
		ciamPartner.EXPECT().JwtInfo(gomock.Any()).Return(jwtInfo, nil)
		transactionProvider.EXPECT().Add(gomock.Any()).Return(nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussH2HCashbackFailed,
			ErrorMessage: apps.ErrMsgBussH2HCashbackFailed,
		})
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/cashbacks", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(fiber.HeaderAuthorization, "Bearer *secret*")
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelB2BClient)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, apps.ErrCodeBussH2HCashbackFailed, m.Meta.Code)
	})

	t.Run("should return 400 failed to apply cashback due error on zero amount", func(t *testing.T) {
		inp.Amount = decimal.Zero
		b, _ := json.Marshal(inp)
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(c), nil)
		ciamPartner.EXPECT().JwtInfo(gomock.Any()).Return(jwtInfo, nil)

		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/cashbacks", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(fiber.HeaderAuthorization, "Bearer *secret*")
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelB2BClient)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Equal(t, apps.ErrCodeBadPayload, m.Meta.Code)
	})
}
