package pg

import (
	"context"
	"log/slog"
	"ownned/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

type UnitOfWork struct {
	tx *sqlx.Tx

	nodeRepository domain.NodeRepository
	usrRepository  domain.UsrRepository
	docRepository  domain.DocRepository
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

type UnitOfWorkFactory struct{
	db *sqlx.DB
	log *slog.Logger
	timeout time.Duration
}

func (f *UnitOfWorkFactory) Do(ctx context.Context, op func(ctx context.Context, tx domain.UnitOfWork) error) error {
	tx, err := f.db.BeginTxx(ctx, nil)
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

	uow := &UnitOfWork{tx:tx}
	txCtx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	err = op(txCtx, uow)
	if err != nil {
		f.log.DebugContext(ctx, "error happened while executing unit of work", "err", err)
		return err
	}

	err = tx.Commit()
	return err
}


func NewUnitOfWorkFactory(db *sqlx.DB, log *slog.Logger, timeout time.Duration) *UnitOfWorkFactory {
	return &UnitOfWorkFactory{db, log, timeout}
}
