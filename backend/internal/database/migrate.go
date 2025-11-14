package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(driver string, databaseURL string, migrationsPath string) error {
	if migrationsPath == "" {
		migrationsPath = "file://migrations"
	}

	switch driver {
	case "postgres":
		m, err := migrate.New(migrationsPath, databaseURL)
		if err != nil {
			return fmt.Errorf("failed to create migration instance: %w", err)
		}
		defer m.Close()

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			return fmt.Errorf("failed to get migration version: %w", err)
		}

		log.Printf("Database migrated to version %d (dirty: %v)", version, dirty)
		return nil

	case "mysql":
		mysqlURL := databaseURL
		if !strings.HasPrefix(mysqlURL, "mysql://") {
			mysqlURL = "mysql://" + strings.TrimPrefix(strings.TrimPrefix(databaseURL, "mysql://"), "mysql+")
		}

		m, err := migrate.New(migrationsPath, mysqlURL)
		if err != nil {
			return fmt.Errorf("failed to create migration instance: %w", err)
		}
		defer m.Close()

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			return fmt.Errorf("failed to get migration version: %w", err)
		}

		log.Printf("Database migrated to version %d (dirty: %v)", version, dirty)
		return nil

	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}
}

func RollbackMigration(driver string, databaseURL string, migrationsPath string, steps int) error {
	if migrationsPath == "" {
		migrationsPath = "file://migrations"
	}

	switch driver {
	case "postgres":
		m, err := migrate.New(migrationsPath, databaseURL)
		if err != nil {
			return fmt.Errorf("failed to create migration instance: %w", err)
		}
		defer m.Close()

		if err := m.Steps(-steps); err != nil {
			return fmt.Errorf("failed to rollback migrations: %w", err)
		}

		log.Printf("Rolled back %d migration(s)", steps)
		return nil

	case "mysql":
		mysqlURL := databaseURL
		if !strings.HasPrefix(mysqlURL, "mysql://") {
			mysqlURL = "mysql://" + strings.TrimPrefix(strings.TrimPrefix(databaseURL, "mysql://"), "mysql+")
		}

		m, err := migrate.New(migrationsPath, mysqlURL)
		if err != nil {
			return fmt.Errorf("failed to create migration instance: %w", err)
		}
		defer m.Close()

		if err := m.Steps(-steps); err != nil {
			return fmt.Errorf("failed to rollback migrations: %w", err)
		}

		log.Printf("Rolled back %d migration(s)", steps)
		return nil

	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}
}
