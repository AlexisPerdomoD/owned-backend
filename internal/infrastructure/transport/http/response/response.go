package response

import (
	"encoding/json"
	"net/http"
	"ownned/internal/infrastructure/transport/http/mapper"
)

func WriteJSON(w http.ResponseWriter, code int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(body)
}

func WriteJSONError(w http.ResponseWriter, err error) error {
	httpErr := mapper.MapError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Code)
	return json.NewEncoder(w).Encode(err)
}
