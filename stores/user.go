package stores

import (
	"fmt"

	"github.com/ugokalp/db"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}

}

func (us *UserStore) FindById(id float64) (db.User, error) {
	var user = db.User{}
	err := us.db.Model(&user).Preload("Urls").Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (us *UserStore) FindByUsername(username string) (db.User, error) {
	var user = db.User{}
	err := us.db.Model(&db.User{}).
		Preload("Urls").
		Select([]string{"Id", "Bio", "Photo", "Username"}).
		Where("username = ?", username).
		Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (us *UserStore) Update(user db.User) (db.User, error) {
	id := fmt.Sprint(user.ID)
	err := us.db.Model(&db.User{}).Where("id = ?", id).Updates(&user).Error
	if err != nil {
		return user, err
	}
	us.db.Model(&user).Association("Urls").Replace(user.Urls)
	return user, nil
}

func (us *UserStore) Delete(userId float64) (db.User, error) {
	var user = db.User{}
	err := us.db.Model(&db.User{}).Where("id = ?", userId).Delete(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
