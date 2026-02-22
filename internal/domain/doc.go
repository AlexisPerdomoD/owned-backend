// Package domain provides the domain layer of the application.
package domain

import (
	"context"
	"time"
)

// DocID is a document version identifier in the system
type DocID = string

// Doc represents a document in the system asociated with a speific NodeID as position in the tree
type Doc struct {
	ID          DocID
	NodeID      NodeID
	UsrID       UsrID
	Description string
	Title       string
	MimeType    string
	SizeInBytes uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DocRepository interface {
	GetByID(ctx context.Context, id DocID) (*Doc, error)

	GetByNodeID(ctx context.Context, id NodeID) (*Doc, error)

	Create(ctx context.Context, d *Doc) error

	Update(ctx context.Context, d *Doc) error

	Delete(ctx context.Context, id DocID) error
}
