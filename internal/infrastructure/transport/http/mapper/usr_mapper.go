package mapper

import (
	"ownned/internal/domain"
	"time"
)

type UsrView struct {
	ID        string         `json:"id"`
	Role      domain.UsrRole `json:"role"`
	RoleTitle string         `json:"roleTitle"`
	Firstname string         `json:"firstname"`
	Lastname  string         `json:"lastname"`
	Username  string         `json:"username"`
	CreatedAt time.Time      `json:"createdAt"`
}

func MapUsrViewFrom(usr *domain.Usr) *UsrView {
	if usr == nil {
		return nil
	}

	return &UsrView{
		ID:        usr.ID,
		Role:      usr.Role,
		RoleTitle: usr.Role.String(),
		Firstname: usr.Firstname,
		Lastname:  usr.Lastname,
		CreatedAt: usr.CreatedAt,
	}
}
