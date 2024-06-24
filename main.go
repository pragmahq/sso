package main

import (
	"log"

	"github.com/pragmahq/sso/database"
	"github.com/pragmahq/sso/web"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer db.Close()
	web.Serve()
}
