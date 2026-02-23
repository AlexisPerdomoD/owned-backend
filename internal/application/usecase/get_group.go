package usecase

import (
	"context"

	"ownned/internal/application/model"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
)

type GetGroupUseCase struct {
	groupRepository domain.GroupRepository
	usrRepository   domain.UsrRepository
	nodeRepository  domain.NodeRepository
}

func (uc *GetGroupUseCase) Execute(ctx context.Context, id domain.GroupID) (*model.PopulateGroup, error) {
	group, err := uc.groupRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, apperror.ErrNotFound(map[string]string{"error": "group entity was not found"})
	}
	// this may be concurrent called with
	// errgroup.WithContext(ctx)
	// "golang.org/x/sync/errgroup"
	usrs, err := uc.usrRepository.GetByGroup(ctx, group.ID)
	if err != nil {
		return nil, err
	}

	nodes, err := uc.nodeRepository.GetByGroup(ctx, group.ID)
	if err != nil {
		return nil, err
	}

	resp := &model.PopulateGroup{
		Group: *group,
		Usrs:  usrs,
		Nodes: nodes,
	}

	return resp, nil
}

func NewGetGroupUseCase(
	gr domain.GroupRepository,
	ur domain.UsrRepository,
	nr domain.NodeRepository,
) *GetGroupUseCase {
	if gr == nil || ur == nil || nr == nil {
		panic("NewGetGroupUseCase received a nil reference as dependency")
	}
	return &GetGroupUseCase{gr, ur, nr}
}
