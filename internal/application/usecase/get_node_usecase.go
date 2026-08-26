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
	*accessChecker
	nr  domain.NodeRepository
	dr  domain.DocRepository
	log *slog.Logger
}

func (uc *GetNodeUseCase) Execute(ctx context.Context, nodeID domain.NodeID) (domain.NodeLike, error) {
	n, err := uc.nr.GetByID(ctx, nodeID)
	if err != nil {
		uc.log.WarnContext(ctx,
			"failed to get node by ID",
			"nodeID", nodeID,
			"error", err)
		return nil, err
	}

	if n == nil {
		detail := make(map[string]string)
		detail["reason"] = fmt.Sprintf("Node with ID=%s was not found", nodeID.String())
		return nil, apperror.ErrNotFound(detail)
	}

	canDo, err := uc.checkNodeAccessTo(ctx, n.Path, domain.GroupReadOnlyAccess)
	if err != nil {
		uc.log.WarnContext(ctx,
			"failed to check if user can access node",
			"nodeID", nodeID,
			"error", err)
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
	nr domain.NodeRepository,
	dr domain.DocRepository,
	mainLogger *slog.Logger,
	ac *accessChecker,
) *GetNodeUseCase {
	helper.NotNilOrPanic(nr, "NodeRepository")
	helper.NotNilOrPanic(dr, "DocRepository")
	helper.NotNilOrPanic(mainLogger, "mainLogger")
	helper.NotNilOrPanic(ac, "accessChecker")
	log := mainLogger.With("usecase", "GetNodeUseCase")
	return &GetNodeUseCase{ac, nr, dr, log}
}
