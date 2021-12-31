package storage

import (
	"context"
	"database/sql"
)

type IDBConnection interface {
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type dbConnection struct {
	sessionName string
	db          *sql.DB
}

func (cnn *dbConnection) Begin(ctx context.Context, opt *sql.TxOptions) (IDBTransaction, error) {
	tx, err := cnn.db.BeginTx(ctx, opt)

	if err != nil {
		return nil, err
	}

	//agrgo la transaccion al contexto
	return &sqlTxDb{
		id: cnn.sessionName,
		tx: tx,
		dbConnection: cnn,
	}, nil
}

func (cnn *dbConnection) Query (ctx context.Context, query string, params ...interface{})  (*sql.Rows, error) {
	//get transaction
	currTx, err := cnn.getOpenedTransaction(ctx)
	if err != nil {
		return nil, err
	}

	if currTx	== nil {
		return 	cnn.db.QueryContext(ctx, query, params...)
	}

	//do the query
	return currTx.Query(ctx, query, params)
}

func (cnn *dbConnection) Exec (ctx context.Context, query string, params ...interface{}) (sql.Result, error) {
	//get transaction
	currTx, err := cnn.getOpenedTransaction(ctx)
	if err != nil {
		return nil, err
	}

	if currTx	== nil {
		return cnn.db.ExecContext(ctx, query, params...)
	}

	//do the query
	return currTx.Exec(ctx, query, params)
}

func (cnn *dbConnection) getOpenedTransaction (ctx context.Context) (IDBTransaction, error) {
	scope :=  TransactionalScopeManagerSingleton().GetCurrentScope(ctx)

	if scope == nil || scope.IsFinished() {
		return nil, nil
	}

	//search for an opened transaction
	tx := scope.GetResource(cnn.sessionName)
	if tx == nil {
		//enlisto a la tx
		scope.EnlistFunc(ctx, func(ctx context.Context, opt *sql.TxOptions) (ITransaction, error) {
			dbTx, err := cnn.Begin(ctx, opt)

			if err != nil {
				return nil, err
			}

			return dbTx.(ITransaction), nil
		})

		tx = scope.GetResource(cnn.sessionName)
	}

	dbTx := tx.(IDBTransaction)
	return dbTx, nil
}


func (cnn *dbConnection) committed(ctx context.Context, transaction ITransaction)  {

}

func (cnn *dbConnection) rollbacked (ctx context.Context, transaction ITransaction)  {

}








