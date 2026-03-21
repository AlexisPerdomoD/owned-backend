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

// hasAccessTo checks if a user has access to a node based on the user's role and the access of the node to the user
func (ac *accessChecker) hasAccessTo(
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
