package storage

import (
	"context"
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



func (mgr *TransactionalScopeManager) Begin(ctx *context.Context)  ITransactionScope {
	tx := NewTransactionScope()

	mgr.sync.Lock()
	defer mgr.sync.Unlock()

	scopeContainer := mgr.getScopes(ctx)
	scopeContainer.scopes = append(scopeContainer.scopes, tx)

	return tx
}

//Commit current tx scope
func (mgr *TransactionalScopeManager) commit(ctx context.Context)   {
	mgr.sync.Lock()
	defer mgr.sync.Unlock()

	scopeContainer := mgr.getScopes(&ctx)

	if scopeContainer == nil || len(scopeContainer.scopes) == 0 {
		return
	}

	//get last scope
	scope := scopeContainer.scopes[len(scopeContainer.scopes)-1]
	scopeContainer.scopes = scopeContainer.scopes[0:len(scopeContainer.scopes)-1]

	tx := scope.(*transactionScope)
	tx.Commit(ctx)
}

//Rollback current tx scope
func (mgr *TransactionalScopeManager) rollback(ctx context.Context)   {
	mgr.sync.Lock()
	defer mgr.sync.Unlock()

	scopeContainer := mgr.getScopes(&ctx)

	if scopeContainer == nil || len(scopeContainer.scopes) == 0 {
		return
	}

	//get last scope
	scope := scopeContainer.scopes[len(scopeContainer.scopes)-1]
	scopeContainer.scopes = scopeContainer.scopes[0:len(scopeContainer.scopes)-1]

	tx := scope.(*transactionScope)
	tx.Rollback(ctx)
}

func (mgr *TransactionalScopeManager) GetCurrentScope(ctx context.Context) ITransactionScope {
	scopes:= mgr.getScopes(&ctx)

	if scopes == nil || len(scopes.scopes) == 0 {
		return nil
	} else {
		return scopes.scopes[len(scopes.scopes) - 1]
	}
}

func (mgr *TransactionalScopeManager) getScopes(ctx *context.Context) *scopesEnvelop {
	if ctx == nil {
		return nil
	}
	scopes := (*ctx).Value(CurrentTransactionScopeKey)

	if scopes == nil {
		scopes = &scopesEnvelop{scopes: make([]ITransactionScope, 0)}
		newCtx := context.WithValue(*ctx, CurrentTransactionScopeKey, scopes)
		ctx = &newCtx
	}

	scopes = (*ctx).Value(CurrentTransactionScopeKey)
	return scopes.(*scopesEnvelop)
}