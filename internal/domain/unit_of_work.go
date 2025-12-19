package domain

import "context"

// UnitOfWork represents an atomic transactional boundary.
//
// A UnitOfWork encapsulates a set of repository operations that must be executed
// atomically. All repository methods accessed through a UnitOfWork are expected
// to operate within the same underlying transaction while Do is executing.
type UnitOfWork interface {

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
	// The context passed to the function should be used for all I/O operations
	// and must respect cancellation and deadlines.
	Do(ctx context.Context, fn func(ctx context.Context) error) error

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
}

type UnitOfWorkFactory interface {
	New() UnitOfWork
}
