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

type Middletrans struct {
	Host   string
	ApiKey string
	Logger *zap.Logger
	Rest
}

type MiddletransAdapter interface {
	WalletTransfer(inp *model.MiddletransWalletTransferRequest) (*model.MiddletransWalletTransferResponse,
		*model.TechnicalError)
}

func NewMiddletrans(mt Middletrans) MiddletransAdapter {
	return &mt
}

func (mt *Middletrans) WalletTransfer(inp *model.MiddletransWalletTransferRequest) (*model.MiddletransWalletTransferResponse, *model.TechnicalError) {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(*inp)
	if err != nil {
		return nil, apps.Exception("failed to build middletrans payload", err, zap.Error(err), mt.Logger)
	}
	req, err := http.NewRequest(fiber.MethodPost, mt.Host+"/ewallet/v1/transfer", payload)
	if err != nil {
		return nil, apps.Exception("failed to create middletrans wallet transfer request", err, zap.Error(err), mt.Logger)
	}
	req.Header.Add(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
	req.Header.Add(apps.HeaderApiKey, mt.ApiKey)
	resp, err := mt.Client().Do(req)
	if err != nil {
		return nil, apps.Exception("failed to wallet transfer using middletrans", err, zap.Error(err), mt.Logger)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			mt.Logger.Error("failed to close the body stream on middletrans adapter", zap.Error(err))
		}
	}(resp.Body)
	var m model.MiddletransWalletTransferResponse
	_ = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, apps.Exception("failed to map response on joso account transfer", err, zap.Error(err), mt.Logger)
	}
	return &m, nil
}
