package pg

import (
	"context"

	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type GroupUsrRepository struct {
	db sqlx.ExtContext
}

func (r *GroupUsrRepository) GetGroupAccess(ctx context.Context, usrID domain.UsrID, groupID domain.GroupID) (*domain.GroupAccess, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *GroupUsrRepository) GetNodeAccess(ctx context.Context, usrID domain.UsrID, nodeID domain.NodeID) (*domain.GroupAccess, error) {
	return nil, apperror.ErrNotImplemented(nil)
}

func (r *GroupUsrRepository) Upsert(ctx context.Context, d *domain.UpsertGroupUsr) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *GroupUsrRepository) UpsertAll(ctx context.Context, d []domain.UpsertGroupUsr) error {
	return apperror.ErrNotImplemented(nil)
}

func (r *GroupUsrRepository) RemoveUsr(ctx context.Context, g domain.GroupID, u domain.UsrID) error {
	return apperror.ErrNotImplemented(nil)
}

func NewGroupUsrRepository(db sqlx.ExtContext) *GroupUsrRepository {
	if db == nil {
		panic("NewGroupUsrRepository received a nil db")
	}

	return &GroupUsrRepository{db}
}
