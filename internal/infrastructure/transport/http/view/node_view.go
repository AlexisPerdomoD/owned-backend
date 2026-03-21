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
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Children    []NodeView      `json:"children,omitempty"`
	Doc         *DocView        `json:"doc,omitempty"`
}

func FolderViewFromDomain(n *domain.Node, chldr []domain.Node) NodeView {
	if n == nil {
		return NodeView{}
	}

	chldrView := make([]NodeView, len(chldr))
	for i, c := range chldr {
		chldrView[i] = FolderViewFromDomain(&c, nil)
	}

	return NodeView{
		ID:          n.ID,
		Name:        n.Name,
		Description: n.Description,
		Type:        n.Type,
		Path:        n.Path,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
		Children:    chldrView,
	}
}

func FileViewFromDomain(n *domain.Node, doc *domain.Doc) NodeView {
	if n == nil {
		return NodeView{}
	}
	docView := DocViewFromDomain(doc)
	return NodeView{
		ID:          n.ID,
		Name:        n.Name,
		Description: n.Description,
		Type:        n.Type,
		Path:        n.Path,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
		Doc:         &docView,
	}
}

type NodeCommentView struct {
	ID        domain.NodeCommentID `json:"id"`
	NodeID    domain.NodeID        `json:"node_id"`
	UsrID     domain.UsrID         `json:"usr_id"`
	Content   string               `json:"content"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

func NodeCommentViewFromDomain(c *domain.NodeComment) NodeCommentView {
	if c == nil {
		return NodeCommentView{}
	}
	return NodeCommentView{
		ID:        c.ID,
		NodeID:    c.NodeID,
		UsrID:     c.UsrID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
