package auth

import (
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type Authenticator struct {
	Secret []byte
}

func NewAuthenticator(secret []byte) *Authenticator {
	return &Authenticator{Secret: secret}
}

const DefaultIssuer = "ng_lu"

// SignUserJWT creates a new JWT token signed by HMAC method
func (a *Authenticator) SignUserJWT(ctx context.Context, username string) (string, error) {
	if a.Secret == nil {
		err := errors.New("cannot sign token without secret")
		logger.Err(ctx, err)
		return "", err
	}

	tokenUUID, err := uuid.NewRandom()
	if err != nil {
		logger.Err(ctx, err)
		return "", err
	}

	claims := jwt.RegisteredClaims{
		Issuer:   DefaultIssuer,
		Subject:  username,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       tokenUUID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		logger.Err(ctx, err)
		return "", err
	}

	return signed, nil
}

// VerifyUserJWT checks if provided token string is valid or not
func (a *Authenticator) VerifyUserJWT(ctx context.Context, tokenString string) (*jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.Secret, nil
	}
	_, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		logger.Err(ctx, err)
		return &claims, err
	}

	return &claims, nil
}
