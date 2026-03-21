package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"ownned/internal/application/dto"
	"ownned/internal/application/storage"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
)

type DeleteDocUseCase struct {
	storage            storage.StorageManager
	docRepository      domain.DocRepository
	nodeRepository     domain.NodeRepository
	usrRepository      domain.UsrRepository
	groupUsrRepository domain.GroupUsrRepository
	log                *slog.Logger
}

func (uc *DeleteDocUseCase) Execute(ctx context.Context, usrID domain.UsrID, docID domain.DocID) (*dto.FileNodeDTO, error) {
	usr, err := uc.usrRepository.GetByID(ctx, usrID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, apperror.ErrUnauthenticated(nil)
	}

	doc, err := uc.docRepository.GetByID(ctx, docID)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		detail := make(map[string]string)
		detail["reason"] = fmt.Sprintf("Doc entity with ID=%s was not found", docID.String())
		return nil, apperror.ErrNotFound(detail)
	}

	node, err := uc.nodeRepository.GetByID(ctx, doc.NodeID)
	if err != nil {
		return nil, err
	}

	if node == nil {
		detail := make(map[string]string)
		detail["reason"] = fmt.Sprintf("Internal state error, node with ID=%s was not found", doc.NodeID.String())
		err := apperror.ErrInternal(detail)
		uc.log.ErrorContext(ctx, "failed to get node by ID", "nodeID", doc.NodeID, "err", err)
		return nil, err
	}

	accss, err := resolveNodeAccess(ctx, uc.groupUsrRepository, usr, node)
	if err != nil {
		return nil, err
	}

	if accss != domain.GroupWriteAccess {
		detail := make(map[string]string)
		detail["reason"] = "user does not have access to remove this doc"
		return nil, apperror.ErrForbidden(detail)
	}

	if err := uc.storage.Delete(ctx, doc.ID.String()); err != nil {
		uc.log.WarnContext(ctx, "failed to delete doc from storage", "docID", docID, "err", err)
		return nil, err
	}

	if err := uc.nodeRepository.Delete(ctx, doc.NodeID); err != nil {
		return nil, err
	}

	return &dto.FileNodeDTO{Node: *node, Doc: *doc}, nil
}

func NewDeleteDocUseCase(
	s storage.StorageManager,
	dr domain.DocRepository,
	nr domain.NodeRepository,
	ur domain.UsrRepository,
	gur domain.GroupUsrRepository,
	mainLogger *slog.Logger,
) *DeleteDocUseCase {
	helper.NotNilOrPanic(s, "StorageManager")
	helper.NotNilOrPanic(dr, "DocRepository")
	helper.NotNilOrPanic(nr, "NodeRepository")
	helper.NotNilOrPanic(ur, "UsrRepository")
	helper.NotNilOrPanic(gur, "GroupUsrRepository")
	helper.NotNilOrPanic(mainLogger, "mainLogger")
	log := mainLogger.With("usecase", "DeleteDocUseCase")
	return &DeleteDocUseCase{s, dr, nr, ur, gur, log}
}
