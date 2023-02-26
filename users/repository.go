package users

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *User) error
	Find(username string, u *User) error
	FindById(id string, u *User) error
	Update(u *User) error
	UsernameAlreadyExists(username string) (bool, error)
	EmailAlreadyExists(email string) (bool, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (ur *UserRepositoryImpl) Create(u *User) error {
	return ur.db.Create(u).Error
}

func (ur *UserRepositoryImpl) Find(username string, u *User) error {
	return ur.db.Where("username = ?", username).Find(u).Error
}

func (ur *UserRepositoryImpl) FindById(id string, u *User) error {
	return ur.db.Where("id = ?", id).Find(u).Error
}

func (ur *UserRepositoryImpl) Update(u *User) error {
	return ur.db.Save(u).Error
}

func (ur *UserRepositoryImpl) UsernameAlreadyExists(username string) (bool, error) {
	var Exists bool
	err := ur.db.Raw("SELECT true FROM users WHERE username = ? AND value = 0", username).Scan(&Exists).Error
	return Exists, err
}

func (ur *UserRepositoryImpl) EmailAlreadyExists(email string) (bool, error) {
	var Exists bool
	err := ur.db.Raw("SELECT true FROM users WHERE email = ? AND value = 0", email).Scan(&Exists).Error
	return Exists, err
}
