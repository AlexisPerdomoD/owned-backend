package handler

import (
	"fmt"
	"net/http"

	"ownned/internal/application/dto"
	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/sctx"
	"ownned/internal/infrastructure/transport/http/decoder"
	"ownned/internal/infrastructure/transport/http/encoder"
	"ownned/internal/infrastructure/transport/http/view"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GroupHandler struct {
	createGroup *usecase.CreateGroupUseCase
	deleteGroup *usecase.DeleteGroupUseCase
}

func (h *GroupHandler) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	usrID, err := sctx.GetUsrID(r.Context())
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	body, err := decoder.ReadFromJSON[dto.CreateGroupDTO](r.Body)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	if err := body.Validate(); err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	group, err := h.createGroup.Execute(r.Context(), usrID, body)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	_ = encoder.WriteJSON(w, http.StatusCreated, view.GroupViewFromDomain(group))
}

func (h *GroupHandler) DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	usrID, err := sctx.GetUsrID(r.Context())
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	groupID, err := uuid.Parse(chi.URLParam(r, "groupID"))
	if err != nil {
		detail := make(map[string]string)
		detail["reason"] = fmt.Sprintf("Group ID was not valid UUID: %s", err)
		_ = encoder.WriteJSONError(w, apperror.ErrBadRequest(detail))
		return
	}

	group, err := h.deleteGroup.Execute(r.Context(), usrID, groupID)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	_ = encoder.WriteJSON(w, http.StatusOK, view.GroupViewFromDomain(group))
}

func NewGroupHandler(cg *usecase.CreateGroupUseCase, dg *usecase.DeleteGroupUseCase) *GroupHandler {
	helper.NotNilOrPanic(cg, "CreateGroupUseCase")
	helper.NotNilOrPanic(dg, "DeleteGroupUseCase")
	return &GroupHandler{cg, dg}
}
