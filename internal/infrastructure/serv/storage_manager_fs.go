package serv

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"ownned/internal/application/storage"
)

type storageManagerFS struct {
	Dir string
}

func (stg *storageManagerFS) Put(ctx context.Context, c storage.StorageUploadCommand) (uint64, error) {
	fname := filepath.Join(stg.Dir, c.Key)
	f, err := os.OpenFile(fname, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return 0, err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}

		if err != nil {
			_ = os.Remove(fname)
		}
	}()

	r := io.LimitReader(c.File, int64(c.MaxSizeBytes)+1)
	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	if n > int64(c.MaxSizeBytes) {
		err = storage.ErrFileTooLarge
	}

	return uint64(n), err
}

func (stg *storageManagerFS) Download(ctx context.Context, key storage.FileKey) (io.ReadCloser, error) {
	fname := filepath.Join(stg.Dir, key)
	f, err := os.Open(fname)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.ErrFileNotFound
		}

		return nil, err
	}

	return f, nil
}

func (stg *storageManagerFS) Delete(ctx context.Context, k storage.FileKey) error {
	fname := filepath.Join(stg.Dir, k)
	err := os.Remove(fname)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func NewStorageManagerFS(dirpath string) storage.StorageManager {
	if err := os.MkdirAll(dirpath, 0o755); err != nil {
		panic(err)
	}

	return &storageManagerFS{dirpath}
}
