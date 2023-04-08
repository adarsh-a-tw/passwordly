package secrets

import "gorm.io/gorm"

type SecretRepository interface {
	CreateCredential(secret *Secret, credential *Credential) error
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
