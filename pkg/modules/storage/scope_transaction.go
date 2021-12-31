package storage

import (
	"context"
	"database/sql"
	"github.com/challenge/pkg/models/errors"
	"sync"
)

// Interfaces
type ITransactionScope interface {
	IsFinished() bool
	GetResource(id string) ITransaction

	Enlist(ctx context.Context, transactional ITransactional) error
	EnlistFunc (ctx context.Context, beginFunc func(ctx context.Context, opt *sql.TxOptions) (ITransaction, error)) error

	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type ITransactional interface {
	Begin(ctx context.Context, opt *sql.TxOptions) (ITransaction, error)
}

// Implementation
type transactionScope struct {
	sync       sync.Mutex
	txOptions  *sql.TxOptions
	isFinished bool
	resources  []ITransaction
	manager    *TransactionalScopeManager
}

func NewTransactionScope(mgr *TransactionalScopeManager, opt *sql.TxOptions) ITransactionScope {
	output := &transactionScope{
		txOptions:  opt,
		isFinished: false,
		resources:  make([]ITransaction, 0),
		manager:    mgr,
	}

	return output
}


func (ts *transactionScope) IsFinished() bool {
	return ts.isFinished
}

func (ts *transactionScope) GetResource(id string) ITransaction {
	ts.sync.Lock()
	defer ts.sync.Unlock()

	for _, tx := range ts.resources {
		if tx.Id() == id {
			return tx
		}
	}

	return nil
}

func (ts *transactionScope) Commit(ctx context.Context) error {
	ts.sync.Lock()
	defer ts.sync.Unlock()

	if ts.isFinished {
		return errors.NewInternalServerErrorMsg("The transaction is finished, can not enlist.")
	}

	ts.isFinished = true

	for _, tx := range ts.resources {
		tx.Commit(ctx)
	}

	ts.manager.Finished(ctx, ts)

	return nil
}

func (ts *transactionScope) Rollback(ctx context.Context) error {
	ts.sync.Lock()
	defer ts.sync.Unlock()

	if ts.isFinished {
		return errors.NewInternalServerErrorMsg("The transaction is finished, can not enlist.")
	}

	ts.isFinished = true

	for _, tx := range ts.resources {
		tx.Rollback(ctx)
	}

	ts.manager.Finished(ctx, ts)
	return nil
}


func (ts *transactionScope) EnlistFunc (ctx context.Context, beginFunc func(ctx context.Context, opt *sql.TxOptions) (ITransaction, error)) error {
	if ts.isFinished {
		return errors.NewInternalServerErrorMsg("The transaction is finished, can not enlist.")
	}

	ts.sync.Lock()
	defer ts.sync.Unlock()

	newTx, err := beginFunc(ctx, nil)

	if err != nil {
		return err
	}

	ts.resources = append(ts.resources, newTx)
	return nil
}

func (ts *transactionScope) Enlist(ctx context.Context, tx ITransactional) error {
	if ts.isFinished {
		return errors.NewInternalServerErrorMsg("The transaction is finished, can not enlist.")
	}

	ts.sync.Lock()
	defer ts.sync.Unlock()

	//start transaction
	newTx, err := tx.Begin(ctx, nil)

	if err != nil {
		return err
	}

	ts.resources = append(ts.resources, newTx)
	return nil
}

