package security

import (
	"context"
)

type noneTokenAuthentication struct {

}

type noneTokenAuthenticationFactory struct {
	
}

func (n noneTokenAuthenticationFactory) Create() ITokenAuthentication {
	panic("implement me")
}

func (t noneTokenAuthentication) GenerateToken(ctx context.Context, user string) string {
	panic("implement me")
}

func (t noneTokenAuthentication) ValidateToken(ctx context.Context, token string) TokenStatus {
	panic("implement me")
}


