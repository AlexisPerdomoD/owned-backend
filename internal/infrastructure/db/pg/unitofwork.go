package pg

import (
	"context"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
)

type UnitOfWork struct {
	nodeRepository domain.NodeRepository
	usrRepository  domain.UsrRepository
	docRepository  domain.DocRepository
}

func (u *UnitOfWork) NodeRepository() domain.NodeRepository {
	if u.nodeRepository == nil {
		u.nodeRepository = NewNodeRepository()
	}

	return u.nodeRepository
}

func (u *UnitOfWork) DocRepository() domain.DocRepository {
	if u.docRepository == nil {
		u.docRepository = NewDocRepository()
	}

	return u.docRepository
}

func (u *UnitOfWork) UsrRepository() domain.UsrRepository {
	if u.usrRepository == nil {
		u.usrRepository = NewUsrRepository()
	}

	return u.usrRepository
}

type UnitOfWorkFactory struct{}

func (f *UnitOfWorkFactory) Do(ctx context.Context, fn func(ctx context.Context, tx domain.UnitOfWork) error) error {
	return apperror.ErrNotImplemented(nil)
}

func NewUnitOfWorkFactory() *UnitOfWorkFactory {
	return &UnitOfWorkFactory{}
}
