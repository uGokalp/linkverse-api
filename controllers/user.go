package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ugokalp/db"
	"github.com/ugokalp/stores"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	Db *gorm.DB
}

func NewController(db *gorm.DB) UserController {
	return UserController{
		Db: db,
	}
}

func (u *UserController) GetUsers(c *gin.Context) {
	var users []db.User
	err := u.Db.Model(&db.User{}).Preload("Urls").Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (u *UserController) GetUserByUsername(c *gin.Context) {
	username := c.Query("username")
	store := stores.NewUserStore(u.Db)
	user, err := store.FindByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *UserController) GetUser(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	store := stores.NewUserStore(u.Db)
	user, err := store.FindById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *UserController) UpdateUser(c *gin.Context) {
	var user db.User
	userId := c.MustGet("userId").(float64)
	store := stores.NewUserStore(u.Db)
	user, err := store.FindById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	err = c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err = store.Update(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	u.Db.Model(&user).Association("Urls").Replace(user.Urls)
	c.JSON(http.StatusOK, user)
}

func (u *UserController) SignUp(c *gin.Context) {
	var user db.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	user.Password = string(hash)
	err = u.Db.Model(&db.User{}).Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (u *UserController) Login(c *gin.Context) {

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email and/or password"})
		return
	}
	var user db.User
	err := u.Db.Model(&db.User{}).Where("email = ?", body.Email).First(&user).Error
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Invalid email and/or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (u *UserController) DeleteUser(c *gin.Context) {
	var user db.User
	userId := c.MustGet("userId").(float64)
	store := stores.NewUserStore(u.Db)
	user, err := store.Delete(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
