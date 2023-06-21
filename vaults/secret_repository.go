package vaults

import "gorm.io/gorm"

type SecretRepository interface {
	CreateCredential(credential *Credential) error
	FindCredentials(credentials *[]Credential, vaultId string) error
}

type SecretRepositoryImpl struct {
	Db *gorm.DB
}

func (sr *SecretRepositoryImpl) CreateCredential(credential *Credential) error {
	return sr.Db.Create(credential).Error
}

func (sr *SecretRepositoryImpl) FindCredentials(credentials *[]Credential, vaultId string) error {
	return sr.Db.Where("vault_refer = ?", vaultId).Order("updated_at DESC, id DESC").Find(credentials).Error
}
