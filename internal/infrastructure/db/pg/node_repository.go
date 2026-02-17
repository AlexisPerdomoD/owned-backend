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

func (r *NodeRepository) GetByID(ctx context.Context, id string) (*domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) GetByIDs(ctx context.Context, ids []string) ([]domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) GetChildren(ctx context.Context, folderID string) ([]domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) GetRoot(ctx context.Context) ([]domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) GetRootByUsr(ctx context.Context, usrID string) ([]domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) Create(ctx context.Context, n *domain.Node) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) Update(ctx context.Context, n *domain.Node) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) Delete(ctx context.Context, id string) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) GetAccess(ctx context.Context, u string, n string) (domain.NodeAccess, error) {
	return 0, apperror.ErrNotImplemented(nil)
}

func (r *NodeRepository) UpdateAccess(ctx context.Context, u string, n string, a domain.NodeAccess) error {
	return apperror.ErrNotImplemented(nil)
}

func NewNodeRepository(db sqlx.ExtContext) *NodeRepository {
	return &NodeRepository{db}
}
