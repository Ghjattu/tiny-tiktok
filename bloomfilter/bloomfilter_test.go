package bloomfilter

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		ClearAll()

		Add(UserBloomFilter, 1)

		buf := make([]byte, binary.MaxVarintLen64)
		binary.PutVarint(buf, 1)

		assert.True(t, userBF.Test(buf))
	})
}

func TestCheckInt64Exist(t *testing.T) {
	ClearAll()
	Add(UserBloomFilter, 1)

	t.Run("exist in bloomfilter", func(t *testing.T) {
		res := CheckInt64Exist(UserBloomFilter, 1)

		assert.True(t, res)
	})

	t.Run("not exist in bloomfilter", func(t *testing.T) {
		res := CheckInt64Exist(UserBloomFilter, 2)

		assert.False(t, res)
	})

	t.Run("unknown bloomfilter", func(t *testing.T) {
		res := CheckInt64Exist(100, 1)

		assert.False(t, res)
	})
}

func TestSelectBloomFilter(t *testing.T) {
	t.Run("user bloomfilter", func(t *testing.T) {
		bf := selectBloomFilter(UserBloomFilter)

		assert.Equal(t, userBF, bf)
	})

	t.Run("video bloomfilter", func(t *testing.T) {
		bf := selectBloomFilter(VideoBloomFilter)

		assert.Equal(t, videoBF, bf)
	})

	t.Run("comment bloomfilter", func(t *testing.T) {
		bf := selectBloomFilter(CommentBloomFilter)

		assert.Equal(t, commentBF, bf)
	})

	t.Run("nil", func(t *testing.T) {
		bf := selectBloomFilter(100)

		assert.Nil(t, bf)
	})
}
