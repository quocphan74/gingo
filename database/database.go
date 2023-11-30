package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/quocphan74/gingo.git/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("Connect database successfully.")
	}
	DB = db
	db.AutoMigrate(
		&models.User{},
		&models.Code{},
		&models.Blog{},
		&models.Comment{},
		&models.Like{},
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // -> Log các câu lệnh truy vấn database trong terminal
	})
}
