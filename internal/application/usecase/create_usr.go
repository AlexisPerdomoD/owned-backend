package usecase

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"ownned/internal/application/dto"
	"ownned/internal/domain"
	"ownned/pkg/apperror"
	"ownned/pkg/col"

	"github.com/google/uuid"
)

type CreateUsrUseCase struct {
	ur     domain.UsrRepository
	nr     domain.NodeRepository
	gur    domain.GroupUsrRepository
	uow    domain.UnitOfWorkFactory
	logger *slog.Logger
}

func (uc *CreateUsrUseCase) Execute(
	ctx context.Context,
	creatorID domain.UsrID,
	args dto.CreateUsrDTO,
) (*domain.Usr, error) {
	creator, err := uc.ur.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if creator.Role != domain.SuperUsrRole {
		return nil, apperror.ErrForbidden(map[string]string{"error": "Usr does not have enought privileges to do this action"})
	}

	usr, err := uc.ur.GetByUsername(ctx, args.Username)
	if err != nil {
		return nil, err
	}

	if usr != nil {
		return nil, apperror.ErrConflic(map[string]string{"error": "Username already in use"})
	}

	usrID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	usrGroupID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	usrNodeRootID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	usr = &domain.Usr{
		ID:        usrID,
		Username:  args.Username,
		Role:      args.Role,
		Firstname: args.Firstname,
		Lastname:  args.Lastname,
	}

	usrNodeRoot := &domain.Node{
		ID:          usrNodeRootID,
		Name:        fmt.Sprintf("root_%s", args.Username),
		Description: "Root node for user " + args.Username,
		Type:        domain.FolderNodeType,
		Path:        domain.NodePathUsrRoot.NewChildPath(usr.ID),
	}

	usrRootGroup := &domain.Group{
		ID:          usrGroupID,
		Name:        fmt.Sprintf("group_%s", args.Username),
		Description: "Group for user " + args.Username,
	}

	nodeRootGroup := &domain.UpsertGroupNode{
		GroupID: usrRootGroup.ID,
		NodeID:  usrNodeRoot.ID,
	}

	usrGroups := col.Set[domain.UpsertGroupUsr]{}
	usrGroups.Add(domain.UpsertGroupUsr{
		GroupID: usrRootGroup.ID,
		UsrID:   usr.ID,
		Access:  domain.GroupWriteAccess,
	})

	if len(args.Access) > 0 {
		for _, v := range args.Access {
			usrGroups.Add(domain.UpsertGroupUsr{
				GroupID: v.GroupID,
				UsrID:   usr.ID,
				Access:  v.Access,
			})
		}
	}

	err = uc.uow.Do(ctx, func(tx domain.UnitOfWork) error {
		txCtx := tx.Ctx()
		if err := tx.UsrRepository().Create(txCtx, usr); err != nil {
			return err
		}

		if err := tx.NodeRepository().Create(txCtx, usrNodeRoot); err != nil {
			return err
		}

		if err := tx.GroupRepository().Create(txCtx, usrRootGroup); err != nil {
			return err
		}

		if err := tx.GroupNodeRepository().Upsert(txCtx, nodeRootGroup); err != nil {
			return err
		}

		if err := tx.GroupUsrRepository().UpsertAll(txCtx, usrGroups.Slice()); err != nil {
			return err
		}

		return nil
	})

	return usr, err
}

func NewCreateUsrUseCase(
	ur domain.UsrRepository,
	nr domain.NodeRepository,
	gur domain.GroupUsrRepository,
	uow domain.UnitOfWorkFactory,
	mainLogger *slog.Logger,
) *CreateUsrUseCase {
	if ur == nil || nr == nil || uow == nil || mainLogger == nil || gur == nil {
		log.Panicln("NewCreateUsrUseCase received a nil reference as dependency")
	}

	logger := mainLogger.With("usecase", "CreateUsrUseCase")
	return &CreateUsrUseCase{ur, nr, gur, uow, logger}
}
