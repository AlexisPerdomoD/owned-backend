package mapper

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"ownned/internal/infrastructure/transport/http/model"
	"ownned/pkg/apperror"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err validator.ValidationErrors) *model.ErrView {
	detail := make(map[string]string)

	for _, fe := range err {
		detail[fe.Field()] = fmt.Sprintf("failed on %s, %s", fe.Tag(), fe.Error())
	}

	return &model.ErrView{
		Code:    http.StatusUnprocessableEntity,
		Message: "validation failed",
		Detail:  detail,
	}

}

func AppError(err *apperror.AppError) *model.ErrView {
	switch err {
	case apperror.ErrNotFoundInstance:
		return &model.ErrView{
			Code:    http.StatusNotFound,
			Message: "resource not found",
			Detail:  err.Detail,
		}

	case apperror.ErrBadRequestInstance:
		return &model.ErrView{
			Code:    http.StatusBadRequest,
			Message: "bad request",
			Detail:  err.Detail,
		}

	case apperror.ErrConflictInstance:
		return &model.ErrView{
			Code:    http.StatusConflict,
			Message: "conflict",
			Detail:  err.Detail,
		}

	case apperror.ErrForbiddenInstance:
		return &model.ErrView{
			Code:    http.StatusForbidden,
			Message: "forbidden",
			Detail:  err.Detail,
		}

	case apperror.ErrUnauthenticatedInstance:
		return &model.ErrView{
			Code:    http.StatusUnauthorized,
			Message: "unauthenticated",
			Detail:  err.Detail,
		}

	case apperror.ErrRateLimitInstance:
		return &model.ErrView{
			Code:    http.StatusTooManyRequests,
			Message: "too many requests",
			Detail:  err.Detail,
		}

	case apperror.ErrInternalInstance:
		return &model.ErrView{
			Code:    http.StatusInternalServerError,
			Message: "internal error",
			Detail:  err.Detail,
		}

	default:
		return &model.ErrView{
			Code:    http.StatusInternalServerError,
			Message: "unknown error",
			Detail:  err.Detail,
		}
	}
}

var errLogger = slog.With("error handler")

func Err(err error) *model.ErrView {
	if err == nil {
		return nil
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return AppError(appErr)
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return ValidationError(validationErrors)
	}

	errLogger.Error("unexpected error happened", "error", err)

	return &model.ErrView{
		Code:    http.StatusInternalServerError,
		Message: "unexpected error",
	}
}
