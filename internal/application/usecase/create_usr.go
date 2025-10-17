package usecase

import (
	"context"
	"log"
	"ownned/internal/application/dto"
	"ownned/internal/domain"
	"ownned/internal/pkg/error_pkg"
	"slices"
)

type CreateUsrUseCase struct {
	ur domain.UsrRepository
	nr domain.NodeRepository
}

func (uc *CreateUsrUseCase) Execute(
	ctx context.Context,
	creatorID domain.UsrID,
	args dto.CreateUsrInputDto,
) (*domain.Usr, error) {
	nodeRepository := uc.nr
	usrRepository := uc.ur

	creator, err := usrRepository.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if creator.Role != domain.SuperUsrRole {
		return nil, error_pkg.ErrForbidden(map[string]string{"general": "usr does not have enought privileges to do this action"})
	}

	usr, err := usrRepository.GetByUsername(ctx, args.Data.Username)
	if err != nil {
		return nil, err
	}

	if usr != nil {
		return nil, error_pkg.ErrConflic(map[string]string{"general": "username already in use"})
	}

	if args.Data.Role != domain.SuperUsrRole && len(args.Access) > 0 {
slices.
	}

	return nil, nil
}

func NewCreateUsrUseCase(
	ur domain.UsrRepository,
	nr domain.NodeRepository,
) *CreateUsrUseCase {
	if ur == nil || nr == nil {
		log.Panicln("NewCreateUsrUseCase received a nil reference as dependency")
	}

	return &CreateUsrUseCase{ur, nr}
}
