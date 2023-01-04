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
	CountByIdentifier(m model.Partner) (*int, *model.TechnicalError)
	FindActiveByCodeAndApiKey(code string, key string) (*model.Partner, *model.TechnicalError)
	FindActiveByEmail(email string) (*model.Partner, *model.TechnicalError)
}

func NewPartner(p Partner) PartnerPersister {
	return &p
}

func (p *Partner) Add(data model.Partner) *model.TechnicalError {
	tx, err := p.Pool.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.Serializable})
	var pid int64
	if err != nil {
		return apps.Exception("failed to begin transaction add partner", err,
			zap.String("code", data.Code.String), p.Logger)
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), `insert into partners (partner, code, api_key, salt, secret, email, 
		msisdn, officer, address, logo, status, is_deleted, created_by, created_date)
		values ($1, $2, $3, $4, $5::bytea, $6, $7, $8, $9, $10, $11, false, $12, now()) returning id`,
		data.Partner.String, data.Code.String, data.ApiKey.String, data.Salt.String, data.Secret, data.Email.String,
		data.Msisdn.String, data.Officer.String, data.Address.String, data.Logo.String, data.Status, data.CreatedBy.Int64).Scan(&pid)
	if err != nil {
		return apps.Exception("failed to insert into partners table", err,
			zap.String("code", data.Code.String), p.Logger)
	}

	if err = tx.Commit(context.Background()); err != nil {
		p.Logger.Panic("transaction add partner failed", zap.Error(err))
	}

	return nil
}

func (p *Partner) CountByIdentifier(data model.Partner) (*int, *model.TechnicalError) {
	total := 0
	rows, err := p.Pool.Query(context.Background(), `select count(id) as total_to_add from partners where 
		(code=$1 or email = $2 or msisdn = $3) AND is_deleted=false`, data.Code.String, data.Email.String,
		data.Msisdn.String)
	if err != nil {
		return nil, apps.Exception("failed to count identifier", err,
			zap.Strings("criteria", []string{data.Code.String, data.Email.String,
				data.Msisdn.String}), p.Logger)
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&total, rows)
	if err != nil {
		return nil, apps.Exception("failed to map count identifier result", err,
			zap.Strings("criteria", []string{data.Code.String, data.Email.String,
				data.Msisdn.String}), p.Logger)
	}

	return &total, nil
}

func (p *Partner) FindActiveByCodeAndApiKey(code string, key string) (*model.Partner, *model.TechnicalError) {
	d := model.Partner{}
	rows, err := p.Pool.Query(context.Background(), ` select id, partner, code, api_key, salt, secret,
			email, msisdn from partners where code = $1 and api_key = $2 
			and status = $3 and is_deleted = false `, code, key, apps.StatusActive)
	if err != nil {
		return nil, apps.Exception("failed to find active by code and api key", err, zap.Strings("", []string{code, key}), p.Logger)
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&d, rows)
	if err != nil {
		return nil, apps.Exception("failed to map active by code and api key", err, zap.Strings("", []string{code, key}), p.Logger)
	}

	return &d, nil
}

func (p *Partner) FindActiveByEmail(email string) (*model.Partner, *model.TechnicalError) {
	d := model.Partner{}
	rows, err := p.Pool.Query(context.Background(), ` select id, partner, code, 
			api_key, salt, secret, email, msisdn, logo, 
			address from partners where email = $1 and status = $2 
			and is_deleted = false `, email, apps.StatusActive)
	if err != nil {
		return nil, apps.Exception("failed to find active by email", err, zap.String("", email), p.Logger)
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&d, rows)
	if err != nil {
		return nil, apps.Exception("failed to map active by email", err, zap.String("", email), p.Logger)
	}

	return &d, nil
}
