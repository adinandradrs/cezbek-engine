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
		Partner: sql.NullString{String: "PT. LinkSaja Indonesia", Valid: true},
		Code:    sql.NullString{String: "LINKSAJA", Valid: true},
		ApiKey:  sql.NullString{String: "api-key-123-abc-456", Valid: true},
		Salt:    sql.NullString{String: "s4lTs3cr3T", Valid: true},
		Secret:  []byte("something"),
		Email:   sql.NullString{String: "kezbeksupport@linksaja.co.id", Valid: true},
		Msisdn:  sql.NullString{String: "628123456789", Valid: true},
		Officer: sql.NullString{String: "Someone", Valid: true},
		Logo:    sql.NullString{String: "/logo/linksaja-1.png", Valid: true},
		Address: sql.NullString{String: "Street A"},
		Status:  apps.StatusActive,
		BaseEntity: model.BaseEntity{
			CreatedBy: sql.NullInt64{Int64: int64(1), Valid: true},
		},
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
		msisdn, officer, address, logo, status, is_deleted, created_by, created_date)
		values ($1, $2, $3, $4, $5::bytea, $6, $7, $8, $9, $10, $11, false, $12, now()) returning id`,
			data.Partner.String, data.Code.String, data.ApiKey.String, data.Salt.String, data.Secret, data.Email.String,
			data.Msisdn.String, data.Officer.String, data.Address.String, data.Logo.String, data.Status, data.CreatedBy.Int64).
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
		msisdn, officer, address, logo, status, is_deleted, created_by, created_date)
		values ($1, $2, $3, $4, $5::bytea, $6, $7, $8, $9, $10, $11, false, $12, now()) returning id`,
			data.Partner.String, data.Code.String, data.ApiKey.String, data.Salt.String, data.Secret, data.Email.String,
			data.Msisdn.String, data.Officer.String, data.Address.String, data.Logo.String, data.Status, data.CreatedBy.Int64).
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
		msisdn, officer, address, logo, status, is_deleted, created_by, created_date)
		values ($1, $2, $3, $4, $5::bytea, $6, $7, $8, $9, $10, $11, false, $12, now()) returning id`,
			data.Partner.String, data.Code.String, data.ApiKey.String, data.Salt.String, data.Secret, data.Email.String,
			data.Msisdn.String, data.Officer.String, data.Address.String, data.Logo.String, data.Status, data.CreatedBy.Int64).
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

func TestPartner_CountByIdentifier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, pool := apps.NewLog(false), pgxpoolmock.NewMockPgxIface(ctrl)
	data := model.Partner{
		Code:   sql.NullString{String: "LINKSAJA"},
		Email:  sql.NullString{String: "someone@linksaja.id"},
		Msisdn: sql.NullString{String: "62811234567"},
	}
	ctx := context.Background()
	persister := NewPartner(Partner{
		Logger: logger,
		Pool:   pool,
	})

	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"total_to_add"}).AddRow(1).ToPgxRows()
		pool.EXPECT().Query(ctx, `select count(id) as total_to_add from partners where 
		(code=$1 or email = $2 or msisdn = $3) AND is_deleted=false`, data.Code.String, data.Email.String,
			data.Msisdn.String).Return(rows, nil)
		total, ex := persister.CountByIdentifier(data)
		assert.Equal(t, 1, *total)
		assert.Nil(t, ex)
	})

	t.Run("should return exception on failed to count query", func(t *testing.T) {
		pool.EXPECT().Query(ctx, `select count(id) as total_to_add from partners where 
		(code=$1 or email = $2 or msisdn = $3) AND is_deleted=false`, data.Code.String, data.Email.String,
			data.Msisdn.String).Return(nil,
			fmt.Errorf("something went wrong on count query"))
		total, ex := persister.CountByIdentifier(data)
		assert.Nil(t, total)
		assert.Equal(t, "something went wrong on count query", ex.Exception)
	})

	t.Run("should return exception on scan count query", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
		pool.EXPECT().Query(ctx, `select count(id) as total_to_add from partners where 
		(code=$1 or email = $2 or msisdn = $3) AND is_deleted=false`, data.Code.String, data.Email.String,
			data.Msisdn.String).Return(rows, nil)
		total, ex := persister.CountByIdentifier(data)
		assert.Nil(t, total)
		assert.NotNil(t, ex)
	})
}

func TestPartner_FindActiveByCodeAndApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, pool := apps.NewLog(false), pgxpoolmock.NewMockPgxIface(ctrl)
	code, key := "LINKSAJA", "api-key-123-abc-456"
	ctx := context.Background()
	persister := NewPartner(Partner{
		Logger: logger,
		Pool:   pool,
	})

	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "partner", "code", "api_key", "salt",
			"secret", "email", "msisdn"}).AddRow(int64(1), sql.NullString{String: "PT. LinkSaja Indonesia Terpadu", Valid: true},
			sql.NullString{String: "LINKSAJA", Valid: true}, sql.NullString{String: "api-key-123-abc-456", Valid: true},
			sql.NullString{String: "s4lTs3cr3T", Valid: true}, []byte("something"),
			sql.NullString{String: "someone@email.net", Valid: true},
			sql.NullString{String: "628123456789", Valid: true}).ToPgxRows()
		pool.EXPECT().Query(ctx, ` select id, partner, code, api_key, salt, secret,
			email, msisdn from partners where code = $1 and api_key = $2 
			and status = $3 and is_deleted = false `, code, key, apps.StatusActive).
			Return(rows, nil)
		data, ex := persister.FindActiveByCodeAndApiKey(code, key)
		assert.Equal(t, int64(1), data.Id)
		assert.NotNil(t, data)
		assert.Nil(t, ex)
	})

	t.Run("should return exception on failed to execute query", func(t *testing.T) {
		pool.EXPECT().Query(ctx, ` select id, partner, code, api_key, salt, secret,
			email, msisdn from partners where code = $1 and api_key = $2 
			and status = $3 and is_deleted = false `, code, key, apps.StatusActive).
			Return(nil, fmt.Errorf("something went wrong on execute query"))
		data, ex := persister.FindActiveByCodeAndApiKey(code, key)
		assert.Nil(t, data)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on map query result", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "partner", "code", "api_key", "salt",
			"secret", "email", "msisdn"}).AddRow(1, sql.NullString{String: "PT. LinkSaja Indonesia Terpadu", Valid: true},
			sql.NullString{String: "LINKSAJA", Valid: true}, sql.NullString{String: "api-key-123-abc-456", Valid: true},
			sql.NullString{String: "s4lTs3cr3T", Valid: true}, []byte("something"),
			sql.NullString{String: "someone@email.net", Valid: true},
			sql.NullString{String: "628123456789", Valid: true}).ToPgxRows()
		pool.EXPECT().Query(ctx, ` select id, partner, code, api_key, salt, secret,
			email, msisdn from partners where code = $1 and api_key = $2 
			and status = $3 and is_deleted = false `, code, key, apps.StatusActive).
			Return(rows, nil)
		data, ex := persister.FindActiveByCodeAndApiKey(code, key)
		assert.Nil(t, data)
		assert.NotNil(t, ex)
	})
}

func TestPartner_FindActiveByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, pool := apps.NewLog(false), pgxpoolmock.NewMockPgxIface(ctrl)
	email := "someone@email.net"
	ctx := context.Background()
	persister := NewPartner(Partner{
		Logger: logger,
		Pool:   pool,
	})

	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "partner", "code", "api_key", "salt",
			"secret", "email", "msisdn", "logo", "address"}).AddRow(int64(1), sql.NullString{String: "PT. LinkSaja Indonesia Terpadu", Valid: true},
			sql.NullString{String: "LINKSAJA", Valid: true}, sql.NullString{String: "api-key-123-abc-456", Valid: true},
			sql.NullString{String: "s4lTs3cr3T", Valid: true}, []byte("something"),
			sql.NullString{String: "someone@email.net", Valid: true},
			sql.NullString{String: "628123456789", Valid: true},
			sql.NullString{String: "/logo/linksaja-1.png", Valid: true},
			sql.NullString{String: "Jl. Nakula Sadewa no. 8B Jakarta Selatan", Valid: true}).ToPgxRows()
		pool.EXPECT().Query(ctx, ` select id, partner, code, 
			api_key, salt, secret, email, msisdn, logo, 
			address from partners where email = $1 and status = $2 
			and is_deleted = false `, email, apps.StatusActive).
			Return(rows, nil)
		data, ex := persister.FindActiveByEmail(email)
		assert.Equal(t, int64(1), data.Id)
		assert.NotNil(t, data)
		assert.Nil(t, ex)
	})

	t.Run("should return exception on failed to execute query", func(t *testing.T) {
		pool.EXPECT().Query(ctx, ` select id, partner, code, 
			api_key, salt, secret, email, msisdn, logo, 
			address from partners where email = $1 and status = $2 
			and is_deleted = false `, email, apps.StatusActive).
			Return(nil, fmt.Errorf("something went wrong on execute query"))
		data, ex := persister.FindActiveByEmail(email)
		assert.Nil(t, data)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on map query result", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "partner", "code", "api_key", "salt",
			"secret", "email", "msisdn", "logo", "address"}).AddRow(1, sql.NullString{String: "PT. LinkSaja Indonesia Terpadu", Valid: true},
			sql.NullString{String: "LINKSAJA", Valid: true}, sql.NullString{String: "api-key-123-abc-456", Valid: true},
			sql.NullString{String: "s4lTs3cr3T", Valid: true}, sql.NullString{String: "s0m3things3creTs!#", Valid: true},
			sql.NullString{String: "someone@email.net", Valid: true},
			sql.NullString{String: "628123456789", Valid: true},
			sql.NullString{String: "/logo/linksaja-1.png", Valid: true},
			sql.NullString{String: "Jl. Nakula Sadewa no. 8B Jakarta Selatan", Valid: true}).ToPgxRows()
		pool.EXPECT().Query(ctx, ` select id, partner, code, 
			api_key, salt, secret, email, msisdn, logo, 
			address from partners where email = $1 and status = $2 
			and is_deleted = false `, email, apps.StatusActive).
			Return(rows, nil)
		data, ex := persister.FindActiveByEmail(email)
		assert.Nil(t, data)
		assert.NotNil(t, ex)
	})
}
