package storage

import (
	"context"
	"database/sql"
)

type ITransaction interface {
	Id() string
	Commit(ctx context.Context) error
	Rollback (ctx context.Context) error
}

type IDBTransaction interface {
	IDBConnection
	ITransaction
}

type sqlTxDb struct {
	id string
	tx *sql.Tx
	dbConnection *dbConnection
}

func (cnn *sqlTxDb) Id ()string {
	return cnn.id
}

func (cnn *sqlTxDb) Query (ctx context.Context, query string, params ...interface{}) (*sql.Rows, error) {
	return cnn.tx.QueryContext(ctx, query, params )
}

func (cnn *sqlTxDb) Exec (ctx context.Context, query string, params ...interface{}) (sql.Result, error) {
	return cnn.tx.ExecContext(ctx, query, params )
}

func (cnn *sqlTxDb) Commit(ctx context.Context) error {
	err := cnn.tx.Commit()

	if err != nil {
		return err
	}

	cnn.dbConnection.committed(ctx, cnn)
	return nil
}

func (cnn *sqlTxDb) Rollback (ctx context.Context) error {
	err := cnn.tx.Commit()

	if err != nil {
		return err
	}

	cnn.dbConnection.rollbacked(ctx, cnn)
	return nil
}
