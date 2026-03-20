package dto

import (
	"io"

	"ownned/internal/domain"
)

type CreateDocInputDTO struct {
	ParentID     domain.NodeID `json:"parentID"`
	Title        string        `json:"title" validate:"required,alphanum,min=1,max=255"`
	Description  string        `json:"description" validate:"max=255"`
	ExpectedSize uint64        `json:"size"`
	Mimetype     string        `json:"mimetype"`
	File         io.ReadCloser `json:"file"`
}

func (dto *CreateDocInputDTO) Validate() error {
	return validate.Struct(dto)
}
