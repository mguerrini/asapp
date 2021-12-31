package security

import (
	"context"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/models/errors"
	"github.com/challenge/pkg/repository"
)

type dbAuthenticationFactory struct {
}

func (d dbAuthenticationFactory) Create(sessionName string) IAuthentication {
	rep := repository.RepositoryFactory().CreateUserRepository(sessionName)
	return &dbAuthentication {
		userRep: rep,
	}
}


type dbAuthentication struct {
	userRep repository.IUserRepository
}

func (d dbAuthentication) Authenticate(ctx context.Context, cred models.Login) error {
	exist, err := d.userRep.ExistUsername(ctx, cred.Username)
	if err != nil {
		return err
	}

	if !exist {
		return errors.NewBadRequestMsg("Invalid username")
	}

	pass, err := d.userRep.GetPassword(ctx, cred.Username)
	if err != nil {
		return err
	}

	if pass == cred.Password {
		return nil
	} else {
		return errors.NewBadRequestMsg("Invalid pass word")
	}
}
