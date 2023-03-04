package vaults

import "gorm.io/gorm"

type VaultRepository interface {
	Create(v *Vault) error
	FetchByUserId(userId string, vaults *[]Vault) error
}

type VaultRepositoryImpl struct {
	Db *gorm.DB
}

func (vr *VaultRepositoryImpl) Create(v *Vault) error {
	return vr.Db.Create(v).Error
}

func (vr *VaultRepositoryImpl) FetchByUserId(userId string, vaults *[]Vault) error {
	return vr.Db.Where("user_refer = ?", userId).Find(vaults).Error
}
