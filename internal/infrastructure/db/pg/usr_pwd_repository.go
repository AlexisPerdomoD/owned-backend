package pg

import (
	"context"

	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type UsrPwdRepository struct {
	db sqlx.ExtContext
}

func (r *UsrPwdRepository) SetPwd(ctx context.Context, id domain.UsrID, pwd []byte) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *UsrPwdRepository) GetPwd(ctx context.Context, id domain.UsrID) ([]byte, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func NewUsrPwdRepository(db sqlx.ExtContext) *UsrPwdRepository {
	return &UsrPwdRepository{db}
}
