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
	Db *gorm.DB
}

func (ur *UserRepositoryImpl) Create(u *User) error {
	return ur.Db.Create(u).Error
}

func (ur *UserRepositoryImpl) Find(username string, u *User) error {
	return ur.Db.Where("username = ?", username).First(u).Error
}

func (ur *UserRepositoryImpl) FindById(id string, u *User) error {
	return ur.Db.Where("id = ?", id).First(u).Error
}

func (ur *UserRepositoryImpl) Update(u *User) error {
	return ur.Db.Save(u).Error
}

func (ur *UserRepositoryImpl) UsernameAlreadyExists(username string) (bool, error) {
	var Exists bool
	err := ur.Db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&Exists).Error
	return Exists, err
}

func (ur *UserRepositoryImpl) EmailAlreadyExists(email string) (bool, error) {
	var Exists bool
	err := ur.Db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&Exists).Error
	return Exists, err
}
