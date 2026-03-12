package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

const getNodeQuery string = `
SELECT
	n.id,
	n.name,
	n.description,
	n.path,
	n.type,
	n.created_at,
	n.updated_at
FROM fs.nodes n`

const insertNodeQuery string = `
INSERT INTO fs.nodes (
	id,
	name,
	description,
	path,
	type,
	created_at,
	updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)`

type nodeRow struct {
	ID          domain.NodeID `db:"id"`
	Name        string        `db:"name"`
	Description string        `db:"description"`
	Path        string        `db:"path"`
	Type        string        `db:"type"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}

func (r *nodeRow) ToDomain() domain.Node {
	return domain.Node{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Path:        domain.NodePath(r.Path),
		Type:        domain.NodeType(r.Type),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

type nodeRepository struct {
	db sqlx.ExtContext
}

func (r *nodeRepository) GetByID(ctx context.Context, id domain.NodeID) (*domain.Node, error) {
	q := fmt.Sprintf("%s\nWHERE n.id=$1", getNodeQuery)
	row := &nodeRow{}
	if err := r.db.QueryRowxContext(ctx, q, id).StructScan(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	d := row.ToDomain()
	return &d, nil
}

func (r *nodeRepository) GetChildren(ctx context.Context, folderID domain.NodePath) ([]domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *nodeRepository) GetRoot(ctx context.Context) ([]domain.Node, error) {
	q := fmt.Sprintf("%s\nWHERE nlevel(n.path)=1", getNodeQuery)

	rows, err := r.db.QueryxContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer safeClose(ctx, rows)
	nodes, err := readSlice[domain.Node, nodeRow](rows)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *nodeRepository) GetRootByGroups(ctx context.Context, groups []domain.GroupID) ([]domain.Node, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *nodeRepository) GetByGroup(ctx context.Context, groupID domain.GroupID) ([]domain.NodeGroupAttach, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *nodeRepository) Create(ctx context.Context, n *domain.Node) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *nodeRepository) Update(ctx context.Context, n *domain.Node) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *nodeRepository) Delete(ctx context.Context, id domain.NodeID) error {
	return apperror.ErrNotImplemented(nil)
}

func NewNodeRepository(db sqlx.ExtContext) domain.NodeRepository {
	return &nodeRepository{db}
}
