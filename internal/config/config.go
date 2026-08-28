package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server  ServerConfig  `toml:"server"`
	Discord DiscordConfig `toml:"discord"`
	DB      DBConfig      `toml:"database"`
}

type ServerConfig struct {
	Addr string `toml:"addr"`
}

type DiscordConfig struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RedirectURI  string `toml:"redirect_uri"`
	GuildID      string `toml:"guild_id"`
}

type DBConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	SSLMode  string `toml:"sslmode"`
}

func Load(path string) (*Config, error) {
	if path == "" {
		path = "arbiter.toml"
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("config: parsing %s: %w", path, err)
	}

	cfg.applyDefaults()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) applyDefaults() {
	if c.Server.Addr == "" {
		c.Server.Addr = ":8080"
	}
	if c.DB.Host == "" {
		c.DB.Host = "localhost"
	}
	if c.DB.Port == 0 {
		c.DB.Port = 5432
	}
	if c.DB.Database == "" {
		c.DB.Database = "arbiter"
	}
	if c.DB.SSLMode == "" {
		c.DB.SSLMode = "disable"
	}
}

func (c *Config) validate() error {
	if c.Discord.ClientID == "" {
		return fmt.Errorf("config: [discord] client_id is required")
	}
	if c.Discord.ClientSecret == "" {
		return fmt.Errorf("config: [discord] client_secret is required")
	}
	if c.Discord.RedirectURI == "" {
		return fmt.Errorf("config: [discord] redirect_uri is required")
	}
	return nil
}

func (db DBConfig) DSN() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(db.User, db.Password),
		Host:     net.JoinHostPort(db.Host, strconv.Itoa(db.Port)),
		Path:     "/" + db.Database,
		RawQuery: url.Values{"sslmode": {db.SSLMode}}.Encode(),
	}

	return u.String()
}
