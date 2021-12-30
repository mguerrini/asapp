package repository

import (
	"context"
	"github.com/challenge/pkg/models"
)

type IUserRepository interface {
	GetPassword(ctx context.Context, userName string) (string, error)
	CreateUser (ctx context.Context, user models.User) (int, error)
	ExistUsername(ctx context.Context, username string) (bool, error)

	GetProfileById(ctx context.Context, userId int) (*models.UserProfile, error)
	GetProfileByUsername(ctx context.Context, username string) (*models.UserProfile, error)
}





