package services

import (
	"context"
	"github.com/challenge/pkg/domain/security"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/models/errors"
	"github.com/challenge/pkg/repository"
)

type UserServices struct {
	userRepository repository.IUserRepository
	authentication security.IAuthentication
}

func NewUserServices (sessionName string) *UserServices {
	rep := repository.RepositoryFactory().CreateUserRepository(sessionName)
	auth := security.AuthenticationFactory().Create(sessionName)

	return &UserServices{
		userRepository: rep,
		authentication: auth,
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

	//create pass
	newPass, err := u.authentication.GeneratePassword(ctx, user.Password)

	if err != nil {
		return nil, err
	}

	user.Password = newPass

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

func (u *UserServices) GetUserProfileByUsername(ctx context.Context, username string) (*models.UserProfile, error) {
	if username == "" {
		return nil, errors.NewBadRequestMsg("The username can not be empty.")
	}

	profile, err := u.userRepository.GetProfileByUsername(ctx, username)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (u *UserServices) GetUserProfileById(ctx context.Context, userId int) (*models.UserProfile, error) {
	if userId <= 0 {
		return nil, errors.NewBadRequestMsg("Invalid user id.")
	}

	profile, err := u.userRepository.GetProfileById(ctx, userId)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

