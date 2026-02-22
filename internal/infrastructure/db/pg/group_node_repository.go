package pg

import (
	"context"

	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type GroupNodeRepository struct {
	db sqlx.ExtContext
}

func (r *GroupNodeRepository) GetByNode(ctx context.Context, nodeID domain.NodeID) ([]domain.GroupNode, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *GroupNodeRepository) GetByGroup(ctx context.Context, groupID domain.GroupID) ([]domain.GroupNode, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *GroupNodeRepository) Upsert(ctx context.Context, d *domain.UpsertGroupNode) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *GroupNodeRepository) UpsertAll(ctx context.Context, d []domain.UpsertGroupNode) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *GroupNodeRepository) RemoveNode(ctx context.Context, g domain.GroupID, n domain.NodeID) error {
	return apperror.ErrNotImplemented(nil)
}

func NewGroupNodeRepository(db sqlx.ExtContext) *GroupNodeRepository {
	if db == nil {
		panic("NewGroupNodeRepository received a nil db")
	}

	return &GroupNodeRepository{db}
}
