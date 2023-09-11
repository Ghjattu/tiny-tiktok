package bloomfilter

import (
	"encoding/binary"

	"github.com/bits-and-blooms/bloom/v3"
)

var (
	userBF    *bloom.BloomFilter
	videoBF   *bloom.BloomFilter
	commentBF *bloom.BloomFilter
)

func init() {
	userBF = bloom.NewWithEstimates(100, 0.01)
	videoBF = bloom.NewWithEstimates(100, 0.01)
	commentBF = bloom.NewWithEstimates(100, 0.01)
}

// Add adds a key to the bloom filter.
//
//	@param bloomFilterType uint
//	@param key int64
func Add(bloomFilterType uint, key int64) {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, key)

	bf := selectBloomFilter(bloomFilterType)
	bf.Add(buf)
}

// CheckInt64Exist checks if a key exists in the bloom filter.
//
//	@param bloomFilterType uint
//	@param key int64
//	@return bool
func CheckInt64Exist(bloomFilterType uint, key int64) bool {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, key)

	bf := selectBloomFilter(bloomFilterType)
	if bf == nil {
		return false
	}

	return bf.Test(buf)
}

// ClearAll clears all the bloom filters.
func ClearAll() {
	userBF.ClearAll()
	videoBF.ClearAll()
	commentBF.ClearAll()
}

// selectBloomFilter selects a bloom filter by the type.
//
//	@param bloomFilterType uint
//	@return *bloom.BloomFilter
func selectBloomFilter(bloomFilterType uint) *bloom.BloomFilter {
	switch bloomFilterType {
	case UserBloomFilter:
		return userBF
	case VideoBloomFilter:
		return videoBF
	case CommentBloomFilter:
		return commentBF
	default:
		return nil
	}
}
