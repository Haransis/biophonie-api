package database

import (
	"database/sql"
	"fmt"
	"os"
)

func InitDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return db, fmt.Errorf("error opening database: %q", err)
	}

	if _, err := db.Exec(
		`CREATE TABLE [IF NOT EXISTS] accounts (
			user_id serial PRIMARY KEY,
			username VARCHAR ( 50 ) UNIQUE NOT NULL,
			password VARCHAR ( 50 ),
			email VARCHAR ( 255 ) UNIQUE,
			created_on TIMESTAMP NOT NULL,
			last_login TIMESTAMP
	)`); err != nil {
		return db, fmt.Errorf("could not create tables: %q", err)
	}

	return db, nil
}
