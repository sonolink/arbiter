package config

import (
	"fmt"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Discord  Discord
	Database Database
	Server   Server
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return &cfg, nil
}

type Discord struct {
	ClientID     string `env:"DISCORD_CLIENT_ID,required"`
	ClientSecret string `env:"DISCORD_CLIENT_SECRET,required"`
	RedirectURI  string `env:"DISCORD_REDIRECT_URI,required"`
}

type Database struct {
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	Name     string `env:"POSTGRES_DB,required"`
	SSLMode  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
}

func (d Database) DSN() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(d.User, d.Password),
		Host:     fmt.Sprintf("%s:%s", d.Host, d.Port),
		Path:     "/" + d.Name,
		RawQuery: "sslmode=" + d.SSLMode,
	}
	return u.String()
}

type Server struct {
	Addr string `env:"SERVER_ADDR" envDefault:":8080"`
}

func LoadDatabase() (Database, error) {
	var cfg Database
	if err := env.Parse(&cfg); err != nil {
		return Database{}, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}
