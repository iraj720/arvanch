package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

type AESTransformer struct {
	aead cipher.AEAD
}

func NewAESTransformer(keyStr string) (*AESTransformer, error) {
	keyArr := sha256.Sum256([]byte(keyStr))

	block, err := aes.NewCipher(keyArr[:])
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &AESTransformer{aead: aesGCM}, nil
}

func (a *AESTransformer) Transform(text string) (string, error) {
	return a.Encrypt(text)
}

func (a *AESTransformer) Encrypt(stringToEncrypt string) (string, error) {
	plaintext := []byte(stringToEncrypt)

	nonce := make([]byte, a.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := a.aead.Seal(nonce, nonce, plaintext, nil)

	return hex.EncodeToString(ciphertext), nil
}

func (a *AESTransformer) Decrypt(encryptedString string) (string, error) {
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	nonceSize := a.aead.NonceSize()

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := a.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
