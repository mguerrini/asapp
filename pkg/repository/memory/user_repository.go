package memory

import (
	"context"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/models/errors"
	"sync"
	"sync/atomic"
)

type memoryUserRepository struct {
	sync  sync.Mutex
	users []models.User
	idSeq int32
}

func NewMemoryUserRepository () *memoryUserRepository {
	return &memoryUserRepository{
		users: make([]models.User, 0),
	}
}

func (m *memoryUserRepository) GetPassword(ctx context.Context, userName string) (string, error) {
	for _, v := range m.users {
		if userName == v.Username {
			return v.Password, nil
		}
	}

	return "", errors.NewBadRequestMsg("Invalid username")
}

func (m *memoryUserRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	id := atomic.AddInt32(&m.idSeq, 1)
	user.Id = int(id)

	m.sync.Lock()
	defer m.sync.Unlock()

	m.users = append(m.users, user)
	return int(id), nil
}


func (m *memoryUserRepository) ExistUsername(ctx context.Context, username string) (bool, error) {
	m.sync.Lock()
	defer m.sync.Unlock()

	for _, u := range m.users {
		if u.Username == username {
			return true, nil
		}
	}

	return false, nil
}

func (m *memoryUserRepository) GetProfileByUsername(username string) (*models.UserProfile, error) {
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
