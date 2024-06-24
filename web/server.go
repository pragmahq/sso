package web

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/pragmahq/sso/database"
)

var router *echo.Echo

func Serve(db *database.DB) {
	router = echo.New()

	registerAuthRoutes(router, db)

	router.Logger.Fatal(router.Start(":" + os.Getenv("SERVER_PORT")))
}
