package models

import "time"

type User struct {
	Id string `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;notNull"`
	Email string `json:"email" gorm:"unique;notNull"`
	Password string `gorm:"notNull"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}