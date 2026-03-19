package serv

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"ownned/internal/application/storage"
)

type storageManagerFS struct {
	DirPath string
}

func (stg *storageManagerFS) Put(ctx context.Context, cmd storage.StorageUploadCommand) (uint64, error) {
	filename := filepath.Join(stg.DirPath, cmd.Key)
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return 0, err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}

		if err != nil {
			_ = os.Remove(filename)
		}
	}()

	r := io.LimitReader(cmd.File, int64(cmd.MaxSizeBytes)+1)
	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	if n > int64(cmd.MaxSizeBytes) {
		err = storage.ErrFileTooLarge
	}

	return uint64(n), err
}

func (stg *storageManagerFS) Download(ctx context.Context, key storage.FileKey) (io.ReadCloser, error) {
	return nil, nil
}

func (stg *storageManagerFS) Delete(ctx context.Context, key storage.FileKey) error {
	return nil
}

func NewStorageManagerFS(dirpath string) storage.StorageManager {
	if err := os.MkdirAll(dirpath, 0o755); err != nil {
		panic(err)
	}

	return &storageManagerFS{dirpath}
}
