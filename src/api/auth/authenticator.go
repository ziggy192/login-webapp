package auth

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
	"bitbucket.org/ziggy192/ng_lu/src/api/redis"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

const defaultIssuer = "ng_lu"

type Authenticator struct {
	secret              []byte
	expiresAfterMinutes int
	tokenBlocker        *TokenBlocker
}

func NewAuthenticator(cfg *config.Config, redisClient *redis.Redis) *Authenticator {
	return &Authenticator{
		secret:              []byte(cfg.AuthSecret),
		expiresAfterMinutes: cfg.JWTExpiresAfterMinutes,
		tokenBlocker:        NewTokenBlocker(redisClient),
	}
}

// SignUserJWT creates a new JWT token signed by HMAC method
func (a *Authenticator) SignUserJWT(ctx context.Context, username string) (string, error) {
	if a.secret == nil {
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
		Issuer:    defaultIssuer,
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.expiresAfterMinutes) * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        tokenUUID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(a.secret)
	if err != nil {
		logger.Err(ctx, err)
		return "", err
	}

	return signed, nil
}

// VerifyUserJWT checks if provided token string is valid or not
func (a *Authenticator) VerifyUserJWT(ctx context.Context, tokenString string) (*jwt.RegisteredClaims, error) {
	claims, err := parseWithClaims(a.secret, tokenString)
	if err != nil {
		logger.Err(ctx, err, tokenString)
		return &claims, err
	}

	blocked, err := a.tokenBlocker.IsBlocked(ctx, tokenString)
	if err != nil {
		return &claims, err
	}
	if blocked {
		return &claims, errors.New("token is logged out")
	}
	return &claims, nil
}

func (a *Authenticator) Logout(ctx context.Context, tokenString string) error {
	claims, _ := parseWithClaims(a.secret, tokenString)
	var duration time.Duration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.After(time.Now()) {
		duration = time.Until(claims.ExpiresAt.Time)
	}
	err := a.tokenBlocker.BlockToken(ctx, tokenString, time.Now(), duration)
	if err != nil {
		logger.Err(ctx, err)
	}
	return nil
}

func parseWithClaims(secret []byte, tokenString string) (jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	}
	_, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	return claims, err
}
