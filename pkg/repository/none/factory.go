package none

import (
	"github.com/challenge/pkg/repository"
)

type noneRepositoryFactory struct {

}

func (n noneRepositoryFactory) CreateUserRepository(sessionName string) repository.IUserRepository {
	panic("implement me")
}

func (n noneRepositoryFactory) CreateMessageRepository(sessionName string) repository.IMessageRepository {
	panic("implement me")
}


