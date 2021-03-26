package db

import (
	"os"

	"github.com/manedurphy/golang-start/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic("failed to connect to postgres")
	}

	db.AutoMigrate(&api.Person{})

	DB = db
}
