// Package auth provides the implementation of the authentication layer.
package auth

import (
	"context"

	authapp "ownned/internal/application/auth"
	"ownned/pkg/apperror"
)

// context keys for user session
type usrSessionKeyType struct{}

// unique instance of usrSessionKeyType to avoid collisions
var usrSessionKey = usrSessionKeyType{}

// GetSession returns the session from the context
func GetSession(ctx context.Context) (*authapp.JWTAccessPayload, error) {
	session, ok := ctx.Value(usrSessionKey).(authapp.JWTAccessPayload)
	if !ok {
		return nil, apperror.ErrInternal(map[string]string{
			"error": "invalid casting of expected type *Session",
		})
	}

	return &session, nil
}

// SetSession sets the session in the context
func SetSession(ctx context.Context, session *authapp.JWTAccessPayload) context.Context {
	return context.WithValue(ctx, usrSessionKey, session)
}
