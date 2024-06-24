package web

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/pragmahq/sso/database"
	"golang.org/x/crypto/bcrypt"
)

type ReqBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerAuthRoutes(router *echo.Echo, db *database.DB) {
	r := router.Group("/auth")
	r.POST("/register", registerUser(db))
	r.POST("/login", login(db))
	r.GET("/logout", logout)

	d := router.Group("/user")
	d.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  SECRET,
		TokenLookup: "cookie:Token",
	}))
	d.GET("/", getUser(db))
}

func registerUser(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req ReqBody
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		// Check if user already exists
		existingUser, _ := database.GetUserByEmail(db, req.Email)
		if existingUser != nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "User already exists"})
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}

		// Create new user
		user := &database.User{
			Id:       uuid.New().String(),
			Email:    req.Email,
			Password: string(hashedPassword),
		}
		err = user.Create(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
	}
}

func login(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req ReqBody
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		// Get user by email
		user, err := database.GetUserByEmail(db, req.Email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}

		// Check password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = user.Id
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token
		t, err := token.SignedString(SECRET)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		// Set cookie
		cookie := new(http.Cookie)
		cookie.Name = "Token"
		cookie.Value = t
		cookie.Expires = time.Now().Add(72 * time.Hour)
		cookie.Path = "/"
		cookie.HttpOnly = true
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]string{"token": t})
	}
}

func logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "Token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1

	c.SetCookie(cookie)
	return c.String(http.StatusOK, "Logged out successfully")
}

func getUser(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string)

		dbUser := &database.User{Id: userID}
		err := dbUser.Read(db)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"id":    dbUser.Id,
			"email": dbUser.Email,
		})
	}
}
