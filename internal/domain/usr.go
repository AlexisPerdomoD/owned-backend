package domain

import (
	"context"
	"time"
)

type UsrID = string

type UsrRole int

const (
	SuperUsrRole UsrRole = iota
	NormalUsrRole
	LimitedUsrRole
)

func (r UsrRole) String() string {
	switch r {
	case SuperUsrRole:
		return "SuperUsrRole"
	case NormalUsrRole:
		return "NormalUsrRole"
	case LimitedUsrRole:
		return "LimitedUsrRole"
	default:
		return "UnknownUsrRole"

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

	Create(ctx context.Context, d *Usr) error

	Update(ctx context.Context, d *Usr) error

	Delete(ctx context.Context, id UsrID) error
}
