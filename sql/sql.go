package sql

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

//go:embed schema/*.sql
var embedMigrations embed.FS

func CreateDB(dbPath string) (*sql.DB, error) {
	goose.SetBaseFS(embedMigrations)

	drv, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := goose.SetDialect("sqlite"); err != nil {
		return nil, err
	}

	// Configure goose logger based on debug flag
	if !viper.GetBool("debug") {
		goose.SetLogger(goose.NopLogger())
	}

	if err := goose.Up(drv, "schema"); err != nil {
		return nil, err
	}

	return drv, nil
}
