package memory

import "github.com/challenge/pkg/repository"

type MemoryRepositoryFactory struct {

}

func (m MemoryRepositoryFactory) CreateUserRepository(sessionName string) repository.IUserRepository {
	return NewMemoryUserRepository()
}

func (m MemoryRepositoryFactory) CreateMessageRepository(sessionName string) repository.IMessageRepository {
	panic("implement me")
}

