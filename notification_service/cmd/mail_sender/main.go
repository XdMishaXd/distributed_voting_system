package main

import (
	"os"

	"notification_service/internal/config"
	"notification_service/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustSetup(cfg.Env)

	log.Info("starting notification_service")

	if err := run(cfg, log); err != nil {
		log.Error("app stopped with error", "error", err)
		os.Exit(1)
	}
}
