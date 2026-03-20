package decoder

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"ownned/internal/application/dto"

	"github.com/google/uuid"
)

func CreateDocInputDTOFromMultipartOnDemand(r *http.Request) (*dto.CreateDocInputDTO, error) {
	form, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	dto := dto.CreateDocInputDTO{}

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
			dto.Title = part.FileName()
			dto.Mimetype = part.Header.Get("Content-Type")
			dto.File = part
			continue
		}

		if err := part.Close(); err != nil {
			return nil, err
		}

	}

	return &dto, nil
}
