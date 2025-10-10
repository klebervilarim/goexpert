package db

import (
	"database/sql"
	"io/ioutil"
	"log"
)

func RunMigrations(db *sql.DB) {
	script, err := ioutil.ReadFile("internal/db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}
}
