package usecase

import (
	"context"
	"fmt"

	"ownned/internal/application/dto"
	"ownned/internal/domain"
	"ownned/pkg/apperror"

	"github.com/google/uuid"
)

type CreateFolderUseCase struct {
	nodeRepository     domain.NodeRepository
	usrRepository      domain.UsrRepository
	groupUsrRepository domain.GroupUsrRepository
}

func (uc *CreateFolderUseCase) Execute(ctx context.Context, creatorID domain.UsrID, args *dto.CreateFolderDTO) (*domain.Node, error) {
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

	if err := args.Validate(); err != nil {
		return nil, err
	}

	parentID, err := uuid.Parse(args.ParentID)
	if err != nil {
		return nil, err
	}

	parent, err := uc.nodeRepository.GetByID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	if parent == nil {
		detail := make(map[string]string)
		detail["reason"] = "parent node was not found"
		return nil, apperror.ErrNotFound(detail)
	}

	if parent.Type != domain.FolderNodeType {
		detail := make(map[string]string)
		detail["reason"] = "parent node is not a folder type"
		return nil, apperror.ErrBadRequest(detail)
	}

	if usr.Role != domain.SuperUsrRole {
		access, err := uc.groupUsrRepository.GetNodeAccess(ctx, usr.ID, parent.ID)
		if err != nil {
			return nil, err
		}

		if access == nil || *access != domain.GroupWriteAccess {
			detail := make(map[string]string)
			detail["reason"] = fmt.Sprintf("Cannot create nodes on this folder=%s with ID=%s", parent.Name, parent.ID)
			return nil, apperror.ErrForbidden(detail)
		}

	}

	folderID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	folder := &domain.Node{
		ID:          folderID,
		Name:        args.Name,
		Description: args.Description,
		Type:        domain.FolderNodeType,
		Path:        parent.Path.NewChildPath(folderID),
	}

	if err := uc.nodeRepository.Create(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
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
