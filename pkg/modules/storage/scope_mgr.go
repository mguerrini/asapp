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

type scopesEnvelop struct {
	scopes []ITransactionScope
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

	newCtx, scopeContainer := mgr.getScopes(ctx)
	scopeContainer.scopes = append(scopeContainer.scopes, tx)

	return newCtx, tx
}

// Finish the last TxScope TODO, search and remove
func (mgr *TransactionalScopeManager) Finished(ctx context.Context, scope ITransactionScope)   {
	mgr.sync.Lock()
	defer mgr.sync.Unlock()

	_, scopeContainer := mgr.getScopes(ctx)

	if scopeContainer == nil || len(scopeContainer.scopes) == 0 {
		return
	}

	//get last scope
	scopeContainer.scopes = scopeContainer.scopes[0:len(scopeContainer.scopes)-1]
}


func (mgr *TransactionalScopeManager) GetCurrentScope(ctx context.Context) ITransactionScope {
	_, scopes := mgr.getScopes(ctx)

	if scopes == nil || len(scopes.scopes) == 0 {
		return nil
	} else {
		return scopes.scopes[len(scopes.scopes) - 1]
	}
}

func (mgr *TransactionalScopeManager) getScopes(ctx context.Context) (context.Context, *scopesEnvelop) {
	if ctx == nil {
		return nil, nil
	}
	scopes := ctx.Value(CurrentTransactionScopeKey)

	if scopes == nil {
		scopes = &scopesEnvelop{scopes: make([]ITransactionScope, 0)}
		newCtx := context.WithValue(ctx, CurrentTransactionScopeKey, scopes)
		ctx = newCtx
	}

	scopes = ctx.Value(CurrentTransactionScopeKey)
	return ctx, scopes.(*scopesEnvelop)
}