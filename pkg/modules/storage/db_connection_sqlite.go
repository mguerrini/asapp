package storage

import (
	"database/sql"
	"github.com/challenge/pkg/models/errors"
	"github.com/challenge/pkg/modules/config"
	_ "github.com/mattn/go-sqlite3"
)


func NewSqliteDbConnection(sessionName string) (IDBConnection, error) {
	dataSource, err := config.ConfigurationSingleton().GetString("root.databases." + sessionName + ".datasource")

	if err != nil {
		return nil, errors.NewInternalServerErrorMsg("The datasource for session " + sessionName + " is not defined")
	}

	return NewSqliteDbConnectionWith(dataSource)
}

func NewSqliteDbConnectionWith(datasource string) (IDBConnection, error) {
	db, err := sql.Open("sqlite3", datasource)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}

	cnnId := "sqlite3_" + datasource
	cnn := &dbConnection{
		sessionName: cnnId,
		db:          db,
	}
	return cnn, nil
}







