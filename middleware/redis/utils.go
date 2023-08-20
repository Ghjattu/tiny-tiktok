package redis

import (
	"math/rand"
	"strconv"
	"time"
)

// RandomDay returns a random time.Duration between 24h and 48h.
//
//	@return time.Duration
func RandomDay() time.Duration {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return time.Hour * time.Duration(24+r.Int31n(25))
}

// HashIncrBy calls HIncrBy command of redis.
//
//	@param key string
//	@param field string
//	@param incr int64
//	@return int64 "the value at field after the increment operation"
//	@return string "status message"
//	@return error
func HashIncrBy(key, field string, incr int64) (int64, string, error) {
	if Rdb.Exists(Ctx, key).Val() == 1 {
		res := Rdb.HIncrBy(Ctx, key, field, incr)

		return res.Val(), "", res.Err()
	}

	return 0, "key does not exist", nil
}

// ConvertStringToInt64 converts a string list to a int64 list.
//
//	@param stringList []string
//	@return []int64
//	@return error
func ConvertStringToInt64(stringList []string) ([]int64, error) {
	intList := make([]int64, 0, len(stringList))

	for _, str := range stringList {
		intNum, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}

		intList = append(intList, intNum)
	}

	return intList, nil
}
