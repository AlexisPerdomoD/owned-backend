package mapper

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"ownned/pkg/apperror"

	"github.com/go-playground/validator/v10"
)

type ErrView struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Detail  map[string]string `json:"detail"`
}

func mapValidationError(err validator.ValidationErrors) *ErrView {
	detail := make(map[string]string)

	for _, fe := range err {
		detail[fe.Field()] = fmt.Sprintf("failed on %s, %s", fe.Tag(), fe.Error())
	}

	return &ErrView{
		Code:    http.StatusUnprocessableEntity,
		Message: "validation failed",
		Detail:  detail,
	}

}

func mapAppError(err *apperror.AppError) *ErrView {
	switch err {
	case apperror.ErrNotFoundInstance:
		return &ErrView{
			Code:    http.StatusNotFound,
			Message: "resource not found",
			Detail:  err.Detail,
		}

	case apperror.ErrBadRequestInstance:
		return &ErrView{
			Code:    http.StatusBadRequest,
			Message: "bad request",
			Detail:  err.Detail,
		}

	case apperror.ErrConflictInstance:
		return &ErrView{
			Code:    http.StatusConflict,
			Message: "conflict",
			Detail:  err.Detail,
		}

	case apperror.ErrForbiddenInstance:
		return &ErrView{
			Code:    http.StatusForbidden,
			Message: "forbidden",
			Detail:  err.Detail,
		}

	case apperror.ErrUnauthenticatedInstance:
		return &ErrView{
			Code:    http.StatusUnauthorized,
			Message: "unauthenticated",
			Detail:  err.Detail,
		}

	case apperror.ErrRateLimitInstance:
		return &ErrView{
			Code:    http.StatusTooManyRequests,
			Message: "too many requests",
			Detail:  err.Detail,
		}

	case apperror.ErrInternalInstance:
		return &ErrView{
			Code:    http.StatusInternalServerError,
			Message: "internal error",
			Detail:  err.Detail,
		}

	default:
		return &ErrView{
			Code:    http.StatusInternalServerError,
			Message: "unknown error",
			Detail:  err.Detail,
		}
	}
}

var errLogger = slog.With("error handler")

func MapError(err error) *ErrView {
	if err == nil {
		return nil
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return mapAppError(appErr)
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return mapValidationError(validationErrors)
	}

	errLogger.Error("unexpected error happened", "error", err)

	return &ErrView{
		Code:    http.StatusInternalServerError,
		Message: "unexpected error",
	}
}
