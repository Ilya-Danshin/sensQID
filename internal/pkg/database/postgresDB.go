package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"sensQID/internal/pkg/cfg"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) Init(ctx context.Context, cfg cfg.DBCfg) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}

	db.conn = conn

	return nil
}
