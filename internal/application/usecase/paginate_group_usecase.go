package usecase

import (
	"context"

	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
	"ownned/pkg/pagination"
)

type PaginateGroupUseCase struct {
	usrRepository   domain.UsrRepository
	groupRepository domain.GroupRepository
}

func (uc *PaginateGroupUseCase) Execute(
	ctx context.Context,
	usrID domain.UsrID,
	page uint,
	limit uint,
	search string,
	onlyMyGroups bool,
) (*pagination.PaginationResult[domain.Group], error) {
	param := domain.GroupPaginateParam{
		UsrID: nil,
		PaginationParam: pagination.PaginationParam{
			Page:   page,
			Limit:  limit,
			Search: search,
		},
	}

	usr, err := uc.usrRepository.GetByID(ctx, usrID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, apperror.ErrUnauthenticated(nil)
	}

	if usr.Role != domain.SuperUsrRole || onlyMyGroups {
		param.UsrID = &usr.ID
	}

	res, err := uc.groupRepository.Paginate(ctx, param)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewPaginateGroupUseCase(
	ur domain.UsrRepository,
	gr domain.GroupRepository,
) *PaginateGroupUseCase {
	helper.NotNilOrPanic(ur, "UsrRepository")
	helper.NotNilOrPanic(gr, "GroupRepository")
	return &PaginateGroupUseCase{ur, gr}
}
