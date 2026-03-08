package pg

import (
	"context"
	"fmt"
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

func (r *groupNodeRow) ToDomain() domain.GroupNode {
	return domain.GroupNode{
		GroupID:    r.GroupID,
		NodeID:     r.NodeID,
		AssignDate: r.AssignedAt,
	}
}

type groupNodeRepository struct {
	db sqlx.ExtContext
}

func (r *groupNodeRepository) GetByNode(ctx context.Context, nodeID domain.NodeID) ([]domain.GroupNode, error) {
	q := fmt.Sprintf("%s\nWHERE gn.node_id=$1", getGroupNodeQuery)
	rows, err := r.db.QueryxContext(ctx, q, nodeID)
	if err != nil {
		return nil, err
	}
	defer safeClose(ctx, rows)
	// por revisar bb si funciona soy dios
	res, err := readSlice[domain.GroupNode, groupNodeRow](rows)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *groupNodeRepository) Upsert(ctx context.Context, d *domain.UpsertGroupNode) error {
	if d == nil {
		return ErrInvalidArgument
	}

	row := groupNodeRow{
		NodeID:     d.NodeID,
		GroupID:    d.GroupID,
		AssignedAt: time.Now().UTC(),
	}

	_, err := sqlx.NamedExecContext(ctx, r.db, upsertGroupNodeQuery, row)

	return err
}

func (r *groupNodeRepository) UpsertAll(ctx context.Context, d []domain.UpsertGroupNode) error {
	if len(d) == 0 {
		return ErrInvalidArgument
	}
	assignedAt := time.Now().UTC()
	rows := make([]domain.GroupNode, len(d))
	for i, gn := range d {
		rows[i].GroupID = gn.GroupID
		rows[i].NodeID = gn.NodeID
		rows[i].AssignDate = assignedAt
	}

	_, err := sqlx.NamedExecContext(ctx, r.db, upsertGroupNodeQuery, rows)

	return err
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
