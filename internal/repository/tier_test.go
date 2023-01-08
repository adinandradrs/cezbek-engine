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
	"time"
)

func TestTier_FindByPartnerMsisdn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	persister := NewTier(Tier{
		Logger: logger,
		Pool:   pool,
	})
	ctx := context.Background()
	pid := int64(1)
	msisdn := "628118770510"
	cmd := `select id, partner_id, msisdn, email,
		current_grade, current_tier, prev_grade, prev_tier, expired_date, transaction_recurring
		from tiers 
		where partner_id = $1 AND 
		msisdn = $2 AND
		is_deleted = false `

	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "partner_id", "msisdn", "email",
			"current_grade", "current_tier", "prev_grade", "prev_tier", "expired_date",
			"transaction_recurring"}).
			AddRow(int64(1), int64(1), sql.NullString{String: "628118770510"}, sql.NullString{String: "someone@email.id"}, 2, sql.NullString{String: "GOLD"},
				1, sql.NullString{String: "BRONZE"}, sql.NullTime{Time: time.Now()}, 3).ToPgxRows()
		pool.EXPECT().Query(ctx, cmd, pid, msisdn).Return(rows, nil)
		v, ex := persister.FindByPartnerMsisdn(pid, msisdn)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on failed to query", func(t *testing.T) {
		pool.EXPECT().Query(ctx, cmd, pid, msisdn).Return(nil, fmt.Errorf("something went wrong"))
		v, ex := persister.FindByPartnerMsisdn(pid, msisdn)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})

	t.Run("should return exception on failed to map the result", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "partner_id", "msisdn", "email",
			"current_grade", "current_tier", "prev_grade", "prev_tier", "expired_date",
			"transaction_recurring"}).
			AddRow(int64(1), 1, sql.NullString{String: "628118770510"}, sql.NullString{String: "someone@email.id"}, 2, sql.NullString{String: "GOLD"},
				1, sql.NullString{String: "BRONZE"}, sql.NullTime{Time: time.Now()}, 3).ToPgxRows()
		pool.EXPECT().Query(ctx, cmd, pid, msisdn).Return(rows, nil)
		v, ex := persister.FindByPartnerMsisdn(pid, msisdn)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}

func TestTier_CountExpire(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	persister := NewTier(Tier{
		Logger: logger,
		Pool:   pool,
	})
	ctx := context.Background()
	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"count"}).AddRow(1).
			ToPgxRows()
		pool.EXPECT().QueryRow(ctx, "select count(id) from tiers where expired_date <= now()").
			Return(rows)
		v, ex := persister.CountExpire()
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on failed to map the result", func(t *testing.T) {
		rows := pgxpoolmock.NewRows(nil).ToPgxRows()
		pool.EXPECT().QueryRow(ctx, "select count(id) from tiers where expired_date <= now()").
			Return(rows)
		v, ex := persister.CountExpire()
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}

func TestTier_Expire(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool, tx := pgxpoolmock.NewMockPgxIface(ctrl), pgxpoolmock.NewMockPgxIface(ctrl)
	persister := NewTier(Tier{
		Logger: logger,
		Pool:   pool,
	})
	cmd := `UPDATE tiers SET
		current_grade = prev_grade, 
		current_tier = prev_tier, 
		prev_grade = current_grade, 
		prev_tier = current_tier, 
		expired_date = $1, 
		transaction_recurring = 1,
		updated_date = NOW(), 
		updated_by = 0 
		WHERE 
			expired_date <= now()`
	ctx := context.Background()
	exp := time.Now()
	t.Run("should success", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(tx, nil)
		tx.EXPECT().Exec(ctx, cmd, exp).Return(nil, nil)
		tx.EXPECT().Commit(ctx).Times(1).Return(nil)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "failed to commit expire tier", r)
			}
		}()
		ex := persister.Expire(exp)
		assert.Nil(t, ex)
	})

	t.Run("should return exception on failed to begin transaction", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(nil, fmt.Errorf("something went wrong"))
		ex := persister.Expire(exp)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on failed to query", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(tx, nil)
		tx.EXPECT().Exec(ctx, cmd, exp).Return(nil, fmt.Errorf("something went wrong"))
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		ex := persister.Expire(exp)
		assert.NotNil(t, ex)
	})

	t.Run("should rollback on failed to commit transaction", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(tx, nil)
		tx.EXPECT().Exec(ctx, cmd, exp).Return(nil, nil)
		tx.EXPECT().Commit(ctx).Times(1).Return(fmt.Errorf("something went wrong"))
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "failed to commit expire tier", r)
			}
		}()
		ex := persister.Expire(exp)
		assert.Nil(t, ex)
	})
}

