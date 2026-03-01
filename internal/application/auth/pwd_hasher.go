package auth

import "errors"

var ErrInvalidHash = errors.New("invalid hash")

// PwdHasher is the interface for password hashing implementations of required algorithms.
type PwdHasher interface {
	// Hash returns the hashed password for the given password.
	Hash(password []byte) ([]byte, error)
	// Compare compares the given hashed password with the given password and returns nil if they match, otherwise ErrInvalidHash.
	// May return error different from ErrInvalidHash that means the comparison failed for some other reason.
	Compare(hashedPassword, password []byte) error
	// RequiredReHash returns true if the hashed password needs to be re-hashed for new policies, algorithm parameters, etc. False otherwise.
	RequiredReHash(hashedPassword []byte) bool
}
