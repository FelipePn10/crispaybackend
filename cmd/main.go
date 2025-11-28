package main

import (
	"log/slog"
	"os"

	"github.com/FelipePn10/crispaybackend/config"
	"github.com/FelipePn10/crispaybackend/internal/database"
	"github.com/FelipePn10/crispaybackend/internal/email"
	"github.com/FelipePn10/crispaybackend/internal/email/service"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := database.NewDB(cfg)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	emailConfig := email.LoadConfigFromEnv()
	emailService := service.NewEmailService(service.EmailConfig(emailConfig))

	api := application{
		config:       cfg,
		logger:       logger,
		db:           db,
		emailService: emailService,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}
