package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSHA1(t *testing.T) {
	tests := []struct {
		text     string
		expected string
	}{
		{
			text:     "hello world :(",
			expected: "c0d061c55eddd1160d06cb04d82bb1f037baf986",
		},
		{
			text:     "سلام جهان :(",
			expected: "0baf9b846efea5c467124859dfa4baf0dfd1fcbd",
		},
		{
			text:     "abc سلام ۱۲۳ #تست",
			expected: "f6f902300051d0e260278151b187becae244abfe",
		},
		{
			text:     `aa bb "" ''`,
			expected: "637b38694e9cd8d87b7de2b3bd1924fa832590a0",
		},
	}
	for i := range tests {
		test := tests[i]

		t.Run(test.text, func(t *testing.T) {
			a := assert.New(t)

			sha1 := SHA1Transformer{}

			r, err := sha1.Transform(test.text)

			a.NoError(err)

			a.Equal(test.expected, r)
		})
	}
}
