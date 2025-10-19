package dto

import (
	"encoding/json"
	"io"
	"log/slog"
	"ownned/internal/domain"

	"github.com/go-playground/validator/v10"
)

type CreateUsrInputDto struct {
	Role      domain.UsrRole  `json:"role" validate:"required"`
	Firstname string          `json:"firstname" validate:"required,min=2,max=50"`
	Lastname  string          `json:"lastname" validate:"required,min=2,max=50"`
	Username  string          `json:"username" validate:"required,email"`
	Access    []domain.NodeID `json:"access" validate:"required,dive,uuid4"`
}

func (dto *CreateUsrInputDto) GetUsrData() *domain.Usr {
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
	return validator.New().Struct(dto)
}

func CreateUsrInputDtoFromJSON(r io.ReadCloser) (*CreateUsrInputDto, error) {
	defer func() {
		if err := r.Close(); err != nil {
			slog.Warn("CreateUsrInputDtoFromJSON could not close properly body", "err", err)
		}
	}()

	var dto CreateUsrInputDto
	if err := json.NewDecoder(r).Decode(&dto); err != nil {
		return nil, err
	}

	return &dto, nil
}
