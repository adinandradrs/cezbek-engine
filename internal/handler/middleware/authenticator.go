package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type PreAuthenticator struct {
	Logger *zap.Logger
}

type JwtAuthenticator struct {
	Logger *zap.Logger
	storage.Cacher
	CiamPartner adaptor.CiamWatcher
}

func NewPreAuthenticator(a *PreAuthenticator) PreAuthenticator {
	return *a
}

func NewJwtAuthenticator(a *JwtAuthenticator) JwtAuthenticator {
	return *a
}

func (a *JwtAuthenticator) forwardClientSession(v string, ctx *fiber.Ctx) (res model.ClientAuthenticationResponse) {
	_ = json.Unmarshal([]byte(v), &res)
	ctx.Request().Header.Add(apps.HeaderSessionId, strconv.FormatInt(*res.Id, 10))
	ctx.Request().Header.Add(apps.HeaderSessionUsername, res.Code)
	ctx.Request().Header.Add(apps.HeaderSessionFullname, res.Company)
	ctx.Request().Header.Add(apps.HeaderSessionRole, "B2BCLIENT")
	return res
}

func ClientSession(ctx *fiber.Ctx) model.SessionRequest {
	id, _ := strconv.ParseInt(ctx.Get(apps.HeaderSessionId), 10, 64)
	return model.SessionRequest{
		Username: ctx.Get(apps.HeaderSessionUsername),
		Email:    ctx.Get(apps.HeaderSessionEmail),
		Msisdn:   ctx.Get(apps.HeaderSessionMsisdn),
		Fullname: ctx.Get(apps.HeaderSessionFullname),
		Id:       id,
		ContextRequest: model.ContextRequest{
			Channel:       ctx.Get(apps.HeaderClientChannel),
			DeviceId:      ctx.Get(apps.HeaderClientDeviceId),
			Authorization: ctx.Get(fiber.HeaderAuthorization),
		},
	}
}

func (a *JwtAuthenticator) ClientFilter() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := validateClientChannel(ctx)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Meta: model.Meta{
					Code:    apps.ErrCodeInvalidChannel,
					Message: apps.ErrMsgInvalidChannel,
				},
			})
		}
		split := strings.Split(ctx.Get(fiber.HeaderAuthorization), "Bearer ")
		if len(split) < 2 {
			return ctx.Status(401).JSON(model.Response{
				Meta: model.Meta{
					Code:    apps.ErrCodeUnauthorized,
					Message: apps.ErrMsgUnauthorized,
				},
			})
		}
		jwt := split[1]

		res, ex := a.CiamPartner.JwtInfo(jwt)
		if ex != nil && ex.Exception == "Token is expired" {
			a.Logger.Error("failed to get expired jwt", zap.Any("", ex))
			return ctx.Status(401).JSON(model.Response{
				Meta: model.Meta{
					Code:    apps.ErrCodeTokenExpired,
					Message: apps.ErrMsgTokenExpired,
				},
			})
		} else if ex != nil && ex.Exception != "Token is expired" {
			a.Logger.Error("failed to get result jwt", zap.Any("", ex))
			return ctx.Status(401).JSON(model.Response{
				Meta: model.Meta{
					Code:    apps.ErrCodeUnauthorized,
					Message: apps.ErrMsgUnauthorized,
				},
			})
		}

		uname := strings.ToUpper(res["cognito:username"].(string))
		v, ex := a.Cacher.Get("CLIENTSESSION", uname)
		if ex != nil {
			a.Logger.Error("failed to get redis data", zap.Any("", ex))
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Meta: model.Meta{
					Code:    apps.ErrCodeUnauthorized,
					Message: apps.ErrMsgUnauthorized,
				},
			})
		}
		a.forwardClientSession(v, ctx)
		return ctx.Next()
	}
}

func validateClientChannel(ctx *fiber.Ctx) (err error) {
	if ctx.Get(apps.HeaderClientChannel) != apps.ChannelB2BClient {
		return fmt.Errorf("invalid channel")
	}
	return nil
}

func (a *PreAuthenticator) ClientFilter() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := validateClientChannel(ctx)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Meta: model.Meta{
					Code:    apps.ErrCodeInvalidChannel,
					Message: apps.ErrMsgInvalidChannel,
				},
			})
		}
		b := struct {
			Code string `json:"code"`
		}{}
		if err := ctx.BodyParser(&b); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err)
		}
		epoch := ctx.Get(apps.HeaderClientTimestamp)
		key := ctx.Get(apps.HeaderApiKey)

		d := string(ctx.Request().Header.Method()) + ":" + strings.ToUpper(b.Code) + ":" + epoch + ":" + strings.ToUpper(key)
		if apps.HMAC(d, b.Code) != ctx.Get(apps.HeaderClientSignature) {
			a.Logger.Error("the given signature is not recognized", zap.String("signature", ctx.Get(apps.HeaderClientSignature)),
				zap.String("code", b.Code))
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Meta: model.Meta{
					Code:    apps.ErrCodeUnauthorized,
					Message: apps.ErrMsgUnauthorized,
				},
			})
		}
		return ctx.Next()
	}
}
