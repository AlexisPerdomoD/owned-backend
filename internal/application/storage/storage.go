package storage

import (
	"io"
	"ownned/internal/domain"
)

type Storage interface {
	Get(identifier domain.DocID) (io.ReadCloser, error)

	Put(identifier domain.DocID, f io.ReadCloser) error

	Remove(identifier domain.DocID) error
}
