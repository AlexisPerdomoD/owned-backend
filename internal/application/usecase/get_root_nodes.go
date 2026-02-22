package usecase

import (
	"context"

	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
)

type GetRootNodesUseCase struct {
	nodeRepository     domain.NodeRepository
	usrRepository      domain.UsrRepository
	groupUsrRepository domain.GroupUsrRepository
}

func (uc *GetRootNodesUseCase) Execute(ctx context.Context, usrID domain.UsrID) ([]domain.Node, error) {
	usr, err := uc.usrRepository.GetByID(ctx, usrID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, apperror.ErrUnauthenticated(nil)
	}

	switch usr.Role {
	case domain.SuperUsrRole:
		return uc.nodeRepository.GetRoot(ctx)

	case domain.LimitedUsrRole, domain.NormalUsrRole:
		groups, err := uc.groupUsrRepository.GetByUsr(ctx, usr.ID)
		if err != nil {
			return nil, err
		}

		groupIDs := make([]domain.GroupID, len(groups))
		for i, g := range groups {
			groupIDs[i] = g.ID
		}

		return uc.nodeRepository.GetRootByGroups(ctx, groupIDs)

	default:
		return nil, apperror.ErrForbidden(nil)
	}
}

func NewGetRootNodesUseCase(
	nr domain.NodeRepository,
	ur domain.UsrRepository,
	gur domain.GroupUsrRepository,
) *GetRootNodesUseCase {
	helper.NotNilOrPanic(nr, "NodeRepository")
	helper.NotNilOrPanic(ur, "UsrRepository")
	helper.NotNilOrPanic(gur, "GroupUsrRepository")
	return &GetRootNodesUseCase{nr, ur, gur}
}
