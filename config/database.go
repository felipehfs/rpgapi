package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Environment is a contant
type Environment string

const (
	// Development connects to the development database
	Development Environment = "development"
	// Test connects to the test database
	Test Environment = "testing"
)

var (
	testUser         = os.Getenv("RPGAPI_TEST_USER")
	testPassword     = os.Getenv("RPGAPI_TEST_PASSWORD")
	testDatabaseName = os.Getenv("RPGAPI_TEST_DATABASE_NAME")

	developmentUser         = os.Getenv("RPGAPI_DEVELOPMENT_USER")
	developmentPassword     = os.Getenv("RPGAPI_DEVELOPMENT_PASSWORD")
	developmentDatabaseName = os.Getenv("RPGAPI_DEVELOPMENT_DATABASE_NAME")
)

// SetupDatabase setups the database
func SetupDatabase(environment Environment) (*sql.DB, error) {

	var dbInfo string

	if environment == Development {
		dbInfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			developmentUser, developmentPassword, developmentDatabaseName)
	}

	if environment == Test {
		dbInfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			testUser, testPassword, testDatabaseName)
	}

	return sql.Open("postgres", dbInfo)
}
