package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type HMACTransformer struct {
	key string
}

func NewHMACTransformer(key string) *HMACTransformer {
	return &HMACTransformer{
		key: key,
	}
}

func (h *HMACTransformer) Transform(text string) (string, error) {
	mac := hmac.New(sha256.New, []byte(h.key))

	_, err := mac.Write([]byte(text))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
