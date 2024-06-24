package database

import (
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testDB *DB

func TestMain(m *testing.M) {
	// Set up test database
	os.Setenv("DATABASE_URL", "postgres://rithulk:postgres@localhost:5432/pragma?sslmode=disable")

	var err error
	testDB, err = InitDB()
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

func TestInitDB(t *testing.T) {
	assert.NotNil(t, testDB)

	var n int
	_, err := testDB.QueryOne(pg.Scan(&n), "SELECT 1")
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}

func TestUserCRUD(t *testing.T) {
	// Test Create
	user := &User{
		Id:          uuid.New().String(),
		Email:       "test@example.com",
		Password:    "password123",
		Permissions: PermissionUser, // Set default permission to User

	}
	err := user.Create(testDB)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Id)

	// Test Read
	readUser := &User{Id: user.Id}
	err = readUser.Read(testDB)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, readUser.Email)

	// Test Update
	user.Email = "updated@example.com"
	err = user.Update(testDB)
	assert.NoError(t, err)

	err = readUser.Read(testDB)
	assert.NoError(t, err)
	assert.Equal(t, "updated@example.com", readUser.Email)

	// Test Delete
	err = user.Delete(testDB)
	assert.NoError(t, err)

	err = readUser.Read(testDB)
	assert.Error(t, err)
	assert.Equal(t, pg.ErrNoRows, err)
}

func TestGetUserByEmail(t *testing.T) {
	// Create a test user
	user := &User{
		Id:       uuid.New().String(),
		Email:    "getbyemail@example.com",
		Password: "password123",
	}
	err := user.Create(testDB)
	assert.NoError(t, err)

	// Test GetUserByEmail
	foundUser, err := GetUserByEmail(testDB, "getbyemail@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Id, foundUser.Id)
	assert.Equal(t, user.Email, foundUser.Email)

	// Test with non-existent email
	_, err = GetUserByEmail(testDB, "nonexistent@example.com")
	assert.Error(t, err)
	assert.Equal(t, pg.ErrNoRows, err)

	// Clean up
	err = user.Delete(testDB)
	assert.NoError(t, err)
}
