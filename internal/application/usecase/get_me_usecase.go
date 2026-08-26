package usecase

import (
	"context"

	"ownned/internal/domain"
	"ownned/pkg/helper"
)

type GetMeUseCase struct {
	identityChecker
	usrRepository domain.UsrRepository
}

func (uc *GetMeUseCase) Execute(ctx context.Context) (*domain.Usr, error) {
	identity, err := uc.getUsrIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return uc.usrRepository.GetByID(ctx, identity.ID)
}

func NewGetMeUseCase(ur domain.UsrRepository) *GetMeUseCase {
	helper.NotNilOrPanic(ur, "UsrRepository")
	return &GetMeUseCase{usrRepository: ur}
}
