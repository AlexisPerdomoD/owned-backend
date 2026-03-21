package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"ownned/internal/domain"
	"ownned/internal/infrastructure/config"
	"ownned/internal/infrastructure/db/pg"

	"github.com/google/uuid"
)

func main() {
	cfg := config.LoadEnvConfig()
	db, err := pg.NewDB(
		cfg.PgDB,
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgUser,
		cfg.PgPassword,
		cfg.PgSsl,
	)
	if err != nil {
		panic(err)
	}

	if err := pg.MigrateUp(db.DB); err != nil {
		panic(err)
	}

	rootusrID, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	rootusrGroupID, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	rootsharedID, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	rootsharedGroupID, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	rootusr := domain.Node{
		ID:          rootusrID,
		Name:        domain.NodePathUsrRoot.String(),
		Description: "Root folder to contain all users root folders.",
		Type:        domain.FolderNodeType,
		Path:        domain.NodePathUsrRoot,
	}

	rootusrGroup := domain.Group{
		ID:          rootusrGroupID,
		Name:        "Root group for usrs root folder",
		Description: "Root group to contain all users root groups.",
	}

	rootusrGroupNode := domain.UpsertGroupNode{
		GroupID: rootusrGroupID,
		NodeID:  rootusrID,
	}

	rootshared := domain.Node{
		ID:          rootsharedID,
		Name:        domain.NodePathSharedRoot.String(),
		Description: "Root folder to contain all shared folders.",
		Type:        domain.FolderNodeType,
		Path:        domain.NodePathSharedRoot,
	}

	rootsharedGroup := domain.Group{
		ID:          rootsharedGroupID,
		Name:        "Root group shared root folder",
		Description: "Root group to contain all shared root groups.",
	}

	rootsharedGroupNode := domain.UpsertGroupNode{
		GroupID: rootsharedGroupID,
		NodeID:  rootsharedID,
	}

	uowFactory := pg.NewUnitOfWorkFactory(db, slog.Default(), time.Second*30)

	ctx := context.Background()
	err = uowFactory.Do(ctx, func(tx domain.UnitOfWork) error {
		if err := tx.NodeRepository().Create(tx.Ctx(), &rootusr); err != nil {
			return err
		}

		if err := tx.GroupRepository().Create(tx.Ctx(), &rootusrGroup); err != nil {
			return err
		}

		if err := tx.GroupNodeRepository().Upsert(tx.Ctx(), &rootusrGroupNode); err != nil {
			return err
		}

		if err := tx.NodeRepository().Create(tx.Ctx(), &rootshared); err != nil {
			return err
		}

		if err := tx.GroupRepository().Create(tx.Ctx(), &rootsharedGroup); err != nil {
			return err
		}

		if err := tx.GroupNodeRepository().Upsert(tx.Ctx(), &rootsharedGroupNode); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Folders roots created successfully, sumary:\n")
	for _, node := range []*domain.Node{&rootusr, &rootshared} {
		fmt.Printf("\nID: %s\nName: %s\nDescription: %s\nPath: %s\nType: %s\nCreatedAt: %s\nUpdatedAt: %s\n",
			node.ID,
			node.Name,
			node.Description,
			node.Path,
			node.Type,
			node.CreatedAt,
			node.UpdatedAt,
		)
	}
}
