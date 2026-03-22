package usecase

import (
	"context"
	"errors"

	"ownned/internal/domain"
)

// accessChecker is an abstraction to access business logic related
type accessChecker struct {
	gur domain.GroupUsrRepository
}

// hasNodeAccessTo checks if a user has access to a node based on the user's role and the access of the node to the user
func (ac *accessChecker) hasNodeAccessTo(
	ctx context.Context,
	u *domain.Usr,
	pth domain.NodePath,
	accs domain.GroupUsrAccess,
) (bool, error) {
	if u.Role == domain.SuperUsrRole {
		return true, nil
	}

	if err := ac.gur.HasAccess(ctx, u.ID, pth, accs); err != nil {
		if errors.Is(err, domain.ErrNoAccess) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// hasGroupAccessTo checks if a user has access to a group based on the user's role and the access of the group to the user
func (ac *accessChecker) hasGroupAccessTo(
	ctx context.Context,
	u *domain.Usr,
	groupID domain.GroupID,
	reqAccs domain.GroupUsrAccess,
) (bool, error) {
	if u.Role == domain.SuperUsrRole {
		return true, nil
	}

	accs, err := ac.gur.GetGroupAccess(ctx, u.ID, groupID)
	if err != nil {
		if errors.Is(err, domain.ErrNoAccess) {
			return false, nil
		}

		return false, err
	}

	return accs.IsEquivalent(reqAccs), nil
}
