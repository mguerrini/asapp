package security

import (
	"context"
)

type noneTokenAuthentication struct {

}

type noneTokenAuthenticationFactory struct {
	
}

func (n noneTokenAuthenticationFactory) Create() ITokenAuthentication {
	return &noneTokenAuthentication{}
}

func (t noneTokenAuthentication) GenerateToken(ctx context.Context, user string) string {
	return "0000000000-0000000000-0000000000-0000000000"
}

func (t noneTokenAuthentication) ValidateToken(ctx context.Context, token string) TokenStatus {
	return SecurityTokenStatus_OK
}


