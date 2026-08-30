package config

import (
	"fmt"
	"net"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Discord  Discord
	Postgres Postgres
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

type Postgres struct {
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	Name     string `env:"POSTGRES_DB,required"`
	SSLMode  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
}

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

type Server struct {
	Addr string `env:"SERVER_ADDR" envDefault:":8080"`
}

func LoadPostgres() (Postgres, error) {
	var cfg Postgres
	if err := env.Parse(&cfg); err != nil {
		return Postgres{}, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}
