package repository

import (
	"context"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type Partner struct {
	Pool   storage.Pooler
	Logger *zap.Logger
}

type PartnerPersister interface {
	Add(m model.Partner) *model.TechnicalError
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
