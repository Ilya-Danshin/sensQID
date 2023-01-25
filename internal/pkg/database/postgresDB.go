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

var anonymizedName = "_anon"

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
	columns, err := db.GetTableColumns(table)
	if err != nil {
		return false, err
	}

	for _, column := range columns {
		if column == columnName {
			return true, nil
		}
	}

	return false, nil
}

func (db *DB) GetTableColumns(table string) ([]string, error) {
	rows, err := db.conn.Query(db.ctx,
		`SELECT column_name
  				FROM information_schema.columns
 				WHERE table_name = $1;`, table)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var columns []string
	for rows.Next() {
		var column string
		err = rows.Scan(&column)
		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	return columns, nil
}

func isContain(i string, arr []string) bool {
	for _, n := range arr {
		if n == i {
			return true
		}
	}
	return false
}

func (db *DB) CreateAnonTable(tableName string, columns []string) error {
	// Get original table columns
	columnsOrigin, err := db.GetTableColumns(tableName)
	if err != nil {
		return err
	}

	queryStr := `CREATE TABLE ` + tableName + anonymizedName + `(`
	for i, column := range columnsOrigin {
		queryStr += column

		// Get type of column from origin table
		columnType, err := db.GetColumnType(tableName, column)
		if err != nil {
			return err
		}
		queryStr += " " + columnType
		// If column in anonymized slice then it should be array type
		if isContain(column, columns) {
			queryStr += "[]"
		}

		if i != len(columnsOrigin)-1 {
			queryStr += ",\n"
		} else {
			queryStr += "\n)"
		}
	}

	rows, err := db.conn.Query(db.ctx, queryStr)
	defer rows.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRandomIntValues(table, column string, n int, value int) ([]int, error) {
	query := "SELECT " + column + " FROM (SELECT DISTINCT " + column + " FROM " + table + " WHERE " + column + " != $1 " + " GROUP BY " + column + ") t ORDER BY random() LIMIT $2"

	//rows, err := db.conn.Query(db.ctx,
	//	`SELECT $1
	//		FROM $2
	//		ORDER BY random()
	//		LIMIT $3`, column, table, n)
	rows, err := db.conn.Query(db.ctx, query, value, n)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var result []int
	for rows.Next() {
		var num int
		err = rows.Scan(&num)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}

	return result, nil
}

func (db *DB) GetRandomStrValues(table, column string, n int, value string) ([]string, error) {
	query := "SELECT " + column + " FROM (SELECT DISTINCT " + column + " FROM " + table + " WHERE " + column + " != $1 " + " GROUP BY " + column + ") t ORDER BY random() LIMIT $2"
	//rows, err := db.conn.Query(db.ctx,
	//	`SELECT $1
	//		FROM $2
	//		ORDER BY random()
	//		LIMIT $3`, column, table, n)
	rows, err := db.conn.Query(db.ctx, query, value, n)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var result []string
	for rows.Next() {
		var str string
		err = rows.Scan(&str)
		if err != nil {
			return nil, err
		}
		result = append(result, str)
	}

	return result, nil
}

func (db *DB) GetColumnType(table, column string) (string, error) {
	row := db.conn.QueryRow(db.ctx,
		`SELECT data_type
  				FROM information_schema.columns
 				WHERE table_name = $1
				AND column_name = $2;`, table, column)

	var columnType string
	err := row.Scan(&columnType)
	if err != nil {
		return "", err
	}

	return columnType, nil
}

func (db *DB) GetTableVolume(table string) (int, error) {
	query := "SELECT count(*) FROM " + table + ";"

	row := db.conn.QueryRow(db.ctx,
		query)

	var volume int
	err := row.Scan(&volume)
	if err != nil {
		return 0, err
	}

	return volume, nil
}

func (db *DB) GetIntValue(table, column string, n int) (int, error) {
	query := "SELECT " + column + " FROM " + table + " OFFSET $1 LIMIT 1;"
	rows, err := db.conn.Query(db.ctx,
		query, n)
	defer rows.Close()
	if err != nil {
		return 0, err
	}
	var value int
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return 0, err
		}
	}
	return value, nil
}

func (db *DB) GetTextValue(table, column string, n int) (string, error) {
	query := "SELECT " + column + " FROM " + table + " OFFSET $1 LIMIT 1;"
	rows, err := db.conn.Query(db.ctx,
		query, n)
	defer rows.Close()
	if err != nil {
		return "", err
	}
	var value string
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return "", err
		}
	}
	return value, nil
}
