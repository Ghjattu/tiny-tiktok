package models

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	mysql_username := os.Getenv("MYSQL_USERNAME")
	mysql_password := os.Getenv("MYSQL_PASSWORD")
	mysql_ip := os.Getenv("MYSQL_IP")
	mysql_port := os.Getenv("MYSQL_PORT")
	mysql_db_name := os.Getenv("MYSQL_DB_NAME")

	dsn := mysql_username + ":" + mysql_password +
		"@tcp(" + mysql_ip + ":" + mysql_port + ")" + "/" + mysql_db_name +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db = d

	db.AutoMigrate(&User{})
}
