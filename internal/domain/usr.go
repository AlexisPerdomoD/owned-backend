package domain

import (
	"context"
	"time"
)

type UsrID = string

type UsrRole string

const (
	SuperUsrRole   UsrRole = "super_usr_role"
	NormalUsrRole  UsrRole = "normal_usr_role"
	LimitedUsrRole UsrRole = "limited_usr_role"
)

func (r UsrRole) IsValid() bool {
	switch r {
	case SuperUsrRole, NormalUsrRole, LimitedUsrRole:
		return true
	default:
		return false
	}
}

type Usr struct {
	ID        UsrID
	Role      UsrRole
	Firstname string
	Lastname  string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UsrRepository interface {
	GetByID(ctx context.Context, id UsrID) (*Usr, error)

	GetByUsername(ctx context.Context, username string) (*Usr, error)

	Create(ctx context.Context, d *Usr) error

	Update(ctx context.Context, d *Usr) error

	Delete(ctx context.Context, id UsrID) error
}
