package usecase

import (
	"context"

	"github.com/google/uuid"
	"ownned/internal/application/dto"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
)

type CreateGroupUseCase struct {
	usrRepository   domain.UsrRepository
	groupRepository domain.GroupRepository
}

func (uc *CreateGroupUseCase) Execute(
	ctx context.Context,
	usrID domain.UsrID,
	args dto.CreateGroupDTO,
) (*domain.Group, error) {
	usr, err := uc.usrRepository.GetByID(ctx, usrID)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, apperror.ErrUnauthenticated(nil)
	}

	if !usr.Role.CanCreateGroup() {
		detail := make(map[string]string)
		detail["reason"] = "Usr does not have permission to do this action."
		return nil, apperror.ErrForbidden(detail)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	group := &domain.Group{
		ID:          id,
		UsrID:       usr.ID,
		Name:        args.Name,
		Description: args.Description,
	}

	if err := uc.groupRepository.Create(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

func NewCreateGroupUseCase(ur domain.UsrRepository, gr domain.GroupRepository) *CreateGroupUseCase {
	helper.NotNilOrPanic(ur, "UsrRepository")
	helper.NotNilOrPanic(gr, "GroupRepository")
	return &CreateGroupUseCase{ur, gr}
}
