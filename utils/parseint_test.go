package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIntWithEmptyString(t *testing.T) {
	s := ""

	statusCode, statusMsg, _ := ParseInt64(s)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "invalid syntax", statusMsg)
}

func TestParseIntWithNonNumberString(t *testing.T) {
	s := "abc"

	statusCode, statusMsg, _ := ParseInt64(s)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "invalid syntax", statusMsg)
}

func TestParseIntWithOutOfRangeString(t *testing.T) {
	s := "922337203685477580888888888"

	statusCode, statusMsg, _ := ParseInt64(s)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the target value out of range", statusMsg)
}

func TestParseIntWithValidString(t *testing.T) {
	s := "123"

	statusCode, statusMsg, i := ParseInt64(s)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "", statusMsg)
	assert.Equal(t, int64(123), i)
}
