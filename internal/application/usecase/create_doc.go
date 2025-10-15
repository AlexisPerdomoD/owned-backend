package usecase

import (
	"ownned/internal/domain"
)

type CreateDocUseCase struct {
	ur domain.UsrRepository
	dr domain.DocRepository
	nr domain.DocRepository
}

func (uc *CreateDocUseCase) Execute() {

}

func NewCreateDocUseCase(
	ur domain.UsrRepository,
	dr domain.DocRepository,
	nr domain.DocRepository,
) *CreateDocUseCase {
	return &CreateDocUseCase{ur, dr, nr}
}
