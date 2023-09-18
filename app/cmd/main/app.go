package main

import (
	"log"
	"small/internal/config"
	"small/internal/server"
	"small/pkg/type/logger"
)

func main() {
	appLogger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	s := server.New(*appLogger, cfg)
	appLogger.Fatal(s.Run())

}
