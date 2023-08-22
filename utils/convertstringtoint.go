package utils

import "strconv"

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
