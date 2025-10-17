package dto

import "ownned/internal/domain"

type CreateUsrInputDto struct {
	Data   domain.Usr
	Access []domain.NodeID
}
