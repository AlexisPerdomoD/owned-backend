package handler

import (
	"net/http"

	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/sctx"
	"ownned/internal/infrastructure/transport/http/response"
	"ownned/internal/infrastructure/transport/http/view"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"

	"github.com/google/uuid"
)

type NodeHandler struct {
	getRoot *usecase.GetRootNodesUseCase
}

func (c *NodeHandler) GetRootHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sctx.GetSession(r.Context())
	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	usrID, err := uuid.Parse(session.UsrID)
	if err != nil {
		_ = response.WriteJSONError(w, apperror.ErrUnauthenticated(nil))
		return
	}

	nodes, err := c.getRoot.Execute(r.Context(), usrID)
	if err != nil {
		_ = response.WriteJSONError(w, err)
		return
	}

	views := make([]view.NodeView, len(nodes))
	for i, n := range nodes {
		views[i] = view.NodeViewFromDomain(&n, nil)
	}

	_ = response.WriteJSON(w, http.StatusOK, nodes)
}

func NewNodeHandler(
	gr *usecase.GetRootNodesUseCase,
) *NodeHandler {
	helper.NotNilOrPanic(gr, "GetRootNodesUseCase")
	return &NodeHandler{gr}
}
