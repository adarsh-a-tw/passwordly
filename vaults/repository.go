package vaults

import "gorm.io/gorm"

type VaultRepository interface {
	Create(v *Vault) error
}

type VaultRepositoryImpl struct {
	Db *gorm.DB
}

func (vr *VaultRepositoryImpl) Create(v *Vault) error {
	return vr.Db.Create(v).Error
}
