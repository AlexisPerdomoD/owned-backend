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

type NodeLike interface {
	GetNode() *Node

	IsFile() bool

	IsFolder() bool

	IsRoot() bool
}

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

func (n *Node) GetNode() *Node {
	return n
}

func (n *Node) IsFile() bool {
	return n.Type == FileNodeType
}

func (n *Node) IsFolder() bool {
	return n.Type == FolderNodeType
}

func (n *Node) IsRoot() bool {
	return n.ParentID == nil
}

type NodeRepository interface {
	GetByID(ctx context.Context, id NodeID) (*Node, error)

	GetByIDs(ctx context.Context, ids []NodeID) ([]Node, error)

	GetChildren(ctx context.Context, folderID NodeID) ([]Node, error)

	GetRoot(ctx context.Context) ([]Node, error)

	GetRootByUsr(ctx context.Context, usrID UsrID) ([]Node, error)

	Create(ctx context.Context, n *Node) error

	Update(ctx context.Context, n *Node) error

	Delete(ctx context.Context, id NodeID) error

	GetAccess(ctx context.Context, u UsrID, n NodeID) (NodeAccess, error)

	UpdateAccess(ctx context.Context, u UsrID, n NodeID, a NodeAccess) error
}
