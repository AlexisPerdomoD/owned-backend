package pg

import (
	"context"

	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type GroupRepository struct {
	db sqlx.ExtContext
}

func (r *GroupRepository) GetByID(ctx context.Context, id domain.GroupID) (*domain.Group, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *GroupRepository) GetByUsr(ctx context.Context, usrID domain.UsrID) ([]domain.GroupUsr, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *GroupRepository) GetByNode(ctx context.Context, nodeID domain.NodeID) ([]domain.GroupNode, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *GroupRepository) Create(ctx context.Context, d *domain.Group) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *GroupRepository) Update(ctx context.Context, d *domain.UpdateGroup) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *GroupRepository) Delete(ctx context.Context, id domain.GroupID) error {
	return apperror.ErrNotImplemented(nil)
}

func NewGroupRepository(db sqlx.ExtContext) *GroupRepository {
	return &GroupRepository{db}
}
