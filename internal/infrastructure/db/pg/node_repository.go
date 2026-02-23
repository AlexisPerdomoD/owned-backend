package pg

import (
	"context"

	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type NodeRepository struct {
	db sqlx.ExtContext
}

func (r *NodeRepository) GetByID(ctx context.Context, id domain.NodeID) (*domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) GetByIDs(ctx context.Context, ids []domain.NodeID) ([]domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) GetChildren(ctx context.Context, folderID domain.NodePath) ([]domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) GetRoot(ctx context.Context) ([]domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) GetRootByGroups(ctx context.Context, groups []domain.GroupID) ([]domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) GetByGroup(ctx context.Context, groupID domain.GroupID) ([]domain.NodeGroupAttach, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) Create(ctx context.Context, n *domain.Node) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) Update(ctx context.Context, n *domain.Node) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) Delete(ctx context.Context, id domain.NodeID) error {
	return apperror.ErrNotImplemented(nil)
}

func NewNodeRepository(db sqlx.ExtContext) *NodeRepository {
	return &NodeRepository{db}
}
