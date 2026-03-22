package auth

import (
	"errors"

	"ownned/internal/domain"
)

var ErrInvalidToken = errors.New("invalid token")

var ErrExpiredToken = errors.New("expired token")

// JWTAccessPayload is the payload of the JWT token for access tokens.
type JWTAccessPayload struct {
	UsrID string         `json:"usr_id"`
	Role  domain.UsrRole `json:"role"`
}

// JWTManager is the interface for JWT management implementations of required algorithms.
// It is responsible for generating and validating JWT tokens with application specific configuration.
type JWTManager interface {
	// GenerateAccessToken generates a JWT token for the given payload.
	// The token is signed with the configured secret.
	// The token is valid for the configured expiration time.
	// Returns the token string or an error.
	GenerateAccessToken(payload JWTAccessPayload) (string, error)

	// ValidateAccessToken validates the given token string.
	// Returns the payload or an error.
	// if the token is invalid or expired, returns ErrInvalidToken or ErrExpiredToken respectively.
	ValidateAccessToken(token string) (*JWTAccessPayload, error)
}
