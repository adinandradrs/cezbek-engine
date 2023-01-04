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

type Gopaid struct {
	Host   string
	ApiKey string
	Logger *zap.Logger
	Rest
}

type GopaidAdapter interface {
	Topup(inp *model.GopaidTopUpRequest) (*model.GopaidTopupResponse, *model.TechnicalError)
}

func NewGopaid(g Gopaid) GopaidAdapter {
	return &g
}

func (g *Gopaid) Topup(inp *model.GopaidTopUpRequest) (*model.GopaidTopupResponse, *model.TechnicalError) {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(*inp)
	if err != nil {
		return nil, apps.Exception("failed to build gopaid payload", err, zap.Error(err), g.Logger)
	}
	req, err := http.NewRequest(fiber.MethodPost, g.Host+"/api/v2/accounts/transfer", payload)
	if err != nil {
		return nil, apps.Exception("failed to create gopaid topup request", err, zap.Error(err), g.Logger)
	}
	req.Header.Add(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
	req.Header.Add(apps.HeaderApiKey, g.ApiKey)
	resp, err := g.Client().Do(req)
	if err != nil {
		return nil, apps.Exception("failed to top up wallet using gopaid", err, zap.Error(err), g.Logger)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			g.Logger.Error("failed to close the body stream on gopaid adapter", zap.Error(err))
		}
	}(resp.Body)
	var m model.GopaidTopupResponse
	_ = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, apps.Exception("failed to map response on gopaid top up balance", err, zap.Error(err), g.Logger)
	}
	g.Logger.Info("success response", zap.Any("", m))
	return &m, nil
}
