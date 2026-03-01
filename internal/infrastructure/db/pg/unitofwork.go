package pg

import (
	"context"
	"log/slog"
	"time"

	"ownned/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UnitOfWork struct {
	tx                  *sqlx.Tx
	ctx                 context.Context
	nodeRepository      *NodeRepository
	usrRepository       *UsrRepository
	usrPwdRepository    *UsrPwdRepository
	docRepository       *DocRepository
	groupRepository     *GroupRepository
	groupUsrRepository  *GroupUsrRepository
	groupNodeRepository *GroupNodeRepository
}

func (u *UnitOfWork) Ctx() context.Context {
	return u.ctx
}

func (u *UnitOfWork) NodeRepository() domain.NodeRepository {
	if u.nodeRepository == nil {
		u.nodeRepository = NewNodeRepository(u.tx)
	}

	return u.nodeRepository
}

func (u *UnitOfWork) DocRepository() domain.DocRepository {
	if u.docRepository == nil {
		u.docRepository = NewDocRepository(u.tx)
	}

	return u.docRepository
}

func (u *UnitOfWork) UsrRepository() domain.UsrRepository {
	if u.usrRepository == nil {
		u.usrRepository = NewUsrRepository(u.tx)
	}

	return u.usrRepository
}

func (u *UnitOfWork) UsrPwdRepository() domain.UsrPwdRepository {
	if u.usrPwdRepository == nil {
		u.usrPwdRepository = NewUsrPwdRepository(u.tx)
	}

	return u.usrPwdRepository
}

func (u *UnitOfWork) GroupRepository() domain.GroupRepository {
	if u.groupRepository == nil {
		u.groupRepository = NewGroupRepository(u.tx)
	}

	return u.groupRepository
}

func (u *UnitOfWork) GroupUsrRepository() domain.GroupUsrRepository {
	if u.groupUsrRepository == nil {
		u.groupUsrRepository = NewGroupUsrRepository(u.tx)
	}

	return u.groupUsrRepository
}

func (u *UnitOfWork) GroupNodeRepository() domain.GroupNodeRepository {
	if u.groupNodeRepository == nil {
		u.groupNodeRepository = NewGroupNodeRepository(u.tx)
	}

	return u.groupNodeRepository
}

type UnitOfWorkFactory struct {
	db      *sqlx.DB
	log     *slog.Logger
	timeout time.Duration
}

func (f *UnitOfWorkFactory) Do(ctx context.Context, op func(tx domain.UnitOfWork) error) error {
	txCtx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	tx, err := f.db.BeginTxx(txCtx, nil)
	if err != nil {
		f.log.WarnContext(ctx, "BeginTxx failed", slog.String("err", err.Error()))
		return err
	}
	defer func() {
		if err == nil {
			return
		}

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			f.log.WarnContext(ctx, "Rollback failed", slog.String("err", rollbackErr.Error()))
		}
	}()
	uow := &UnitOfWork{tx: tx, ctx: txCtx}
	err = op(uow)
	if err != nil {
		f.log.DebugContext(txCtx, "error happened while executing unit of work", "err", err)
		return err
	}

	err = tx.Commit()
	return err
}

func NewUnitOfWorkFactory(db *sqlx.DB, log *slog.Logger, timeout time.Duration) *UnitOfWorkFactory {
	return &UnitOfWorkFactory{db, log, timeout}
}
