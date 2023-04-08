package secrets

import (
	"time"

	"github.com/adarsh-a-tw/passwordly/vaults"
)

type Secret struct {
	Id         string `json:"id" gorm:"primaryKey"`
	Name       string `json:"name" gorm:"notNull"`
	VaultRefer string
	Vault      vaults.Vault `gorm:"foreignKey:VaultRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	Type       string       `json:"type" gorm:"notNull"`
}

func (s *Secret) SecretType() SecretType {
	// Need to find better way to do this
	if s.Type == "CREDENTIAL" {
		return TypeCredential
	} else if s.Type == "KEY" {
		return TypeKey
	}
	return TypeDocument
}

type Securable interface {
	Type() SecretType
}

type Credential struct {
	Id          string `json:"id" gorm:"primaryKey"`
	Username    string `json:"username" gorm:"notNull"`
	Password    string `json:"password" gorm:"notNull"`
	SecretRefer string
	Secret      Secret `gorm:"foreignKey:SecretRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (c *Credential) Type() SecretType { return TypeCredential }

type Key struct {
	Id          string `json:"id" gorm:"primaryKey"`
	Value       string `json:"value" gorm:"notNull"`
	SecretRefer string
	Secret      Secret `gorm:"foreignKey:SecretRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (k *Key) Type() SecretType { return TypeKey }

type Document struct {
	Id          string `json:"id" gorm:"primaryKey"`
	Content     string `json:"content" gorm:"notNull"`
	SecretRefer string
	Secret      Secret `gorm:"foreignKey:SecretRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (d *Document) Type() SecretType { return TypeDocument }
