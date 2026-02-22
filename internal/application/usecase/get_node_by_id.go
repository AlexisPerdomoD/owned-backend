package usecase

import (
	"context"
	"log"

	"ownned/internal/application/model"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
)

type GetNodeByIDUseCase struct {
	ur  domain.UsrRepository
	nr  domain.NodeRepository
	dr  domain.DocRepository
	gur domain.GroupUsrRepository
}

func (uc *GetNodeByIDUseCase) Execute(ctx context.Context, usrID domain.UsrID, nodeID domain.NodeID) (domain.NodeLike, error) {
	nodeRepository := uc.nr
	usrRepository := uc.ur
	docRepository := uc.dr

	usr, err := usrRepository.GetByID(ctx, usrID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, apperror.ErrForbidden(nil)
	}

	node, err := nodeRepository.GetByID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	if node == nil {
		return nil, apperror.ErrNotFound(nil)
	}

	if usr.Role != domain.SuperUsrRole {
		access, err := uc.gur.GetNodeAccess(ctx, usr.ID, node.ID)
		if err != nil {
			return nil, err
		}

		if access == nil {
			return nil, apperror.ErrForbidden(nil)
		}
	}

	if node.Type == domain.FileNodeType {
		doc, err := docRepository.GetByNodeID(ctx, node.ID)
		if err != nil {
			return nil, err
		}

		if doc == nil {
			return nil, apperror.ErrNotFound(map[string]string{"error": "doc entity was not found"})
		}

		return &model.FileNodeDTO{Node: *node, Doc: *doc}, nil
	}

	children, err := nodeRepository.GetChildren(ctx, node.Path)
	if err != nil {
		return nil, err
	}

	return &model.FolderNodeDTO{Node: *node, Children: children}, nil
}

func NewGetNodeByIDUseCase(
	ur domain.UsrRepository,
	nr domain.NodeRepository,
	dr domain.DocRepository,
	gur domain.GroupUsrRepository,
) *GetNodeByIDUseCase {
	if ur == nil || nr == nil || dr == nil {
		log.Panicln("NewGetNodeByIDUseCase received a nil reference as dependency")
	}

	return &GetNodeByIDUseCase{ur, nr, dr, gur}
}
