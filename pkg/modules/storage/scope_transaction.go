package storage

import (
	"context"
	"errors"
	"sync"
)

// Interfaces
type ITransactionScope interface {
	IsFinished() bool
	Enlist(ctx context.Context, transactional ITransactional) error
}

type ITransactional interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context)
	Rollback (ctx context.Context)
}

// Implementation
type transactionScope struct {
	sync sync.Mutex
	isFinished bool
	resources  []ITransactional
}

func NewTransactionScope() ITransactionScope {
	output := &transactionScope{
		isFinished: false,
		resources:  make([]ITransactional, 0),
	}

	return output
}


func (ts *transactionScope) IsFinished() bool {
	return ts.isFinished
}

func (ts *transactionScope) Commit(ctx context.Context) error {
	if ts.isFinished {
		return errors.New("The transaction is finished, can not enlist.")
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
		return errors.New("The transaction is finished, can not enlist.")
	}

	ts.sync.Lock()
	defer ts.sync.Unlock()

	for _, tx := range ts.resources {
		tx.Rollback(ctx)
	}

	return nil
}

func (ts *transactionScope) Enlist(ctx context.Context, tx ITransactional) error {
	if ts.isFinished {
		return errors.New("The transaction is finished, can not enlist.")
	}

	ts.sync.Lock()
	defer ts.sync.Unlock()

	//start transaction
	err := tx.Begin(ctx)

	if err != nil {
		return err
	}

	ts.resources = append(ts.resources, tx)
	return nil
}

