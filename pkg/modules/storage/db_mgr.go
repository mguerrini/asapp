package storage

import (
	"context"
	"errors"
	"github.com/challenge/pkg/modules/config"
	"sync"
)

var dbManagerInstanceOnce sync.Once
var dbManagerInstance *DBManager

type DBManager struct {
	sync        sync.Mutex
	connections map[string]IDBConnection
}


func DBManagerSingleton() *DBManager {
	dbManagerInstanceOnce.Do(func() {
		if dbManagerInstance != nil {
			return
		}

		dbManagerInstance = &DBManager{
			connections:     make(map[string]IDBConnection, 0),
		}
	})

	return dbManagerInstance
}

func (mgr *DBManager) Open(ctx context.Context, sessionName string) (IDBConnection, error) {
	if len(sessionName) == 0 {
		return nil, errors.New("The database connection can not be empty")
	}

	mgr.sync.Lock()
	defer mgr.sync.Unlock()

	if db, ok := mgr.connections[sessionName]; ok {
		return db, nil
	}

	dbType, err := config.ConfigurationSingleton().GetString("root.databases." + sessionName + ".type")
	if err != nil {
		return nil, err
	}

	switch dbType {
		case "memory":
			return nil, nil

		case "sql":
			return nil, nil
	}

	return nil, errors.New("Invalid db type in configuration for connections '" + sessionName + "'")
}

func (mgr *DBManager) CloseAll() {
	mgr.sync.Lock()
	defer mgr.sync.Unlock()

	for _, v := range mgr.connections {
		v.(Closable).Close()
	}

	mgr.connections = make(map[string]IDBConnection, 0)
}

func (mgr *DBManager) enlist(ctx context.Context, connection ITransactional) {
	scope := TransactionalScopeManagerSingleton().GetCurrentScope(ctx)

	if scope == nil {
		return
	}

	scope.Enlist(ctx, connection)
}