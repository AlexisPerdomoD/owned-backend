package model

import "ownned/internal/domain"

type PopulateGroup struct {
	domain.Group
	Nodes []domain.NodeGroupAttach
	Usrs  []domain.UsrGroupAccess
}
