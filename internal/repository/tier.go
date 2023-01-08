package repository

import (
	"context"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"time"
)

type Tier struct {
	Pool   storage.Pooler
	Logger *zap.Logger
}

type TierPersister interface {
	FindByPartnerMsisdn(pid int64, msisdn string) (*model.Tier, *model.TechnicalError)
	Add(tier model.Tier) *model.TechnicalError
	Update(tier model.Tier) *model.TechnicalError
	Expire(expired time.Time) *model.TechnicalError
	CountExpire() (*int, *model.TechnicalError)
}

func NewTier(t Tier) TierPersister {
	return &t
}

func (t *Tier) FindByPartnerMsisdn(pid int64, msisdn string) (*model.Tier, *model.TechnicalError) {
	var d model.Tier
	rows, err := t.Pool.Query(context.Background(), `select id, partner_id, msisdn, email,
		current_grade, current_tier, prev_grade, prev_tier, expired_date, transaction_recurring
		from tiers 
		where partner_id = $1 AND 
		msisdn = $2 AND
		is_deleted = false `, pid, msisdn)
	if err != nil {
		t.Logger.Info("", zap.Int64("partner_id", pid), zap.String("msisdn", msisdn))
		return nil, apps.Exception("failed to find tier by partner msisdn", err, zap.Any("", nil), t.Logger)
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&d, rows)
	if err != nil {
		t.Logger.Info("", zap.Int64("partner_id", pid), zap.String("msisdn", msisdn))
		return nil, apps.Exception("failed to map find tier by partner msisdn", err, zap.Any("", nil), t.Logger)
	}
	return &d, nil
}

func (t *Tier) addJourney(j model.TierJourney, tx pgx.Tx) *model.TechnicalError {
	_, err := tx.Exec(context.Background(), `INSERT INTO tier_journeys 
		(last_transaction_id, current_grade, current_tier, notes, is_deleted, created_by, created_date, tier_id)
		VALUES ($1, $2, $3, $4, FALSE, $5, NOW(), $6)`,
		j.LastTransactionId, j.CurrentGrade, j.CurrentTier.String,
		j.Notes.String, j.CreatedBy.Int64, j.TierId,
	)
	if err != nil {
		return apps.Exception("failed to add tier journey tx", err, zap.Any("", j), t.Logger)
	}
	return nil
}

func (t *Tier) Update(tier model.Tier) *model.TechnicalError {
	tx, err := t.Pool.BeginTx(context.Background(),
		pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return apps.Exception("failed to begin add tier tx", err, zap.Any("", tier), t.Logger)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `UPDATE tiers SET 
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
			msisdn = $10 and partner_id = $11`,
		tier.NextGrade, tier.NextTier.String, tier.CurrentGrade, tier.CurrentTier.String,
		tier.PrevGrade, tier.PrevTier.String, tier.ExpiredDate.Time, tier.TransactionRecurring,
		tier.UpdatedBy.Int64, tier.Msisdn.String, tier.PartnerId,
	)
	if err != nil {
		return apps.Exception("failed to update tier tx", err, zap.Any("", tier), t.Logger)
	}
	tier.Journey.TierId = tier.Id
	ex := t.addJourney(tier.Journey, tx)
	if ex != nil {
		return ex
	}
	if err = tx.Commit(context.Background()); err != nil {
		t.Logger.Panic("failed to commit update tier", zap.Any("tier", tier))
	}
	return nil
}

func (t *Tier) Add(tier model.Tier) *model.TechnicalError {
	tx, err := t.Pool.BeginTx(context.Background(),
		pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return apps.Exception("failed to begin add tier tx", err, zap.Any("", tier), t.Logger)
	}
	defer tx.Rollback(context.Background())

	var tid int64
	err = tx.QueryRow(context.Background(), `INSERT INTO tiers 
		(partner_id, msisdn, email, current_grade, current_tier,
		prev_grade, prev_tier, expired_date, transaction_recurring, is_deleted, 
		created_by, created_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, FALSE, $10, NOW()) RETURNING ID`,
		tier.PartnerId, tier.Msisdn.String, tier.Email.String,
		tier.CurrentGrade, tier.CurrentTier.String, tier.PrevGrade, tier.PrevTier.String, tier.ExpiredDate.Time,
		tier.TransactionRecurring, tier.CreatedBy.Int64,
	).Scan(&tid)
	if err != nil {
		return apps.Exception("failed to add tier tx", err, zap.Any("", tier), t.Logger)
	}
	tier.Journey.TierId = tid
	ex := t.addJourney(tier.Journey, tx)
	if ex != nil {
		return ex
	}
	if err = tx.Commit(context.Background()); err != nil {
		t.Logger.Panic("failed to commit add tier", zap.Any("tier", tier))
	}
	return nil
}

func (t *Tier) CountExpire() (*int, *model.TechnicalError) {
	var count int
	row := t.Pool.QueryRow(context.Background(), "select count(id) from tiers where expired_date <= now()")
	err := row.Scan(&count)
	if err != nil {
		return nil, apps.Exception("failed to count expire today", err, zap.Time("", time.Now()), t.Logger)
	}
	return &count, nil
}

func (t *Tier) Expire(expired time.Time) *model.TechnicalError {
	tx, err := t.Pool.BeginTx(context.Background(),
		pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return apps.Exception("failed to begin expire tier tx", err, zap.Time("", time.Now()), t.Logger)
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `UPDATE tiers SET
		current_grade = prev_grade, 
		current_tier = prev_tier, 
		prev_grade = current_grade, 
		prev_tier = current_tier, 
		expired_date = $1, 
		transaction_recurring = 1,
		updated_date = NOW(), 
		updated_by = 0 
		WHERE 
			expired_date <= now()`,
		expired,
	)
	if err != nil {
		return apps.Exception("failed to expire tier tx", err, zap.Time("", time.Now()), t.Logger)
	}
	if err = tx.Commit(context.Background()); err != nil {
		t.Logger.Panic("failed to commit expire tier", zap.Time("", time.Now()))
	}
	return nil
}
