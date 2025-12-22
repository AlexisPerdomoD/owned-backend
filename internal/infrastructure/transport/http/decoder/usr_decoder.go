package decoder

import (
	"encoding/json"
	"io"
	"ownned/internal/application/model"
)

func CreateUsrDTOFromJSON(r io.Reader) (*model.CreateUsrInputDTO, error) {
	var dto model.CreateUsrInputDTO
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		return nil, err
	}

	return &dto, nil
}
