package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"ownned/internal/application/dto"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
)

type GetNodeUseCase struct {
	accessChecker
	ur  domain.UsrRepository
	nr  domain.NodeRepository
	dr  domain.DocRepository
	log *slog.Logger
}

func (uc *GetNodeUseCase) Execute(ctx context.Context, usrID domain.UsrID, nodeID domain.NodeID) (domain.NodeLike, error) {
	usr, err := uc.ur.GetByID(ctx, usrID)
	if err != nil {
		uc.log.WarnContext(ctx, "failed to get user by ID", "usrID", usrID, "error", err)
		return nil, err
	}

	if usr == nil {
		return nil, apperror.ErrForbidden(nil)
	}

	n, err := uc.nr.GetByID(ctx, nodeID)
	if err != nil {
		uc.log.WarnContext(ctx, "failed to get node by ID", "nodeID", nodeID, "error", err)
		return nil, err
	}

	if n == nil {
		detail := make(map[string]string)
		detail["reason"] = fmt.Sprintf("Node with ID=%s was not found", nodeID.String())
		return nil, apperror.ErrNotFound(detail)
	}

	canDo, err := uc.hasAccessTo(ctx, usr, n.Path, domain.GroupReadOnlyAccess)
	if err != nil {
		uc.log.WarnContext(ctx, "failed to check if user can access node", "nodeID", nodeID, "error", err)
		return nil, err
	}

	if !canDo {
		detail := make(map[string]string)
		detail["reason"] = fmt.Sprintf("User does not have access to specified node ID=%s", nodeID.String())
		return nil, apperror.ErrForbidden(detail)
	}

	if n.Type == domain.FileNodeType {
		doc, err := uc.dr.GetByNodeID(ctx, n.ID)
		if err != nil {
			uc.log.WarnContext(ctx, "failed to get doc by node ID", "nodeID", nodeID, "error", err)
			return nil, err
		}

		if doc == nil {
			detail := make(map[string]string)
			detail["reason"] = fmt.Sprintf("Doc entity associated with node at path %s was not found", n.Path)
			return nil, apperror.ErrNotFound(map[string]string{"error": "doc entity was not found"})
		}

		return &dto.FileNodeDTO{Node: *n, Doc: *doc}, nil
	}

	chld, err := uc.nr.GetChildren(ctx, n.Path)
	if err != nil {
		uc.log.WarnContext(ctx, "failed to get children", "nodeID", nodeID, "error", err)
		return nil, err
	}

	return &dto.FolderNodeDTO{Node: *n, Children: chld}, nil
}

func NewGetNodeByIDUseCase(
	ur domain.UsrRepository,
	nr domain.NodeRepository,
	dr domain.DocRepository,
	gur domain.GroupUsrRepository,
	mainLogger *slog.Logger,
) *GetNodeUseCase {
	helper.NotNilOrPanic(ur, "UsrRepository")
	helper.NotNilOrPanic(nr, "NodeRepository")
	helper.NotNilOrPanic(dr, "DocRepository")
	helper.NotNilOrPanic(gur, "GroupUsrRepository")
	helper.NotNilOrPanic(mainLogger, "mainLogger")
	log := mainLogger.With("usecase", "GetNodeUseCase")
	ac := accessChecker{gur}
	return &GetNodeUseCase{ac, ur, nr, dr, log}
}
