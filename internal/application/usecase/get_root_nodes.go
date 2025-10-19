package usecase

import (
	"context"
	"ownned/internal/domain"
	"ownned/internal/pkg/error_pkg"
	h "ownned/internal/pkg/helper_pkg"
)

type GetRootNodesUseCase struct {
	nodeRepository domain.NodeRepository
	usrRepository  domain.UsrRepository
}

func (uc *GetRootNodesUseCase) Execute(ctx context.Context, usrID domain.UsrID) ([]domain.Node, error) {
	usr, err := uc.usrRepository.GetByID(ctx, usrID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, error_pkg.ErrUnauthenticated(nil)
	}

	switch usr.Role {
	case domain.SuperUsrRole:
		return uc.nodeRepository.GetRoot(ctx)

	case domain.LimitedUsrRole, domain.NormalUsrRole:
		return uc.nodeRepository.GetRootByUsr(ctx, usr.ID)

	default:
		return nil, error_pkg.ErrForbidden(nil)
	}
}

func NewGetRootNodesUseCase(
	nr domain.NodeRepository,
	ur domain.UsrRepository,
) *GetRootNodesUseCase {
	h.AssertNotNil(nr, "NodeRepository")
	h.AssertNotNil(ur, "UsrRepository")
	return &GetRootNodesUseCase{nr, ur}
}
