package dto

import "ownned/internal/domain"

type CreateGroupDTO struct {
	Name        string `json:"name" validate:"required,min=1,max=255,excludes=\\/"`
	Description string `json:"description" validate:"max=1000"`
}

func (dto *CreateGroupDTO) Validate() error {
	return validate.Struct(dto)
}

type PopulateGroup struct {
	domain.Group
	Nodes []domain.NodeGroupAttach
	Usrs  []domain.UsrGroupAccess
}
