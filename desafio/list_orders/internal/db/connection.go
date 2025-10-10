package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	connStr := "host=db port=5432 user=postgres password=postgres dbname=orders sslmode=disable"
	return sql.Open("postgres", connStr)
}
