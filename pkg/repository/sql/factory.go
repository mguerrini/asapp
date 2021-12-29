package sql

import (
	"github.com/challenge/pkg/repository"
)

type SqlRepositoryFactory struct {

}

func (s SqlRepositoryFactory) CreateUserRepository(sessionName string) repository.IUserRepository {
	panic("implement me")
}

func (s SqlRepositoryFactory) CreateMessageRepository(sessionName string) repository.IMessageRepository {
	panic("implement me")
}


