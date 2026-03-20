package view

import (
	"time"

	"ownned/internal/domain"
)

type NodeView struct {
	ID          domain.NodeID   `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        domain.NodeType `json:"type"`
	Path        domain.NodePath `json:"path"`
	Children    []NodeView      `json:"children,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func NodeViewFromDomain(n *domain.Node, chldr []NodeView) NodeView {
	if n == nil {
		return NodeView{}
	}

	return NodeView{
		ID:          n.ID,
		Name:        n.Name,
		Description: n.Description,
		Type:        n.Type,
		Path:        n.Path,
		Children:    chldr,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
	}
}
