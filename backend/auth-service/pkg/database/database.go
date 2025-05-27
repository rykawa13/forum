package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// Config represents database configuration
type Config struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DefaultConfig returns default database configuration
func DefaultConfig() Config {
	return Config{
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

func NewPostgresConnectionFromURL(url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	// Apply configuration
	cfg := DefaultConfig()
	config.MaxConns = int32(cfg.MaxOpenConns)
	config.MaxConnLifetime = cfg.ConnMaxLifetime

	// Create pool
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func RunMigrations(dbURL string, migrationsPath string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// HealthCheck проверяет доступность базы данных
func HealthCheck(ctx context.Context, db *pgxpool.Pool) error {
	return db.Ping(ctx)
}
