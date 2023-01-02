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
	Dao         repository.PartnerPersister
	CiamWatcher adaptor.CiamWatcher
	S3Watcher   adaptor.S3Watcher
	PathS3      *string
	Logger      *zap.Logger
}

type PartnerManager interface {
	Add(inp *model.AddPartnerRequest) (*model.TransactionResponse, *model.BusinessError)
}

func NewPartner(p Partner) PartnerManager {
	return &p
}

func (p *Partner) uploadLogo(tid string, logo multipart.FileHeader) (*string, *model.TechnicalError) {
	fname, nsplit := tid, strings.Split(logo.Filename, ".")
	fext := nsplit[len(nsplit)-1]
	floc := *p.PathS3 + "logo/" + fname + "." + fext
	f, _ := logo.Open()
	_, err := p.S3Watcher.Upload(&model.S3UploadRequest{
		Destination: floc,
		Source:      f,
		ContentType: logo.Header.Get(fiber.HeaderContentType),
	})
	if err != nil {
		return nil, err
	}
	return &floc, nil
}

func (p *Partner) generateSecret(inp *model.AddPartnerRequest, pass *string) (*string, []byte, *model.TechnicalError) {
	salt := apps.Hash(inp.Code + ":" + uuid.NewString())
	secret, ex := apps.Encrypt(*pass, salt, p.Logger)
	if ex != nil {
		return nil, nil, ex
	}
	return &salt, secret, nil
}

func (p *Partner) ciamRegistration(data *model.Partner, pass *string) (*model.CiamUserResponse, *model.BusinessError) {
	resp, ex := p.CiamWatcher.OnboardPartner(model.CiamOnboardPartnerRequest{
		Email:       data.Email.String,
		PhoneNumber: data.Msisdn.String,
		Username:    data.Code.String,
		Name:        data.Partner.String,
		Picture:     data.Logo.String,
		Password:    *pass,
	})
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeESBUnavailable,
			ErrorMessage: apps.ErrMsgESBUnavailable,
		}
	}
	return resp, nil
}

func (p *Partner) Add(inp *model.AddPartnerRequest) (*model.TransactionResponse, *model.BusinessError) {
	tid := apps.TransactionId(inp.Code + apps.DefaultTrxId)
	data := model.Partner{
		Address: sql.NullString{String: inp.Address, Valid: true},
		Officer: sql.NullString{String: inp.Officer, Valid: true},
		Partner: sql.NullString{String: inp.Partner, Valid: true},
		Code:    sql.NullString{String: inp.Code, Valid: true},
		Email:   sql.NullString{String: inp.Email, Valid: true},
		Msisdn:  sql.NullString{String: inp.Msisdn, Valid: true},
		Status:  apps.StatusActive,
		BaseEntity: model.BaseEntity{
			CreatedBy: sql.NullInt64{Int64: inp.SessionRequest.Id, Valid: true},
		},
	}

	count, err := p.Dao.CountByIdentifier(data)
	if *count > 0 || err != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussPartnerExists,
			ErrorMessage: apps.ErrMsgBussPartnerExists,
		}
	}

	gpass, ex := apps.RandomPassword(12, 5, 3, p.Logger)
	p.Logger.Info("generate password", zap.String("code", inp.Code), zap.String("genpass", gpass))
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	salt, secret, ex := p.generateSecret(inp, &gpass)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	data.Salt = sql.NullString{String: *salt, Valid: true}
	data.Secret = secret

	floc, ex := p.uploadLogo(tid, inp.Logo)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	data.Logo = sql.NullString{String: *floc, Valid: true}

	ciam, cex := p.ciamRegistration(&data, &gpass)
	if cex != nil {
		return nil, cex
	}
	data.ApiKey = sql.NullString{String: apps.Hash(inp.Code + uuid.NewString() + ciam.SubId), Valid: true}
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
