package usecase

import (
	"ownned/internal/domain"
)

type CreateDocUseCase struct {
	ur domain.UsrRepository
	dr domain.DocRepository
	nr domain.NodeRepository
}

func (uc *CreateDocUseCase) Execute(creatorID domain.UsrID, newUsr *domain.Usr) {

}

func NewCreateDocUseCase(
	ur domain.UsrRepository,
	dr domain.DocRepository,
	nr domain.NodeRepository,
) *CreateDocUseCase {
	return &CreateDocUseCase{ur, dr, nr}
}
