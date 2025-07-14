package helper_db

import (
	"context"

	"gorm.io/gorm"
)

type ctxKey string

const txKey ctxKey = "gorm_tx"

type transactionBundler struct {
	db *gorm.DB
}

type ITransactionBundler interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if ok {
		return tx
	}
	return nil
}

func NewTransactionBundler(db *gorm.DB) ITransactionBundler {
	return &transactionBundler{
		db: db,
	}
}

func (t *transactionBundler) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := t.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if e := recover(); e != nil {
			tx.Rollback()
		}
	}()

	ctxWithTx := NewContextWithTx(ctx, tx)

	if err := fn(ctxWithTx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
