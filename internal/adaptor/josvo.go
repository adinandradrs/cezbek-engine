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

type Josvo struct {
	Host   string
	ApiKey string
	Logger *zap.Logger
	Rest
}

type JosvoAdapter interface {
	AccountTransfer(inp *model.JosvoAccountTransferRequest) (*model.JosvoAccountTransferResponse, *model.TechnicalError)
}

func NewJosvo(j Josvo) JosvoAdapter {
	return &j
}

func (j *Josvo) AccountTransfer(inp *model.JosvoAccountTransferRequest) (*model.JosvoAccountTransferResponse, *model.TechnicalError) {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(*inp)
	if err != nil {
		return nil, apps.Exception("failed to build josvo payload", err, zap.Error(err), j.Logger)
	}
	req, err := http.NewRequest(fiber.MethodPost, j.Host+"/api/v2/accounts/transfer", payload)
	if err != nil {
		return nil, apps.Exception("failed to create josvo account transfer request", err, zap.Error(err), j.Logger)
	}
	req.Header.Add(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
	req.Header.Add(apps.HeaderApiKey, j.ApiKey)
	resp, err := j.Client().Do(req)
	if err != nil {
		return nil, apps.Exception("failed to fund transfer using josvo", err, zap.Error(err), j.Logger)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			j.Logger.Error("failed to close the body stream on josvo adapter", zap.Error(err))
		}
	}(resp.Body)
	var m model.JosvoAccountTransferResponse
	_ = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, apps.Exception("failed to map response on joso account transfer", err, zap.Error(err), j.Logger)
	}
	j.Logger.Info("success response", zap.Any("", m))
	return &m, nil
}
