package models

import (
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func getDatabaseName(isTest bool) string {
	if isTest {
		return os.Getenv("MYSQL_DB_NAME_TEST")
	}

	return os.Getenv("MYSQL_DB_NAME")
}

// InitDatabase initializes the database,
// isTest is true if the database is for testing
// and false if the database is for production.
//
//	@param isTest bool
func InitDatabase(isTest bool) {
	// When running tests in Go, the test environment is isolated from
	// the system environment variables by design.
	// This means that the test environment does not inherit
	// the environment variables set in our shell or system.
	// As a result, the os.Getenv() function in the test file will not be able to
	// access the environment variables set outside the test execution context.
	// So we need to load the environment variables manually.
	godotenv.Load("../.env")

	mysql_username := os.Getenv("MYSQL_USERNAME")
	mysql_password := os.Getenv("MYSQL_PASSWORD")
	mysql_ip := os.Getenv("MYSQL_IP")
	mysql_port := os.Getenv("MYSQL_PORT")
	mysql_db_name := getDatabaseName(isTest)

	dsn := mysql_username + ":" + mysql_password +
		"@tcp(" + mysql_ip + ":" + mysql_port + ")" + "/" + mysql_db_name +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = d

	// Drop table if isTest is true to avoid error
	// when running tests multiple times in a row.
	if isTest {
		db.Migrator().DropTable(&User{}, &Video{}, &FavoriteRel{})
	}

	db.AutoMigrate(&User{}, &Video{}, &FavoriteRel{})
}
