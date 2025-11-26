package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg := config{
		addr: ":4070",
	}

	api := application{
		config: cfg,
	}
	if err := api.run(api.mount()); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}

}
