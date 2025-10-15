package storage

import (
	"io"
	"ownned/internal/domain"
)

type Storage interface {
	GetDoc(identifier domain.DocID) (io.Reader, error)

	UploadDoc(identifier domain.DocID, f io.Reader) error

	Delete(identifier domain.DocID) error
}
