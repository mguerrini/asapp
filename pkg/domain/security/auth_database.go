package security

import (
	"context"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/models/errors"
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type dbAuthenticationFactory struct {
}

func (d dbAuthenticationFactory) Create(sessionName string) IAuthentication {
	rep := repository.RepositoryFactory().CreateUserRepository(sessionName)
	mode, err := config.ConfigurationSingleton().GetString("root.authentication.password")

	if err != nil {
		panic(err)
	}

	return &dbAuthentication {
		userRep: rep,
		passwordMode: mode,
	}
}


type dbAuthentication struct {
	userRep      repository.IUserRepository
	passwordMode string
}

func (d dbAuthentication) GeneratePassword(ctx context.Context, password string) (string, error) {
	if d.passwordMode == "plain" {
		return password, nil
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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

	if d.passwordMode == "hash" {
		err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(cred.Password))
		return err
	} else {
		if pass == cred.Password {
			return nil
		} else {
			return errors.NewBadRequestMsg("Invalid pass word")
		}
	}
}
