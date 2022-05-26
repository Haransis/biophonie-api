package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const schema = `--sql
	CREATE TABLE IF NOT EXISTS accounts (
		user_id serial PRIMARY KEY,
		username VARCHAR ( 50 ) UNIQUE NOT NULL,
		created_on TIMESTAMP NOT NULL,
		last_login TIMESTAMP
	)`

func InitDb() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return db, fmt.Errorf("error opening database: %q", err)
	}

	db.MustExec(schema)

	tx := db.MustBegin()
	tx.Exec("INSERT INTO accounts (username, created_on, last_login) VALUES ($1,now(),now())", "alice")
	tx.Exec("INSERT INTO accounts (username, created_on, last_login) VALUES ($1,now(),now())", "bob")
	tx.Commit()

	return db, nil
}
