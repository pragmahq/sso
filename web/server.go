package web

import (
	"os"

	"github.com/labstack/echo/v4"
)

var router *echo.Echo

func Serve() {
	router = echo.New()

	router.Logger.Fatal(router.Start(":" + os.Getenv("SERVER_PORT")))
}
