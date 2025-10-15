package domain

import (
	"context"
	"time"
)

type NodeType string

const (
	FolderNodeType NodeType = "folder"
	FileNodeType   NodeType = "file"
)

type NodeID = string

type Node struct {
	ID          NodeID
	ParentID    *NodeID
	Name        string
	Description string
	Path        string
	Type        NodeType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NodeRepository interface {
	GetByID(ctx context.Context, id NodeID) (*Node, error)

	Create(ctx context.Context, n *Node) error

	Update(ctx context.Context, n *Node) error

	Delete(ctx context.Context, id NodeID) error
}
