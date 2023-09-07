package models

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func getDatabaseName(isTest bool) string {
	if isTest {
		MYSQL_DB_NAME_TEST := os.Getenv("MYSQL_DB_NAME_TEST")
		if MYSQL_DB_NAME_TEST == "" {
			MYSQL_DB_NAME_TEST = "tiktok_test"
		}

		return MYSQL_DB_NAME_TEST
	}

	MYSQL_DB_NAME := os.Getenv("MYSQL_DB_NAME")
	if MYSQL_DB_NAME == "" {
		MYSQL_DB_NAME = "tiktok"
	}

	return MYSQL_DB_NAME
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

	MySQL_USERNAME := os.Getenv("MYSQL_USERNAME")
	if MySQL_USERNAME == "" {
		MySQL_USERNAME = "root"
	}

	MySQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")

	MySQL_IP := os.Getenv("MYSQL_IP")
	if MySQL_IP == "" {
		MySQL_IP = "127.0.0.1"
	}

	MySQL_PORT := os.Getenv("MYSQL_PORT")
	if MySQL_PORT == "" {
		MySQL_PORT = "3306"
	}

	MySQL_DB_NAME := getDatabaseName(isTest)

	dsn := MySQL_USERNAME + ":" + MySQL_PASSWORD +
		"@tcp(" + MySQL_IP + ":" + MySQL_PORT + ")" + "/" + MySQL_DB_NAME +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,
		},
	)

	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}

	db = d

	// Drop table if isTest is true to avoid error
	// when running tests multiple times in a row.
	if isTest {
		db.Migrator().
			DropTable(&User{}, &Video{}, &FavoriteRel{}, &Comment{}, &FollowRel{}, &Message{})
	}

	db.AutoMigrate(&User{}, &Video{}, &FavoriteRel{}, &Comment{}, &FollowRel{}, &Message{})
}
