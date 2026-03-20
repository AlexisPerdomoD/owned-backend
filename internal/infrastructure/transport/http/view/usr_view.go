package view

import (
	"time"

	"ownned/internal/domain"
)

type UsrView struct {
	ID        string         `json:"id"`
	Role      domain.UsrRole `json:"role"`
	RoleTitle string         `json:"roleTitle"`
	Firstname string         `json:"firstname"`
	Lastname  string         `json:"lastname"`
	Username  string         `json:"username"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func UsrViewFromDomain(usr *domain.Usr) UsrView {
	if usr == nil {
		return UsrView{}
	}
	return UsrView{
		ID:        usr.ID.String(),
		Role:      usr.Role,
		RoleTitle: string(usr.Role),
		Firstname: usr.Firstname,
		Lastname:  usr.Lastname,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
	}
}
