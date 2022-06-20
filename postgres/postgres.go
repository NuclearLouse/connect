package postgres

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	User         string `cfg:"user"`
	Pass         string `cfg:"password"`
	Host         string `cfg:"host"`
	Port         int    `cfg:"port"`
	Database     string `cfg:"database"`
	Schema       string `cfg:"schema"`
	SSLMode      string `cfg:"sslmode"`
	PoolMaxConns int    `cfg:"max_open_conns"`
}

//DefaultConfig return
// Host:         "localhost",
// Port:         5432,
// Database:     "postgres",
// User:         "postgres",
// Pass:         "postgres",
// SSLMode:      "disable",
// PoolMaxConns: 10,
func DefaultConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         5432,
		Database:     "postgres",
		User:         "postgres",
		Pass:         "postgres",
		Schema:       "public",
		SSLMode:      "disable",
		PoolMaxConns: 10,
	}
}

func Connect(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {

	val := url.Values{}
	val.Set("sslmode", cfg.SSLMode)
	val.Set("pool_max_conns", fmt.Sprintf("%d", cfg.PoolMaxConns))
	url := url.URL{
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		User:     url.UserPassword(cfg.User, cfg.Pass),
		Path:     cfg.Database,
		RawQuery: val.Encode(),
	}

	pool, err := pgxpool.Connect(ctx, url.String())
	if err != nil {
		return nil, err
	}

	// pool.Config().MaxConnIdleTime = 60 * time.Second
	// pool.Config().MaxConnLifetime = 5 * time.Minute

	return pool, pool.Ping(ctx)
}
