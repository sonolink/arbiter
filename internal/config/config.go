package config

import (
	"fmt"
	"net"
	"net/url"

	"github.com/caarlos0/env/v11"
)

// Config holds all runtime settings for the application.
type Config struct {
	Discord  Discord
	Postgres Postgres
	Server   Server
}

// Load reads the full configuration from environment variables.
func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return &cfg, nil
}

// Discord holds the OAuth credentials for the Discord application.
type Discord struct {
	ClientID     string `env:"DISCORD_CLIENT_ID,required"`
	ClientSecret string `env:"DISCORD_CLIENT_SECRET,required"`
	RedirectURI  string `env:"DISCORD_REDIRECT_URI,required"`
}

// Postgres holds the settings used to build a database connection.
type Postgres struct {
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	Name     string `env:"POSTGRES_DB,required"`
	SSLMode  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
}

// DSN builds a Postgres connection string from the settings.
func (p Postgres) DSN() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(p.User, p.Password),
		Host:     net.JoinHostPort(p.Host, p.Port),
		Path:     "/" + p.Name,
		RawQuery: url.Values{"sslmode": {p.SSLMode}}.Encode(),
	}
	return u.String()
}

// Server holds the HTTP listener settings.
type Server struct {
	Addr string `env:"SERVER_ADDR" envDefault:":8080"`
}

// LoadPostgres reads only the Postgres settings, without the rest of the app config.
func LoadPostgres() (Postgres, error) {
	var cfg Postgres
	if err := env.Parse(&cfg); err != nil {
		return Postgres{}, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}
