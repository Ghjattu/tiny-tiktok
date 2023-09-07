package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertInterfaceToInt64(t *testing.T) {
	t.Run("interface type is float64", func(t *testing.T) {
		value := interface{}(float64(1))

		got, err := ConvertInterfaceToInt64(value)

		assert.Equal(t, int64(1), got)
		assert.Nil(t, err)
	})

	t.Run("interface type is int", func(t *testing.T) {
		value := interface{}(int(1))

		got, err := ConvertInterfaceToInt64(value)

		assert.Equal(t, int64(1), got)
		assert.Nil(t, err)
	})

	t.Run("unhandled type", func(t *testing.T) {
		value := interface{}("1")

		got, err := ConvertInterfaceToInt64(value)

		assert.Equal(t, int64(0), got)
		assert.NotNil(t, err)
	})
}

func TestConvertInterfaceToInt64Slice(t *testing.T) {
	t.Run("interface type is float64 slice", func(t *testing.T) {
		value := interface{}([]float64{1, 2, 3})

		got, err := ConvertInterfaceToInt64Slice(value)

		assert.Equal(t, []int64{1, 2, 3}, got)
		assert.Nil(t, err)
	})

	t.Run("interface type is int64 slice", func(t *testing.T) {
		value := interface{}([]int64{1, 2, 3})

		got, err := ConvertInterfaceToInt64Slice(value)

		assert.Equal(t, []int64{1, 2, 3}, got)
		assert.Nil(t, err)
	})

	t.Run("unhandled type", func(t *testing.T) {
		value := interface{}("1")

		got, err := ConvertInterfaceToInt64Slice(value)

		assert.Equal(t, []int64{}, got)
		assert.NotNil(t, err)
	})
}
