package vaults

import "gorm.io/gorm"

type VaultRepository interface {
	Create(v *Vault) error
	FetchByUserId(userId string, vaults *[]Vault) error
	FetchById(id string, v *Vault) error
	Update(v *Vault) error
	Delete(v *Vault) error
}

type VaultRepositoryImpl struct {
	Db *gorm.DB
}

func (vr *VaultRepositoryImpl) Create(v *Vault) error {
	return vr.Db.Create(v).Error
}

func (vr *VaultRepositoryImpl) FetchByUserId(userId string, vaults *[]Vault) error {
	return vr.Db.Where("user_refer = ?", userId).Order("updated_at DESC, id DESC").Find(vaults).Error
}

func (vr *VaultRepositoryImpl) FetchById(id string, v *Vault) error {
	return vr.Db.Where("Id = ?", id).First(v).Error
}

func (vr *VaultRepositoryImpl) Update(v *Vault) error {
	return vr.Db.Save(v).Error
}

func (vr *VaultRepositoryImpl) Delete(v *Vault) error {
	tx := vr.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("vault_refer = ?", v.Id).Delete(&Credential{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(v).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
