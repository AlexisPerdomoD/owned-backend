package pg

import (
	"context"
	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type UsrRepository struct {
	db sqlx.ExtContext
}

func (r *UsrRepository) GetByID(ctx context.Context, id domain.UsrID) (*domain.Usr, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *UsrRepository) GetByUsername(ctx context.Context, username string) (*domain.Usr, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *UsrRepository) Create(ctx context.Context, usr *domain.Usr) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *UsrRepository) Update(ctx context.Context, usr *domain.Usr) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *UsrRepository) Delete(ctx context.Context, id domain.UsrID) error {
	return apperror.ErrNotImplemented(nil)
}

func NewUsrRepository(db sqlx.ExtContext) *UsrRepository {
	return &UsrRepository{db}
}
