package storage

import (
	"github.com/challenge/pkg/models/errors"
	"github.com/challenge/pkg/modules/config"
	"sync"
)

var dbManagerInstanceOnce sync.Once
var dbManagerInstance *DBManager

type Closable interface {
	Close()
}

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

func (mgr *DBManager) CreateConnection(sessionName string) (IDBConnection, error) {
	if len(sessionName) == 0 {
		return nil, errors.NewInternalServerErrorMsg("The database connection can not be empty")
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
		case "sql":
			cnn, err := mgr.createSqlConnection(sessionName)

			if err != nil {
				return nil, err
			}

			mgr.connections[sessionName] = cnn
			return cnn, nil

		case "nosql":
			return nil, errors.NewInternalServerErrorMsg("Invalid db type in configuration for connections '" + sessionName + "'")

	}

	return nil, errors.NewInternalServerErrorMsg("Invalid db type in configuration for connections '" + sessionName + "'")
}

func (mgr *DBManager) createSqlConnection(sessionName string) (IDBConnection, error) {
	driver, err := config.ConfigurationSingleton().GetString("root.databases." + sessionName + ".driver")

	if err != nil {
		return nil, err
	}

	if driver != "sqlite3" {
		return nil, errors.NewInternalServerErrorMsg("Invalid driver for db type in configuration for connections '" + sessionName + "'")
	}
	//must use default sql driver

	return NewSqliteDbConnection(sessionName)
}


func (mgr *DBManager) Dispose() {
	mgr.sync.Lock()
	defer mgr.sync.Unlock()

	for _, v := range mgr.connections {
		c, ok := v.(Closable)

		if ok {
			c.Close()
		}
	}

	mgr.connections = make(map[string]IDBConnection, 0)
}
