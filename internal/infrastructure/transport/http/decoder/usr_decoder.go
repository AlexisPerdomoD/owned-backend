// Package decoder contains the decoders for the transport layer
package decoder

import (
	"encoding/json"
	"io"

	"ownned/internal/application/dto"
)

func CreateUsrDTOFromJSON(r io.Reader) (*dto.CreateUsrInputDTO, error) {
	var dto dto.CreateUsrInputDTO
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		return nil, err
	}

	return &dto, nil
}
