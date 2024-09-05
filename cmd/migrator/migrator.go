package main

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrationPath = "./schema"
)

func main() {
	migrat, err := migrate.New("file://"+migrationPath, "postgresql://user:password@localhost:5432/jpsdb")
	if err != nil {
		slog.Error(fmt.Sprintf("Error migrate: %s", err))
		return
	}

	if err = migrat.Up(); err != nil {
		slog.Error(fmt.Sprintf("Error migrate: %s", err))
		return
	}
}
