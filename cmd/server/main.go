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

func logRoutes(r chi.Router) {
	methodColors := map[string]string{
		"GET":    "\033[32m", // verde
		"POST":   "\033[34m", // azul
		"PUT":    "\033[33m", // amarillo
		"PATCH":  "\033[33m", // amarillo
		"DELETE": "\033[31m", // rojo
	}

	const (
		reset = "\033[0m"
		bold  = "\033[1m"
	)

	type routeEntry struct {
		methods []string
		path    string
	}
	grouped := make(map[string]*routeEntry)
	order := []string{}

	_ = chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if _, exists := grouped[route]; !exists {
			grouped[route] = &routeEntry{path: route}
			order = append(order, route)
		}
		grouped[route].methods = append(grouped[route].methods, method)
		return nil
	})

	fmt.Println(bold + "registered routes:" + reset)
	for idx, path := range order {
		entry := grouped[path]
		coloredMethods := make([]string, len(entry.methods))
		for i, m := range entry.methods {
			color, ok := methodColors[m]
			if !ok {
				color = "\033[37m"
			}
			coloredMethods[i] = color + bold + m + reset
		}
		fmt.Printf("  %2d. %s%-45s%s %s\n",
			idx+1,
			bold, entry.path, reset,
			strings.Join(coloredMethods, " "),
		)
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
	l := slog.New(slog.
		NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

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
	nodeCommentRepository := pg.NewNodeCommentRepository(db)
	groupRepository := pg.NewGroupRepository(db)
	groupUsrRepository := pg.NewGroupUsrRepository(db)
	docRepository := pg.NewDocRepository(db)
	unitOfWorkFactory := pg.NewUnitOfWorkFactory(db, l, time.Second*30)

	// GROUPS
	getGroup := usecase.
		NewGetGroupUseCase(
			usrRepository,
			nodeRepository,
			groupRepository,
			groupUsrRepository)
	paginateGroup := usecase.
		NewPaginateGroupUseCase(
			usrRepository,
			groupRepository)
	createGroup := usecase.
		NewCreateGroupUseCase(
			usrRepository,
			unitOfWorkFactory)
	deleteGroup := usecase.
		NewDeleteGroupUseCase(
			usrRepository,
			groupRepository,
			groupUsrRepository)

	// GROUPS ROUTES
	groupHandler := handler.
		NewGroupHandler(
			getGroup,
			paginateGroup,
			createGroup,
			deleteGroup)
	groupR := chi.NewRouter()

	groupR.Post("/", authmiddleware.
		IsAuthenticated(groupHandler.CreateGroupHandler))

	groupR.Get("/{groupID}", authmiddleware.
		IsAuthenticated(groupHandler.GetGroupHandler))

	groupR.Get("/paginate", authmiddleware.
		IsAuthenticated(groupHandler.PaginateGroupHandler))

	groupR.Delete("/{groupID}", authmiddleware.
		IsAuthenticated(groupHandler.DeleteGroupHandler))

	// USERS
	createUsr := usecase.
		NewCreateUsrUseCase(
			usrRepository,
			unitOfWorkFactory,
			pwdHasher,
			l)
	getUsr := usecase.
		NewGetUsrUseCase(usrRepository)
	paginateUsr := usecase.
		NewPaginateUsrUseCase(usrRepository)
	loginUsr := usecase.
		NewLoginUsrUseCase(
			usrRepository,
			usrPwdRepository,
			pwdHasher,
			jwtService)
	// USR ROUTES
	usrHandler := handler.
		NewUsrHandler(
			loginUsr,
			createUsr,
			getUsr,
			paginateUsr,
			handler.UsrHandlerConfig{
				Secure:   cfg.Mode != "local",
				SameSite: http.SameSiteLaxMode,
			})
	usrR := chi.NewRouter()
	usrR.Get("/{usrID}", authmiddleware.
		IsAuthenticated(usrHandler.GetUsrHandler))
	usrR.Post("/", authmiddleware.
		IsSuperUsr(usrHandler.CreateUsrHandler))
	usrR.Get("/paginate", authmiddleware.
		IsAuthenticated(usrHandler.PaginateUsrHandler))
	usrR.Post("/login", usrHandler.LoginUsrHandler)
	usrR.Delete("/logout", usrHandler.LogoutUsrHandler)

	// NODES
	getRoot := usecase.
		NewGetRootNodesUseCase(
			nodeRepository,
			usrRepository,
			groupRepository,
			l)
	createFolder := usecase.
		NewCreateFolderUseCase(
			nodeRepository,
			usrRepository,
			groupUsrRepository)
	getNode := usecase.
		NewGetNodeByIDUseCase(
			usrRepository,
			nodeRepository,
			docRepository,
			groupUsrRepository,
			l)
	// NODES ROUTES
	nodeHandler := handler.
		NewNodeHandler(
			getRoot,
			createFolder,
			getNode)
	nodeR := chi.NewRouter()
	nodeR.Get("/", authmiddleware.
		IsAuthenticated(nodeHandler.GetRootHandler))
	nodeR.Post("/", authmiddleware.
		IsAuthenticated(nodeHandler.CreateFolderHandler))
	nodeR.Get("/{nodeID}", authmiddleware.
		IsAuthenticated(nodeHandler.GetNodeHandler))

	// NODE COMMENTS
	getNodeComments := usecase.
		NewGetNodeCommentsUseCase(
			usrRepository,
			nodeRepository,
			nodeCommentRepository,
			groupUsrRepository)
	createNodeComment := usecase.
		NewCreateNodeCommentUseCase(
			usrRepository,
			nodeRepository,
			nodeCommentRepository,
			groupUsrRepository,
			l)
	updateNodeComment := usecase.
		NewUpdateNodeCommentUseCase(
			usrRepository,
			nodeCommentRepository)
	deleteNodeComment := usecase.
		NewDeleteNodeCommentUseCase(
			usrRepository,
			nodeRepository,
			nodeCommentRepository,
			groupUsrRepository)
	// NODE COMMENTS ROUTES
	nodeCommentHandler := handler.
		NewNodeCommentHandler(
			getNodeComments,
			createNodeComment,
			updateNodeComment,
			deleteNodeComment)
	nodeCommentR := chi.NewRouter()
	nodeR.Get("/:nodeID/comments", authmiddleware.
		IsAuthenticated(nodeCommentHandler.GetNodeCommentsHandler))
	nodeR.Post("/:nodeID/comments", authmiddleware.
		IsAuthenticated(nodeCommentHandler.CreateNodeCommentHandler))
	nodeCommentR.Patch("/{nodeCommentID}", authmiddleware.
		IsAuthenticated(nodeCommentHandler.UpdateNodeCommentHandler))
	nodeCommentR.Delete("/{nodeCommentID}", authmiddleware.
		IsAuthenticated(nodeCommentHandler.DeleteNodeCommentHandler))

	// DOCS
	createDoc := usecase.
		NewCreateDocUseCase(
			usrRepository,
			docRepository,
			nodeRepository,
			groupUsrRepository,
			unitOfWorkFactory,
			storage,
			l)
	deleteDoc := usecase.
		NewDeleteDocUseCase(
			storage,
			docRepository,
			nodeRepository,
			usrRepository,
			groupUsrRepository,
			l)

	// DOCS ROUTES
	docHandler := handler.
		NewDocHandler(
			createDoc,
			deleteDoc)
	docR := chi.NewRouter()
	docR.Post("/", authmiddleware.
		IsAuthenticated(docHandler.CreateDocHandler))
	docR.Delete("/{docID}", authmiddleware.
		IsAuthenticated(docHandler.DeleteDocHandler))

	// SERVER ROUTES

	r := chi.NewRouter()
	r.Mount("/api/v1/groups", groupR)
	r.Mount("/api/v1/usrs", usrR)
	r.Mount("/api/v1/nodes", nodeR)
	r.Mount("/api/v1/comments", nodeCommentR)
	r.Mount("/api/v1/docs", docR)
	logRoutes(r)

	l.Info("server starting at:", "port", cfg.Port)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}
