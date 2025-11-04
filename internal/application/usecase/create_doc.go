package usecase

import (
	"context"
	"log/slog"
	"ownned/internal/application/dto"
	"ownned/internal/application/storage"
	"ownned/internal/domain"
	"ownned/internal/pkg/error_pkg"
)

type CreateDocUseCaseResponse struct {
	Doc  *domain.Doc
	Node *domain.Node
}

type CreateDocUseCase struct {
	usrRepository     domain.UsrRepository
	docRepository     domain.DocRepository
	nodeRepository    domain.NodeRepository
	unitOfWorkFactory domain.UnitOfWorkFactory
	storage           storage.Storage
	logger            *slog.Logger
}

func (uc *CreateDocUseCase) Execute(ctx context.Context, creatorID domain.UsrID, arg *dto.CreateDocInputDto) (*CreateDocUseCaseResponse, error) {
	usr, err := uc.usrRepository.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if usr == nil || usr.Role == domain.LimitedUsrRole {
		return nil, error_pkg.ErrForbidden(nil)
	}

	folder, err := uc.nodeRepository.GetByID(ctx, arg.ParentID)
	if err != nil {
		return nil, err
	}

	if folder == nil {
		return nil, error_pkg.ErrNotFound(map[string]string{"parentID": "folder was not found"})
	}

	if folder.Type != domain.FolderNodeType {
		return nil, error_pkg.ErrBadRequest(map[string]string{"parentID": "does not point to a folder"})
	}

	access, err := uc.nodeRepository.GetAccess(ctx, usr.ID, folder.ID)
	if err != nil {
		return nil, err
	}

	if access != domain.WriteAccess {
		return nil, error_pkg.ErrForbidden(map[string]string{"parentID": "usr does not have enought access"})
	}

	uploadArgs := arg.GetUploadArgs()

	if err = uc.storage.Upload(ctx, uploadArgs); err != nil {
		return nil, err
	}

	tx := uc.unitOfWorkFactory.New()
	txOut, err := tx.Do(
		ctx,
		func(txCtx context.Context, tx domain.UnitOfWork) (any, error) {
			txNodeRepository := tx.NodeRepository()
			node := &domain.Node{
				ParentID:    &arg.ParentID,
				Type:        domain.FileNodeType,
				Description: arg.Description,
				Name:        arg.Title,
			}

			if err := txNodeRepository.Create(txCtx, node); err != nil {
				return nil, err
			}

			doc := &domain.Doc{
				ID:          uploadArgs.ID,
				MimeType:    uploadArgs.Mimetype,
				Title:       arg.Title,
				Description: arg.Description,
				NodeID:      node.ID,
				UsrID:       usr.ID,
				SizeInBytes: uint64(uploadArgs.Size),
			}

			if err := tx.DocRepository().Create(txCtx, doc); err != nil {
				return nil, err
			}

			return &CreateDocUseCaseResponse{Doc: doc, Node: node}, nil
		},
	)

	if err != nil {
		deleteErr := uc.storage.Remove(ctx, uploadArgs.ID)
		if deleteErr != nil {
			uc.logger.Warn("error deleting file after document creation err",
				"err",
				deleteErr,
			)
		}

		return nil, err
	}

	res, ok := txOut.(*CreateDocUseCaseResponse)
	if !ok {
		uc.logger.Warn("unexpected result type from tx.Do", "got", txOut)
		return nil, error_pkg.ErrInternal(nil)
	}

	return res, nil
}

func NewCreateDocUseCase(
	ur domain.UsrRepository,
	dr domain.DocRepository,
	nr domain.NodeRepository,
	uowf domain.UnitOfWorkFactory,
	storage storage.Storage,
	mainLogger *slog.Logger,
) *CreateDocUseCase {
	if ur == nil || dr == nil || nr == nil || uowf == nil || storage == nil {
		panic("missing dependencies for NewCreateDocUseCase")
	}

	logger := mainLogger.With("usecase", "CreateDocUseCase")
	return &CreateDocUseCase{ur, dr, nr, uowf, storage, logger}
}
