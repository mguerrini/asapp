package security

import (
	"context"
	"github.com/challenge/pkg/models"
)

type noneAuthentication struct {
	
}

func (n noneAuthentication) Authenticate(ctx context.Context, cred models.Login) error {
	return nil
}

type noneAuthenticationFactory struct {

}

func (n noneAuthenticationFactory) Create(sessionName string) IAuthentication {
	return &noneAuthentication{}
}

