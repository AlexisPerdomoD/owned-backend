package usecase

import (
	"context"

	"ownned/internal/application/storage"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
)

type DeleteDoceUseCase struct {
	storage            storage.Storage
	docRepository      domain.DocRepository
	nodeRepository     domain.NodeRepository
	usrRepository      domain.UsrRepository
	groupUsrRepository domain.GroupUsrRepository
}

func (uc *DeleteDoceUseCase) validateUsrAccess(ctx context.Context, usrID domain.UsrID, nodeID domain.NodeID) error {
	access, err := uc.groupUsrRepository.GetNodeAccess(ctx, usrID, nodeID)
	if err != nil {
		return err
	}

	if access != domain.GroupWriteAccess {
		return apperror.ErrForbidden(nil)
	}

	return nil
}

func (uc *DeleteDoceUseCase) Execute(ctx context.Context, userID domain.UsrID, docID domain.DocID) (*domain.Doc, error) {
	usr, err := uc.usrRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, apperror.ErrNotFound(map[string]string{"usrID": "usr entity was not found"})
	}

	doc, err := uc.docRepository.GetByID(ctx, docID)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, apperror.ErrNotFound(map[string]string{"docID": "doc entity was not found"})
	}

	if usr.Role != domain.SuperUsrRole {
		if err := uc.validateUsrAccess(ctx, usr.ID, doc.NodeID); err != nil {
			return nil, err
		}
	}

	if err := uc.storage.Remove(ctx, doc.ID); err != nil {
		return nil, err
	}

	if err := uc.nodeRepository.Delete(ctx, doc.NodeID); err != nil {
		return nil, err
	}

	return doc, nil
}

func NewDeleteDocUseCase(
	s storage.Storage,
	dr domain.DocRepository,
	nr domain.NodeRepository,
	ur domain.UsrRepository,
	gur domain.GroupUsrRepository,
) *DeleteDoceUseCase {
	if s == nil || gur == nil || dr == nil || nr == nil || ur == nil {
		panic("NewDeleteDocUseCase has been provided with nil dependencies")
	}

	return &DeleteDoceUseCase{s, dr, nr, ur, gur}
}
