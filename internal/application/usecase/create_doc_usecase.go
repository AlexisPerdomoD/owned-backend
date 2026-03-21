package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"ownned/internal/application/dto"
	"ownned/internal/application/storage"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
)

type CreateDocUseCaseResponse struct {
	Doc  *domain.Doc
	Node *domain.Node
}

type CreateDocUseCase struct {
	accessChecker
	usrRepository     domain.UsrRepository
	docRepository     domain.DocRepository
	nodeRepository    domain.NodeRepository
	unitOfWorkFactory domain.UnitOfWorkFactory
	storage           storage.StorageManager
	logger            *slog.Logger
}

func (uc *CreateDocUseCase) Execute(ctx context.Context, creatorID domain.UsrID, arg *dto.CreateDocInputDTO) (*CreateDocUseCaseResponse, error) {
	usr, err := uc.usrRepository.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, apperror.ErrUnauthenticated(nil)
	}

	if usr.Role == domain.LimitedUsrRole {
		detail := make(map[string]string)
		detail["reason"] = "user does not have permission to create documents"
		return nil, apperror.ErrForbidden(detail)
	}

	folder, err := uc.nodeRepository.GetByID(ctx, arg.ParentID)
	if err != nil {
		return nil, err
	}

	if folder == nil {
		detail := make(map[string]string)
		detail["reason"] = "Folder was not found."
		return nil, apperror.ErrNotFound(detail)
	}

	if folder.Type != domain.FolderNodeType {
		detail := make(map[string]string)
		detail["reason"] = "ParentID does not point to a folder."
		return nil, apperror.ErrBadRequest(detail)
	}

	canDo, err := uc.hasAccessTo(ctx, usr, folder.Path, domain.GroupWriteAccess)
	if err != nil {
		uc.logger.WarnContext(ctx, "failed to check if user can access node", "nodeID", folder.ID, "error", err)
		return nil, err
	}
	if !canDo {
		detail := make(map[string]string)
		detail["reason"] = fmt.Sprintf("User does not have access to specified node ID=%s to create documents", folder.ID.String())
		return nil, apperror.ErrForbidden(detail)
	}

	nodeID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	docID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	uploadCommand := storage.StorageUploadCommand{
		Key:          docID.String(),
		MaxSizeBytes: arg.ExpectedSize,
		File:         arg.File,
	}
	size, err := uc.storage.Put(ctx, uploadCommand)
	if err != nil {
		return nil, err
	}

	response := &CreateDocUseCaseResponse{
		Node: &domain.Node{
			ID:          nodeID,
			Type:        domain.FileNodeType,
			Description: arg.Description,
			Name:        arg.Title,
			Path:        folder.Path.NewChildPath(nodeID),
		},
		Doc: &domain.Doc{
			ID:          docID,
			NodeID:      nodeID,
			MimeType:    arg.Mimetype,
			Title:       arg.Title,
			Description: arg.Description,
			UsrID:       usr.ID,
			SizeInBytes: size,
		},
	}

	if err := uc.saveDoc(ctx, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (uc *CreateDocUseCase) saveDoc(ctx context.Context, response *CreateDocUseCaseResponse) error {
	err := uc.unitOfWorkFactory.Do(ctx, func(tx domain.UnitOfWork) error {
		txCtx := tx.Ctx()
		if err := tx.NodeRepository().Create(txCtx, response.Node); err != nil {
			return err
		}

		if err := tx.DocRepository().Create(txCtx, response.Doc); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		deleteErr := uc.storage.Delete(ctx, response.Doc.ID.String())
		if deleteErr != nil {
			uc.logger.Warn("error deleting file after document creation err",
				"err",
				deleteErr,
			)
		}
	}

	return err
}

func NewCreateDocUseCase(
	ur domain.UsrRepository,
	dr domain.DocRepository,
	nr domain.NodeRepository,
	gur domain.GroupUsrRepository,
	uowf domain.UnitOfWorkFactory,
	storage storage.StorageManager,
	mainLogger *slog.Logger,
) *CreateDocUseCase {
	helper.NotNilOrPanic(ur, "UsrRepository")
	helper.NotNilOrPanic(dr, "DocRepository")
	helper.NotNilOrPanic(nr, "NodeRepository")
	helper.NotNilOrPanic(uowf, "UnitOfWorkFactory")
	helper.NotNilOrPanic(storage, "StorageManager")
	helper.NotNilOrPanic(mainLogger, "mainLogger")
	logger := mainLogger.With("usecase", "CreateDocUseCase")
	ac := accessChecker{gur}
	return &CreateDocUseCase{ac, ur, dr, nr, uowf, storage, logger}
}
