package zrepository

import (
	"context"
	"database/sql"
)

func Transaction(ctx context.Context, db *sql.DB, fn func(context.Context) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	txCtx := WithTx(ctx, tx)
	if err := fn(txCtx); err != nil {
		return err
	}
	return tx.Commit()
}

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, "tx", tx)
}

func getTx(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value("tx").(*sql.Tx)
	return tx
}

func queryRow(ctx context.Context, db *sql.DB, query string, args ...any) *sql.Row {
	if tx := getTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return db.QueryRowContext(ctx, query, args...)
}

func exec(ctx context.Context, db *sql.DB, query string, args ...any) (sql.Result, error) {
	if tx := getTx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return db.ExecContext(ctx, query, args...)
}

func query(ctx context.Context, db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	if tx := getTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}
	return db.QueryContext(ctx, query, args...)
}
