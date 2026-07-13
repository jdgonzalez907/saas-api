package main

import (
	"log"

	"jdgonzalez907/users-api/internal/infrastructure"
	_ "time/tzdata" // embed timezone DB for distroless images (no /usr/share/zoneinfo)
)

func main() {
	app, err := infrastructure.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	defer app.Close()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
