package dto

import (
	"encoding/json"
	"io"
	"ownned/internal/domain"
)

type CreateFolderInputDto struct {
	ParentID    *domain.NodeID `json:"parentID"`
	Name        string         `json:"name" validate:"required,alphanum,min=1,max=255,excludes=\\/"`
	Description string         `json:"description" validate:"max=255"`
}

func (dto *CreateFolderInputDto) GetData() *domain.Node {
	return &domain.Node{
		ParentID:    dto.ParentID,
		Name:        dto.Name,
		Description: dto.Description,
		Type:        domain.FolderNodeType,
	}
}

func (dto *CreateFolderInputDto) Validate() error {
	return validate.Struct(dto)
}

func ValidateCreateFolderInputDtoFromJSON(r io.Reader) (*CreateFolderInputDto, error) {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	var dto CreateFolderInputDto
	if err := decoder.Decode(&dto); err != nil {
		return nil, err
	}

	if err := dto.Validate(); err != nil {
		return nil, err
	}

	return &dto, nil
}
