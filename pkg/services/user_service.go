package services

import (
	"context"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/models/errors"
	"github.com/challenge/pkg/repository"
)

type UserServices struct {
	userRepository repository.IUserRepository
}

func NewUserServices (sessionName string) *UserServices {
	rep := repository.RepositoryFactory().CreateUserRepository(sessionName)

	return &UserServices{
		userRepository: rep,
	}
}


func (u *UserServices) CreateUser(ctx context.Context, user models.User) (*models.UserProfile, error)  {
	//validate user data
	if user.Username == "" {
		return nil, errors.NewBadRequestMsg("The username can not be empty.")
	}

	if user.Password == "" {
		return nil, errors.NewBadRequestMsg("The password can not be empty.")
	}

	//validate existent user
	exist, err := u.userRepository.ExistUsername(ctx, user.Username)

	if err != nil {
		return nil, err
	}

	if exist {
		return nil, errors.NewBadRequestMsg("The username already exists.")
	}

	//create
	id, errC := u.userRepository.CreateUser(ctx, user)
	if errC != nil {
		return nil, errC
	}

	return &models.UserProfile{
		Id:       id,
		Username: user.Username,
	}, nil
}

func (u *UserServices) GetUserProfile(username string) (*models.UserProfile, error) {
	if username == "" {
		return nil, errors.NewBadRequestMsg("The username can not be empty.")
	}

	profile, err := u.userRepository.GetProfileByUsername(username)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

