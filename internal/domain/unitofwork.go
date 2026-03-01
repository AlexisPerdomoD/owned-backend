package domain

import "context"

// UnitOfWork represents an atomic transactional boundary.
//
// A UnitOfWork encapsulates a set of repository operations that must be executed
// atomically. All repository methods accessed through a UnitOfWork are expected
// to operate within the same underlying transaction while Do is executing.
type UnitOfWork interface {
	// ctx returns the context associated with the current UnitOfWork.
	//
	// The returned context is used for all I/O operations and must respect
	// cancellation and deadlines.
	Ctx() context.Context

	// NodeRepository returns a NodeRepository bound to the current UnitOfWork.
	//
	// All operations performed through the returned repository while Do is
	// executing must participate in the active transaction.
	NodeRepository() NodeRepository
	// DocRepository returns a DocRepository bound to the current UnitOfWork.
	//
	// All operations performed through the returned repository while Do is
	// executing must participate in the active transaction.
	DocRepository() DocRepository

	// UsrRepository returns a UsrRepository bound to the current UnitOfWork.
	//
	// All operations performed through the returned repository while Do is
	// executing must participate in the active transaction.
	UsrRepository() UsrRepository

	// UsrPwdRepository returns a UsrPwdRepository bound to the current UnitOfWork.
	//
	// All operations performed through the returned repository while Do is
	// executing must participate in the active transaction.
	UsrPwdRepository() UsrPwdRepository

	// GroupRepository returns a GroupRepository bound to the current UnitOfWork.
	//
	// All operations performed through the returned repository while Do is
	// executing must participate in the active transaction.
	GroupRepository() GroupRepository
	// GroupUsrRepository returns a GroupUsrRepository bound to the current UnitOfWork.
	//
	// All operations performed through the returned repository while Do is
	// executing must participate in the active transaction.
	GroupUsrRepository() GroupUsrRepository
	// GroupNodeRepository returns a GroupNodeRepository bound to the current UnitOfWork.
	//
	// All operations performed through the returned repository while Do is
	// executing must participate in the active transaction.
	GroupNodeRepository() GroupNodeRepository
}

type UnitOfWorkFactory interface {
	// Do executes the given function within a transactional context.
	//
	// The provided function is executed inside a single transaction. If the
	// function returns a non-nil error, the transaction is rolled back and the
	// error is returned to the caller. If the function returns nil, the
	// transaction is committed.
	//
	// Implementations must guarantee that all repository instances obtained
	// from this UnitOfWork during the execution of the function operate on the
	// same transactional state.
	//

	Do(ctx context.Context, fn func(uow UnitOfWork) error) error
}