func TestTier_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	ctx := context.Background()
	pool, tx := pgxpoolmock.NewMockPgxIface(ctrl), pgxpoolmock.NewMockPgxIface(ctrl)
	persister := NewTier(Tier{
		Logger: logger,
		Pool:   pool,
	})
	mcmd := `INSERT INTO tiers 
		(partner_id, msisdn, email, current_grade, current_tier,
		prev_grade, prev_tier, expired_date, transaction_recurring, is_deleted, 
		created_by, created_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, FALSE, $10, NOW()) RETURNING ID`
	ccmd := `INSERT INTO tier_journeys 
		(last_transaction_id, current_grade, current_tier, notes, is_deleted, created_by, created_date, tier_id)
		VALUES ($1, $2, $3, $4, FALSE, $5, NOW(), $6)`
	tier := model.Tier{
		PartnerId:            1,
		Msisdn:               sql.NullString{String: "628118770510"},
		Email:                sql.NullString{String: "adinandra.dharmasurya@gmail.com"},
		CurrentGrade:         2,
		CurrentTier:          sql.NullString{String: "SILVER"},
		PrevGrade:            1,
		PrevTier:             sql.NullString{String: "GOLD"},
		ExpiredDate:          sql.NullTime{Time: time.Now()},
		TransactionRecurring: 1,
		BaseEntity: model.BaseEntity{
			UpdatedBy: sql.NullInt64{Int64: 1},
		},
		Journey: model.TierJourney{
			CurrentGrade: 2,
			CurrentTier:  sql.NullString{String: "SILVER"},
			Notes:        sql.NullString{String: "Something"},
			BaseEntity: model.BaseEntity{
				CreatedBy: sql.NullInt64{Int64: 1},
			},
		},
	}

	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"ID"}).AddRow(int64(1)).ToPgxRows()
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).Return(tx, nil)
		tx.EXPECT().QueryRow(ctx, mcmd, tier.PartnerId, tier.Msisdn.String, tier.Email.String,
			tier.CurrentGrade, tier.CurrentTier.String, tier.PrevGrade, tier.PrevTier.String, tier.ExpiredDate.Time,
			tier.TransactionRecurring, tier.CreatedBy.Int64).Return(rows)
		tx.EXPECT().Exec(context.Background(), ccmd,
			tier.Journey.LastTransactionId, tier.Journey.CurrentGrade, tier.Journey.CurrentTier.String,
			tier.Journey.Notes.String, tier.Journey.CreatedBy.Int64, tier.Journey.TierId,
		).Return(nil, nil)
		tx.EXPECT().Commit(ctx).Times(1).Return(nil)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "failed to commit add tier", r)
			}
		}()
		ex := persister.Add(tier)
		assert.Nil(t, ex)
	})

	t.Run("should return exception on failed to begin transaction", func(t *testing.T) {
		pool.EXPECT().BeginTx(context.Background(),
			pgx.TxOptions{IsoLevel: pgx.Serializable}).Return(nil, fmt.Errorf("something went wrong"))
		ex := persister.Add(tier)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on failed to map query result", func(t *testing.T) {
		rows := pgxpoolmock.NewRows(nil).ToPgxRows()
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).Return(tx, nil)
		tx.EXPECT().QueryRow(ctx, mcmd, tier.PartnerId, tier.Msisdn.String, tier.Email.String,
			tier.CurrentGrade, tier.CurrentTier.String, tier.PrevGrade, tier.PrevTier.String, tier.ExpiredDate.Time,
			tier.TransactionRecurring, tier.CreatedBy.Int64).Return(rows)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		ex := persister.Add(tier)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on failed execute child command", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"ID"}).AddRow(int64(1)).ToPgxRows()
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).Return(tx, nil)
		tx.EXPECT().QueryRow(ctx, mcmd, tier.PartnerId, tier.Msisdn.String, tier.Email.String,
			tier.CurrentGrade, tier.CurrentTier.String, tier.PrevGrade, tier.PrevTier.String, tier.ExpiredDate.Time,
			tier.TransactionRecurring, tier.CreatedBy.Int64).Return(rows)
		tx.EXPECT().Exec(context.Background(), ccmd,
			tier.Journey.LastTransactionId, tier.Journey.CurrentGrade, tier.Journey.CurrentTier.String,
			tier.Journey.Notes.String, tier.Journey.CreatedBy.Int64, tier.Journey.TierId,
		).Return(nil, fmt.Errorf("something went wrong"))
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		ex := persister.Add(tier)
		assert.NotNil(t, ex)
	})
}

