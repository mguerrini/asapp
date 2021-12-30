package security

import (
	"context"
	"strings"
)

type noneTokenAuthentication struct {

}

type noneTokenAuthenticationFactory struct {
	
}

func (n noneTokenAuthenticationFactory) Create() ITokenAuthentication {
	return &noneTokenAuthentication{}
}

func (t noneTokenAuthentication) GenerateToken(ctx context.Context, user string) (string, error) {
	return "0000000000-0000000000-0000000000-0000000000", nil
}

func (t noneTokenAuthentication) ValidateToken(ctx context.Context, token string) TokenStatus {
	auxToken := strings.ToLower(token)
	auxToken = strings.TrimPrefix(auxToken, "bearer")

	trim := len(token) - len(auxToken)
	if trim > 0 {
		token = token[trim:]
		token = strings.TrimSpace(token)
	}

	if token == "" || token == "0000000000-0000000000-0000000000-0000000000"{
		return SecurityTokenStatus_OK
	} else {
		return SecurityTokenStatus_Invalid
	}
}


