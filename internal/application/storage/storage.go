// Package storage provides the storage layer of the application.
package storage

import (
	"context"
	"io"
)

type UploadArgs struct {
	ID       string
	Mimetype string
	Size     uint64
	File     io.Reader
}

type Storage interface {
	Get(ctx context.Context, id string) (io.ReadCloser, error)

	Upload(ctx context.Context, args *UploadArgs) error

	Remove(ctx context.Context, id string) error
}
