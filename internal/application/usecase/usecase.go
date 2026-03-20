// Package usecase provides the usecase layer of the application.
// It contains the application's business logic and orchestrates the interactions between the domain and infrastructure layers.
package usecase

import (
	"context"

	"ownned/internal/domain"
)

// resolveNodeAccess checks if the user has access to the given node.
func resolveNodeAccess(ctx context.Context, gur domain.GroupUsrRepository, usr *domain.Usr, node *domain.Node) (domain.GroupUsrAccess, error) {
	if usr.Role == domain.SuperUsrRole {
		return domain.GroupWriteAccess, nil
	}

	accss, err := gur.GetNodeAccess(ctx, usr.ID, node.ID)
	if err != nil {
		return domain.GroupNoneAccess, err
	}

	return accss, nil
}
