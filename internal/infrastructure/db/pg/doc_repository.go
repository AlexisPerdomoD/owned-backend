package pg

import (
	"context"

	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type DocRepository struct {
	db sqlx.ExtContext
}

func (r *DocRepository) GetByID(ctx context.Context, id domain.DocID) (*domain.Doc, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *DocRepository) GetByNodeID(ctx context.Context, id domain.DocID) (*domain.Doc, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *DocRepository) GetAllFromNodeID(ctx context.Context, id domain.NodeID) ([]domain.Doc, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *DocRepository) Create(ctx context.Context, d *domain.Doc) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *DocRepository) Update(ctx context.Context, d *domain.Doc) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *DocRepository) Delete(ctx context.Context, id domain.DocID) error {
	return apperror.ErrNotImplemented(nil)
}

func NewDocRepository(db sqlx.ExtContext) *DocRepository {
	return &DocRepository{db}
}
