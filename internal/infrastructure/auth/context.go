package auth

import (
	"context"
	"ownned/pkg/apperror"
)

type Session struct {
	UserID  string `json:"user_id"`
	Role    int    `json:"role"`
	DateISO int    `json:"date_iso"`
}

// context keys for user session
type usrSessionKeyType struct{}

// unique instance of usrSessionKeyType to avoid collisions
var usrSessionKey = usrSessionKeyType{}

// GetSession returns the session from the context
func GetSession(ctx context.Context) (*Session, error) {
	session, ok := ctx.Value(usrSessionKey).(Session)
	if !ok {
		return nil, apperror.ErrInternal(map[string]string{
			"error": "invalid casting of expected type *Session",
		})
	}

	return &session, nil
}

// SetSession sets the session in the context
func SetSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, usrSessionKey, session)
}
