package usecase

import (
	"context"
	"log"
	"log/slog"
	"ownned/internal/application/storage"
	"ownned/internal/domain"
	"ownned/internal/pkg/helper_pkg"
	"time"
)

type DeleteNodeUseCase struct {
	usrRepository  domain.UsrRepository
	nodeRepository domain.NodeRepository
	docRepository  domain.DocRepository
	storage        storage.Storage
	logger         *slog.Logger
}

func (uc *DeleteNodeUseCase) deleteDocsAsync(docs []domain.Doc) {

	logger := uc.logger.With("scope", "deleteDocsAsync")
	storage := uc.storage
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic at deleteDocsAsync", "recover", r)
		}
	}()

	deletions := helper_pkg.MapConcurrent(
		docs,
		func(doc domain.Doc) (*domain.Doc, error) { return &doc, storage.Remove(ctx, doc.ID) },
		20,
	)

	for _, deletion := range deletions {
		if deletion.IsOk() {
			continue
		}

		logger.Warn("failed to delete doc from storage",
			"docID", deletion.Value.ID,
			"docTitle", deletion.Value.Title,
			"nodeID", deletion.Value.NodeID,
			"err", deletion.Error,
		)
	}

}

func (uc *DeleteNodeUseCase) Execute(ctx context.Context, usrID domain.UsrID, nodeID domain.NodeID) error {
	usr, err := uc.usrRepository.GetByID(ctx, usrID)
	if err != nil {
		return err
	}

	if usr == nil {
		return domain.ErrUnauthenticated(nil)
	}

	if usr.Role != domain.LimitedUsrRole {
		return domain.ErrForbidden(map[string]string{"general": "usr does not have access to do this action"})
	}

	node, err := uc.nodeRepository.GetByID(ctx, nodeID)
	if err != nil {
		return err
	}

	if node == nil {
		return domain.ErrNotFound(map[string]string{"nodeID": "entity was not found by id=" + nodeID})
	}

	if usr.Role != domain.SuperUsrRole {
		access, err := uc.nodeRepository.GetAccess(ctx, usr.ID, node.ID)
		if err != nil {
			return err
		}

		if access != domain.WriteAccess {
			return domain.ErrForbidden(
				map[string]string{
					"general": "Not enought privileges to delete node",
					"nodeID":  "usr does not have permission to do this action over this resource nodeID=" + node.ID,
				},
			)
		}
	}

	docs, err := uc.docRepository.GetByNodeID(ctx, node.ID)
	if err != nil {
		return err
	}

	// IS EXPECTED TO NODE DELETIONS TO DELETE EVERY CHILDREN AND ALSO DOCS ASSOCIATE ON CASCADE
	if err := uc.nodeRepository.Delete(ctx, node.ID); err != nil {
		return err
	}

	if len(docs) > 0 {
		go uc.deleteDocsAsync(docs)
	}

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
