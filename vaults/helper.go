package vaults

func validateVaultOwner(vr VaultRepository, vaultId string, userId string) (bool, error) {
	var vault Vault
	if err := vr.FetchById(vaultId, &vault); err != nil {
		return false, err
	}

	if vault.UserRefer != userId {
		return false, nil
	}

	return true, nil
}
