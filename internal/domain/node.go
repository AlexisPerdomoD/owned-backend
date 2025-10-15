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

type NodeAccess uint8

const (
	NoAccess NodeAccess = iota
	ReadOnlyAccess
	WriteAccess
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

type FileNode struct {
	Node
	Docs []Doc
}

type FolderNode struct {
	Node
	Children []Node
}

type ChildNode interface {
	FileNode | FolderNode
}

type NodeRepository interface {
	GetByID(ctx context.Context, id NodeID) (*Node, error)

	Create(ctx context.Context, n *Node) error

	Update(ctx context.Context, n *Node) error

	Delete(ctx context.Context, id NodeID) error

	GetAccess(u UsrID, n NodeID) (NodeAccess, error)

	UpdateAccess(u UsrID, n NodeID, a NodeAccess) error
}
