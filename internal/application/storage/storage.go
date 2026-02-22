// Package storage provides the storage layer of the application.
package storage

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type UploadArgs struct {
	Mimetype string
	Size     uint64
	File     io.Reader
}

type Storage interface {
	Get(ctx context.Context, key uuid.UUID) (io.ReadCloser, error)

	Upload(ctx context.Context, key uuid.UUID, args *UploadArgs) error

	Remove(ctx context.Context, key uuid.UUID) error
}