func TestTier_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	ctx := context.Background()
	pool, tx := pgxpoolmock.NewMockPgxIface(ctrl), pgxpoolmock.NewMockPgxIface(ctrl)
	persister := NewTier(Tier{
		Logger: logger,
		Pool:   pool,
	})
	mcmd := `UPDATE tiers SET 
		next_grade = $1, 
		next_tier = $2, 
		current_grade = $3, 
		current_tier = $4, 
		prev_grade = $5, 
		prev_tier = $6, 
		expired_date = $7, 
		transaction_recurring = $8,
		updated_date = NOW(), 
		updated_by = $9 
		WHERE 
			msisdn = $10 and partner_id = $11`
	ccmd := `INSERT INTO tier_journeys 
		(last_transaction_id, current_grade, current_tier, notes, is_deleted, created_by, created_date, tier_id)
		VALUES ($1, $2, $3, $4, FALSE, $5, NOW(), $6)`
	tier := model.Tier{
		PartnerId:            1,
		Msisdn:               sql.NullString{String: "628118770510"},
		Email:                sql.NullString{String: "adinandra.dharmasurya@gmail.com"},
		CurrentGrade:         2,
		CurrentTier:          sql.NullString{String: "SILVER"},
		PrevGrade:            1,
		PrevTier:             sql.NullString{String: "GOLD"},
		ExpiredDate:          sql.NullTime{Time: time.Now()},
		TransactionRecurring: 1,
		BaseEntity: model.BaseEntity{
			UpdatedBy: sql.NullInt64{Int64: 1},
		},
		Journey: model.TierJourney{
			CurrentGrade: 2,
			CurrentTier:  sql.NullString{String: "SILVER"},
			Notes:        sql.NullString{String: "Something"},
			BaseEntity: model.BaseEntity{
				CreatedBy: sql.NullInt64{Int64: 1},
			},
		},
	}
	t.Run("should success", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).Return(tx, nil)
		tx.EXPECT().Exec(ctx, mcmd, tier.NextGrade, tier.NextTier.String, tier.CurrentGrade, tier.CurrentTier.String,
			tier.PrevGrade, tier.PrevTier.String, tier.ExpiredDate.Time, tier.TransactionRecurring,
			tier.UpdatedBy.Int64, tier.Msisdn.String, tier.PartnerId).Return(nil, nil)
		tx.EXPECT().Exec(context.Background(), ccmd,
			tier.Journey.LastTransactionId, tier.Journey.CurrentGrade, tier.Journey.CurrentTier.String,
			tier.Journey.Notes.String, tier.Journey.CreatedBy.Int64, tier.Journey.TierId,
		).Return(nil, nil)
		tx.EXPECT().Commit(ctx).Times(1).Return(nil)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "failed to commit update tier", r)
			}
		}()
		ex := persister.Update(tier)
		assert.Nil(t, ex)
	})

	t.Run("should return exception on failed to begin transaction", func(t *testing.T) {
		pool.EXPECT().BeginTx(context.Background(),
			pgx.TxOptions{IsoLevel: pgx.Serializable}).Return(nil, fmt.Errorf("something went wrong"))
		ex := persister.Update(tier)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on failed execute child command", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).Return(tx, nil)
		tx.EXPECT().Exec(ctx, mcmd, tier.NextGrade, tier.NextTier.String, tier.CurrentGrade, tier.CurrentTier.String,
			tier.PrevGrade, tier.PrevTier.String, tier.ExpiredDate.Time, tier.TransactionRecurring,
			tier.UpdatedBy.Int64, tier.Msisdn.String, tier.PartnerId).Return(nil, nil)
		tx.EXPECT().Exec(context.Background(), ccmd,
			tier.Journey.LastTransactionId, tier.Journey.CurrentGrade, tier.Journey.CurrentTier.String,
			tier.Journey.Notes.String, tier.Journey.CreatedBy.Int64, tier.Journey.TierId,
		).Return(nil, fmt.Errorf("something went wrong"))
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		ex := persister.Update(tier)
		assert.NotNil(t, ex)
	})
}
