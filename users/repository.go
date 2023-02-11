package users

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *User) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (ur *UserRepositoryImpl) Create(u *User) error {
	return ur.db.Create(u).Error
}
