package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pragmahq/sso/database"
)

func RequirePermission(permission int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*database.User)
			if user.Permissions&permission == 0 {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient permissions"})
			}
			return next(c)
		}
	}
}

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return RequirePermission(database.PermissionAdmin)(next)
}

func RequireEditor(next echo.HandlerFunc) echo.HandlerFunc {
	return RequirePermission(database.PermissionEditor)(next)
}
