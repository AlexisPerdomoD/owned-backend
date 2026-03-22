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
	ur  domain.UsrRepository
	uow domain.UnitOfWorkFactory
}

func (uc *CreateGroupUseCase) Execute(
	ctx context.Context,
	usrID domain.UsrID,
	args *dto.CreateGroupDTO,
) (*domain.Group, error) {
	usr, err := uc.ur.GetByID(ctx, usrID)
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

	groupID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	group := &domain.Group{
		ID:          groupID,
		UsrID:       usr.ID,
		Name:        args.Name,
		Description: args.Description,
	}

	groupUsr := &domain.UpsertGroupUsr{
		GroupID: group.ID,
		UsrID:   usr.ID,
		Access:  domain.GroupOwnerAccess,
	}

	if err := uc.uow.Do(ctx, func(tx domain.UnitOfWork) error {
		if err := tx.GroupRepository().Create(tx.Ctx(), group); err != nil {
			return err
		}

		if err := tx.GroupUsrRepository().Upsert(tx.Ctx(), groupUsr); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return group, nil
}

func NewCreateGroupUseCase(ur domain.UsrRepository, uowf domain.UnitOfWorkFactory) *CreateGroupUseCase {
	helper.NotNilOrPanic(ur, "UsrRepository")
	helper.NotNilOrPanic(uowf, "UnitOfWorkFactory")
	return &CreateGroupUseCase{ur, uowf}
}
