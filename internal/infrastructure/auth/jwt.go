package auth

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
