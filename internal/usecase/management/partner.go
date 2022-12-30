package management

import (
	"database/sql"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"mime/multipart"
	"strings"
	"time"
)

type Partner struct {
	Dao       repository.PartnerPersister
	Logger    *zap.Logger
	S3Watcher adaptor.S3Watcher
	CDN       string
	PathS3    string
}

type PartnerManager interface {
	Add(inp *model.AddPartnerRequest) (*model.TransactionResponse, *model.BusinessError)
}

func NewPartner(p Partner) PartnerManager {
	return &p
}

func (p Partner) uploadLogo(tid string, logo multipart.FileHeader) (*string, *model.TechnicalError) {
	fname, nsplit := tid, strings.Split(logo.Filename, ".")
	fext := nsplit[len(nsplit)-1]
	floc := p.PathS3 + "logo/" + fname + "." + fext
	f, _ := logo.Open()
	_, err := p.S3Watcher.Upload(&adaptor.S3UploadRequest{
		Destination: floc,
		Source:      f,
		ContentType: logo.Header.Get(fiber.HeaderContentType),
	})
	if err != nil {
		return nil, err
	}
	return &floc, nil
}

func (p Partner) generateSecret(inp *model.AddPartnerRequest) (*string, []byte, *model.TechnicalError) {
	salt := apps.Hash(inp.Code + ":" + uuid.NewString())
	gpass, ex := apps.RandomPassword(12, 5, 3, p.Logger)
	if ex != nil {
		return nil, nil, ex
	}
	secret, ex := apps.Encrypt(gpass, salt, p.Logger)
	if ex != nil {
		return nil, nil, ex
	}
	return &salt, secret, nil
}

func (p Partner) Add(inp *model.AddPartnerRequest) (*model.TransactionResponse, *model.BusinessError) {
	tid := apps.TransactionId(inp.Code + apps.DefaultTrxId)
	data := model.Partner{
		Code:   sql.NullString{String: inp.Code, Valid: true},
		Email:  sql.NullString{String: inp.Email, Valid: true},
		Msisdn: sql.NullString{String: inp.Msisdn, Valid: true},
	}

	count, err := p.Dao.CountByIdentifier(data)
	if *count == 0 || err != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussPartnerExists,
			ErrorMessage: apps.ErrMsgBussPartnerExists,
		}
	}

	floc, err := p.uploadLogo(tid, inp.Logo)
	if err != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}

	salt, secret, ex := p.generateSecret(inp)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}

	data.Partner = sql.NullString{String: inp.Partner, Valid: true}
	data.ApiKey = sql.NullString{String: apps.Hash(inp.Code + uuid.NewString() + inp.Email), Valid: true}
	data.Salt = sql.NullString{String: *salt, Valid: true}
	data.Secret = secret
	data.Officer = sql.NullString{String: inp.Officer, Valid: true}
	data.PartnerLogo = sql.NullString{String: *floc, Valid: true}
	data.Status = apps.StatusActive
	ex = p.Dao.Add(data)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSubmitted,
			ErrorMessage: apps.ErrMsgSubmitted,
		}
	}

	return &model.TransactionResponse{
		TransactionId:        tid,
		TransactionTimestamp: time.Now().Unix(),
	}, nil
}
