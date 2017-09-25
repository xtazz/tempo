package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
)

func initDb(filename string) (*sql.DB, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		bytes, err := ioutil.ReadFile("db.sql")

		if err != nil {
			return nil, err
		}

		dbInitStatement := string(bytes[:])

		db, err := sql.Open("sqlite3", filename)

		if err != nil {
			return nil, err
		}

		_, err = db.Exec(dbInitStatement)

		if err != nil {
			return nil, err
		}

		return db, nil
	} else {
		return sql.Open("sqlite3", filename)
	}
}
