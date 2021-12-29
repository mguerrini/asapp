package security

import "context"

type jwsTokenAuthenticationFactory struct {

}

type jwsTokenAuthentication struct {

}


func (j jwsTokenAuthenticationFactory) Create() ITokenAuthentication {
	panic("implement me")
}

func (j jwsTokenAuthentication) GenerateToken(ctx context.Context, user string) string {
	panic("implement me")
}

func (j jwsTokenAuthentication) ValidateToken(ctx context.Context, token string) TokenStatus {
	panic("implement me")
}
