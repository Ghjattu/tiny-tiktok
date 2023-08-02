package models

import "github.com/joho/godotenv"

func beforeTest() {
	// When running tests in Go, the test environment is isolated from
	// the system environment variables by design.
	// This means that the test environment does not inherit
	// the environment variables set in our shell or system.
	// As a result, the os.Getenv() function in the test file will not be able to
	// access the environment variables set outside the test execution context.
	// So we need to load the environment variables manually.
	godotenv.Load("../.env")
}
