package repository

import (
	"github.com/challenge/pkg/repository/memory"
	"sync"
)

type memoryRepositoryFactory struct {
	sync sync.Mutex
	userRepos map[string] IUserRepository
	msgsRepos map[string] IMessageRepository
	healthRepos map[string] IHealthRepository
}

func NewMemoryRepositoryFactory() IRepositoryFactory {
	return &memoryRepositoryFactory{
		sync:      sync.Mutex{},
		userRepos: make(map[string] IUserRepository, 0),
		msgsRepos: make(map[string] IMessageRepository, 0),
		healthRepos: make(map[string] IHealthRepository, 0),
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
	m.sync.Lock()
	defer m.sync.Unlock()

	if repo, ok := m.msgsRepos[sessionName]; ok {
		return repo
	}
	output := memory.NewMemoryMessageRepository()
	m.msgsRepos[sessionName] = output

	return output
}

func (m memoryRepositoryFactory) CreateHealthRepository(sessionName string)  IHealthRepository {
	m.sync.Lock()
	defer m.sync.Unlock()

	if repo, ok := m.healthRepos[sessionName]; ok {
		return repo
	}
	output := memory.NewMemoryHealthRepositoryRepository()
	m.healthRepos[sessionName] = output

	return output
}

