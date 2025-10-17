package domain

import "context"

type UnitOfWork interface {
	Do(ctx context.Context, tx func(uow UnitOfWork) (any, error)) (any, error)

	NodeRepository() NodeRepository

	DocRepository() DocRepository

	UsrRepository() UsrRepository
}

type UnitOfWorkFactory interface {
	New() UnitOfWork
}
