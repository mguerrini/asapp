package security

import (
	"context"
	"github.com/challenge/pkg/models/errors"
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/modules/logger"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

type jwtTokenAuthenticationFactory struct {

}

type jwtTokenAuthentication struct {
	expirationTime int //minutes
	key []byte
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (j jwtTokenAuthenticationFactory) Create() ITokenAuthentication {
	expiration, err := config.ConfigurationSingleton().GetInt("root.token_auth.expiration")
	if err != nil {
		logger.Warn("JWS Authentication expiration's is not defined")
		expiration = 5
	}

	key, err := config.ConfigurationSingleton().GetString("root.token_auth.jwt_key")
	if err != nil {
		logger.Warn("JWS Authentication key's is not defined")
		key = "123456789987654312"
	}

	return &jwtTokenAuthentication{
		expirationTime: expiration,
		key: []byte(key),
	}
}

func (j jwtTokenAuthentication) GenerateToken(ctx context.Context, user string) (string, error) {
	tokenDuration := time.Duration(j.expirationTime) * time.Minute
	expirationTime := time.Now().Add(tokenDuration)

	claims := &Claims{
		Username:       user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.key)

	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", errors.NewInternalServerError(err)
	}

	return tokenString, nil
}

func (j jwtTokenAuthentication) ValidateToken(ctx context.Context, token string) TokenStatus {
	//remove bearer Prefix
	auxToken := strings.ToLower(token)
	auxToken = strings.TrimPrefix(auxToken, "bearer")

	trim := len(token) - len(auxToken)
	if trim > 0 {
		token = token[trim:]
		token = strings.TrimSpace(token)
	}

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})

	if err != nil {
		if !tkn.Valid && err != nil && strings.Contains(err.Error(), "token is expired by") {
			return SecurityTokenStatus_Expired
		} else {
			return SecurityTokenStatus_Invalid
		}
	}

	return SecurityTokenStatus_OK
}
