package security

import (
	"context"
	"github.com/challenge/pkg/models"
)

type dbAuthentication struct {
	
}

func (d dbAuthentication) Authenticate(ctx context.Context, cred models.Login) error {
	panic("implement me")
}

type dbAuthenticationFactory struct {

}

func (d dbAuthenticationFactory) Create(sessionName string) IAuthentication {
	panic("implement me")
}
