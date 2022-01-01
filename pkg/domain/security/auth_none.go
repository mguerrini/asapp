package security

import (
	"context"
	"github.com/challenge/pkg/models"
)

type noneAuthenticationFactory struct {

}

func (n noneAuthenticationFactory) Create(sessionName string) IAuthentication {
	return &noneAuthentication{}
}


type noneAuthentication struct {

}


func (n noneAuthentication) Authenticate(ctx context.Context, cred models.Login) error {
	return nil
}

func (d noneAuthentication) GeneratePassword(ctx context.Context, password string) (string, error) {
	return password, nil
}


