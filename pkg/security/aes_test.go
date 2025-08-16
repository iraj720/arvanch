package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {
	cases := []struct {
		name   string
		text   string
		encKey string
		decKey string
		err    bool
	}{
		{
			name:   "normal",
			text:   "Extraordinary Attorney Woo",
			encKey: "secret",
			decKey: "secret",
		},
		{
			name:   "very long text",
			text:   strings.Repeat("i", 20000),
			encKey: "secret",
			decKey: "secret",
		},
		{
			name:   "very long key",
			text:   "Extraordinary Attorney Woo",
			encKey: strings.Repeat("i", 20000),
			decKey: strings.Repeat("i", 20000),
		},
		{
			name:   "key mismatch",
			text:   "Extraordinary Attorney Woo",
			encKey: "secret1",
			decKey: "secret2",
			err:    true,
		},
	}

	for i := range cases {
		tc := cases[i]

		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)

			encryptor, err := NewAESTransformer(tc.encKey)
			a.NoError(err)

			encrypted, err := encryptor.Encrypt(tc.text)

			a.NoError(err)

			decryptor, err := NewAESTransformer(tc.decKey)
			a.NoError(err)

			decrypted, err := decryptor.Decrypt(encrypted)

			if tc.err {
				a.Error(err)
			} else {
				a.Equal(tc.text, decrypted)
			}
		})
	}
}
