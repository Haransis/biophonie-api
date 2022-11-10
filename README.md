# biophonie-api

## Environment variables
The execution of biophonie API requires some environment variables:
* DATABASE_URL: the url of the postgres server
(example: "postgres://postgres:example@localhost:5432/postgres?sslmode=disable")
* PUBLIC_PATH: the public assets folder
(example: "$HOME/go/src/github.com/haran/biophonie-api/public")
* SECRETS_FOLDER: the folder containing rsa keys and admin password (example: "$HOME/go/src/github.com/haran/biophonie-api/testassets")
* PORT: opened port of the API