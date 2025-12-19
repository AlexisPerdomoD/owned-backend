package usecase

import (
	"context"
	"ownned/internal/application/storage"
	"ownned/internal/domain"
	"ownned/internal/pkg/error_pkg"
)

type DeleteDoceUseCase struct {
	storage           storage.Storage
	docRepository     domain.DocRepository
	nodeRepository    domain.NodeRepository
	usrRepository     domain.UsrRepository
	unitOfWorkFactory domain.UnitOfWorkFactory
}

func (uc *DeleteDoceUseCase) validateUsrAccess(ctx context.Context, usrID domain.UsrID, nodeID domain.NodeID) error {
	access, err := uc.nodeRepository.GetAccess(ctx, usrID, nodeID)
	if err != nil {
		return err
	}

	if access != domain.WriteAccess {
		return error_pkg.ErrForbidden(nil)
	}

	return nil
}

func (uc *DeleteDoceUseCase) Execute(ctx context.Context, userID domain.UsrID, docID domain.DocID) (*domain.Doc, error) {
	channel := make(chan any, 2)
	defer close(channel)

	go func() {
		usr, err := uc.usrRepository.GetByID(ctx, userID)
		if err != nil {
			channel <- err
			return
		}

		if usr == nil {
			channel <- error_pkg.ErrNotFound(map[string]string{"usrID": "usr entity was not found"})
			return
		}

		channel <- usr
	}()

	go func() {
		doc, err := uc.docRepository.GetByID(ctx, docID)

		if err != nil {
			channel <- err
			return
		}

		if doc == nil {
			channel <- error_pkg.ErrNotFound(map[string]string{"docID": "doc entity was not found"})
			return
		}

		channel <- doc
	}()

	var usr *domain.Usr
	var doc *domain.Doc

	for val := range channel {
		err, isErr := val.(error)
		if isErr {
			return nil, err
		}

		usrOut, isUsr := val.(*domain.Usr)
		if isUsr {
			usr = usrOut
			continue
		}

		docOut, isDoc := val.(*domain.Doc)
		if isDoc {
			doc = docOut
			continue
		}

		return nil, error_pkg.ErrInternal(nil)
	}

	if usr.Role != domain.SuperUsrRole {
		if err := uc.validateUsrAccess(ctx, usr.ID, doc.NodeID); err != nil {
			return nil, err
		}
	}

	if err := uc.storage.Remove(ctx, doc.ID); err != nil {
		return nil, err
	}

	tx := uc.unitOfWorkFactory.New()
	_, err := tx.Do(
		ctx,
		func(txCtx context.Context, tx domain.UnitOfWork) (any, error) {
			if err := tx.DocRepository().Delete(txCtx, doc.ID); err != nil {
				return nil, err
			}

			versions, err := tx.DocRepository().GetByNodeID(txCtx, doc.NodeID)
			if err != nil {
				return nil, err
			}

			if len(versions) == 0 {
				err := tx.NodeRepository().Delete(txCtx, doc.NodeID)
				if err != nil {
					return nil, err
				}
			}

			return nil, nil
		})

	return doc, err
}

func NewDeleteDocUseCase(
	s storage.Storage,
	dr domain.DocRepository,
	nr domain.NodeRepository,
	ur domain.UsrRepository,
	uoff domain.UnitOfWorkFactory,
) *DeleteDoceUseCase {
	if s == nil || uoff == nil || dr == nil || nr == nil || ur == nil {
		panic("NewDeleteDocUseCase has been provided with nil dependencies")
	}

	return &DeleteDoceUseCase{s, dr, nr, ur, uoff}
}
