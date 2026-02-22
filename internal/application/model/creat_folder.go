package model

import (
	"encoding/json"
	"io"

	"ownned/internal/domain"
)

type CreateFolderInputDTO struct {
	ParentID    *domain.NodeID `json:"parent_id"`
	Name        string         `json:"name" validate:"required,alphanum,min=1,max=255,excludes=\\/"`
	Description string         `json:"description" validate:"max=255"`
}

func (dto *CreateFolderInputDTO) GetData() (parentID *domain.NodeID, node *domain.Node) {
	return dto.ParentID, &domain.Node{
		Name:        dto.Name,
		Description: dto.Description,
		Type:        domain.FolderNodeType,
	}
}

func (dto *CreateFolderInputDTO) Validate() error {
	return validate.Struct(dto)
}

func NewCreateFolderInputDtoFromJSON(r io.Reader) (*CreateFolderInputDTO, error) {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	var dto CreateFolderInputDTO
	if err := decoder.Decode(&dto); err != nil {
		return nil, err
	}

	return &dto, dto.Validate()
}
