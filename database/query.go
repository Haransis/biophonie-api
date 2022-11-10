package database

const (
	initTables = `--sql
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
			amplitudes FLOAT [],
			picture VARCHAR ( 42 ) NOT NULL,
			sound VARCHAR ( 42 ) NOT NULL,
			available BOOLEAN NOT NULL DEFAULT FALSE
		);
	`

	createAdmin = `--sql
		INSERT INTO accounts (name, created_on, password, admin) 
		VALUES ($1,now(),$2,'t') 
		ON CONFLICT DO NOTHING
	`

	GetClosestGeoId = `--sql
		WITH excluded(id) AS ( SELECT UNNEST($2::int[])) 
		SELECT geo.id FROM geopoints geo 
		WHERE NOT EXISTS(SELECT 1 FROM excluded e WHERE geo.id = e.id) AND available = TRUE
		ORDER BY geo.location <-> GeomFromEWKB($1)
		LIMIT 1;
	`

	GetGeoPoint = `--sql
		SELECT * FROM geopoints WHERE id = $1
	`

	GetUserByName = `--sql
		SELECT * FROM accounts WHERE name = $1
	`

	GetUserById = `--sql
		SELECT * FROM accounts WHERE id = $1
	`

	PostUser = `--sql
		INSERT INTO accounts (name, created_on, password) 
		VALUES ($1,now(),$2) 
		RETURNING id
	`

	PostGeoPoint = `--sql
		INSERT INTO geopoints (title, user_id, location, amplitudes, picture, sound, created_on) 
		VALUES (:title,:user_id,GeomFromEWKB(:location),:amplitudes,:picture,:sound,:created_on) 
		RETURNING id
	`

	EnableGeoPoint = `--sql
		UPDATE geopoints SET available = TRUE WHERE id = $1 AND available = FALSE
	`

	DeleteGeoPoint = `--sql
		DELETE FROM geopoints WHERE id = $1
	`

	MakeAdmin = `--sql
		UPDATE accounts SET admin = TRUE WHERE id = $1
	`
)
