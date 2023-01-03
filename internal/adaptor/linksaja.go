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

type Linksaja struct {
	Host     string
	Logger   *zap.Logger
	Username string
	Password string
	Rest
}

type LinksajaAdapter interface {
	Authorization() (*model.LinksajaAuthorizationResponse, *model.TechnicalError)
	FundTransfer(inp *model.LinksajaFundTransferRequest) (*model.LinksajaFundTransferResponse, *model.TechnicalError)
}

func NewLinksaja(l Linksaja) LinksajaAdapter {
	return &l
}

func (l *Linksaja) Authorization() (*model.LinksajaAuthorizationResponse, *model.TechnicalError) {
	inp, payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{l.Username, l.Password}, new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(inp)
	if err != nil {
		return nil, apps.Exception("failed to build linksaja payload", err, zap.Error(err), l.Logger)
	}
	req, err := http.NewRequest(fiber.MethodPost, l.Host+"/api/v1/thirdparty/authorization", payload)
	if err != nil {
		return nil, apps.Exception("failed to create linksaja auth request", err, zap.Error(err), l.Logger)
	}
	req.Header.Add(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
	resp, err := l.Rest.Client().Do(req)
	if err != nil {
		return nil, apps.Exception("failed to auth linksaja", err, zap.Error(err), l.Logger)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			l.Logger.Error("failed to close the body stream on linksaja adapter auth", zap.Error(err))
		}
	}(resp.Body)
	var m model.LinksajaAuthorizationResponse
	_ = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, apps.Exception("failed to map response on linksaja auth", err, zap.Error(err), l.Logger)
	}
	return &m, nil
}

func (l *Linksaja) FundTransfer(inp *model.LinksajaFundTransferRequest) (*model.LinksajaFundTransferResponse, *model.TechnicalError) {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(*inp)
	if err != nil {
		return nil, apps.Exception("failed to build linksaja payload", err, zap.Error(err), l.Logger)
	}
	req, err := http.NewRequest(fiber.MethodPost, l.Host+"/api/v1/transfer/fund", payload)
	if err != nil {
		return nil, apps.Exception("failed to create linksaja fund transfer request", err, zap.Error(err), l.Logger)
	}
	req.Header.Add(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
	req.Header.Add(fiber.HeaderAuthorization, "Bearer "+inp.Bearer)
	resp, err := l.Rest.Client().Do(req)
	if err != nil {
		return nil, apps.Exception("failed to fund transfer using linksaja", err, zap.Error(err), l.Logger)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			l.Logger.Error("failed to close the body stream on linksaja adapter", zap.Error(err))
		}
	}(resp.Body)
	var m model.LinksajaFundTransferResponse
	_ = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, apps.Exception("failed to map response on linksaja fund transfer", err, zap.Error(err), l.Logger)
	}
	return &m, nil
}
