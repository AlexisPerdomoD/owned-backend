package error_pkg

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type HTTPError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Detail  map[string]string `json:"detail"`
}

func mapValidationError(err validator.ValidationErrors) *HTTPError {
	detail := make(map[string]string)

	for _, fe := range err {
		detail[fe.Field()] = fmt.Sprintf("failed on %s, %s", fe.Tag(), fe.Error())
	}

	return &HTTPError{
		Code:    http.StatusUnprocessableEntity,
		Message: "validation failed",
		Detail:  detail,
	}

}

func mapAppError(err *AppError) *HTTPError {
	switch err.error {
	case ErrNotFoundInstance:
		return &HTTPError{
			Code:    http.StatusNotFound,
			Message: "resource not found",
			Detail:  err.Detail,
		}

	case ErrBadRequestInstance:
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "bad request",
			Detail:  err.Detail,
		}

	case ErrConflictInstance:
		return &HTTPError{
			Code:    http.StatusConflict,
			Message: "conflict",
			Detail:  err.Detail,
		}

	case ErrForbiddenInstance:
		return &HTTPError{
			Code:    http.StatusForbidden,
			Message: "forbidden",
			Detail:  err.Detail,
		}

	case ErrUnauthenticatedInstance:
		return &HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "unauthenticated",
			Detail:  err.Detail,
		}

	case ErrRateLimitInstance:
		return &HTTPError{
			Code:    http.StatusTooManyRequests,
			Message: "too many requests",
			Detail:  err.Detail,
		}

	case ErrInternalInstance:
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "internal error",
			Detail:  err.Detail,
		}

	default:
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "unknown error",
			Detail:  err.Detail,
		}
	}
}

var errLogger = slog.With("error handler")

func MapError(err error) *HTTPError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return mapAppError(appErr)
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return mapValidationError(validationErrors)
	}

	errLogger.Error("unexpected error happened", "error", err)

	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: "unexpected error",
	}
}
