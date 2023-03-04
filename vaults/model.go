package vaults

import (
	"time"

	"github.com/adarsh-a-tw/passwordly/users"
)

type Vault struct {
	Id        string `gorm:"primaryKey"`
	Name      string `gorm:"notNull"`
	UserRefer string
	User      users.User `gorm:"foreignKey:UserRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
