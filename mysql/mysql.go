package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB //use the map to stores sessions
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) OpenDB() error {
	var err error
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	ip := os.Getenv("DB_IP")
	constring := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, ip, dbname)
	db.db, err = sql.Open("mysql", constring)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

type Tx struct {
	*sql.Tx
	db *DB
}

func (db *DB) BeginTx(opts *sql.TxOptions) (*Tx, error) {
	ctx := context.Background()
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx: tx,
		db: db,
	}, nil
}
