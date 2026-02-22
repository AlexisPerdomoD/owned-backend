package usecase

import (
	"context"
	"log/slog"

	"ownned/internal/application/model"
	"ownned/internal/application/storage"
	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/google/uuid"
)

type CreateDocUseCaseResponse struct {
	Doc  *domain.Doc
	Node *domain.Node
}

type CreateDocUseCase struct {
	usrRepository      domain.UsrRepository
	docRepository      domain.DocRepository
	nodeRepository     domain.NodeRepository
	groupUsrRepository domain.GroupUsrRepository
	unitOfWorkFactory  domain.UnitOfWorkFactory
	storage            storage.Storage
	logger             *slog.Logger
}

func (uc *CreateDocUseCase) Execute(ctx context.Context, creatorID domain.UsrID, arg *model.CreateDocInputDTO) (*CreateDocUseCaseResponse, error) {
	usr, err := uc.usrRepository.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if usr == nil || usr.Role == domain.LimitedUsrRole {
		return nil, apperror.ErrForbidden(nil)
	}

	folder, err := uc.nodeRepository.GetByID(ctx, arg.ParentID)
	if err != nil {
		return nil, err
	}

	if folder == nil {
		return nil, apperror.ErrNotFound(map[string]string{"error": "Folder was not found"})
	}

	if folder.Type != domain.FolderNodeType {
		return nil, apperror.ErrBadRequest(map[string]string{"error": "parentID does not point to a folder"})
	}

	access, err := uc.groupUsrRepository.GetNodeAccess(ctx, usr.ID, folder.ID)
	if err != nil {
		return nil, err
	}

	if access == nil || *access != domain.GroupWriteAccess {
		return nil, apperror.ErrForbidden(map[string]string{"error": "Usr does not have enought access"})
	}

	nodeID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	docID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	uploadArgs := arg.GetUploadArgs()
	if err = uc.storage.Upload(ctx, docID, uploadArgs); err != nil {
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
			MimeType:    uploadArgs.Mimetype,
			Title:       arg.Title,
			Description: arg.Description,
			UsrID:       usr.ID,
			SizeInBytes: uint64(uploadArgs.Size),
		},
	}

	if err := uc.saveDoc(ctx, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (uc *CreateDocUseCase) saveDoc(ctx context.Context, response *CreateDocUseCaseResponse) error {
	err := uc.unitOfWorkFactory.Do(ctx, func(txCtx context.Context, tx domain.UnitOfWork) error {
		if err := tx.NodeRepository().Create(txCtx, response.Node); err != nil {
			return err
		}

		if err := tx.DocRepository().Create(txCtx, response.Doc); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		deleteErr := uc.storage.Remove(ctx, response.Doc.ID)
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
	storage storage.Storage,
	mainLogger *slog.Logger,
) *CreateDocUseCase {
	if ur == nil || dr == nil || nr == nil || uowf == nil || storage == nil {
		panic("missing dependencies for NewCreateDocUseCase")
	}

	logger := mainLogger.With("usecase", "CreateDocUseCase")
	return &CreateDocUseCase{ur, dr, nr, gur, uowf, storage, logger}
}
