package config

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Log      Log
	Discord  Discord
	Postgres Postgres
	Server   Server
}

func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}

type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

func (f *LogFormat) UnmarshalText(text []byte) error {
	switch format := LogFormat(text); format {
	case LogFormatText, LogFormatJSON:
		*f = format

		return nil
	default:
		return fmt.Errorf(
			"invalid log format %q, want %q or %q",
			format,
			LogFormatText,
			LogFormatJSON,
		)
	}
}

type Log struct {
	Format LogFormat  `env:"LOG_FORMAT" envDefault:"text"`
	Level  slog.Level `env:"LOG_LEVEL"  envDefault:"info"`
}

func (l Log) Handler(w io.Writer) slog.Handler {
	opts := &slog.HandlerOptions{Level: l.Level}

	switch l.Format {
	case LogFormatJSON:
		return slog.NewJSONHandler(w, opts)
	default:
		return slog.NewTextHandler(w, opts)
	}
}

func LoadLog() (Log, error) {
	cfg, err := env.ParseAs[Log]()
	if err != nil {
		return Log{}, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
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
	Host            string        `env:"SERVER_HOST"             envDefault:"127.0.0.1"`
	Port            int           `env:"SERVER_PORT"             envDefault:"8080"`
	ReadTimeout     time.Duration `env:"SERVER_READ_TIMEOUT"     envDefault:"5s"`
	WriteTimeout    time.Duration `env:"SERVER_WRITE_TIMEOUT"    envDefault:"30s"`
	IdleTimeout     time.Duration `env:"SERVER_IDLE_TIMEOUT"     envDefault:"120s"`
	ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

func (s Server) Addr() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

func LoadPostgres() (Postgres, error) {
	cfg, err := env.ParseAs[Postgres]()
	if err != nil {
		return Postgres{}, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}
