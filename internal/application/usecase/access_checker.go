package usecase

import (
	"context"
	"errors"

	"ownned/internal/application/auth"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/helper"

	"github.com/google/uuid"
)

// usrIdentity is the user identity
type usrIdentity struct {
	ID   domain.UsrID
	Role domain.UsrRole
}

type identityChecker struct{}

// getUsrIdentity returns the user identity from the context
func (ic identityChecker) getUsrIdentity(ctx context.Context) (i usrIdentity, err error) {
	s, err := auth.GetSession(ctx)
	if err != nil {
		return i, err
	}

	usrID, err := uuid.Parse(s.UsrID)
	if err != nil {
		detail := make(map[string]string)
		detail["reason"] = "invalid user ID"
		return i, apperror.ErrInternal(detail)
	}

	if !s.Role.IsValid() {
		detail := make(map[string]string)
		detail["reason"] = "invalid user role"
		return i, apperror.ErrForbidden(detail)
	}

	i.ID = usrID
	i.Role = s.Role
	return i, nil
}

// accessChecker is an abstraction to access business logic related
type accessChecker struct {
	identityChecker
	groupUsrRepository domain.GroupUsrRepository
}

// checkNodeAccessTo checks if a user has access to a node based on the user's role and the access of the node to the user
func (ic *accessChecker) checkNodeAccessTo(
	ctx context.Context,
	pth domain.NodePath,
	accs domain.GroupUsrAccess,
) (bool, error) {
	u, err := ic.getUsrIdentity(ctx)
	if err != nil {
		return false, err
	}

	if u.Role == domain.SuperUsrRole {
		return true, nil
	}

	if err := ic.groupUsrRepository.HasAccess(ctx, u.ID, pth, accs); err != nil {
		if errors.Is(err, domain.ErrNoAccess) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// checkGroupAccessTo checks if a user has access to a group based on the user's role and the access of the group to the user
func (ic *accessChecker) checkGroupAccessTo(
	ctx context.Context,
	groupID domain.GroupID,
	reqAccs domain.GroupUsrAccess,
) (bool, error) {
	u, err := ic.getUsrIdentity(ctx)
	if err != nil {
		return false, err
	}

	if u.Role == domain.SuperUsrRole {
		return true, nil
	}

	accs, err := ic.groupUsrRepository.GetGroupAccess(ctx, u.ID, groupID)
	if err != nil {
		if errors.Is(err, domain.ErrNoAccess) {
			return false, nil
		}

		return false, err
	}

	return accs.IsEquivalent(reqAccs), nil
}

func NewAccessChecker(gur domain.GroupUsrRepository) *accessChecker {
	helper.NotNilOrPanic(gur, "GroupUsrRepository")
	return &accessChecker{groupUsrRepository: gur}
}
