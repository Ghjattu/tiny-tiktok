package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertInt64ToString(t *testing.T) {
	list := []int64{1, 2, 3}

	strList, err := ConvertInt64ToString(list)

	assert.Nil(t, err)
	assert.Equal(t, []string{"1", "2", "3"}, strList)
}
