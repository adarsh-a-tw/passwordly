package vaults

import "time"

type CreateVaultRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateVaultRequest CreateVaultRequest

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

type CreateSecretRequest struct {
	Name     string     `json:"name" binding:"required"`
	Type     SecretType `json:"type" binding:"required,secret_type"`
	Username string     `json:"username,omitempty"`
	Password string     `json:"password,omitempty"`
	Value    string     `json:"value,omitempty"`
	Document string     `json:"document,omitempty"`
}

type SecretResponse struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Type      SecretType `json:"type"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Username  string     `json:"username,omitempty"`
	Password  string     `json:"password,omitempty"`
	Value     string     `json:"value,omitempty"`
	Document  string     `json:"document,omitempty"`
}
