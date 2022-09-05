package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// TODO retrieve it from secret config
var adminPassword, _ = bcrypt.GenerateFromPassword([]byte("57aba9df-969f-4871-a095-e916d06ba38b"), bcrypt.DefaultCost)

const schema = `--sql
	CREATE TABLE IF NOT EXISTS accounts (
		id serial PRIMARY KEY,
		name VARCHAR ( 20 ) UNIQUE NOT NULL,
		password VARCHAR ( 60 ) UNIQUE NOT NULL,
		admin BOOLEAN NOT NULL DEFAULT FALSE,
		created_on TIMESTAMP NOT NULL
	);
	CREATE TABLE IF NOT EXISTS geopoints (
		id serial PRIMARY KEY,
		title VARCHAR ( 30 ) NOT NULL,
		user_id serial NOT NULL,
		location geography ( POINT , 4326 ),
		created_on TIMESTAMP NOT NULL,
		amplitudes INT [],
		picture VARCHAR ( 42 ) NOT NULL,
		sound VARCHAR ( 42 ) NOT NULL,
		available BOOLEAN NOT NULL DEFAULT FALSE
	);`

/*const addIndexing = `--sql
CREATE INDEX idx_geopoints_geom ON geopoints USING gist ((location));
`*/ //creates a gist index to help location research (to use ?)

func InitDb() (*sqlx.DB, error) {
	fmt.Println(os.Getenv("DATABASE_URL"))
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return db, fmt.Errorf("error opening database: %q", err)
	}

	db.MustExec(schema)
	db.MustExec("INSERT INTO accounts (name, created_on, password, admin) VALUES ($1,now(),$2,'t') ON CONFLICT DO NOTHING", "admin", adminPassword)

	return db, nil
}
