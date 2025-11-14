package database

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    _ "github.com/go-sql-driver/mysql"
    _ "github.com/jackc/pgx/v5/stdlib"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func InitDB() (*gorm.DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://postgres:password@localhost:5432/appdb?sslmode=disable"
    }

    driver := os.Getenv("DATABASE_DRIVER")
    if driver == "" {
        driver = detectDriver(dsn)
    }

    migrationPath := os.Getenv("MIGRATIONS_PATH")
    if migrationPath == "" {
        wd, err := os.Getwd()
        if err != nil {
            return nil, fmt.Errorf("failed to get working directory: %w", err)
        }
        migrationPath = fmt.Sprintf("file://%s", filepath.Join(wd, "migrations"))
    } else if !strings.HasPrefix(migrationPath, "file://") {
        absPath, err := filepath.Abs(migrationPath)
        if err != nil {
            return nil, fmt.Errorf("failed to resolve migrations path: %w", err)
        }
        migrationPath = fmt.Sprintf("file://%s", absPath)
    }

    if err := RunMigrations(driver, dsn, migrationPath); err != nil {
        return nil, fmt.Errorf("failed to run migrations: %w", err)
    }

    var dialector gorm.Dialector
    switch strings.ToLower(driver) {
    case "postgres", "postgresql":
        dialector = postgres.Open(dsn)
    case "mysql":
        gormDSN := strings.TrimPrefix(strings.TrimPrefix(dsn, "mysql://"), "mysql+")
        dialector = mysql.Open(gormDSN)
    default:
        return nil, fmt.Errorf("unsupported database driver: %s", driver)
    }

    db, err := gorm.Open(dialector, &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get database instance: %w", err)
    }

    // Configure connection pool
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(25)

    return db, nil
}

func detectDriver(dsn string) string {
    dsnLower := strings.ToLower(dsn)
    switch {
    case strings.HasPrefix(dsnLower, "mysql://"), strings.Contains(dsnLower, "@tcp("), strings.Contains(dsnLower, "mysql+"):
        return "mysql"
    case strings.HasPrefix(dsnLower, "postgres://"), strings.HasPrefix(dsnLower, "postgresql://"), strings.HasPrefix(dsnLower, "pgx://"):
        return "postgres"
    default:
        return "postgres"
    }
}
