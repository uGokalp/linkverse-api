package main

import (
	models "github.com/ugokalp/db"
	"github.com/ugokalp/server"
	"github.com/ugokalp/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	utils.LoadEnv()
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Url{})
	server.InitServer(db)
}
