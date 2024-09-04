package main

import (
	"context"
	"fmt"
	"jps/internal/config"
	"jps/internal/handler"
	"jps/internal/server"
	storage "jps/internal/storage/postgresql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustReadConfig()

	log := setupLogger(cfg.Env)

	log.Info("App started", slog.Any("config", cfg))

	db, err := storage.NewSqlxDB(cfg)
	if err != nil {
		log.Error(fmt.Sprintf("db doesn't init: %s", err))
		panic("db doesn't init")
	}

	log.Info("App started", slog.Any("config", cfg))

	storage := storage.NewPostgreDB(db)

	handlers := handler.NewHandler(storage)

	serv := new(server.Server)

	go func() {
		if err := serv.Run(cfg, handlers.InitRouter()); err != nil {
			log.Error(fmt.Sprintf("cannot run server: %s", err))
			panic("cannot run server")
		}
	}()

	cnl := make(chan os.Signal, 1)
	signal.Notify(cnl, syscall.SIGTERM, syscall.SIGINT)
	<-cnl
	log.Info("App stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		log.Error(fmt.Sprintf("an error occurred while executing graceful shutdown: %s", err))
	}

	if err := db.Close(); err != nil {
		log.Error(fmt.Sprintf("an error occurred while executing close db: %s", err))
	}
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
