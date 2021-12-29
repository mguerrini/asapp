package repository

import (
	"github.com/challenge/pkg/repository/memory"
	"sync"
)

type memoryRepositoryFactory struct {
	sync sync.Mutex
	userRepos map[string] IUserRepository
	msgsRepos map[string] IMessageRepository
}

func NewMemoryRepositoryFactory() IRepositoryFactory {
	return &memoryRepositoryFactory{
		sync:      sync.Mutex{},
		userRepos: make(map[string] IUserRepository, 0),
		msgsRepos: make(map[string] IMessageRepository, 0),
	}
}

func (m memoryRepositoryFactory) CreateUserRepository(sessionName string) IUserRepository {
	m.sync.Lock()
	defer m.sync.Unlock()

	if repo, ok := m.userRepos[sessionName]; ok {
		return repo
	}
	output := memory.NewMemoryUserRepository()
	m.userRepos[sessionName] = output

	return output
}

func (m memoryRepositoryFactory) CreateMessageRepository(sessionName string) IMessageRepository {
	panic("implement me")
}

