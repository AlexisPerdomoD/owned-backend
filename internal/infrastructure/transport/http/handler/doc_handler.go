package handler

import (
	"net/http"

	"github.com/google/uuid"
	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/sctx"
	"ownned/internal/infrastructure/transport/http/decoder"
	"ownned/internal/infrastructure/transport/http/encoder"
	"ownned/internal/infrastructure/transport/http/view"
	"ownned/pkg/helper"
)

type DocHandler struct {
	createDoc *usecase.CreateDocUseCase
}

func (h *DocHandler) CreateDocHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sctx.GetSession(r.Context())
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	usrID, err := uuid.Parse(session.UsrID)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	dto, err := decoder.CreateDocInputDTOFromMultipartOnDemand(r)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	defer func() {
		if dto.File != nil {
			_ = dto.File.Close()
		}
	}()

	fileNode, err := h.createDoc.Execute(r.Context(), usrID, dto)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	resp := view.FileViewFromDomain(fileNode.Node, fileNode.Doc)
	_ = encoder.WriteJSON(w, http.StatusCreated, resp)
}

func NewDocHandler(cduc *usecase.CreateDocUseCase) *DocHandler {
	helper.NotNilOrPanic(cduc, "CreateDocUseCase")
	return &DocHandler{cduc}
}
