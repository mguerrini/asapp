package services

import (
	"context"
	"github.com/challenge/pkg/domain/security"
	"github.com/challenge/pkg/models"
)

type AuthServices struct {
	authentication security.IAuthentication
	authorization  security.ITokenAuthentication
}



func NewAuthServices (sessionName string) *AuthServices {
	tauth := security.TokenAuthenticationFactory().Create()
	auth := security.AuthenticationFactory().Create(sessionName)

	return &AuthServices{
		authentication: auth,
		authorization:  tauth,
	}
}


func (s *AuthServices) ValidateUser(ctx context.Context, cred models.Login) error {
	err := s.authentication.Authenticate(ctx, cred)
	return err
}

func (s *AuthServices) GenerateToken(ctx context.Context, userName string) (token string) {
	return s.authorization.GenerateToken(ctx, userName)
}

func (s *AuthServices) ValidateToken(ctx context.Context, token string) security.TokenStatus {
	return s.authorization.ValidateToken(ctx, token)
}
