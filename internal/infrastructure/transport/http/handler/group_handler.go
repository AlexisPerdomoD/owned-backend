package handler

import (
	"net/http"

	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/sctx"
	"ownned/internal/infrastructure/transport/http/decoder"
	"ownned/internal/infrastructure/transport/http/encoder"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"

	"github.com/google/uuid"
)

type GroupHandler struct {
	createGroup *usecase.CreateGroupUseCase
	deleteGroup *usecase.DeleteGroupUseCase
}

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sctx.GetSession(r.Context())
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
	}

	usrID, err := uuid.Parse(session.UsrID)
	if err != nil {
		detail := make(map[string]string)
		detail["reason"] = "Invalid user ID provided."
		_ = encoder.WriteJSONError(w, apperror.ErrBadRequest(detail))
	}

	body, err := decoder.CreateDTOFromJSON(r.Body)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
	}

	if err := body.Validate(); err != nil {
		_ = encoder.WriteJSONError(w, err)
	}
}

func NewGroupHandler(cg *usecase.CreateGroupUseCase, dg *usecase.DeleteGroupUseCase) *GroupHandler {
	helper.NotNilOrPanic(cg, "CreateGroupUseCase")
	helper.NotNilOrPanic(dg, "DeleteGroupUseCase")
	return &GroupHandler{cg, dg}
}
