package pg

import (
	"context"
	"time"

	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

const getGroupNodeQuery string = `
SELECT 
	gn.group_id,
	gn.node_id,
	gn.assigned_at
FROM fs.group_nodes gn`

const upsertGroupNodeQuery = `
INSERT INTO fs.group_usrs(
	group_id,
	usr_id,
	assigned_at)
VALUES (
	:group_id,
	:usr_id,
	:assigned_at) 
ON CONFLICT(group_id, usr_id) DO NOTHING`

type groupNodeRow struct {
	GroupID    domain.GroupID `db:"group_id"`
	NodeID     domain.NodeID  `db:"node_id"`
	AssignedAt time.Time      `db:"assigned_at"`
}

type groupNodeRepository struct {
	db sqlx.ExtContext
}

func (r *groupNodeRepository) GetByNode(ctx context.Context, nodeID domain.NodeID) ([]domain.GroupNode, error) {
	return nil, nil
}

func (r *groupNodeRepository) Upsert(ctx context.Context, d *domain.UpsertGroupNode) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *groupNodeRepository) UpsertAll(ctx context.Context, d []domain.UpsertGroupNode) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *groupNodeRepository) RemoveNode(ctx context.Context, g domain.GroupID, n domain.NodeID) error {
	return apperror.ErrNotImplemented(nil)
}

func NewGroupNodeRepository(db sqlx.ExtContext) domain.GroupNodeRepository {
	if db == nil {
		panic("NewGroupNodeRepository received a nil db")
	}

	return &groupNodeRepository{db}
}
