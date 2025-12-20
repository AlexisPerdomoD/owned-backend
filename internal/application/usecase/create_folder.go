package usecase

import (
	"context"
	"fmt"
	"ownned/internal/application/dto"
	"ownned/internal/domain"
)

type CreateFolderUseCase struct {
	nodeRepository    domain.NodeRepository
	usrRepository     domain.UsrRepository
	unitOfWorkFactory domain.UnitOfWorkFactory
}

func (uc *CreateFolderUseCase) Execute(ctx context.Context, creatorID domain.UsrID, dto *dto.CreateFolderInputDTO) (*domain.Node, error) {
	usr, err := uc.usrRepository.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, domain.ErrUnauthenticated(nil)
	}

	if usr.Role == domain.LimitedUsrRole {
		return nil, domain.ErrForbidden(nil)
	}

	folder := dto.GetData()

	if folder.ParentID != nil {
		parent, err := uc.nodeRepository.GetByID(ctx, *folder.ParentID)
		if err != nil {
			return nil, err
		}

		if parent == nil {
			return nil, domain.ErrNotFound(
				map[string]string{
					"ParentID": "parent node was not found",
				})
		}

		if parent.Type != domain.FolderNodeType {
			return nil, domain.ErrBadRequest(
				map[string]string{
					"parentID": "parent node is not a folder type",
				})
		}

		if usr.Role != domain.SuperUsrRole {
			access, err := uc.nodeRepository.GetAccess(ctx, usr.ID, parent.ID)
			if err != nil {
				return nil, err
			}

			if access != domain.WriteAccess {
				return nil, domain.ErrForbidden(
					map[string]string{
						"parentID": fmt.Sprintf("usr cannot create nodes on this folder=%s with ID=%s", parent.Name, parent.ID),
					})
			}

		}
	}

	tx := uc.unitOfWorkFactory.New()
	err = tx.Do(ctx, func(txCtx context.Context) error {
		nr := tx.NodeRepository()

		if err := nr.Create(txCtx, folder); err != nil {
			return err
		}

		if usr.Role == domain.SuperUsrRole {
			return nil
		}

		return nr.UpdateAccess(txCtx, usr.ID, folder.ID, domain.WriteAccess)
	})

	return folder, err
}

func NewCreateFolderUseCase(
	nr domain.NodeRepository,
	ur domain.UsrRepository,
	uowf domain.UnitOfWorkFactory,
) *CreateFolderUseCase {

	if nr == nil || ur == nil || uowf == nil {
		panic("NewCreateFolderUseCase receive nil dependencies")
	}

	return &CreateFolderUseCase{nr, ur, uowf}
}
