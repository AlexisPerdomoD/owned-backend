package dto

import "ownned/internal/domain"

type FileNodeDTO struct {
	domain.Node
	Doc domain.Doc
}

type FolderNodeDTO struct {
	domain.Node
	Children []domain.Node
}

type CreateFolderDTO struct {
	ParentID    string `json:"parent_id" validate:"required,uuid"`
	Name        string `json:"name" validate:"required,min=1,max=255,excludes=\\/"`
	Description string `json:"description" validate:"max=255"`
}

func (dto *CreateFolderDTO) Validate() error {
	return validate.Struct(dto)
}
