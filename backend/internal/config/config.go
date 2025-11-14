package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port        string
	Environment string
}

type DatabaseConfig struct {
	URL string
}

type StorageConfig struct {
	Type          string
	MinIOEndpoint string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket   string
	MinIOUseSSL   bool
}

type AuthConfig struct {
	JWTSecret string
}

func Load() (*Config, error) {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Try to read config file, but don't fail if it doesn't exist
	_ = viper.ReadInConfig()

	cfg := &Config{
		Server: ServerConfig{
			Port:        viper.GetString("PORT"),
			Environment: viper.GetString("ENVIRONMENT"),
		},
		Database: DatabaseConfig{
			URL: viper.GetString("DATABASE_URL"),
		},
		Storage: StorageConfig{
			Type:          viper.GetString("STORAGE_TYPE"),
			MinIOEndpoint: viper.GetString("MINIO_ENDPOINT"),
			MinIOAccessKey: viper.GetString("MINIO_ACCESS_KEY"),
			MinIOSecretKey: viper.GetString("MINIO_SECRET_KEY"),
			MinIOBucket:   viper.GetString("MINIO_BUCKET"),
			MinIOUseSSL:   viper.GetBool("MINIO_USE_SSL"),
		},
		Auth: AuthConfig{
			JWTSecret: viper.GetString("JWT_SECRET"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("DATABASE_URL", "postgres://postgres:password@localhost:5432/appdb?sslmode=disable")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("STORAGE_TYPE", "minio")
	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin")
	viper.SetDefault("MINIO_BUCKET", "media")
	viper.SetDefault("MINIO_USE_SSL", false)
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("PORT environment variable is required")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}
