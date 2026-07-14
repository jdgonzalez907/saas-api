package main

import (
	"jdgonzalez907/saas-api/internal/configuration"
	"log"

	_ "time/tzdata" // embed timezone DB for distroless images (no /usr/share/zoneinfo)
)

func main() {
	app, err := configuration.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	defer app.Close()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
