package pg

import (
	"context"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
)

type DocRepository struct {
}

func (r *DocRepository) GetByID(ctx context.Context, id string) (*domain.Doc, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *DocRepository) GetByNodeID(ctx context.Context, id string) ([]domain.Doc, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *DocRepository) Create(ctx context.Context, d *domain.Doc) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *DocRepository) Update(ctx context.Context, d *domain.Doc) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *DocRepository) Delete(ctx context.Context, id string) error {
	return apperror.ErrNotImplemented(nil)
}

func NewDocRepository() *DocRepository {
	return &DocRepository{}
}
