package usecase

import (
	"context"
	"fmt"

	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"
)

type DeleteNodeCommentUseCase struct {
	identityChecker
	nodeCommentRepository domain.NodeCommentRepository
}

func (uc *DeleteNodeCommentUseCase) Execute(ctx context.Context, commentID domain.NodeCommentID) (*domain.NodeComment, error) {
	comment, err := uc.nodeCommentRepository.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if comment == nil {
		detail := make(map[string]string)
		detail["reason"] = fmt.Sprintf("NodeComment with ID=%s was not found", commentID)
		return nil, apperror.ErrNotFound(detail)
	}

	identity, err := uc.getUsrIdentity(ctx)
	if err != nil {
		return nil, err
	}

	if identity.Role != domain.SuperUsrRole && identity.ID != comment.UsrID {
		detail := make(map[string]string)
		detail["reason"] = fmt.Sprintf("User does not have access to NodeComment with ID='%s'", commentID)
		return nil, apperror.ErrForbidden(detail)
	}

	err = uc.nodeCommentRepository.Delete(ctx, comment.ID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func NewDeleteNodeCommentUseCase(ncr domain.NodeCommentRepository) *DeleteNodeCommentUseCase {
	helper.NotNilOrPanic(ncr, "NodeCommentRepository")
	return &DeleteNodeCommentUseCase{nodeCommentRepository: ncr}
}
