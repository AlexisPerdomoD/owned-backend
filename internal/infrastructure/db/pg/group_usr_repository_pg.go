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

const getGroupUsrQuery string = `
SELECT 
	gu.group_id,
	gu.usr_id,
	gu.access,
	gu.assigned_at
FROM fs.group_usrs gu`

const getGroupAccessQuery string = "SELECT gu.access from fs.group_usrs gu\nWHERE gu.usr_id=$1 AND gu.group_id=$2"

const getNodeAccessQuery string = `
SELECT 
	gu.access 
FROM fs.group_usrs gu
INNER JOIN fs.group_nodes gn 
	ON gn.group_id = gu.group_id
WHERE 
	gn.node_id = $1
	AND gu.usr_id = $2
ORDER BY 
	CASE gu.access
    	WHEN 'write_access'     THEN 1
    	WHEN 'read_only_access' THEN 2
	END
LIMIT 1`

const upsertGroupUsrQuery = `
INSERT INTO fs.group_usrs(
	group_id,
	usr_id,
	access,
	assigned_at)
VALUES (
	:group_id,
	:usr_id,
	:access,
	:assigned_at) 
ON CONFLICT(group_id, usr_id) DO UPDATE SET 
	access = EXCLUDED.access, 
	assigned_at = EXCLUDED.assigned_at`

type groupUsrRow struct {
	GroupID    domain.GroupID        `db:"group_id"`
	UsrID      domain.UsrID          `db:"usr_id"`
	Access     domain.GroupUsrAccess `db:"access"`
	AssignedAt time.Time             `db:"assigned_at"`
}

func (r *groupUsrRow) ToDomain() domain.GroupUsr {
	return domain.GroupUsr{
		GroupID:    r.GroupID,
		UsrID:      r.UsrID,
		Access:     r.Access,
		AssignDate: r.AssignedAt,
	}
}

type groupUsrRepository struct {
	db sqlx.ExtContext
}

func (r *groupUsrRepository) GetGroupAccess(
	ctx context.Context,
	usrID domain.UsrID,
	groupID domain.GroupID,
) (domain.GroupUsrAccess, error) {
	q := fmt.Sprintf("%s\nWHERE gu.usr_id=$1 AND gu.group_id=$2", getGroupAccessQuery)
	var access domain.GroupUsrAccess
	err := r.db.QueryRowxContext(ctx, q, usrID, groupID).Scan(&access)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.GroupNoneAccess, nil
		}

		return domain.GroupNoneAccess, err
	}

	return access, nil
}

func (r *groupUsrRepository) GetNodeAccess(ctx context.Context, usrID domain.UsrID, nodeID domain.NodeID) (domain.GroupUsrAccess, error) {
	var access domain.GroupUsrAccess
	err := r.db.QueryRowxContext(ctx, getNodeAccessQuery, usrID, nodeID).Scan(&access)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.GroupNoneAccess, nil
		}

		return domain.GroupNoneAccess, err
	}

	return access, nil
}

func (r *groupUsrRepository) GetByUsr(ctx context.Context, usrID domain.UsrID) ([]domain.GroupUsr, error) {
	q := fmt.Sprintf("%s\nWHERE gu.usr_id=$1", getGroupUsrQuery)
	rows, err := r.db.QueryxContext(ctx, q, usrID)
	if err != nil {
		return nil, err
	}

	defer safeClose(ctx, rows)
	// por revisar bb si funciona soy dios
	res, err := readSlice[domain.GroupUsr, groupUsrRow](rows)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *groupUsrRepository) Upsert(ctx context.Context, d *domain.UpsertGroupUsr) error {
	if d == nil {
		return ErrInvalidArgument
	}
	assignedAt := time.Now().UTC()
	row := groupUsrRow{
		UsrID:      d.UsrID,
		GroupID:    d.GroupID,
		Access:     d.Access,
		AssignedAt: assignedAt,
	}

	_, err := sqlx.NamedExecContext(ctx, r.db, upsertGroupUsrQuery, row)
	if err != nil {
		return err
	}

	return nil
}

func (r *groupUsrRepository) UpsertAll(ctx context.Context, d []domain.UpsertGroupUsr) error {
	if len(d) == 0 {
		return ErrInvalidArgument
	}

	rows := make([]groupUsrRow, len(d))
	assignedAt := time.Now().UTC()
	for i, gu := range d {
		row := &rows[i]
		row.UsrID = gu.UsrID
		row.GroupID = gu.GroupID
		row.Access = gu.Access
		row.AssignedAt = assignedAt
	}

	_, err := sqlx.NamedExecContext(ctx, r.db, upsertGroupUsrQuery, rows)
	if err != nil {
		return err
	}

	return nil
}

func (r *groupUsrRepository) RemoveUsr(ctx context.Context, g domain.GroupID, u domain.UsrID) error {
	return apperror.ErrNotImplemented(nil)
}

func NewGroupUsrRepository(db sqlx.ExtContext) domain.GroupUsrRepository {
	if db == nil {
		panic("NewGroupUsrRepository received a nil db")
	}

	return &groupUsrRepository{db}
}
