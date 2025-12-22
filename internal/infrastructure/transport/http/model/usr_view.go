package model

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
