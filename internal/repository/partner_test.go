package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPartner_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, pool, tx := apps.NewLog(false), pgxpoolmock.NewMockPgxIface(ctrl),
		pgxpoolmock.NewMockPgxIface(ctrl)
	data := model.Partner{
		Partner:     sql.NullString{String: "PT. LinkSaja Indonesia", Valid: true},
		Code:        sql.NullString{String: "LINKSAJA", Valid: true},
		ApiKey:      sql.NullString{String: "api-key-123-abc-456", Valid: true},
		Salt:        sql.NullString{String: "s4lTs3cr3T", Valid: true},
		Secret:      sql.NullString{String: "s0m3things3creTs!#", Valid: true},
		Email:       sql.NullString{String: "kezbeksupport@linksaja.co.id", Valid: true},
		Msisdn:      sql.NullString{String: "628123456789", Valid: true},
		Officer:     sql.NullString{String: "Someone", Valid: true},
		PartnerLogo: sql.NullString{String: "/logo/linksaja-1.png", Valid: true},
		Status:      apps.StatusActive,
	}
	ctx := context.Background()
	persister := NewPartner(Partner{
		Logger: logger,
		Pool:   pool,
	})

	t.Run("should success", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(tx, nil)
		rows := pgxpoolmock.NewRows([]string{"id"}).AddRow(1).ToPgxRows()
		tx.EXPECT().QueryRow(context.Background(), `insert into partners (partner, code, api_key, salt, secret, email, 
		msisdn, email, officer, address, partner_logo, status, is_deleted, created_by, created_date)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, false, $13, now()) returning id`,
			data.Partner.String, data.Code.String, data.ApiKey.String, data.Salt.String, data.Secret.String, data.Email.String,
			data.Msisdn.String, data.Email.String, data.Officer.String, data.PartnerLogo.String, data.Status).
			Return(rows)
		tx.EXPECT().Commit(ctx).Times(1).Return(nil)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "transaction add partner failed", r)
			}
		}()
		ex := persister.Add(data)
		assert.Nil(t, ex)
	})

	t.Run("should return exception on failed to begin transaction", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(nil, fmt.Errorf("something went wrong on begin transaction"))
		ex := persister.Add(data)
		assert.Equal(t, "something went wrong on begin transaction", ex.Exception)
	})

	t.Run("should return exception on failed to insert", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(tx, nil)
		rows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
		tx.EXPECT().QueryRow(context.Background(), `insert into partners (partner, code, api_key, salt, secret, email, 
		msisdn, email, officer, address, partner_logo, status, is_deleted, created_by, created_date)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, false, $13, now()) returning id`,
			data.Partner.String, data.Code.String, data.ApiKey.String, data.Salt.String, data.Secret.String, data.Email.String,
			data.Msisdn.String, data.Email.String, data.Officer.String, data.PartnerLogo.String, data.Status).
			Return(rows)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		ex := persister.Add(data)
		assert.NotNil(t, ex)
	})

	t.Run("should rollback on commit failure", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(tx, nil)
		rows := pgxpoolmock.NewRows([]string{"id"}).AddRow(1).ToPgxRows()
		tx.EXPECT().QueryRow(context.Background(), `insert into partners (partner, code, api_key, salt, secret, email, 
		msisdn, email, officer, address, partner_logo, status, is_deleted, created_by, created_date)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, false, $13, now()) returning id`,
			data.Partner.String, data.Code.String, data.ApiKey.String, data.Salt.String, data.Secret.String, data.Email.String,
			data.Msisdn.String, data.Email.String, data.Officer.String, data.PartnerLogo.String, data.Status).
			Return(rows)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		tx.EXPECT().Commit(ctx).Times(1).Return(fmt.Errorf("something went wrong on commit insert partner tx"))
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "transaction add partner failed", r)
			}
		}()
		ex := persister.Add(data)
		assert.Equal(t, "something went wrong on commit insert partner tx", ex.Exception)
	})

}
