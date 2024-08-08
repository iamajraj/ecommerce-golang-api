package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("ecommerce.db"), &gorm.Config{})
	return err
}
