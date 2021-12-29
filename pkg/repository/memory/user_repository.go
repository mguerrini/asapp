package memory

import (
	"context"
	"errors"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/repository"
	"sync"
	"sync/atomic"
)

type MemoryUserRepository struct {
	sync  sync.Mutex
	users []models.User
	idSeq int32
}

func NewMemoryUserRepository () repository.IUserRepository {
	return &MemoryUserRepository{
		users: make([]models.User, 0),
	}
}

func (m MemoryUserRepository) GetPassword(ctx context.Context, userName string) (string, error) {
	for _, v := range m.users {
		if userName == v.Username {
			return v.Password, nil
		}
	}

	return "", errors.New("Invalid username")
}

func (m MemoryUserRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	id := atomic.AddInt32(&m.idSeq, 1)
	user.Id = int(id)

	m.sync.Lock()
	defer m.sync.Unlock()

	m.users = append(m.users, user)
	return int(id), nil
}


func (m MemoryUserRepository) ExistUsername(ctx context.Context, username string) (bool, error) {
	m.sync.Lock()
	defer m.sync.Unlock()

	for _, u := range m.users {
		if u.Username == username {
			return true, nil
		}
	}

	return false, nil
}

func (m MemoryUserRepository) GetProfileByUsername(username string) (*models.UserProfile, error) {
	m.sync.Lock()
	defer m.sync.Unlock()

	for _, u := range m.users {
		if u.Username == username {
			profile := models.UserProfile{
				Id:       u.Id,
				Username: u.Username,
			}
			return &profile, nil
		}
	}

	return nil, nil
}
