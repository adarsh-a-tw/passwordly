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

type Credential struct {
	Id         string `json:"id" gorm:"primaryKey"`
	Name       string `json:"name" gorm:"notNull"`
	Username   string `json:"username" gorm:"notNull"`
	Password   string `json:"password" gorm:"notNull"`
	VaultRefer string
	Vault      Vault     `gorm:"foreignKey:VaultRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Key struct {
	Id         string `json:"id" gorm:"primaryKey"`
	Name       string `gorm:"notNull"`
	Value      string `json:"value" gorm:"notNull"`
	VaultRefer string
	Vault      Vault     `gorm:"foreignKey:VaultRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Document struct {
	Id         string `json:"id" gorm:"primaryKey"`
	Name       string `gorm:"notNull"`
	Content    string `json:"content" gorm:"notNull"`
	VaultRefer string
	Vault      Vault     `gorm:"foreignKey:VaultRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Securable interface {
	Type() SecretType
}

func (Credential) Type() SecretType {
	return TypeCredential
}

func (Key) Type() SecretType {
	return TypeKey
}

func (Document) Type() SecretType {
	return TypeDocument
}
