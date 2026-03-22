package dto

import "ownned/internal/domain"

type CreateGroupDTO struct {
	Name        string `json:"name" validate:"required,min=1,max=255,excludes=\\/"`
	Description string `json:"description" validate:"max=1000"`
}

type PopulateGroup struct {
	domain.Group
	Nodes []domain.NodeGroupAttach
	Usrs  []domain.UsrGroupAccess
}
