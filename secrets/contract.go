package secrets

import "time"

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
