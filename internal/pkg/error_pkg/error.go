package error_pkg

type AppError struct {
	error
	Detail map[string]string
}

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
	return &AppError{error: ErrConflicInstance, Detail: detail}
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

func ErrUnknown(
	detail map[string]string,
) *AppError {
	return &AppError{error: ErrUnknownInstance, Detail: detail}
}
