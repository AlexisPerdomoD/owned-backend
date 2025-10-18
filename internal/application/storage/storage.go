package storage

import (
	"context"
	"io"
	"ownned/internal/domain"
)

type Storage interface {
	Get(ctx context.Context, docID domain.DocID) (io.ReadCloser, error)

	Put(ctx context.Context, docID domain.DocID, f io.Reader) error

	Remove(ctx context.Context, docID domain.DocID) error
}
