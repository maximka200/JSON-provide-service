package main

import (
	"jps/internal/config"
	"jps/internal/handler"
	storage "jps/internal/storage/postgresql"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustReadConfig()

	log := setupLogger(cfg.Env)

	log.Info("App started", slog.Any("config", cfg))

	db, err := storage.NewSqlxDB(cfg)
	if err != nil {
		log.Error("db doesn't init")
		panic("db doesn't init")
	}

	log.Info("App started", slog.Any("config", cfg))

	storage := storage.NewPostgreDB(db)

	handlers := handler.NewHandler(storage)

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
