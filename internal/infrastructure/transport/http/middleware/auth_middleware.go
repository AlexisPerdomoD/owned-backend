package middleware

import (
	"net/http"
	"ownned/internal/infrastructure/auth"
	"ownned/internal/infrastructure/transport/http/response"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
	"strings"
)

type AuthMiddleware struct {
	jwtValidator auth.JWTValidator
}

func (m *AuthMiddleware) IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			_ = response.WriteJSONError(w, apperror.ErrUnauthenticated(
				map[string]string{
					"general": "session token not found",
				}))
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			_ = response.WriteJSONError(w, apperror.ErrUnauthenticated(
				map[string]string{
					"general": "malformed session token",
				}))
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		session, err := m.jwtValidator.Validate(token)
		if err != nil {
			_ = response.WriteJSONError(w, err)
			return
		}

		ctx := auth.SetSession(r.Context(), session)

		next(w, r.WithContext(ctx))
	})

}

func NewAuthMiddleware(jwtValidator auth.JWTValidator) *AuthMiddleware {
	helper.NotNilOrPanic(jwtValidator, "jwtValidator")
	return &AuthMiddleware{jwtValidator}
}
