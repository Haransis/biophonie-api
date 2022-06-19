package main

import (
	"log"

	"github.com/haran/biophonie-api/controller"
	_ "github.com/haran/biophonie-api/docs"
)

// @title Swagger for biophonie-api
// @version         1.0
// @description     API of biophonie (https://secret-garden-77375.herokuapp.com/). Files are located in "assets/"
// @termsOfService  TODO

// @contact.name   TODO
// @contact.url    TODO
// @contact.email  TODO

// @license.name  GPL-3.0 license
// @license.url   https://www.gnu.org/licenses/gpl-3.0.en.html

// @BasePath /api/v1

func main() {
	c := controller.NewController()
	r := controller.SetupRouter(c)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Stopping server: %q", err)
	}
}
