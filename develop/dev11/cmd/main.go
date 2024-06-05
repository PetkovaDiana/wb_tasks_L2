package main

import (
	"main.go/internal/config"
	"main.go/internal/storage"
	"main.go/server"
)

func main() {

	cfg, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	routes := handler.Newhandler

	srv := server.NewServer(cfg.ServerConfig)
	if err = srv.Run(routes.InitRoutes()); err != nil {
		panic(err)
	}

	storage.NewEventStorage()
}
