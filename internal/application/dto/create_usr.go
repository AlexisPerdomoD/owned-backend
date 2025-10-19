package dto

import (
	"encoding/json"
	"io"
	"ownned/internal/domain"
)

type CreateUsrInputDto struct {
	Role      domain.UsrRole  `json:"role" validate:"required,min=0,max=2"`
	Firstname string          `json:"firstname" validate:"required,min=2,max=50"`
	Lastname  string          `json:"lastname" validate:"required,min=2,max=50"`
	Username  string          `json:"username" validate:"required,email"`
	Access    []domain.NodeID `json:"access" validate:"required,dive,uuid4"`
}

func (dto *CreateUsrInputDto) GetData() *domain.Usr {
	return &domain.Usr{
		Role:      dto.Role,
		Firstname: dto.Firstname,
		Lastname:  dto.Lastname,
		Username:  dto.Username,
	}
}

func (dto *CreateUsrInputDto) GetUsrAccess() []domain.NodeID {
	return dto.Access
}

func (dto *CreateUsrInputDto) Validate() error {
	return validate.Struct(dto)
}

func ValidateUsrInputDtoFromJSON(r io.Reader) (*CreateUsrInputDto, error) {
	var dto CreateUsrInputDto
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		return nil, err
	}

	if err := dto.Validate(); err != nil {
		return nil, err
	}

	return &dto, nil
}
