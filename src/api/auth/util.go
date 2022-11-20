package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
)

type ctxKey string

const contextKeyClaims ctxKey = "claims"

// SaveClaims saves claims to the context
func SaveClaims(ctx context.Context, claims *jwt.RegisteredClaims) context.Context {
	return context.WithValue(ctx, contextKeyClaims, claims)
}

// GetClaims returns the user's claims from context
func GetClaims(ctx context.Context) *jwt.RegisteredClaims {
	if ctx == nil {
		return nil
	}
	claims := ctx.Value(contextKeyClaims)
	if claims == nil {
		return nil
	}

	return claims.(*jwt.RegisteredClaims)
}

// GetUsername returns the authenticated username from context
func GetUsername(ctx context.Context) string {
	claims := GetClaims(ctx)
	if claims == nil {
		return ""
	}
	return claims.Subject
}
