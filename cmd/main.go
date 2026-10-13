package main

import (
	"log"
	"log/slog"

	grpc "github.com/Lykeion-org/demes/internal/grpc"
	models "github.com/Lykeion-org/demes/internal/models"
	"github.com/Lykeion-org/go-shared/pkg/config"
	"github.com/Lykeion-org/go-shared/pkg/db"
	utils "github.com/Lykeion-org/go-shared/pkg/helpers"
)

func main(){
	cfg, err := config.LoadConfigFile[Config]("../config.yaml")
	if err != nil {
		panic(err)
	}

	utils.InitializeStandardLogger("debug")

	connString := db.CreateConnectionString(cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName, cfg.DatabaseUserName, cfg.DatabasePassword)
	db, err:= db.ConnectAndMigrate(connString, &models.LykeionUser{})
	if err!=nil{
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}

	slog.Info("Connected to database and migrated models")

	// One account per email, compared case-insensitively. This fails if the
	// table already holds duplicates; demes still starts and blocks new ones,
	// but the duplicates have to be cleaned up before the index can exist.
	err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_lykeion_users_email_address_lower ON lykeion_users (lower(email_address))`).Error
	if err != nil {
		slog.Error("Could not create unique email index, remove duplicate emails from lykeion_users and restart", "error", err)
	}

	server := grpc.NewDemesServer(db)
	server.StartServer(cfg.UserServicePort)
	log.Println("GRPC server started", "port", cfg.UserServicePort)
	defer server.StopServer()

	select{}
}