package handler

import (
	"bytes"
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/handler/middleware"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/usecase/client"
	"github.com/adinandradrs/cezbek-engine/mock/usecase/partner"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestAuthorizationHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	partnerOnboardProvider, clientOnboardProvider := partner.NewMockOnboardProvider(ctrl),
		client.NewMockOnboardProvider(ctrl)

	api := fiber.New()
	authorization := api.Group("/api/v1/authorization")

	preAuthenticator := middleware.NewPreAuthenticator(&middleware.PreAuthenticator{
		Logger: logger,
	})
	preAuthClientFilter := preAuthenticator.ClientFilter()
	AuthorizationHandler(authorization, Authorization{
		PartnerOnboardProvider: partnerOnboardProvider,
		ClientOnboardProvider:  clientOnboardProvider,
		ClientFilter:           preAuthClientFilter,
	})

	t.Run("should return 200 success to auth b2b", func(t *testing.T) {
		inp := model.OfficerAuthenticationRequest{
			Email: "someone@email.net",
		}
		b, _ := json.Marshal(inp)
		partnerOnboardProvider.EXPECT().Authenticate(gomock.Any()).Return(&model.OfficerAuthenticationResponse{
			RemainingSeconds: 300,
			TransactionResponse: model.TransactionResponse{
				TransactionId:        "TRX-001",
				TransactionTimestamp: time.Now().Unix(),
			},
		}, nil)
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/authorization/b2b", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
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

	t.Run("should return 400 invalid payload failed to auth b2b", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/authorization/b2b", nil)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelEBizKezbek)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("should return 401 failed to auth b2b", func(t *testing.T) {
		inp := model.OfficerAuthenticationRequest{
			Email: "someone@email.net",
		}
		b, _ := json.Marshal(inp)
		partnerOnboardProvider.EXPECT().Authenticate(gomock.Any()).Return(nil, &model.BusinessError{
			ErrorMessage: apps.ErrMsgUnauthorized,
			ErrorCode:    apps.ErrCodeUnauthorized,
		})
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/authorization/b2b", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelEBizKezbek)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
		assert.Equal(t, apps.ErrCodeUnauthorized, m.Meta.Code)
	})

	t.Run("should return 500 failed to auth b2b", func(t *testing.T) {
		inp := model.OfficerAuthenticationRequest{
			Email: "someone@email.net",
		}
		b, _ := json.Marshal(inp)
		partnerOnboardProvider.EXPECT().Authenticate(gomock.Any()).Return(nil, &model.BusinessError{
			ErrorMessage: apps.ErrMsgSomethingWrong,
			ErrorCode:    apps.ErrCodeSomethingWrong,
		})
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/authorization/b2b", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelEBizKezbek)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, apps.ErrCodeSomethingWrong, m.Meta.Code)
	})

	t.Run("should return 200 success to auth OTP b2b", func(t *testing.T) {
		inp := model.OfficerValidationRequest{
			TransactionId: "TRX-001",
			Otp:           "123456",
		}
		b, _ := json.Marshal(inp)
		partnerOnboardProvider.EXPECT().Validate(gomock.Any()).Return(&model.OfficerValidationResponse{
			Id:      1,
			Code:    "CODE_A",
			Email:   "someone@email.net",
			Msisdn:  "628118770510",
			Company: "Company A",
			UrlLogo: "https://img.com/something.jpeg",
			SessionResponse: model.SessionResponse{
				Token:        "**secret**",
				RefreshToken: "**secret**",
			},
		}, nil)
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/authorization/otp", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelEBizKezbek)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		req.Header.Add(apps.HeaderClientTrxId, "TRX-001")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.NotNil(t, m.Data)
	})

	t.Run("should return 400 failed to auth OTP b2b forgot header transaction_id", func(t *testing.T) {
		inp := model.OfficerValidationRequest{
			Otp: "123456",
		}
		b, _ := json.Marshal(inp)
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/authorization/otp", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelEBizKezbek)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("should return 400 failed to auth OTP b2b wrong otp", func(t *testing.T) {
		inp := model.OfficerValidationRequest{
			TransactionId: "TRX-001",
			Otp:           "123456",
		}
		b, _ := json.Marshal(inp)
		partnerOnboardProvider.EXPECT().Validate(gomock.Any()).Return(nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussPartnerOTPInvalid,
			ErrorMessage: apps.ErrMsgBussPartnerOTPInvalid,
		})
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/authorization/otp", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelEBizKezbek)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		req.Header.Add(apps.HeaderClientTrxId, "TRX-001")
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Equal(t, apps.ErrCodeBussPartnerOTPInvalid, m.Meta.Code)
	})

	t.Run("should return 200 to auth client", func(t *testing.T) {
		inp := model.ClientAuthenticationRequest{
			Code: "LAJADA",
		}
		cid := int64(1)
		clientOnboardProvider.EXPECT().Authenticate(gomock.Any()).Return(&model.ClientAuthenticationResponse{
			Code:    "LAJADA",
			Id:      &cid,
			Company: "PT. LAJADA COMMERCE",
			SessionResponse: model.SessionResponse{
				Token:        "**secret**",
				RefreshToken: "**secret**",
			},
		}, nil)
		b, _ := json.Marshal(inp)
		epoch := time.Now().Unix()
		s := apps.HMAC(fiber.MethodPost+":"+strings.ToUpper(inp.Code)+":"+strconv.FormatInt(epoch, 10)+":"+strings.ToUpper("api-key-123-456"), inp.Code)
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/authorization/client", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelB2BClient)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		req.Header.Add(apps.HeaderApiKey, "api-key-123-456")
		req.Header.Add(apps.HeaderClientTimestamp, strconv.FormatInt(epoch, 10))
		req.Header.Add(apps.HeaderClientSignature, s)
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})

	t.Run("should return 401 to auth client", func(t *testing.T) {
		inp := model.ClientAuthenticationRequest{
			Code: "LAJADA",
		}
		clientOnboardProvider.EXPECT().Authenticate(gomock.Any()).Return(nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeUnauthorized,
			ErrorMessage: apps.ErrMsgUnauthorized,
		})
		b, _ := json.Marshal(inp)
		epoch := time.Now().Unix()
		s := apps.HMAC(fiber.MethodPost+":"+strings.ToUpper(inp.Code)+":"+strconv.FormatInt(epoch, 10)+":"+strings.ToUpper("api-key-123-456"), inp.Code)
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/authorization/client", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelB2BClient)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		req.Header.Add(apps.HeaderApiKey, "api-key-123-456")
		req.Header.Add(apps.HeaderClientTimestamp, strconv.FormatInt(epoch, 10))
		req.Header.Add(apps.HeaderClientSignature, s)
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("should return 500 to auth client", func(t *testing.T) {
		inp := model.ClientAuthenticationRequest{
			Code: "LAJADA",
		}
		clientOnboardProvider.EXPECT().Authenticate(gomock.Any()).Return(nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		})
		b, _ := json.Marshal(inp)
		epoch := time.Now().Unix()
		s := apps.HMAC(fiber.MethodPost+":"+strings.ToUpper(inp.Code)+":"+strconv.FormatInt(epoch, 10)+":"+strings.ToUpper("api-key-123-456"), inp.Code)
		req := httptest.NewRequest(fiber.MethodPost, "/api/v1/authorization/client", bytes.NewBuffer(b))
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(apps.HeaderClientChannel, apps.ChannelB2BClient)
		req.Header.Add(apps.HeaderClientDeviceId, "f-123-456")
		req.Header.Add(apps.HeaderClientOs, "Android 10")
		req.Header.Add(apps.HeaderClientVersion, "1.0.0")
		req.Header.Add(apps.HeaderApiKey, "api-key-123-456")
		req.Header.Add(apps.HeaderClientTimestamp, strconv.FormatInt(epoch, 10))
		req.Header.Add(apps.HeaderClientSignature, s)
		res, _ := api.Test(req, 100)
		m := model.Response{}
		_ = json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})
}
