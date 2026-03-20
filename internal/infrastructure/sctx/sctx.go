// Package sctx provides the implementation of the authentication layer to centralize session information handling from contexts.
package sctx

import (
	"context"

	"ownned/internal/application/auth"
	"ownned/pkg/apperror"
)

// context keys for user session
type usrSessionKeyType struct{}

// unique instance of usrSessionKeyType to avoid collisions
var usrSessionKey = usrSessionKeyType{}

// GetSession returns the session from the context
func GetSession(ctx context.Context) (*auth.JWTAccessPayload, error) {
	session, ok := ctx.Value(usrSessionKey).(*auth.JWTAccessPayload)
	if !ok {
		detail := make(map[string]string)
		detail["reason"] = "invalid casting of expected type *Session"
		return nil, apperror.ErrInternal(detail)
	}

	return session, nil
}

// SetSession sets the session in the context
func SetSession(ctx context.Context, session *auth.JWTAccessPayload) context.Context {
	return context.WithValue(ctx, usrSessionKey, session)
}
