package vaults

type CreateVaultRequest struct {
	Name string `json:"name" binding:"required"`
}

type VaultResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
