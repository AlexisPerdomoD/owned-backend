package view

import (
	"time"

	"github.com/google/uuid"
	"ownned/internal/domain"
)

type GroupView struct {
	ID          uuid.UUID `json:"id"`
	UsrID       uuid.UUID `json:"usr_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GroupViewFromDomain(group *domain.Group) GroupView {
	if group == nil {
		return GroupView{}
	}

	return GroupView{
		ID:          group.ID,
		UsrID:       group.UsrID,
		Name:        group.Name,
		Description: group.Description,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
}
