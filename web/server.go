package web

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pragmahq/sso/database"
)

var router *echo.Echo

func Serve(db *database.DB) {
	router = echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	registerAuthRoutes(router, db)

	router.Logger.Fatal(router.Start(":" + os.Getenv("SERVER_PORT")))
}
