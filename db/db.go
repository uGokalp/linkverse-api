package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Bio      string `json:"bio"`
	Photo    string `json:"photo"`
	Username string `json:"username" binding:"required" gorm:"unique;not null;type:varchar(50);default:null"`
	Password string `json:"password" binding:"required" gorm:"not null;type:varchar(100);default:null"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null;type:varchar(100);default:null"`
	Urls     []Url  `json:"urls"`
}

type UserLogin struct {
	gorm.Model
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Url struct {
	gorm.Model
	UserId   uint   `json:"userId" binding:"required"`
	Url      string `json:"url" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Image    string `json:"image"`
	Position int    `json:"position"`
}
