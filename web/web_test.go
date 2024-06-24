package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pragmahq/sso/database"
	"github.com/stretchr/testify/assert"
)

var testDB *database.DB

func TestMain(m *testing.M) {
	// Set up test database
	os.Setenv("DATABASE_URL", "postgres://rithulk:postgres@localhost:5432/pragma?sslmode=disable")

	var err error
	testDB, err = database.InitDB()
	if err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Clean up
	_, err = testDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		panic(err)
	}
	testDB.Close()

	os.Exit(code)
}

func TestRegisterUser(t *testing.T) {
	e := echo.New()
	reqBody := `{"email":"test@example.com","password":"password123"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := registerUser(testDB)
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User created successfully", response["message"])
	}
}

func TestLogin(t *testing.T) {
	e := echo.New()
	reqBody := `{"email":"test@example.com","password":"password123"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := login(testDB)
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["token"])
	}
}

func TestLogout(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, logout(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Logged out successfully", rec.Body.String())

		cookie := rec.Header().Get("Set-Cookie")
		assert.Contains(t, cookie, "Token=")
		assert.Contains(t, cookie, "Max-Age=0")
	}
}

func TestGetUser(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a test user
	testUser := &database.User{
		Id:       uuid.New().String(),
		Email:    "testuser@example.com",
		Password: "password123",
	}
	err := testUser.Create(testDB)
	assert.NoError(t, err)

	// Mock JWT token
	c.Set("user", &jwt.Token{
		Claims: jwt.MapClaims{
			"user_id": testUser.Id,
		},
	})

	h := getUser(testDB)
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.Id, response["id"])
		assert.Equal(t, testUser.Email, response["email"])
	}

	// Clean up
	err = testUser.Delete(testDB)
	assert.NoError(t, err)
}
