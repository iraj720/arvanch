//nolint:gosec
package security

import (
	"crypto/sha1"
	"encoding/hex"
)

type SHA1Transformer struct{}

func (SHA1Transformer) Transform(text string) (string, error) {
	s := sha1.Sum([]byte(text))
	return hex.EncodeToString(s[:]), nil
}
