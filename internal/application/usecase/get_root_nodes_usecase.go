package usecase

import (
	"context"
	"log/slog"

	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
)

type GetRootNodesUseCase struct {
	nodeRepository  domain.NodeRepository
	usrRepository   domain.UsrRepository
	groupRepository domain.GroupRepository
	log             *slog.Logger
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
		groups, err := uc.groupRepository.GetByUsrAssigned(ctx, usr.ID)
		if err != nil {
			return nil, err
		}

		if len(groups) == 0 {
			detail := make(map[string]string)
			detail["reason"] = "no groups associated with user"
			return nil, apperror.ErrForbidden(detail)
		}

		gIDs := make([]domain.GroupID, len(groups))
		for i, g := range groups {
			gIDs[i] = g.ID
		}

		return uc.nodeRepository.GetRootByGroups(ctx, gIDs)

	default:
		uc.log.Warn("user role not supported for getting root nodes", "role", usr.Role)
		return nil, apperror.ErrForbidden(nil)
	}
}

func NewGetRootNodesUseCase(
	nr domain.NodeRepository,
	ur domain.UsrRepository,
	gr domain.GroupRepository,
	mainLogger *slog.Logger,
) *GetRootNodesUseCase {
	helper.NotNilOrPanic(nr, "NodeRepository")
	helper.NotNilOrPanic(ur, "UsrRepository")
	helper.NotNilOrPanic(gr, "GroupUsrRepository")
	helper.NotNilOrPanic(mainLogger, "mainLogger")
	log := mainLogger.With("usecase", "GetRootNodesUseCase")
	return &GetRootNodesUseCase{nr, ur, gr, log}
}
