package persistence

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"telemetry.ex.gosuda.org/telemetry/internal/persistence/database"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

type PersistenceClientConfig struct {
	DSN             string
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

type PersistenceClient struct {
	pool *sql.DB
	db   *database.Queries
}

var _ types.PersistenceService = (*PersistenceClient)(nil)

func NewPersistenceClient(ctx context.Context, config *PersistenceClientConfig) (*PersistenceClient, error) {
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	dbtx := database.New(db)

	return &PersistenceClient{pool: db, db: dbtx}, nil
}

func (g *PersistenceClient) Close() error {
	err := g.pool.Close()
	if err != nil {
		return err
	}

	return nil
}
