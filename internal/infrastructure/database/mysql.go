package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(dsn string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return
}
