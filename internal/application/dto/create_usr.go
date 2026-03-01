package dto

import "ownned/internal/domain"

type AccessDTO struct {
	GroupID domain.GroupID     `json:"group_id" validate:"required,uuid7"`
	Access  domain.GroupAccess `json:"access" validate:"required,oneof=read_only_access write_access"`
}

type CreateUsrInputDTO struct {
	Role      domain.UsrRole `json:"role" validate:"required,oneof=super_usr_role normal_usr_role limited_usr_role"`
	Firstname string         `json:"firstname" validate:"required,min=2,max=50"`
	Lastname  string         `json:"lastname" validate:"required,min=2,max=50"`
	Username  string         `json:"username" validate:"required,email"`
	Pwd       string         `json:"pwd" validate:"required,min=8,max=255,alphanum,excludes=\\/"`
	Access    []AccessDTO    `json:"access" validate:"required,dive"`
}

func (dto *CreateUsrInputDTO) Validate() error {
	return validate.Struct(dto)
}
