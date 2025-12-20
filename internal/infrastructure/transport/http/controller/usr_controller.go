package controller

import (
	"net/http"
	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/transport/http/mapper"
	"ownned/internal/infrastructure/transport/http/response"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UsrController struct {
	createUsrUseCase *usecase.CreateUsrUseCase
	getUsrUseCase    *usecase.GetUsrUseCase
}

func (c *UsrController) GetUsrHandler(w http.ResponseWriter, r *http.Request) {
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

	view := mapper.MapUsrViewFrom(usr)
	_ = response.WriteJSON(w, http.StatusOK, view)
}

func (c *UsrController) CreateUsrHandler(w http.ResponseWriter, r *http.Request) {
	_ = response.WriteJSONError(w, apperror.ErrNotImplemented(nil))
}

func NewUsrController(
	cu *usecase.CreateUsrUseCase,
	gu *usecase.GetUsrUseCase,
) *UsrController {
	helper.AssertNotNil(cu, "CreateUsrUseCase")
	helper.AssertNotNil(gu, "GetUsrUseCase")
	return &UsrController{cu, gu}
}
