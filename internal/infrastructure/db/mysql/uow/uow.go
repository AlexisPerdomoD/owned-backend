package uow

import (
	"context"
	"ownned/internal/domain"
	"ownned/internal/infrastructure/db/mysql/repo"
	"ownned/pkg/apperror"
)

type UnitOfWork struct {
	nodeRepository domain.NodeRepository
	usrRepository  domain.UsrRepository
	docRepository  domain.DocRepository
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return apperror.ErrNotImplemented(nil)
}

func (u *UnitOfWork) NodeRepository() domain.NodeRepository {
	if u.nodeRepository == nil {
		u.nodeRepository = repo.NewNodeRepository()
	}

	return u.nodeRepository
}

func (u *UnitOfWork) DocRepository() domain.DocRepository {
	if u.docRepository == nil {
		u.docRepository = repo.NewDocRepository()
	}

	return u.docRepository
}

func (u *UnitOfWork) UsrRepository() domain.UsrRepository {
	if u.usrRepository == nil {
		u.usrRepository = repo.NewUsrRepository()
	}

	return u.usrRepository
}

type UnitOfWorkFactory struct{}

func (f *UnitOfWorkFactory) New() domain.UnitOfWork {
	return &UnitOfWork{}
}

func New() *UnitOfWorkFactory {
	return &UnitOfWorkFactory{}
}
