package adaptor

import (
	"bytes"
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Xenit struct {
	Host               string
	BasicAuthorization string
	Logger             *zap.Logger
	Rest
}

type XenitAdapter interface {
	WalletTopup(inp *model.XenitWalletTopupRequest) (*model.XenitWalletTopupResponse, *model.TechnicalError)
}

func NewXenit(x Xenit) XenitAdapter {
	return &x
}

func (x *Xenit) WalletTopup(inp *model.XenitWalletTopupRequest) (*model.XenitWalletTopupResponse, *model.TechnicalError) {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(*inp)
	if err != nil {
		return nil, apps.Exception("failed to build xenit payload", err, zap.Error(err), x.Logger)
	}
	req, err := http.NewRequest(fiber.MethodPost, x.Host+"/api/v1/e-wallet/top-up", payload)
	if err != nil {
		return nil, apps.Exception("failed to create xenit top-up wallet request", err, zap.Error(err), x.Logger)
	}
	req.Header.Add(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
	req.Header.Add(apps.HeaderApiKey, "Basic "+x.BasicAuthorization)
	resp, err := x.Client().Do(req)
	if err != nil {
		return nil, apps.Exception("failed to wallet top-up using xenit", err, zap.Error(err), x.Logger)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			x.Logger.Error("failed to close the body stream on xenit adapter", zap.Error(err))
		}
	}(resp.Body)
	if resp.StatusCode != fiber.StatusOK {
		return nil, apps.Exception("bad response on xenit", err, zap.Any("", resp.Body), x.Logger)
	}
	var m model.XenitWalletTopupResponse
	_ = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, apps.Exception("failed to map response on xenit wallet top-up", err,
			zap.Error(err), x.Logger)
	}
	x.Logger.Info("success response", zap.Any("", m))
	return &m, nil
}
