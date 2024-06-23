package main

import (
	"log"

	"github.com/pragmahq/sso/db"
	"github.com/pragmahq/sso/web"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	web.Serve()
}
