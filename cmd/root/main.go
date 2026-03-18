package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"ownned/internal/application/dto"
	"ownned/internal/application/usecase"
	"ownned/internal/domain"
	"ownned/internal/infrastructure/config"
	"ownned/internal/infrastructure/db/pg"
	"ownned/internal/infrastructure/serv"
)

func main() {
	cfg := config.LoadEnvConfig()
	// SERVICES

	usrname := flag.String("usrname", "", "Unique username (email) of user")
	pwd := flag.String("pwd", "", "Password of the new user")

	flag.Parse()

	if *usrname == "" {
		fmt.Fprintln(os.Stderr, "error: -usrname is required")
		flag.Usage()
		os.Exit(1)
	}

	if *pwd == "" {
		fmt.Fprintln(os.Stderr, "error: -pwd is required")
		flag.Usage()
		os.Exit(1)
	}

	usrDTO := dto.CreateUsrDTO{
		Firstname: "admin",
		Lastname:  "admin",
		Username:  *usrname,
		Pwd:       *pwd,
		Role:      domain.SuperUsrRole,
		Access:    make([]dto.CreateAccessDTO, 0),
	}

	if err := usrDTO.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %e", err)
		flag.Usage()
		os.Exit(1)
	}

	pwdHasher := serv.NewPwdHasherArgon2(
		cfg.PwdTime,
		cfg.PwdMemKiB,
		cfg.PwdThreads,
		cfg.PwdHashLen,
		cfg.PwdSaltLen,
	)

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
	ur := pg.NewUsrRepository(db)
	uow := pg.NewUnitOfWorkFactory(db, slog.Default(), time.Second*30)
	createUsr := usecase.NewCreateUsrUseCase(ur, uow, pwdHasher, slog.Default())
	ctx := context.Background()
	usr, err := createUsr.Execute(ctx, usrDTO)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s creado exitosamente", usr.Username)
}
