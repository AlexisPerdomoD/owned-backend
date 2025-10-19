package usecase

import (
	"context"
	"fmt"
	"log"
	"ownned/internal/application/dto"
	"ownned/internal/domain"
	"ownned/internal/pkg/error_pkg"
)

type CreateFolderUseCase struct {
	nodeRepository    domain.NodeRepository
	usrRepository     domain.UsrRepository
	unitOfWorkFactory domain.UnitOfWorkFactory
}

func (uc *CreateFolderUseCase) Execute(ctx context.Context, creatorID domain.UsrID, dto *dto.CreateFolderInputDto) (*domain.FolderNode, error) {
	usr, err := uc.usrRepository.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, error_pkg.ErrUnauthenticated(nil)
	}

	if usr.Role == domain.LimitedUsrRole {
		return nil, error_pkg.ErrForbidden(nil)
	}

	node := dto.GetData()

	if node.ParentID != nil {
		parent, err := uc.nodeRepository.GetByID(ctx, *node.ParentID)
		if err != nil {
			return nil, err
		}

		if parent == nil {
			return nil, error_pkg.ErrNotFound(map[string]string{"ParentID": "parent node was not found"})
		}

		if parent.Type != domain.FolderNodeType {
			return nil, error_pkg.ErrBadRequest(map[string]string{"parentID": "parent node is not a folder type"})
		}

		if usr.Role != domain.SuperUsrRole {
			access, err := uc.nodeRepository.GetAccess(ctx, usr.ID, node.ID)
			if err != nil {
				return nil, err
			}

			if access != domain.WriteAccess {
				return nil, error_pkg.ErrForbidden(
					map[string]string{
						"parentID": fmt.Sprintf("usr cannot create nodes on this folder=%s with ID=%s", parent.Name, parent.ID),
					})
			}

		}
	}

	tx, err := uc.
		unitOfWorkFactory.New().
		Do(ctx, func(
			ctx context.Context,
			uow domain.UnitOfWork,
		) (any, error) {
			return nil, nil
		})

	return nil, nil
}

func NewCreateFolderUseCase(
	nr domain.NodeRepository,
	ur domain.UsrRepository,
	uowf domain.UnitOfWorkFactory,
) *CreateFolderUseCase {

	if nr == nil || ur == nil || uowf == nil {
		log.Panicln("NewCreateFolderUseCase receive nil dependencies")
	}

	return &CreateFolderUseCase{nr, ur, uowf}
}
