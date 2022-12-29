package repository

import (
	"context"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type Partner struct {
	Pool   storage.Pooler
	Logger *zap.Logger
}

type PartnerPersister interface {
	Add(m model.Partner) *model.TechnicalError
	CountByCode(code string) (*int, *model.TechnicalError)
	FindActiveByCodeAndApiKey(code string, key string) (*model.Partner, *model.TechnicalError)
}

func NewPartner(p Partner) PartnerPersister {
	return &p
}

func (p Partner) Add(data model.Partner) *model.TechnicalError {
	tx, err := p.Pool.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.Serializable})
	var pid int64
	if err != nil {
		return apps.Exception("failed to begin transaction add partner", err,
			zap.String("code", data.Code.String), p.Logger)
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), `insert into partners (partner, code, api_key, salt, secret, email, 
		msisdn, email, officer, address, partner_logo, status, is_deleted, created_by, created_date)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, false, $13, now()) returning id`,
		data.Partner.String, data.Code.String, data.ApiKey.String, data.Salt.String, data.Secret.String, data.Email.String,
		data.Msisdn.String, data.Email.String, data.Officer.String, data.PartnerLogo.String, data.Status).Scan(&pid)
	if err != nil {
		return apps.Exception("failed to insert into partners table", err,
			zap.String("code", data.Code.String), p.Logger)
	}

	if err = tx.Commit(context.Background()); err != nil {
		p.Logger.Panic("transaction add partner failed", zap.Error(err))
	}

	return nil
}

func (p Partner) CountByCode(code string) (*int, *model.TechnicalError) {
	total := 0
	rows, err := p.Pool.Query(context.Background(), `select count(id) as total from partners where 
		code=$1 AND is_deleted=false`, code)
	if err != nil {
		return nil, apps.Exception("failed to count by code", err,
			zap.String("code", code), p.Logger)
	}

	err = pgxscan.ScanOne(&total, rows)
	if err != nil {
		return nil, apps.Exception("failed to map count own by code", err,
			zap.String("code", code), p.Logger)
	}

	return &total, nil
}

func (p Partner) FindActiveByCodeAndApiKey(code string, key string) (*model.Partner, *model.TechnicalError) {
	d := model.Partner{}
	query, err := p.Pool.Query(context.Background(), ` select id, partner, code, api_key, salt, secret,
			email, msisdn from partners where code = $1 and api_key = $2 
			and status = $3 and is_deleted = false `, code, key, apps.StatusActive)
	if err != nil {
		return nil, apps.Exception("failed to find active by code and api key", err, zap.Strings("", []string{code, key}), p.Logger)
	}
	err = pgxscan.ScanOne(&d, query)
	if err != nil {
		return nil, apps.Exception("failed to map active by code and api key", err, zap.Strings("", []string{code, key}), p.Logger)
	}
	return &d, nil

}
