package domain

import "errors"

type AppError struct {
	error
	Detail map[string]string
}

func (e *AppError) GetInstance() error {
	return e.error
}

var (
	ErrNotFoundInstance        = errors.New("not_found")
	ErrBadRequestInstance      = errors.New("bad_request")
	ErrConflictInstance        = errors.New("conflict")
	ErrUnauthenticatedInstance = errors.New("unauthenticated")
	ErrForbiddenInstance       = errors.New("forbidden")
	ErrAbortedInstance         = errors.New("aborted")
	ErrRateLimitInstance       = errors.New("rate_limit")
	ErrExternalServiceInstance = errors.New("external_service")
	ErrInternalInstance        = errors.New("internal")
	ErrNotImplementedInstance  = errors.New("not_implemented")
	ErrUnknownInstance         = errors.New("unknown")
)

func ErrNotFound(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrNotFoundInstance, Detail: detail}
}

func ErrBadRequest(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrBadRequestInstance, Detail: detail}
}

func ErrConflic(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrConflictInstance, Detail: detail}
}

func ErrUnauthenticated(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrUnauthenticatedInstance, Detail: detail}
}

func ErrForbidden(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrForbiddenInstance, Detail: detail}
}

func ErrAborted(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrAbortedInstance, Detail: detail}
}

func ErrRateLimit(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrRateLimitInstance, Detail: detail}
}

func ErrExternalService(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrExternalServiceInstance, Detail: detail}
}

func ErrInternal(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrInternalInstance, Detail: detail}
}

func ErrNotImplemented(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrNotImplementedInstance, Detail: detail}
}

func ErrUnknown(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrUnknownInstance, Detail: detail}
}
