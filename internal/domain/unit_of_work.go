package domain

import "context"

type UnitOfWork interface {
	Do(ctx context.Context, tx func(ctx context.Context, uow UnitOfWork) (any, error)) (any, error)

	NodeRepository() NodeRepository

	DocRepository() DocRepository

	UsrRepository() UsrRepository
}

type UnitOfWorkFactory interface {
	New() UnitOfWork
}
