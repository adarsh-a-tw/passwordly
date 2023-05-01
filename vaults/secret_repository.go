package vaults

import "gorm.io/gorm"

type SecretRepository interface {
	CreateCredential(secret *Secret, credential *Credential) error
	FindSecretById(secret *Secret, secretId string) error
	FindCredentialById(credential *Credential, credentialId string) error
}

type SecretRepositoryImpl struct {
	Db *gorm.DB
}

func (sr *SecretRepositoryImpl) CreateCredential(secret *Secret, credential *Credential) error {
	tx := sr.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(secret).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(credential).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (sr *SecretRepositoryImpl) FindSecretById(secret *Secret, secretId string) error {
	return sr.Db.Where("id = ?", secretId).First(secret).Error
}

func (sr *SecretRepositoryImpl) FindCredentialById(credential *Credential, credentialId string) error {
	return sr.Db.Where("id = ?", credentialId).First(credential).Error
}
