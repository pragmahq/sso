package web

import (
	"net/http"

	"github.com/labstack/echo"
)

type ReqBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerAuthRoutes() {
	r := router.Group("/auth")
	// r.POST("/", login)
	r.GET("/", logout)
	// d := router.Group("/user")
	// d.Use(echojwt.WithConfig(echojwt.Config{
	// 	SigningKey:  SECRET,
	// 	TokenLookup: "cookie:Token",
	// }))
	// d.GET("/", getuser)
}

// func login(c echo.Context) error {}

func logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "Token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1

	c.SetCookie(cookie)
	return c.String(http.StatusOK, "OK")
}
