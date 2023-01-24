package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"sensQID/internal/pkg/cfg"
)

type DB struct {
	ctx  context.Context
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

	db.ctx = ctx
	db.conn = conn

	return nil
}

func (db *DB) IsTableExist(tableName string) (bool, error) {
	rows, err := db.conn.Query(db.ctx,
		`SELECT tablename FROM pg_catalog.pg_tables;`)
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return false, err
		}

		if table == tableName {
			return true, nil
		}
	}

	return false, nil
}

func (db *DB) IsColumnExist(table, columnName string) (bool, error) {
	rows, err := db.conn.Query(db.ctx,
		`SELECT column_name
  				FROM information_schema.columns
 				WHERE table_name   = $1;`, table)
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() {
		var column string
		err := rows.Scan(&column)
		if err != nil {
			return false, err
		}

		if column == columnName {
			return true, nil
		}
	}

	return false, nil
}
