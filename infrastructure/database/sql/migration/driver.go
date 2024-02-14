package database

import (
	"database/sql"
	"embed"
	"ffxvi-bard/config"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"log"
	"net/http"

	"ffxvi-bard/port/contract"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed query/*.sql
var migrations embed.FS

// Convert embedded files to http.FileSystem
var httpFS http.FileSystem = http.FS(migrations)

type migrationDriver struct {
	database string
	path     string
}

// NewMigrationDriver creates a new instance of the migration driver
func NewMigrationDriver(config *config.DatabaseConfig) contract.SqlMigrationDriverInterface {
	return &migrationDriver{database: config.Database, path: config.Path}
}

func (d *migrationDriver) connection() (*sql.DB, error) {
	db, err := sql.Open(d.database, d.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	return db, nil
}

func (d *migrationDriver) Execute(commandType string) {
	db, err := d.connection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatalf("Failed to create sqlite instance: %v", err)
	}

	sourceDriver, err := httpfs.New(http.FS(migrations), "query")
	if err != nil {
		log.Fatalf("Failed to create source driver: %v", err)
	}

	m, err := migrate.NewWithInstance("httpfs", sourceDriver, "sqlite3", driver)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if commandType == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply up migrations: %v", err)
		}
	} else if commandType == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply down migrations: %v", err)
		}
	} else {
		log.Fatalf("Unsupported command type: %s", commandType)
	}
}
