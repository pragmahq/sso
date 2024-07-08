package web

import (
	"net/http"
	"time"

	"github.com/go-pg/pg/v10"
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

type RegisterBody struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	InviteCode string `json:"inviteCode"`
}

func registerAuthRoutes(router *echo.Echo, db *database.DB) {
	r := router.Group("/api/auth")
	r.POST("/register", registerUser(db))
	r.POST("/login", login(db))
	r.GET("/logout", logout)
	r.GET("/validate", validateToken(db))
	r.GET("/validate-invite/:invite", validateInvite(db))

	d := router.Group("/api/user")
	d.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  SECRET,
		TokenLookup: "cookie:Token",
	}))
	d.GET("/", getUser(db))
}

func validateInvite(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.Param("invite")

		if code == "" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}

		inviteCode, err := database.GetInviteCode(db, code)
		if err != nil || inviteCode == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid invite code"})
		}
		if inviteCode.UsedBy != "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invite code already used"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "valid"})
	}
}

func registerUser(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req RegisterBody
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		existingUser := &database.User{Email: req.Email}
		err := db.Model(existingUser).Where("email = ?", req.Email).Select()
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "User already exists"})
		} else if err != pg.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}

		inviteCode, err := database.GetInviteCode(db, req.InviteCode)
		if err != nil || inviteCode == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid invite code"})
		}
		if inviteCode.UsedBy != "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invite code already used"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}

		user := &database.User{
			Id:       uuid.New().String(),
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		err = db.RunInTransaction(c.Request().Context(), func(tx *pg.Tx) error {
			_, err := tx.Model(user).Insert()
			if err != nil {
				return err
			}

			inviteCode.UsedBy = user.Id
			if inviteCode.UsedAt == nil {
				inviteCode.UsedAt = new(time.Time)
			}
			*inviteCode.UsedAt = time.Now()

			_, err = tx.Model(inviteCode).
				Set("used_by = ?used_by").
				Set("used_at = ?used_at").
				Where("id = ?id").
				Update()
			if err != nil {
				return err
			}

			return nil
		})

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

func validateToken(db *database.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No token provided"})
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			return []byte(SECRET), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token has expired"})
				}
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
			}

			user := &database.User{Id: userID}
			err := db.Model(user).WherePK().Select()
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"valid": true,
				"user": map[string]interface{}{
					"id":    user.Id,
					"email": user.Email,
				},
			})
		}

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}
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
