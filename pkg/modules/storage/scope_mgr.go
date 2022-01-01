package storage

import (
	"context"
	"database/sql"
	"sync"
)

const (
	CurrentTransactionScopeKey string = "TransactionScope"
)

var transactionalScopeManagerOnce sync.Once
var transactionalScopeManagerInstance *TransactionalScopeManager

type TransactionalScopeManager struct {
	sync sync.Mutex
}

type txBoundary struct {
	tx     ITransactionScope
	parent context.Context
}

func TransactionalScopeManagerSingleton() *TransactionalScopeManager {
	transactionalScopeManagerOnce.Do(func() {
		if transactionalScopeManagerInstance != nil {
			return
		}

		transactionalScopeManagerInstance = &TransactionalScopeManager{}
	})

	return transactionalScopeManagerInstance
}



func (mgr *TransactionalScopeManager) Begin(ctx context.Context, opt *sql.TxOptions) (context.Context, ITransactionScope) {
	tx := NewTransactionScope(mgr, opt)

	mgr.sync.Lock()
	defer mgr.sync.Unlock()

	scope := txBoundary{
		tx:     tx,
		parent: ctx,
	}
	newCtx := context.WithValue(ctx, CurrentTransactionScopeKey, scope)

	return newCtx, tx
}

// Finish the last TxScope TODO, search and remove
func (mgr *TransactionalScopeManager) Finished(ctx context.Context, scope ITransactionScope) context.Context  {
	mgr.sync.Lock()
	defer mgr.sync.Unlock()

	txScope := ctx.Value(CurrentTransactionScopeKey)

	if txScope == nil{
		return ctx
	}

	aux := txScope.(txBoundary)

	return aux.parent
}


func (mgr *TransactionalScopeManager) GetCurrentScope(ctx context.Context) ITransactionScope {
	txScope := ctx.Value(CurrentTransactionScopeKey)

	if txScope == nil{
		return nil
	}
	aux := txScope.(txBoundary)
	return aux.tx
}
