package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringToInt64(t *testing.T) {
	t.Run("convert successfully", func(t *testing.T) {
		list := []string{"1", "2", "3"}

		res, err := ConvertStringToInt64(list)

		assert.Equal(t, []int64{1, 2, 3}, res)
		assert.Nil(t, err)
	})

	t.Run("convert failed", func(t *testing.T) {
		list := []string{"1", "2", "test"}

		res, err := ConvertStringToInt64(list)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
