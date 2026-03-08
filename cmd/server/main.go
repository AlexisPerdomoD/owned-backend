package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"ownned/internal/application/auth"
	"ownned/internal/application/usecase"
	"ownned/internal/domain"
	"ownned/internal/infrastructure/config"
	"ownned/internal/infrastructure/db/pg"
	"ownned/internal/infrastructure/transport/http/handler"
	"ownned/internal/infrastructure/transport/http/middleware"

	"github.com/go-chi/chi/v5"
)

// start point baby
func main() {
	cfg := config.LoadEnvConfig()
	// services
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	var jwtService auth.JWTManager
	// DB
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

	usrRepository := pg.NewUsrRepository(db)
	var nodeRepository domain.NodeRepository = pg.NewNodeRepository(db)
	groupUsrRepository := pg.NewGroupUsrRepository(db)
	// var docRepository domain.DocRepository = repo.NewDocRepository()
	unitOfWorkFactory := pg.NewUnitOfWorkFactory(db, l, time.Second*30)

	createUsr := usecase.NewCreateUsrUseCase(usrRepository, nodeRepository, groupUsrRepository, unitOfWorkFactory, l)
	getUsr := usecase.NewGetUsrUseCase(usrRepository)
	usrHandler := handler.NewUsrHandler(createUsr, getUsr)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello Ownned"))
		if err != nil {
			l.Warn("some error happend", "err", err)
		}
	})

	// midlewares
	authmiddleware := middleware.NewAuthMiddleware(jwtService)

	// Usr Routes
	usrR := chi.NewRouter()
	usrR.Get("/{id}", authmiddleware.IsAuthenticated(usrHandler.GetUsrHandler))
	usrR.Post("/", authmiddleware.IsAuthenticated(usrHandler.CreateUsrHandler))

	r.Mount("/api/v1/usr", usrR)
	logRoutes(r, l)

	l.Info("server starting at:", "port", cfg.Port)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}

func logRoutes(r chi.Router, l *slog.Logger) {
	for idx, route := range r.Routes() {
		l.Info("registered route", "idx", idx+1, "path", route.Pattern)
		if route.SubRoutes == nil {
			continue
		}

		if len(route.SubRoutes.Routes()) == 0 {
			continue
		}

		for subIdx, subRoute := range route.SubRoutes.Routes() {
			l.Info("registered sub route",
				"idx", fmt.Sprintf("%d.%d", idx+1, subIdx+1),
				"path", fmt.Sprintf("%s%s",
					strings.TrimSuffix(
						route.Pattern,
						"/*",
					),
					subRoute.Pattern,
				),
			)
		}

	}
}
