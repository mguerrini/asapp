package sql

import "github.com/challenge/pkg/modules/storage"

type sqlHealthRepository struct {
	dbCnn storage.IDBConnection
}

func NewHealthRepository (sessionName string) (*sqlHealthRepository, error) {
	//get connection string
	dbConn, err := storage.DBManagerSingleton().CreateConnection(sessionName)

	if err != nil {
		return nil, err
	}
	return &sqlHealthRepository{
		dbCnn: dbConn,
	}, nil
}


func (r *sqlHealthRepository) Ping() error {
	return r.dbCnn.Ping()
}