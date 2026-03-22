package usecase

import (
	"context"
	"fmt"

	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
)

type DeleteGroupUseCase struct {
	accessChecker
	usrRepository   domain.UsrRepository
	groupRepository domain.GroupRepository
}

func (uc *DeleteGroupUseCase) Execute(ctx context.Context, usrID domain.UsrID, groupID domain.GroupID) (*domain.Group, error) {
	usr, err := uc.usrRepository.GetByID(ctx, usrID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, apperror.ErrUnauthenticated(nil)
	}

	if !usr.Role.CanDeleteGroup() {
		detail := make(map[string]string)
		detail["reason"] = "User can not do this action."
		return nil, apperror.ErrForbidden(detail)
	}

	group, err := uc.groupRepository.GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if group == nil {
		detail := make(map[string]string)
		detail["reason"] = "Group does not exist with ID=" + groupID.String()
		return nil, apperror.ErrNotFound(detail)
	}

	canDo, err := uc.hasGroupAccessTo(ctx, usr, group.ID, domain.GroupOwnerAccess)
	if err != nil {
		return nil, err
	}

	if !canDo {
		detail := make(map[string]string)
		detail["reason"] = fmt.Sprintf("User does not have access to specified group ID=%s", groupID)
		return nil, apperror.ErrForbidden(detail)
	}

	if err := uc.groupRepository.Delete(ctx, group.ID); err != nil {
		return nil, err
	}

	return group, nil
}

func NewDeleteGroupUseCase(ur domain.UsrRepository, gr domain.GroupRepository, gur domain.GroupUsrRepository) *DeleteGroupUseCase {
	helper.NotNilOrPanic(ur, "usrRepository")
	helper.NotNilOrPanic(gr, "groupRepository")
	ac := accessChecker{gur}
	return &DeleteGroupUseCase{ac, ur, gr}
}
