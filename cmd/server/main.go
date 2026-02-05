package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"ownned/internal/application/usecase"
	"ownned/internal/domain"
	"ownned/internal/infrastructure/auth"
	"ownned/internal/infrastructure/db/pg"
	"ownned/internal/infrastructure/transport/http/handler"
	"ownned/internal/infrastructure/transport/http/middleware"
	"strings"

	"github.com/go-chi/chi/v5"
)

// start point baby
func main() {
	//services
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	jwtService := auth.NewJWTService()
	// DB
	var usrRepository domain.UsrRepository = pg.NewUsrRepository()
	var nodeRepository domain.NodeRepository = pg.NewNodeRepository()
	// var docRepository domain.DocRepository = repo.NewDocRepository()
	var unitOfWorkFactory domain.UnitOfWorkFactory = pg.NewUnitOfWorkFactory()

	createUsr := usecase.NewCreateUsrUseCase(usrRepository, nodeRepository, unitOfWorkFactory, l)
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
	log_routes(r, l)

	PORT := 9090
	l.Info("server starting at:", "port", PORT)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
}

func log_routes(r chi.Router, l *slog.Logger) {

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
