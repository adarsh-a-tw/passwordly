package vaults

type CreateVaultRequest struct {
	Name string `json:"name" binding:"required"`
}

type VaultResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type VaultListResponse struct {
	Vaults []VaultResponse `json:"vaults"`
}

func (vlr *VaultListResponse) load(vaults []Vault) {
	var vaultResponses []VaultResponse
	for _, vault := range vaults {
		vaultResponses = append(vaultResponses, VaultResponse{Id: vault.Id, Name: vault.Name})
	}
	vlr.Vaults = vaultResponses
}
