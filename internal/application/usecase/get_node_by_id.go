package usecase

import (
	"context"
	"log"
	"ownned/internal/domain"
	error_pkg "ownned/internal/pkg/error_pkg"
)

type GetNodeByIDUseCase struct {
	ur domain.UsrRepository
	nr domain.NodeRepository
	dr domain.DocRepository
}

type GetNodeByIDUseCaseResponse interface {
	*domain.FolderNode | *domain.FileNode
}

func (uc *GetNodeByIDUseCase) Execute(ctx context.Context, usrID domain.UsrID, nodeID domain.NodeID) (GetNodeByIDUseCaseResponse, error) {
	nodeRepository := uc.nr
	usrRepository := uc.ur
	docRepository := uc.dr

	usr, err := usrRepository.GetByID(ctx, usrID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, error_pkg.ErrForbidden(nil)
	}

	node, err := nodeRepository.GetByID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	if node == nil {
		return nil, error_pkg.ErrNotFound(nil)
	}

	if usr.Role != domain.SuperUsrRole {
		access, err := nodeRepository.GetAccess(usr.ID, node.ID)
		if err != nil {
			return nil, err
		}

		if access == domain.NoAccess {
			return nil, error_pkg.ErrForbidden(nil)
		}
	}

	if node.Type == domain.FileNodeType {
		docs, err := docRepository.GetByNodeID(ctx, node.ID)
		if err != nil {
			return nil, err
		}

		return &domain.FileNode{Node: *node, Docs: docs}, nil
	}

	// children, err := nodeRepository.GetChildren()

	return &domain.FolderNode{Node: *node, Children: []domain.Node{}}, nil
}

func NewGetNodeByIDUseCase(
	ur domain.UsrRepository,
	nr domain.NodeRepository,
	dr domain.DocRepository,
) *GetNodeByIDUseCase {

	if ur == nil || nr == nil || dr == nil {
		log.Panicln("NewGetNodeByIDUseCase received a nil reference as dependency")
	}

	return &GetNodeByIDUseCase{ur, nr, dr}
}
