package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/adarsh-a-tw/passwordly/common"
)

type AesEncryptionProvider struct {
	gcm cipher.AEAD
}

func NewEncryptionProvider() (*AesEncryptionProvider, error) {
	c, err := aes.NewCipher([]byte(common.Cfg.EncryptionKey))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	return &AesEncryptionProvider{gcm}, nil
}

func (aep *AesEncryptionProvider) Encrypt(plainText string) (string, error) {
	nonce := make([]byte, aep.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return string(aep.gcm.Seal(nonce, nonce, []byte(plainText), nil)), nil
}

func (aep *AesEncryptionProvider) Decrypt(cipherText string) (string, error) {
	cipherTextBytes := []byte(cipherText)
	nonceSize := aep.gcm.NonceSize()

	nonce, cipherTextBytes := cipherTextBytes[:nonceSize], cipherTextBytes[nonceSize:]
	plaintextBytes, err := aep.gcm.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}
	return string(plaintextBytes), nil
}
