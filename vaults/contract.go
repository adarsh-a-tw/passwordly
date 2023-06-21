package vaults

type CreateVaultRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateVaultRequest CreateVaultRequest

type VaultResponse struct {
	Id        string           `json:"id"`
	Name      string           `json:"name"`
	Secrets   []SecretResponse `json:"secrets,omitempty"`
	CreatedAt int64            `json:"created_at,omitempty"`
	UpdatedAt int64            `json:"updated_at,omitempty"`
}

func (vr *VaultResponse) load(v Vault, creds []Credential) {
	vr.Id = v.Id
	vr.Name = v.Name

	var secretResponses = make([]SecretResponse, 0)
	for _, cred := range creds {
		sr := SecretResponse{}
		sr.load(cred)
		secretResponses = append(secretResponses, sr)
	}

	vr.Secrets = secretResponses
}

type VaultListResponse struct {
	Vaults []VaultResponse `json:"vaults"`
}

func (vlr *VaultListResponse) load(vaults []Vault) {
	var vaultResponses = make([]VaultResponse, 0)
	for _, vault := range vaults {
		vaultResponses = append(vaultResponses, VaultResponse{Id: vault.Id, Name: vault.Name, CreatedAt: vault.CreatedAt.Unix(), UpdatedAt: vault.UpdatedAt.Unix()})
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
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
	Username  string     `json:"username,omitempty"`
	Password  string     `json:"password,omitempty"`
	Value     string     `json:"value,omitempty"`
	Document  string     `json:"document,omitempty"`
}

func (sr *SecretResponse) load(s Securable) {
	switch s.Type() {
	case TypeCredential:
		cred := s.(Credential)
		sr.Id = cred.Id
		sr.Name = cred.Name
		sr.Type = TypeCredential
		sr.CreatedAt = cred.CreatedAt.Unix()
		sr.UpdatedAt = cred.UpdatedAt.Unix()
		sr.Username = cred.Username
		sr.Password = cred.Password
	case TypeKey:
		key := s.(Key)
		sr.Id = key.Id
		sr.Name = key.Name
		sr.Type = TypeKey
		sr.CreatedAt = key.CreatedAt.Unix()
		sr.UpdatedAt = key.UpdatedAt.Unix()
		sr.Value = key.Value
	case TypeDocument:
		document := s.(Document)
		sr.Id = document.Id
		sr.Name = document.Name
		sr.Type = TypeDocument
		sr.CreatedAt = document.CreatedAt.Unix()
		sr.UpdatedAt = document.UpdatedAt.Unix()
		sr.Document = document.Content
	}
}
