package handler

import (
	"net/http"

	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/transport/http/decoder"
	"ownned/internal/infrastructure/transport/http/mapper"
	"ownned/internal/infrastructure/transport/http/response"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UsrHandler struct {
	loginUsrUseCase  *usecase.LoginUsrUseCase
	createUsrUseCase *usecase.CreateUsrUseCase
	getUsrUseCase    *usecase.GetUsrUseCase
	secure           bool
}

func (c *UsrHandler) GetUsrHandler(w http.ResponseWriter, r *http.Request) {
	usrID, err := uuid.Parse(chi.URLParam(r, "usrID"))
	if err != nil {
		_ = response.WriteJSONError(w, apperror.ErrBadRequest(map[string]string{"usrID": "invalido"}))
		return
	}

	usr, err := c.getUsrUseCase.Execute(r.Context(), usrID)
	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	view := mapper.UsrViewFromDomain(usr)
	_ = response.WriteJSON(w, http.StatusOK, view)
}

func (c *UsrHandler) CreateUsrHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: mejorar
	defer r.Body.Close()

	body, err := decoder.CreateUsrDTOFromJSON(r.Body)
	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	ctx := r.Context()
	usr, err := c.createUsrUseCase.Execute(ctx, *body)
	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	view := mapper.UsrViewFromDomain(usr)
	_ = response.WriteJSON(w, http.StatusCreated, view)
}

func (c *UsrHandler) LoginUsrHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: mejorar
	defer r.Body.Close()

	body, err := decoder.LoginUsrDTOFromJSON(r.Body)
	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	ctx := r.Context()
	sessionToken, err := c.loginUsrUseCase.Execute(ctx, *body)
	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		MaxAge:   3600,
	})

	resp := make(map[string]string)
	resp["message"] = "logged properly"
	_ = response.WriteJSON(w, http.StatusCreated, resp)
}

func NewUsrHandler(
	lu *usecase.LoginUsrUseCase,
	cu *usecase.CreateUsrUseCase,
	gu *usecase.GetUsrUseCase,
	secure bool,
) *UsrHandler {
	helper.NotNilOrPanic(lu, "LoginUsrUseCase")
	helper.NotNilOrPanic(cu, "CreateUsrUseCase")
	helper.NotNilOrPanic(gu, "GetUsrUseCase")
	return &UsrHandler{lu, cu, gu, secure}
}
