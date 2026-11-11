package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string            `yaml:"environment"`
	Server      ServerConfig      `yaml:"server"`
	Application ApplicationConfig `yaml:"application"`
	Databases   DatabasesConfig   `yaml:"databases"`
	Security    SecurityConfig    `yaml:"security"`
}

type SecurityConfig struct {
	// EncryptionKey is a base64-encoded 32-byte key (AES-256) used to encrypt
	// secrets at rest (e.g. identity provider client secrets).
	EncryptionKey string `yaml:"encryption_key"`
}

type ServerConfig struct {
	Port                   string   `yaml:"port"`
	ShutdownTimeoutSeconds int      `yaml:"shutdown_timeout"`
	AllowedOrigins         []string `yaml:"allowed_origins"`
	// CookieOrigins are the origins (the backoffice frontend) whose requests
	// are authenticated with the access token cookie instead of a Bearer token.
	CookieOrigins []string `yaml:"cookie_origins"`
}

func (s ServerConfig) ShutdownTimeout() time.Duration {
	return time.Duration(s.ShutdownTimeoutSeconds) * time.Second
}

func (s *ServerConfig) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		Port                   string `yaml:"port"`
		ShutdownTimeoutSeconds int    `yaml:"shutdown_timeout"`
		AllowedOrigins         string `yaml:"allowed_origins"`
		CookieOrigins          string `yaml:"cookie_origins"`
	}
	if err := value.Decode(&raw); err != nil {
		return err
	}
	s.Port = raw.Port
	s.ShutdownTimeoutSeconds = raw.ShutdownTimeoutSeconds
	s.AllowedOrigins = splitOrigins(raw.AllowedOrigins)
	s.CookieOrigins = splitOrigins(raw.CookieOrigins)
	return nil
}

func splitOrigins(raw string) []string {
	var origins []string
	for _, o := range strings.Split(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}

type ApplicationConfig struct {
	Name                          string `yaml:"name"`
	SeedToken                     string `yaml:"seed_token"`
	BackofficeClientId            string `yaml:"backoffice_client_id"`
	BackofficeProtectedResourceId string `yaml:"backoffice_protected_resource_id"`
	BackofficeClientRedirectUri   string `yaml:"backoffice_client_redirect_uri"`
	BackofficeClientSecret        string `yaml:"backoffice_client_secret"`
	ProtocolUrl                   string `yaml:"protocol_url"`
}

type DatabasesConfig struct {
	Postgres      string `yaml:"postgres"`
	MigrationsDSN string `yaml:"migrations_dsn"`
}

func Load() (*Config, error) {
	env := EnvOr("ENV", "local")
	if env == "local" {
		_ = godotenv.Load(".env.local") // Local environment file is optional.
	}

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
