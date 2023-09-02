package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateString(t *testing.T) {
	s := GenerateRandomString(3)

	assert.Equal(t, 3, len(s))
	for i := 0; i < len(s); i++ {
		assert.True(t, s[i] >= 'a' && s[i] <= 'z')
	}
}
