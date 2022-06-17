package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

const schema = `--sql
	CREATE TABLE IF NOT EXISTS accounts (
		id serial PRIMARY KEY,
		name VARCHAR ( 50 ) UNIQUE NOT NULL,
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
		sound VARCHAR ( 42 ) NOT NULL
	);`

func InitDb() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return db, fmt.Errorf("error opening database: %q", err)
	}

	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO accounts (name, created_on) VALUES ($1,now()) ON CONFLICT DO NOTHING", "alice")
	tx.MustExec("INSERT INTO accounts (name, created_on) VALUES ($1,now()) ON CONFLICT DO NOTHING", "bob")
	tx.MustExec("INSERT INTO geopoints (title, user_id, location, created_on, amplitudes, picture, sound) VALUES ($1,$2,ST_GeomFromText($3),now(),$4,$5,$6) ON CONFLICT DO NOTHING", "Forest by night", "1", "Point(0.0 0.0)", "{1,2,3,4}", "https://example.com/image.jpg", "https://example.com/sound.mp3")
	tx.Commit()

	return db, nil
}
