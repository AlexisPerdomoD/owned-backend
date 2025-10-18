package usecase

import (
	"context"
	"log"
	"log/slog"
	"ownned/internal/application/storage"
	"ownned/internal/domain"
	"ownned/internal/pkg/error_pkg"
	"ownned/internal/pkg/helper_pkg"
)

type DeleteNodeUseCase struct {
	usrRepository  domain.UsrRepository
	nodeRepository domain.NodeRepository
	docRepository  domain.DocRepository
	storage        storage.Storage
	logger         *slog.Logger
}

func (uc *DeleteNodeUseCase) Execute(ctx context.Context, usrID domain.UsrID, nodeID domain.NodeID) error {
	usr, err := uc.usrRepository.GetByID(ctx, usrID)
	if err != nil {
		return err
	}

	if usr == nil {
		return error_pkg.ErrUnauthenticated(nil)
	}

	if usr.Role != domain.LimitedUsrRole {
		return error_pkg.ErrForbidden(map[string]string{"general": "usr does not have access to do this action"})
	}

	node, err := uc.nodeRepository.GetByID(ctx, nodeID)
	if err != nil {
		return err
	}

	if node == nil {
		return error_pkg.ErrNotFound(map[string]string{"nodeID": "entity was not found by id=" + nodeID})
	}

	if usr.Role != domain.SuperUsrRole {
		access, err := uc.nodeRepository.GetAccess(ctx, usr.ID, node.ID)
		if err != nil {
			return err
		}

		if access != domain.WriteAccess {
			return error_pkg.ErrForbidden(
				map[string]string{
					"general": "Not enought privileges to delete node",
					"nodeID":  "usr does not have permission to do this action over this resource nodeID=" + node.ID,
				},
			)
		}
	}
	// this may be a private handler
	if node.Type == domain.FileNodeType {
		docs, err := uc.docRepository.GetByNodeID(ctx, node.ID)
		if err != nil {
			return err
		}

		if err := uc.nodeRepository.Delete(ctx, node.ID); err != nil {
			return err
		}

		if len(docs) == 0 {
			return nil
		}

		go func() {
			deletions := helper_pkg.MapConcurrent(docs, func(doc domain.Doc) (*domain.Doc, error) {
				return &doc, uc.storage.Remove(doc.ID)
			}, 10)

			for _, deletion := range deletions {

				if deletion.IsOk() {
					continue
				}

				uc.logger.Warn("failed to delete doc from storage",
					"docID", deletion.Value.ID,
					"docTitle", deletion.Value.Title,
					"nodeID", node.ID,
					"err", deletion.Error,
				)
			}
		}()

	}

	// this may be a private handler
	// folder node
	// todo
	return nil
}

func NewDeleteNodeUseCase(
	ur domain.UsrRepository,
	nr domain.NodeRepository,
	dr domain.DocRepository,
	storage storage.Storage,
	mainLogger *slog.Logger,
) *DeleteNodeUseCase {
	if ur == nil || nr == nil || dr == nil || storage == nil || mainLogger == nil {
		log.Panicln("DeleteNodeUseCase some dependencies provided were nil")
	}

	logger := mainLogger.With("usecase", "DeleteNodeUseCase")

	return &DeleteNodeUseCase{ur, nr, dr, storage, logger}

}
