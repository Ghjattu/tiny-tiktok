package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInt64(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		s := ""

		statusCode, statusMsg, _ := ParseInt64(s)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "invalid syntax", statusMsg)
	})

	t.Run("non number string", func(t *testing.T) {
		s := "abc"

		statusCode, statusMsg, _ := ParseInt64(s)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "invalid syntax", statusMsg)
	})

	t.Run("out of range string", func(t *testing.T) {
		s := "922337203685477580888888888"

		statusCode, statusMsg, _ := ParseInt64(s)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "the target value out of range", statusMsg)
	})

	t.Run("parse successfully", func(t *testing.T) {
		s := "123"

		statusCode, statusMsg, i := ParseInt64(s)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "", statusMsg)
		assert.Equal(t, int64(123), i)
	})
}
