package pg

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const getGroupQuery string = `
SELECT 
	g.id,
	g.name,
	g.description,
	g.created_at,
	g.updated_at
FROM fs.group g`

const insertGroupQuery string = `
INSERT INTO fs.groups(
	id,
	name,
	description,
	created_at,
	updated_at
) VALUES ($1, $2, $3, $4, $5)`

const updateGroupQuery string = `
UPDATE fs.groups SET
	name = $1, 
	description = $2
WHERE id=$3`

type groupRow struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (r *groupRow) ToDomain() domain.Group {
	return domain.Group{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

type groupRepository struct {
	db sqlx.ExtContext
}

func (r *groupRepository) GetByID(ctx context.Context, id domain.GroupID) (*domain.Group, error) {
	q := fmt.Sprintf("%s\nWHERE id=$1", getGroupQuery)
	row := &groupRow{}
	err := r.db.QueryRowxContext(ctx, q, id).StructScan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	d := row.ToDomain()
	return &d, nil
}

func (r *groupRepository) GetByIDs(ctx context.Context, ids []domain.GroupID) (map[domain.GroupID]*domain.Group, error) {
	res := make(map[domain.GroupID]*domain.Group)
	for _, id := range ids {
		res[id] = nil
	}

	if len(ids) == 0 {
		return res, nil
	}

	q := fmt.Sprintf("%s\nWHERE id=ANY($1)", getGroupQuery)
	rows, err := r.db.QueryxContext(ctx, q, ids)
	if err != nil {
		return nil, err
	}
	defer safeClose(ctx, rows)

	for rows.Next() {
		row := &groupRow{}
		if err := rows.StructScan(row); err != nil {
			return nil, err
		}
		d := row.ToDomain()
		res[row.ID] = &d
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *groupRepository) GetByUsr(ctx context.Context, usrID domain.UsrID) ([]domain.Group, error) {
	q := fmt.Sprintf("%s\nINNER JOIN fs.group_usrs gu ON gu.group_id = g.id\nWHERE gu.usrID=$1", getGroupQuery)
	rows, err := r.db.QueryxContext(ctx, q, usrID)
	if err != nil {
		return nil, err
	}

	defer safeClose(ctx, rows)
	res, err := readSlice[domain.Group, groupRow](rows)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *groupRepository) Create(ctx context.Context, d *domain.Group) error {
	if d == nil {
		return ErrInvalidArgument
	}

	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, insertGroupQuery,
		d.ID,
		d.Name,
		d.Description,
		now,
		now)
	if err != nil {
		return err
	}

	d.CreatedAt = now
	d.UpdatedAt = now
	return nil
}

func (r *groupRepository) Update(ctx context.Context, d *domain.Group) error {
	if d == nil {
		return ErrInvalidArgument
	}

	res, err := r.db.ExecContext(ctx, updateGroupQuery, d.Name, d.Description, d.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count < 1 {
		return apperror.ErrNotFound(nil)
	}

	d.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *groupRepository) Delete(ctx context.Context, id domain.GroupID) error {
	return apperror.ErrNotImplemented(nil)
}

func NewGroupRepository(db sqlx.ExtContext) domain.GroupRepository {
	return &groupRepository{db}
}
