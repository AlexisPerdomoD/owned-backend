package handler

import (
	"net/http"
	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/auth"
	"ownned/internal/infrastructure/transport/http/decoder"
	"ownned/internal/infrastructure/transport/http/mapper"
	"ownned/internal/infrastructure/transport/http/response"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UsrHandler struct {
	createUsrUseCase *usecase.CreateUsrUseCase
	getUsrUseCase    *usecase.GetUsrUseCase
}

func (c *UsrHandler) GetUsrHandler(w http.ResponseWriter, r *http.Request) {
	usrID, err := uuid.Parse(chi.URLParam(r, "usrID"))
	if err != nil {
		_ = response.WriteJSONError(w, apperror.ErrBadRequest(map[string]string{"usrID": "invalido"}))
		return
	}

	usr, err := c.getUsrUseCase.Execute(r.Context(), usrID.String())
	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	view := mapper.UsrViewFromDomain(usr)
	_ = response.WriteJSON(w, http.StatusOK, view)
}

func (c *UsrHandler) CreateUsrHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: mejorar
	defer func() { _ = r.Body.Close() }()

	ctx := r.Context()
	session, err := auth.GetSession(ctx)
	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	body, err := decoder.CreateUsrDTOFromJSON(r.Body)

	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	if err := body.Validate(); err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	usr, err := c.createUsrUseCase.Execute(ctx, session.UserID, *body)
	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	view := mapper.UsrViewFromDomain(usr)
	_ = response.WriteJSON(w, http.StatusCreated, view)
}

func NewUsrHandler(
	cu *usecase.CreateUsrUseCase,
	gu *usecase.GetUsrUseCase,
) *UsrHandler {
	helper.AssertNotNil(cu, "CreateUsrUseCase")
	helper.AssertNotNil(gu, "GetUsrUseCase")
	return &UsrHandler{cu, gu}
}
