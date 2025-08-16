package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHMAC(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		key            string
		expectedResult string
	}{
		{
			name:           "normal",
			text:           "Extraordinary Attorney Woo",
			key:            "secret",
			expectedResult: "pQaI9BsmlDUauAyThqfHFxSaXFfJ8jEQpYA90/aazFU=",
		},
		{
			name:           "very long key",
			text:           "Extraordinary Attorney Woo",
			key:            strings.Repeat("i", 20000),
			expectedResult: "6cWt1MrArvWCh2PBAnE++vhk7Lz3OIcw9nzsL8qDBZ4=",
		},
		{
			name:           "very long text",
			text:           strings.Repeat("i", 20000),
			key:            "secret",
			expectedResult: "/TF3CAITR7N8RuiIeCmHkd1AYfYZ7PgVdyp0TopHc2Y=",
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			a := assert.New(t)

			h := NewHMACTransformer(test.key)

			r, err := h.Transform(test.text)

			a.NoError(err)

			a.Equal(test.expectedResult, r)
		})
	}
}
