package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string            `yaml:"environment"`
	Server      ServerConfig      `yaml:"server"`
	Application ApplicationConfig `yaml:"application"`
	Databases   DatabasesConfig   `yaml:"databases"`
}

type ServerConfig struct {
	Host                   string `yaml:"host"`
	Port                   string `yaml:"port"`
	ShutdownTimeoutSeconds int    `yaml:"shutdown_timeout"`
}

func (s ServerConfig) ShutdownTimeout() time.Duration {
	return time.Duration(s.ShutdownTimeoutSeconds) * time.Second
}

type ApplicationConfig struct {
	Name             string `yaml:"name"`
	SeedAdminName    string `yaml:"seed_admin_name"`
	BackofficePRName string `yaml:"backoffice_pr_name"`
}

type DatabasesConfig struct {
	Postgres       string `yaml:"postgres"`
	MigrationsDSN  string `yaml:"migrations_dsn"`
	MigrationsPath string `yaml:"migrations_path"`
}

func Load() (*Config, error) {
	configPath := EnvOr("CONFIG_FILE_PATH", "config/config.yaml")

	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", configPath, err)
	}

	expanded := os.ExpandEnv(string(raw))

	cfg := &Config{}
	if err = yaml.Unmarshal([]byte(expanded), cfg); err != nil {
		return nil, fmt.Errorf("error parsing config.yaml: %w", err)
	}

	return cfg, nil
}
