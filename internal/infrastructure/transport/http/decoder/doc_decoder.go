package decoder

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"ownned/internal/application/dto"

	"github.com/google/uuid"
)

func CreateDocInputDTOFromMultipartOnDemandOld(r *http.Request) (*dto.CreateDocInputDTO, error) {
	form, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	dto := dto.CreateDocInputDTO{}

label:
	for {
		part, err := form.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		switch part.FormName() {
		case "parent_id":
			data, err := io.ReadAll(part)
			if err != nil {
				_ = part.Close()
				return nil, err
			}

			dto.ParentID, err = uuid.Parse(string(data))
			if err != nil {
				_ = part.Close()
				return nil, err
			}

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
			dto.Filename = part.FileName()
			dto.Mimetype = part.Header.Get("Content-Type")
			dto.File = part
			continue label
		}

		if err := part.Close(); err != nil {
			return nil, err
		}

	}

	return &dto, nil
}

func CreateDocInputDTOFromMultipartOnDemand(r *http.Request) (*dto.CreateDocInputDTO, error) {
	form, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	dto := &dto.CreateDocInputDTO{}

reading:
	for {
		part, err := form.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch part.FormName() {
		case "file":
			// dejamos el part directo como dto.File
			dto.File = part
			dto.Filename = part.FileName()
			dto.Mimetype = part.Header.Get("Content-Type")
			break reading
		default:
			data, err := io.ReadAll(part)
			if err != nil {
				_ = part.Close()
				return nil, err
			}
			switch part.FormName() {
			case "parent_id":
				dto.ParentID, err = uuid.Parse(string(data))
				if err != nil {
					_ = part.Close()
					return nil, err
				}
			case "description":
				dto.Description = string(data)
			case "size":
				size, err := strconv.ParseUint(string(data), 10, 64)
				if err != nil {
					_ = part.Close()
					return nil, fmt.Errorf("invalid size provided %w", err)
				}
				dto.ExpectedSize = size
			}
			_ = part.Close()
		}
	}

	return dto, nil
}
