package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/config"
	"ownned/internal/infrastructure/db/pg"
	"ownned/internal/infrastructure/serv"
	"ownned/internal/infrastructure/transport/http/handler"
	"ownned/internal/infrastructure/transport/http/middleware"
)

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

// start point baby
func main() {
	cfg := config.LoadEnvConfig()
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

	// SERVICES
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	jwtService := serv.NewJWTManagerST(
		[]byte(cfg.SessionSecret),
		time.Hour,
		"ownned",
	)

	pwdHasher := serv.NewPwdHasherArgon2(
		cfg.PwdTime,
		cfg.PwdMemKiB,
		cfg.PwdThreads,
		cfg.PwdHashLen,
		cfg.PwdSaltLen,
	)

	storage := serv.NewStorageManagerFS(cfg.LocalStorageDir)

	// MIDLEWARES
	authmiddleware := middleware.NewAuthMiddleware(jwtService)

	usrRepository := pg.NewUsrRepository(db)
	usrPwdRepository := pg.NewUsrPwdRepository(db)
	nodeRepository := pg.NewNodeRepository(db)
	groupRepository := pg.NewGroupRepository(db)
	groupUsrRepository := pg.NewGroupUsrRepository(db)
	docRepository := pg.NewDocRepository(db)
	unitOfWorkFactory := pg.NewUnitOfWorkFactory(db, l, time.Second*30)

	// USERS
	createUsr := usecase.NewCreateUsrUseCase(usrRepository, unitOfWorkFactory, pwdHasher, l)
	getUsr := usecase.NewGetUsrUseCase(usrRepository)
	loginUsr := usecase.NewLoginUsrUseCase(usrRepository, usrPwdRepository, pwdHasher, jwtService)
	// USR ROUTES
	usrHandler := handler.NewUsrHandler(loginUsr, createUsr, getUsr, handler.UsrHandlerConfig{
		Secure:   cfg.Mode != "local",
		SameSite: http.SameSiteLaxMode,
	})
	usrR := chi.NewRouter()
	usrR.Get("/{id}", authmiddleware.IsAuthenticated(usrHandler.GetUsrHandler))
	usrR.Post("/", authmiddleware.IsAuthenticated(usrHandler.CreateUsrHandler))
	usrR.Post("/login", usrHandler.LoginUsrHandler)

	// NODES
	getRoot := usecase.NewGetRootNodesUseCase(nodeRepository, usrRepository, groupRepository, l)
	createFolder := usecase.NewCreateFolderUseCase(nodeRepository, usrRepository, groupUsrRepository)
	getNode := usecase.NewGetNodeByIDUseCase(usrRepository, nodeRepository, docRepository, groupUsrRepository, l)
	// NODES ROUTES
	nodeHandler := handler.NewNodeHandler(getRoot, createFolder, getNode)
	nodeR := chi.NewRouter()
	nodeR.Get("/", authmiddleware.IsAuthenticated(nodeHandler.GetRootHandler))
	nodeR.Post("/", authmiddleware.IsAuthenticated(nodeHandler.CreateFolderHandler))
	nodeR.Get("/{nodeID}", authmiddleware.IsAuthenticated(nodeHandler.GetNodeHandler))

	// DOCS
	createDoc := usecase.NewCreateDocUseCase(usrRepository, docRepository, nodeRepository, groupUsrRepository, unitOfWorkFactory, storage, l)
	// DOCS ROUTES
	docHandler := handler.NewDocHandler(createDoc)
	docR := chi.NewRouter()
	docR.Post("/", authmiddleware.IsAuthenticated(docHandler.CreateDocHandler))

	// SERVER ROUTES
	r := chi.NewRouter()
	r.Mount("/api/v1/usrs", usrR)
	r.Mount("/api/v1/nodes", nodeR)
	r.Mount("/api/v1/docs", docR)
	logRoutes(r, l)

	l.Info("server starting at:", "port", cfg.Port)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}
