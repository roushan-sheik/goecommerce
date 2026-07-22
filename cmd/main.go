package main

import (
	"goecommerce/internal/config"
	"goecommerce/internal/server"
)

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)
	server.Start(db, cfg)
}
