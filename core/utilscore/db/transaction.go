package dbutils

import (
	errm "authentication_service/core/errmodule"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
)

func RollbackTransactionDB(ctx context.Context, tx pgx.Tx) {
	if rErr := tx.Rollback(ctx); rErr != nil && !errors.Is(rErr, pgx.ErrTxClosed) {
		log.Printf("error: failed to rollback transaction: %v", rErr)
	}
}

func BeginTransaction(ctx context.Context, databasePull *pgxpool.Pool) (*pgxpool.Conn, pgx.Tx, error) {
	conn, err := databasePull.Acquire(ctx)
	if err != nil {
		logrus.Errorf("🔴 error: %s: %+v", "BeginTransaction-BeginTransaction", err)
		return nil, nil, err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		logrus.Errorf("🔴 error: %s: %+v", "BeginTransaction-Begin", err)
		conn.Release()
		return nil, nil, err
	}
	return conn, tx, nil
}

// ExecuteTx выполняет функцию fn в транзакции.
// Если tx не nil, то используется переданная транзакция.
// Иначе создается новая транзакция.
func ExecuteTx(ctx context.Context, db *pgxpool.Pool, tx pgx.Tx, fn func(tx pgx.Tx) error) *errm.Error {
	var err error
	var transactionStarted bool

	if tx == nil {
		tx, err = db.Begin(ctx)
		if err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "ExecuteTx-Begin", err)
			return errm.NewError("failed to begin transaction", err)
		}

		transactionStarted = true
	}

	err = fn(tx)

	if transactionStarted {
		if err != nil {
			logrus.Errorf("🔴 error: %s: %+v", "ExecuteTx-transactionStarted", err)
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return errm.NewError("failed to rollback transaction", err)
			}
			return errm.NewError("failed to execute transaction", err)
		}

		if commitErr := tx.Commit(ctx); commitErr != nil {
			return errm.NewError("failed to commit transaction", commitErr)
		}
	}

	if err != nil {
		return errm.NewError("failed to execute transaction", err)
	}
	return nil
}
