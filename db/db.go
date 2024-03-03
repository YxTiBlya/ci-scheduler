package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pingcap/errors"

	"github.com/YxTiBlya/ci-core/logger"
)

type DB struct {
	*pgxpool.Pool
	log *logger.Logger
	cfg Config
}

func New(cfg Config) *DB {
	return &DB{
		log: logger.New("db"),
		cfg: cfg,
	}
}

func (db *DB) Start(ctx context.Context) error {
	cfg, err := pgxpool.ParseConfig(db.cfg.String())
	if err != nil {
		return errors.Wrap(err, "pgxpool parse config")
	}

	db.Pool, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "pgxpool new with config")
	}
	if err := db.Ping(ctx); err != nil {
		return errors.Wrap(err, "db ping")
	}

	return nil
}

func (db *DB) Stop(ctx context.Context) error {
	db.Pool.Close()
	db.log.Sync()
	return nil
}
