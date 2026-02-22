package domain

import (
	"context"
	"time"
)

type (

	// GroupID is a group identifier in the system
	GroupID string

	// Group represents a identifier to be tag to a Folder or file in the system and be associated with a user(s)
	Group struct {
		ID          GroupID
		Name        string
		Description string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
)

type GroupAccess string

const (
	GroupReadOnlyAccess GroupAccess = "read_only_access"
	GroupWriteAccess    GroupAccess = "write_access"
)

type CreateGroup struct {
	Name        string
	Description string
}

type UpdateGroup struct {
	Name        *string
	Description *string
}

func (g *UpdateGroup) IsEmpty() bool {
	return g.Name == nil && g.Description == nil
}

// GroupRepository is the interface to interact with the group repository
type GroupRepository interface {
	// GetByID returns a group by its identifier
	GetByID(ctx context.Context, id GroupID) (*Group, error)
	// Create a new group in the system
	// Returns an error if nil data is provided
	Create(ctx context.Context, d *CreateGroup) error
	// Update a group in the system
	// Returns an error if nil data is provided
	Update(ctx context.Context, d *UpdateGroup) error
	// Delete a group in the system by its identifier if it exists
	Delete(ctx context.Context, id GroupID) error
}

type GroupUsr struct {
	Group
	UsrID      UsrID
	Access     GroupAccess
	AssignDate time.Time
}

type UpsertGroupUsr struct {
	GroupID GroupID
	UsrID   UsrID
	Access  GroupAccess
}

// GroupUsrRepository is the interface to interact with usr - group relations
type GroupUsrRepository interface {
	// GetByGroup returns a list of access of users to a group
	GetByGroup(ctx context.Context, g GroupID) ([]GroupUsr, error)
	// GetByUsr returns a list of groups accessed by a user
	GetByUsr(ctx context.Context, usrID UsrID) ([]GroupUsr, error)
	// Upsert access to a groups for a users
	// Returns an error if nil data is provided
	Upsert(ctx context.Context, d *UpsertGroupUsr) error
	// Upsert access to a groups for a users
	// Returns an error if nil data is provided
	UpsertAll(ctx context.Context, d []UpsertGroupUsr) error
	// RemoveUsr removes a user from a group if it exists on it
	RemoveUsr(ctx context.Context, g GroupID, u UsrID) error
}

type GroupNode struct {
	Group
	NodeID     NodeID
	AssignDate time.Time
}

type UpsertGroupNode struct {
	GroupID GroupID
	NodeID  NodeID
}

// GroupNodeRepository is the interface to interact with node - group relations
type GroupNodeRepository interface {
	// GetByNode returns a list of groups attached to a node
	GetByNode(ctx context.Context, nodeID NodeID) ([]GroupNode, error)
	// GetByGroup returns a list of nodes attached to a group
	GetByGroup(ctx context.Context, groupID GroupID) ([]GroupNode, error)
	// Upsert attach a node to a group, does nothing if node is already attached
	// Returns an error if nil data is provided
	Upsert(ctx context.Context, d *UpsertGroupNode) error
	// Upsert attach to nodes to groups
	// Returns an error if nil data is provided
	UpsertAll(ctx context.Context, d []UpsertGroupNode) error
	// RemoveNode removes a node from a group if it exists on it
	RemoveNode(ctx context.Context, g GroupID, n NodeID) error
}
