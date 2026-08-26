package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"ownned/internal/application/usecase"
	"ownned/internal/infrastructure/config"
	"ownned/internal/infrastructure/db/pg"
	"ownned/internal/infrastructure/srv"
	"ownned/internal/infrastructure/transport/http/handler"
	"ownned/internal/infrastructure/transport/http/middleware"
	"ownned/pkg/http_log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// START POINT BABY
func main() {
	// =========================================================================
	// CONFIG
	// =========================================================================

	cfg := config.LoadEnvConfig()

	// =========================================================================
	// DB
	// =========================================================================

	db, err := pg.
		NewDB(
			cfg.PgDB,
			cfg.PgHost,
			cfg.PgPort,
			cfg.PgUser,
			cfg.PgPassword,
			cfg.PgSsl)
	if err != nil {
		panic(err)
	}
	if err := pg.MigrateUp(db.DB); err != nil {
		panic(err)
	}
	// =========================================================================
	// SERVICES
	// =========================================================================

	lg := slog.New(slog.
		NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	jwtManager := srv.
		NewJWTManagerST(
			[]byte(cfg.SessionSecret),
			time.Hour,
			"ownned")
	pwdHasher := srv.
		NewPwdHasherArgon2(
			cfg.PwdTime,
			cfg.PwdMemKiB,
			cfg.PwdThreads,
			cfg.PwdHashLen,
			cfg.PwdSaltLen)
	storage := srv.
		NewStorageManagerFS(cfg.LocalStorageDir)

	// =========================================================================
	// REPOSITORIES
	// =========================================================================

	usrRepository := pg.
		NewUsrRepository(db)
	usrPwdRepository := pg.
		NewUsrPwdRepository(db)
	nodeRepository := pg.
		NewNodeRepository(db)
	nodeCommentRepository := pg.
		NewNodeCommentRepository(db)
	groupRepository := pg.
		NewGroupRepository(db)
	groupNodeRepository := pg.
		NewGroupNodeRepository(db)
	groupUsrRepository := pg.
		NewGroupUsrRepository(db)
	docRepository := pg.
		NewDocRepository(db)
	unitOfWorkFactory := pg.
		NewUnitOfWorkFactory(db, lg, time.Second*30)

	accessChecker := usecase.
		NewAccessChecker(groupUsrRepository)

	// =========================================================================
	// MIDLEWARES
	// =========================================================================

	authM := middleware.
		NewAuthMiddleware(jwtManager)

	// =========================================================================
	// GROUPS ROUTES
	// =========================================================================

	getGroup := usecase.
		NewGetGroupUseCase(
			usrRepository,
			nodeRepository,
			groupRepository,
			accessChecker)
	paginateGroup := usecase.
		NewPaginateGroupUseCase(groupRepository)
	updateGroup := usecase.
		// TODO: add swagger doc here
		NewUpdateGroupUseCase(groupRepository)
	createGroup := usecase.
		NewCreateGroupUseCase(unitOfWorkFactory)
	deleteGroup := usecase.
		NewDeleteGroupUseCase(
			groupRepository,
			accessChecker)
	createGroupNode := usecase.
		NewCreateGroupNodeUseCase(
			groupRepository,
			groupNodeRepository,
			nodeRepository,
			accessChecker)

	deleteGroupNode := usecase.
		NewDeleteGroupNodeUseCase(
			groupRepository,
			groupNodeRepository,
			accessChecker)
	upsertGroupUsr := usecase.
		NewUpsertGroupUsrUseCase(
			usrRepository,
			accessChecker)
	deleteGroupUsr := usecase.
		NewDeleteGroupUsrUseCase(accessChecker)

	// ROUTES
	groupH := handler.
		NewGroupHandler(
			getGroup,
			paginateGroup,
			createGroup,
			updateGroup,
			deleteGroup,
			createGroupNode,
			deleteGroupNode,
			upsertGroupUsr,
			deleteGroupUsr)
	groupR := chi.NewRouter()

	groupR.Post("/", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(groupH.CreateGroupHandler))

	groupR.Get("/{groupID}", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(groupH.GetGroupHandler))

	groupR.Patch("/{groupID}", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(groupH.UpdateGroupHandler))

	groupR.Get("/paginate", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(groupH.PaginateGroupHandler))

	groupR.Delete("/{groupID}", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(groupH.DeleteGroupHandler))

	groupR.Post("/{groupID}/nodes", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(groupH.CreateGroupNodeHandler))

	groupR.Delete("/{groupID}/nodes/{nodeID}", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(groupH.DeleteGroupNodeHandler))

	groupR.Post("/{groupID}/users", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(groupH.UpsertGroupUsrHandler))

	groupR.Delete("/{groupID}/users/{usrID}", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(groupH.DeleteGroupUsrHandler))

	// =========================================================================
	// USERS ROUTES
	// =========================================================================

	createUsr := usecase.
		NewCreateUsrUseCase(
			usrRepository,
			unitOfWorkFactory,
			pwdHasher,
			lg)
	getUsr := usecase.
		NewGetUsrUseCase(usrRepository)
	getMe := usecase.
		NewGetMeUseCase(usrRepository)
	paginateUsr := usecase.
		NewPaginateUsrUseCase(usrRepository)
	loginUsr := usecase.
		NewLoginUsrUseCase(
			usrRepository,
			usrPwdRepository,
			pwdHasher,
			jwtManager)

	// ROUTES
	secure := cfg.Mode != "local"
	sameSite := http.SameSiteStrictMode
	if cfg.Mode == "local" {
		sameSite = http.SameSiteLaxMode
	}
	usrH := handler.
		NewUsrHandler(
			loginUsr,
			createUsr,
			getMe,
			getUsr,
			paginateUsr,
			handler.UsrHandlerConfig{
				Secure:   secure,
				SameSite: sameSite,
			})
	usrR := chi.NewRouter()
	usrR.Get("/me", authM.
		IsAuthenticated(usrH.GetMeHandler))
	usrR.Get("/{usrID}", authM.
		IsAuthenticated(usrH.GetUsrHandler))
	usrR.Get("/paginate", authM.
		IsAuthenticated(usrH.PaginateUsrHandler))
	usrR.Post("/", authM.
		IsSuperUsr(usrH.CreateUsrHandler))
	usrR.Post("/login", usrH.LoginUsrHandler)
	usrR.Delete("/logout", usrH.LogoutUsrHandler)

	// =========================================================================
	// NODES ROUTES
	// =========================================================================

	getRoot := usecase.
		NewGetRootNodesUseCase(
			nodeRepository,
			groupRepository,
			lg)
	createFolder := usecase.
		NewCreateFolderUseCase(
			nodeRepository,
			accessChecker)
	getNode := usecase.
		NewGetNodeByIDUseCase(
			nodeRepository,
			docRepository,
			lg,
			accessChecker)

	// NODES
	nodeH := handler.
		NewNodeHandler(
			getRoot,
			createFolder,
			getNode)
	nodeR := chi.NewRouter()
	nodeR.Get("/", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(nodeH.GetRootHandler))
	nodeR.Post("/", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(nodeH.CreateFolderHandler))
	nodeR.Get("/{nodeID}", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(nodeH.GetNodeHandler))

	// =========================================================================
	// NODE COMMENTS
	// =========================================================================

	getNodeComments := usecase.
		NewGetNodeCommentsUseCase(
			nodeRepository,
			nodeCommentRepository,
			accessChecker)
	createNodeComment := usecase.
		NewCreateNodeCommentUseCase(
			nodeRepository,
			nodeCommentRepository,
			lg,
			accessChecker)
	updateNodeComment := usecase.
		NewUpdateNodeCommentUseCase(
			nodeCommentRepository)
	deleteNodeComment := usecase.
		NewDeleteNodeCommentUseCase(
			nodeCommentRepository)
	// ROUTES
	nodeCommentH := handler.
		NewNodeCommentHandler(
			getNodeComments,
			createNodeComment,
			updateNodeComment,
			deleteNodeComment)
	nodeCommentR := chi.NewRouter()
	nodeR.Get("/{nodeID}/comments", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(nodeCommentH.GetNodeCommentsHandler))
	nodeR.Post("/{nodeID}/comments", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(nodeCommentH.CreateNodeCommentHandler))
	nodeCommentR.Patch("/{nodeCommentID}", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(nodeCommentH.UpdateNodeCommentHandler))
	nodeCommentR.Delete("/{nodeCommentID}", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(nodeCommentH.DeleteNodeCommentHandler))

	// =========================================================================
	// DOCS ROUTES
	// =========================================================================

	createDoc := usecase.
		NewCreateDocUseCase(
			docRepository,
			nodeRepository,
			unitOfWorkFactory,
			storage,
			lg,

			accessChecker)
	deleteDoc := usecase.
		NewDeleteDocUseCase(
			storage,
			docRepository,
			nodeRepository,
			lg,
			accessChecker)
	downloadDoc := usecase.
		NewDownloadDocUseCase(
			nodeRepository,
			docRepository,
			storage,
			accessChecker)

	// ROUTES
	docH := handler.
		NewDocHandler(
			createDoc,
			deleteDoc,
			downloadDoc)
	docR := chi.NewRouter()
	docR.Post("/", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(docH.CreateDocHandler))
	docR.Delete("/{docID}", authM.
		// TODO: Add swagger doc here
		IsAuthenticated(docH.DeleteDocHandler))
	docR.Get("/{docID}/download", authM.
		// TODO: add swagger doc here
		IsAuthenticated(docH.DownloadDocHandler))

	// =========================================================================
	// SERVER START POINT
	// =========================================================================

	r := chi.NewRouter()
	if cfg.Mode == "local" {
		// Basic CORS
		// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"http://localhost:5173"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders: []string{"Accept", "Content-Type"},
			// ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	// TODO: Add swagger doc here
	r.Mount("/api/v1/groups", groupR)
	// TODO: Add swagger doc here
	r.Mount("/api/v1/usrs", usrR)
	// TODO: Add swagger doc here
	r.Mount("/api/v1/nodes", nodeR)
	// TODO: Add swagger doc here
	r.Mount("/api/v1/comments", nodeCommentR)
	// TODO: Add swagger doc here
	r.Mount("/api/v1/docs", docR)

	// =========================================================================
	// SERVE WEB APP
	// =========================================================================
	if cfg.ServeWebApp {
		// Static assets
		// TODO: make this path a config option maybe
		r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/dist/assets"))))

		// SPA fallback
		// TODO: make this path a config option maybe
		r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/dist/index.html")
		}))
	}

	http_log.ChiRouterLog(r)
	lg.Info("server starting at:",
		"Mode", cfg.Mode,
		"port", cfg.Port,
		"ServeWebApp", cfg.ServeWebApp)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}
