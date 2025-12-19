package usecase

import (
	"context"
	"log"
	"log/slog"
	"ownned/internal/application/dto"
	"ownned/internal/domain"
	"ownned/internal/pkg/error_pkg"
	"ownned/internal/pkg/helper_pkg"
)

type CreateUsrUseCase struct {
	ur     domain.UsrRepository
	nr     domain.NodeRepository
	uow    domain.UnitOfWorkFactory
	logger *slog.Logger
}

func (uc *CreateUsrUseCase) Execute(
	ctx context.Context,
	creatorID domain.UsrID,
	args dto.CreateUsrInputDto,
) (*domain.Usr, error) {
	usrRepository := uc.ur
	unitOfWorkFactory := uc.uow

	creator, err := usrRepository.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if creator.Role != domain.SuperUsrRole {
		return nil, error_pkg.ErrForbidden(map[string]string{"general": "usr does not have enought privileges to do this action"})
	}

	usr, err := usrRepository.GetByUsername(ctx, args.Username)
	if err != nil {
		return nil, err
	}

	if usr != nil {
		return nil, error_pkg.ErrConflic(map[string]string{"general": "username already in use"})
	}

	usr = args.GetData()
	tx := unitOfWorkFactory.New()
	if err = tx.Do(ctx, func(txCtx context.Context) error {

		if err := tx.UsrRepository().Create(txCtx, usr); err != nil {
			return err
		}

		if usr.Role != domain.SuperUsrRole && len(args.Access) > 0 {
			nodes, err := tx.NodeRepository().GetByIDs(txCtx, args.Access)
			if err != nil {
				return err
			}

			result := helper_pkg.MapConcurrent(nodes, func(n domain.Node) (any, error) {
				access := domain.ReadOnlyAccess

				if usr.Role == domain.NormalUsrRole {
					access = domain.WriteAccess
				}

				return nil, tx.NodeRepository().UpdateAccess(txCtx, usr.ID, n.ID, access)
			}, 1000)

			for _, v := range result {
				if v.IsOk() {
					continue
				}

				return v.Error
			}

		}

		return nil
	}); err != nil {
		uc.logger.Log(ctx, slog.LevelDebug, "transaction failed:", "err", err)
		return nil, err
	}

	return usr, nil
}

func NewCreateUsrUseCase(
	ur domain.UsrRepository,
	nr domain.NodeRepository,
	uow domain.UnitOfWorkFactory,
	mainLogger *slog.Logger,
) *CreateUsrUseCase {
	if ur == nil || nr == nil || uow == nil || mainLogger == nil {
		log.Panicln("NewCreateUsrUseCase received a nil reference as dependency")
	}

	logger := mainLogger.With("usecase", "CreateUsrUseCase")
	return &CreateUsrUseCase{ur, nr, uow, logger}
}
