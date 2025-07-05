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

	server := grpc.NewDemesServer(db)
	server.StartServer(cfg.UserServicePort)
	log.Println("GRPC server started", "port", cfg.UserServicePort)
	defer server.StopServer()

	select{}
}