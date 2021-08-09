package config

import (
	"fmt"
	"os"

	"github.com/sing3demons/go-echo-product/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	// database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok", host, user, pass, name, port)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	database.AutoMigrate(&models.Products{})
	database.AutoMigrate(&models.User{})
	// database.Migrator().DropTable(&models.Products{})

	db = database
}
func GetDB() *gorm.DB {
	return db
}
