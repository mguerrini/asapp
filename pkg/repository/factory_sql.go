package repository

import "github.com/challenge/pkg/repository/sql"

type sqlRepositoryFactory struct {

}

func NewSqlRepositoryFactory() IRepositoryFactory {
	return &sqlRepositoryFactory{}
}

func (s sqlRepositoryFactory) CreateUserRepository(sessionName string) IUserRepository {
	rep, err := sql.NewSqlUserRepository(sessionName)

	if err != nil {
		panic(err)
	}

	return rep
}

func (s sqlRepositoryFactory) CreateMessageRepository(sessionName string) IMessageRepository {
	rep, err := sql.NewSqlMessageRepository(sessionName)

	if err != nil {
		panic(err)
	}

	return rep
}

func (s sqlRepositoryFactory) CreateHealthRepository(sessionName string)  IHealthRepository {
	rep, err := sql.NewHealthRepository(sessionName)

	if err != nil {
		panic(err)
	}

	return rep
}
