package database

import (
	"fmt"
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

	adminPassword := os.Getenv("SECRETS_FOLDER") + string(os.PathSeparator) + adminKey
	adminHash, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	db.MustExec(initTables)
	db.MustExec(createAdmin, "admin", adminHash)

	return db, nil
}
