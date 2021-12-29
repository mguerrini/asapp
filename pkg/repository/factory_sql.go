package repository

type sqlRepositoryFactory struct {

}

func NewSqlRepositoryFactory() IRepositoryFactory {
	return &sqlRepositoryFactory{}
}

func (s sqlRepositoryFactory) CreateUserRepository(sessionName string) IUserRepository {
	panic("implement me")
}

func (s sqlRepositoryFactory) CreateMessageRepository(sessionName string) IMessageRepository {
	panic("implement me")
}

