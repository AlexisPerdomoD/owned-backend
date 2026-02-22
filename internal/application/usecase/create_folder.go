package usecase

import (
	"context"
	"fmt"

	"ownned/internal/application/model"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
)

type CreateFolderUseCase struct {
	nodeRepository     domain.NodeRepository
	usrRepository      domain.UsrRepository
	groupUsrRepository domain.GroupUsrRepository
}

func (uc *CreateFolderUseCase) Execute(ctx context.Context, creatorID domain.UsrID, dto *model.CreateFolderInputDTO) (*domain.Node, error) {
	usr, err := uc.usrRepository.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, apperror.ErrUnauthenticated(nil)
	}

	if usr.Role == domain.LimitedUsrRole {
		return nil, apperror.ErrForbidden(nil)
	}

	parentID, folder := dto.GetData()

	parent, err := uc.nodeRepository.GetByID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	if parent == nil {
		return nil, apperror.ErrNotFound(
			map[string]string{
				"error": "parent node was not found",
			})
	}

	if parent.Type != domain.FolderNodeType {
		return nil, apperror.ErrBadRequest(
			map[string]string{
				"error": "parent node is not a folder type",
			})
	}

	folder.Path = parent.Path.NewChildPath(folder.ID)

	if usr.Role != domain.SuperUsrRole {
		access, err := uc.groupUsrRepository.GetNodeAccess(ctx, usr.ID, parent.ID)
		if err != nil {
			return nil, err
		}

		if access == nil || *access != domain.GroupWriteAccess {
			return nil, apperror.ErrForbidden(
				map[string]string{
					"error": fmt.Sprintf("usr cannot create nodes on this folder=%s with ID=%s", parent.Name, parent.ID),
				})
		}

	}

	if err := uc.nodeRepository.Create(ctx, folder); err != nil {
		return nil, err
	}

	return folder, err
}

func NewCreateFolderUseCase(
	nr domain.NodeRepository,
	ur domain.UsrRepository,
	gur domain.GroupUsrRepository,
) *CreateFolderUseCase {
	if nr == nil || ur == nil || gur == nil {
		panic("NewCreateFolderUseCase receive nil dependencies")
	}

	return &CreateFolderUseCase{nr, ur, gur}
}
