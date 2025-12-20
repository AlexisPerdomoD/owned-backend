package controller

import (
	"log"
	"net/http"
	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/transport/http/mapper"
	"ownned/internal/infrastructure/transport/http/response"
	"ownned/pkg/apperror"

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
		httpErr := &mapper.ErrView{
			Code:    http.StatusBadRequest,
			Message: "usrID invalido",
		}
		_ = response.WriteJSON(w, httpErr.Code, httpErr)
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

func (c *UsrController) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", c.CreateUsrHandler)

	r.Get("/{usrID}", c.GetUsrHandler)

	return r
}

func NewUsrController(
	cu *usecase.CreateUsrUseCase,
	gu *usecase.GetUsrUseCase,
) *UsrController {

	if cu == nil || gu == nil {
		log.Panic("missing dependencies for NewUsrController")
	}

	return &UsrController{cu, gu}
}
