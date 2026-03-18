package middleware

import (
	"net/http"
	"strings"

	"ownned/internal/application/auth"
	"ownned/internal/infrastructure/sctx"
	"ownned/internal/infrastructure/transport/http/response"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
)

type AuthMiddleware struct {
	jwtValidator auth.JWTManager
}

func (m *AuthMiddleware) getSessionFromBearer(h http.Header) (*auth.JWTAccessPayload, error) {
	token := h.Get("Authorization")
	if token == "" {
		return nil, apperror.ErrUnauthenticated(
			map[string]string{
				"general": "session token not found",
			})
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return nil, apperror.ErrUnauthenticated(
			map[string]string{
				"general": "malformed session token",
			})
	}

	token = strings.TrimPrefix(token, "Bearer ")
	session, err := m.jwtValidator.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (m *AuthMiddleware) IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.getSessionFromBearer(r.Header)
		if err != nil {
			_ = response.WriteJSONError(w, err)
			return
		}

		ctx := sctx.SetSession(r.Context(), session)
		next(w, r.WithContext(ctx))
	})
}

func NewAuthMiddleware(jwtValidator auth.JWTManager) *AuthMiddleware {
	helper.NotNilOrPanic(jwtValidator, "jwtValidator")
	return &AuthMiddleware{jwtValidator}
}
