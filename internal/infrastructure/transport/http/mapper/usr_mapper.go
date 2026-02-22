package mapper

import (
	"ownned/internal/domain"
	"ownned/internal/infrastructure/transport/http/model"
)

func UsrViewFromDomain(usr *domain.Usr) *model.UsrView {
	if usr == nil {
		return nil
	}
	return &model.UsrView{
		ID:        usr.ID,
		Role:      usr.Role,
		RoleTitle: string(usr.Role),
		Firstname: usr.Firstname,
		Lastname:  usr.Lastname,
		CreatedAt: usr.CreatedAt,
	}
}
