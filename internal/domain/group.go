package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GroupID is a group identifier in the system
type GroupID = uuid.UUID

// Group represents a identifier to be tag to a Folder or file in the system and be associated with a user(s)
type Group struct {
	ID          GroupID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GroupAccess string

const (
	GroupReadOnlyAccess GroupAccess = "read_only_access"
	GroupWriteAccess    GroupAccess = "write_access"
)

type UpdateGroup struct {
	Name        *string
	Description *string
}

func (g *UpdateGroup) IsEmpty() bool {
	return g.Name == nil && g.Description == nil
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

type GroupNode struct {
	Group
	NodeID     NodeID
	AssignDate time.Time
}

type UpsertGroupNode struct {
	GroupID GroupID
	NodeID  NodeID
}

// GroupRepository is the interface to interact with the group repository
type GroupRepository interface {
	// GetByID returns a group by its identifier
	GetByID(ctx context.Context, id GroupID) (*Group, error)
	// GetByUsr returns a list of groups attached to a user and their access
	GetByUsr(ctx context.Context, usrID UsrID) ([]GroupUsr, error)
	// GetByNode returns a list of groups attached to a node and their access
	GetByNode(ctx context.Context, nodeID NodeID) ([]GroupNode, error)
	// Create a new group in the system
	// Returns an error if nil data is provided
	Create(ctx context.Context, d *Group) error
	// Update a group in the system
	// Returns an error if nil data is provided
	Update(ctx context.Context, d *UpdateGroup) error
	// Delete a group in the system by its identifier if it exists
	Delete(ctx context.Context, id GroupID) error
}

// GroupUsrRepository is the interface to interact with usr - group relations
type GroupUsrRepository interface {
	// GetGroupAccess returns the access of a user to a Group based on usrs group access, if no access is found it returns nil
	GetGroupAccess(ctx context.Context, usrID UsrID, groupID GroupID) (*GroupAccess, error)
	// GetNodeAccess returns the access of a user to a Node based on usrs group access, if no access is found it returns nil
	GetNodeAccess(ctx context.Context, usrID UsrID, nodeID NodeID) (*GroupAccess, error)
	// Upsert access to a groups for a users
	// Returns an error if nil data is provided
	Upsert(ctx context.Context, d *UpsertGroupUsr) error
	// Upsert access to a groups for a users
	// Returns an error if nil data is provided
	UpsertAll(ctx context.Context, d []UpsertGroupUsr) error
	// RemoveUsr removes a user from a group if it exists on it
	RemoveUsr(ctx context.Context, g GroupID, u UsrID) error
}

// GroupNodeRepository is the interface to interact with node - group relations
type GroupNodeRepository interface {
	// Upsert attach a node to a group, does nothing if node is already attached
	// Returns an error if nil data is provided
	Upsert(ctx context.Context, d *UpsertGroupNode) error
	// Upsert attach to nodes to groups
	// Returns an error if nil data is provide
	UpsertAll(ctx context.Context, d []UpsertGroupNode) error
	// RemoveNode removes a node from a group if it exist on it
	RemoveNode(ctx context.Context, g GroupID, n NodeID) error
}
