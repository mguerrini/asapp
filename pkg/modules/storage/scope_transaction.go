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
	Enlist(ctx context.Context, transactional ITransactional) error
	EnlistFunc (ctx context.Context, beginFunc func(ctx context.Context, opt *sql.TxOptions) (ITransaction, error)) error
	GetResource(id string) ITransaction
}

type ITransactional interface {
	Begin(ctx context.Context, opt *sql.TxOptions) (ITransaction, error)
}

// Implementation
type transactionScope struct {
	sync sync.Mutex
	isFinished bool
	resources  []ITransaction
}


func NewTransactionScope() ITransactionScope {
	output := &transactionScope{
		isFinished: false,
		resources:  make([]ITransaction, 0),
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
	if ts.isFinished {
		return errors.NewInternalServerErrorMsg("The transaction is finished, can not enlist.")
	}

	ts.sync.Lock()
	defer ts.sync.Unlock()

	for _, tx := range ts.resources {
		tx.Commit(ctx)
	}

	return nil
}

func (ts *transactionScope) Rollback(ctx context.Context) error {
	if ts.isFinished {
		return errors.NewInternalServerErrorMsg("The transaction is finished, can not enlist.")
	}

	ts.sync.Lock()
	defer ts.sync.Unlock()

	for _, tx := range ts.resources {
		tx.Rollback(ctx)
	}

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

