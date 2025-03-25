package main

import (
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetConnectionString(t *testing.T) {
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_PASSWORD", "testpassword")
	os.Setenv("POSTGRES_DB", "testdb")

	expected := "host=localhost port=5432 user=testuser password=testpassword dbname=testdb sslmode=disable"
	actual := getConnectionString()

	assert.Equal(t, expected, actual, "Connection string does not match expected value")
}

func TestGetConnectionString_MissingEnvVars(t *testing.T) {
	os.Clearenv()

	assert.Panics(t, func() {
		getConnectionString()
	}, "Expected getConnectionString to panic when environment variables are missing")
}

func TestGetConnectionString_portNotInt(t *testing.T) {
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "notAnInt")
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_PASSWORD", "testpassword")
	os.Setenv("POSTGRES_DB", "testdb")
	assert.Panics(t, func() {
		getConnectionString()
	}, "Expected getConnectionString to panic when POSTGRES_PORT is not an integer")
}

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	assert.NotNil(t, router, "Router should not be nil")
	assert.IsType(t, &gin.Engine{}, router.Engine, "Router.Engine should be of type *gin.Engine")
}

func TestSetupRouter(t *testing.T) {
	dbm, _, err := sqlmock.New()
	assert.NoError(t, err)
	db := sqlx.NewDb(dbm, "postgres")

	assert.NoError(t, err, "Failed to open database connection")
	defer db.Close()

	router := setupRouter(db)

	assert.NotNil(t, router, "Router should not be nil")
	assert.IsType(t, &Router{}, router, "setupRouter should return a *Router")

	// Check if routes are correctly set up
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/exercises/:uuid"},
		{"POST", "/exercises"},
		{"PUT", "/exercises/:uuid"},
		{"DELETE", "/exercises/:uuid"},
	}

	for _, route := range routes {
		assert.True(t, routeExists(router.Engine, route.method, route.path),
			"Route %s %s should exist", route.method, route.path)
	}
}

func routeExists(engine *gin.Engine, method, path string) bool {
	for _, route := range engine.Routes() {
		if route.Method == method && route.Path == path {
			return true
		}
	}
	return false
}
