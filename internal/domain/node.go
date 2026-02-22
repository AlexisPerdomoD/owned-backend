package domain

import (
	"context"
	"time"
)

type NodePath string

type NodeType string

const (
	FolderNodeType NodeType = "folder"
	FileNodeType   NodeType = "file"
)

type NodeLike interface {
	GetNode() *Node

	IsFile() bool

	IsFolder() bool

	IsRoot() bool
}

type NodeID = string

type Node struct {
	ID          NodeID
	Name        string
	Description string
	Path        NodePath
	Type        NodeType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (n *Node) GetNode() *Node {
	return n
}

func (n *Node) IsFile() bool {
	return n.Type == FileNodeType
}

func (n *Node) IsFolder() bool {
	return n.Type == FolderNodeType
}

type NodeRepository interface {
	GetByID(ctx context.Context, id NodeID) (*Node, error)

	GetByIDs(ctx context.Context, ids []NodeID) ([]Node, error)

	GetChildren(ctx context.Context, path NodePath) ([]Node, error)

	GetRoot(ctx context.Context) ([]Node, error)

	GetRootByGroups(ctx context.Context, groups []GroupID) ([]Node, error)

	Create(ctx context.Context, n *Node) error

	Update(ctx context.Context, n *Node) error

	Delete(ctx context.Context, id NodeID) error
}
