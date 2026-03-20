package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/transport/http/decoder"
	"ownned/internal/infrastructure/transport/http/encoder"
	"ownned/internal/infrastructure/transport/http/view"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
)

type UsrHandlerConfig struct {
	Secure   bool
	SameSite http.SameSite
}

type UsrHandler struct {
	loginUsr  *usecase.LoginUsrUseCase
	createUsr *usecase.CreateUsrUseCase
	getUsr    *usecase.GetUsrUseCase
	cfg       UsrHandlerConfig
}

func (c *UsrHandler) GetUsrHandler(w http.ResponseWriter, r *http.Request) {
	usrID, err := uuid.Parse(chi.URLParam(r, "usrID"))
	if err != nil {
		detail := make(map[string]string)
		detail["usrID"] = "invalid uuid provided"
		_ = encoder.WriteJSONError(w, apperror.ErrBadRequest(detail))
		return
	}

	usr, err := c.getUsr.Execute(r.Context(), usrID)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	_ = encoder.WriteJSON(w, http.StatusOK, view.UsrViewFromDomain(usr))
}

func (c *UsrHandler) CreateUsrHandler(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	body, err := decoder.CreateUsrDTOFromJSON(r.Body)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	ctx := r.Context()
	usr, err := c.createUsr.Execute(ctx, *body)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	_ = encoder.WriteJSON(w, http.StatusCreated, view.UsrViewFromDomain(usr))
}

func (c *UsrHandler) LoginUsrHandler(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	body, err := decoder.LoginUsrDTOFromJSON(r.Body)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	sessionToken, err := c.loginUsr.Execute(r.Context(), *body)
	if err != nil {
		_ = encoder.WriteJSONError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   c.cfg.Secure,
		SameSite: c.cfg.SameSite,
		Path:     "/",
		MaxAge:   3600,
	})

	resp := make(map[string]string)
	resp["message"] = "logged properly"
	_ = encoder.WriteJSON(w, http.StatusCreated, resp)
}

func NewUsrHandler(
	lu *usecase.LoginUsrUseCase,
	cu *usecase.CreateUsrUseCase,
	gu *usecase.GetUsrUseCase,
	cfg UsrHandlerConfig,
) *UsrHandler {
	helper.NotNilOrPanic(lu, "LoginUsrUseCase")
	helper.NotNilOrPanic(cu, "CreateUsrUseCase")
	helper.NotNilOrPanic(gu, "GetUsrUseCase")
	return &UsrHandler{lu, cu, gu, cfg}
}
