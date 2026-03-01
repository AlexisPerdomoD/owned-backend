package dto

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

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

func NewCreateDocInputDtoFromMultipartOnDemand(r *http.Request) (*CreateDocInputDTO, error) {
	form, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	dto := CreateDocInputDTO{}

	for {
		part, err := form.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		switch part.FormName() {
		case "description":
			data, err := io.ReadAll(part)
			if err != nil {
				_ = part.Close()
				return nil, err
			}

			dto.Description = string(data)
		case "size":
			data, err := io.ReadAll(part)
			if err != nil {
				_ = part.Close()
				return nil, err
			}

			size, err := strconv.ParseUint(string(data), 10, 0)
			if err != nil {
				_ = part.Close()
				return nil, fmt.Errorf("invalid size provided %w", err)
			}

			dto.ExpectedSize = uint64(size)

		case "file":
			dto.Title = part.FileName()
			dto.Mimetype = part.Header.Get("Content-Type")
			dto.File = part
			continue
		}

		if err := part.Close(); err != nil {
			return nil, err
		}

	}

	return &dto, dto.Validate()
}
