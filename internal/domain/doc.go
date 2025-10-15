package domain

import (
	"context"
	"time"
)

type DocID = string

type Doc struct {
	ID          DocID
	NodeID      NodeID
	UsrID       UsrID
	Description string
	Title       string
	MimeType    string
	SizeInBytes int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DocRepository interface {
	GetByID(ctx context.Context, id DocID) (*Doc, error)

	Create(ctx context.Context, d *Doc) error

	Update(ctx context.Context, d *Doc) error

	Delete(ctx context.Context, id DocID) error
}
