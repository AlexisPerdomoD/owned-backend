package dto

import (
	"ownned/internal/domain"
	"ownned/pkg/apperror"
)

type CreateAccessDTO struct {
	GroupID domain.GroupID        `json:"group_id" validate:"required,uuid7"`
	Access  domain.GroupUsrAccess `json:"access" validate:"required,oneof=read_only_access write_access"`
}

type CreateUsrDTO struct {
	Role      domain.UsrRole    `json:"role" validate:"required,oneof=super_usr_role normal_usr_role limited_usr_role"`
	Firstname string            `json:"firstname" validate:"required,min=2,max=50"`
	Lastname  string            `json:"lastname" validate:"required,min=2,max=50"`
	Username  string            `json:"username" validate:"required,email"`
	Pwd       string            `validate:"required,min=8,max=255"`
	Access    []CreateAccessDTO `json:"access" validate:"required,dive"`
}

func (dto *CreateUsrDTO) Validate() error {
	if err := validate.Struct(dto); err != nil {
		return err
	}

	if !isValidPwd(dto.Pwd) {
		detail := make(map[string]string)
		detail["invalid_password"] = "password does not contain any letter, digit or symbol"
		return apperror.ErrBadRequest(detail)
	}

	return nil
}

type LoginUsrDTO struct {
	Username string `json:"username" validate:"required,email"`
	Pwd      string `json:"password" validate:"required,min=8,max=255"`
}

func (dto *LoginUsrDTO) Validate() error {
	return validate.Struct(dto)
}
