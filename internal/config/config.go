package config

import (
	"fmt"

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
	URL string `env:"DATABASE_URL,required"`
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
