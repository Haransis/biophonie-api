package database

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

/*const addIndexing = `--sql
CREATE INDEX idx_geopoints_geom ON geopoints USING gist ((location));
`*/ //creates a gist index to help location research (to use ?)

const adminKey = "admin.key"

func InitDb() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return db, fmt.Errorf("error opening database: %q", err)
	}

	adminFile := os.Getenv("SECRETS_FOLDER") + string(os.PathSeparator) + adminKey
	adminPassword, err := ioutil.ReadFile(adminFile)
	if err != nil {
		return db, fmt.Errorf("error opening admin file: %q", err)
	}

	adminHash, _ := bcrypt.GenerateFromPassword(adminPassword, bcrypt.DefaultCost)
	db.MustExec(initTables)
	db.MustExec(createAdmin, "admin", adminHash)

	return db, nil
}
