package vaults

type SecretType string

const (
	TypeCredential SecretType = "CREDENTIAL"
	TypeKey        SecretType = "KEY"
	TypeDocument              = "DOCUMENT"
)

func (st SecretType) IsValid() bool {
	switch st {
	case TypeCredential, TypeKey, TypeDocument:
		return true
	}
	return false
}
