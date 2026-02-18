package auth

import "ownned/pkg/apperror"

// JWTGenerator is an interface for generating JWT tokens
// using internal app configurations
type JWTGenerator interface {
	Generate(usr *Session) (string, error)
}

// JWTValidator is an interface for validating JWT tokens
// using internal app configurations
type JWTValidator interface {
	Validate(token string) (*Session, error)
}

type JWTService struct {
	secret string
}

func (j *JWTService) Generate(usr *Session) (string, error) {
	return "", apperror.ErrNotImplemented(nil)
}

func (j *JWTService) Validate(token string) (*Session, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret}
}
