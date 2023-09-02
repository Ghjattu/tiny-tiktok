package utils

import "strconv"

// ConvertInt64ToString converts a int64 list to a string list.
//
//	@param intList []int64
//	@return []string
//	@return error
func ConvertInt64ToString(intList []int64) ([]string, error) {
	stringList := make([]string, 0, len(intList))

	for _, intNum := range intList {
		str := strconv.FormatInt(intNum, 10)

		stringList = append(stringList, str)
	}

	return stringList, nil
}
