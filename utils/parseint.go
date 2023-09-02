package utils

import "strconv"

// ParseInt64 parse string to int64.
//
//	@param s string
//	@return int32 "status_code"
//	@return string "status_msg"
//	@return int64 "int64 value"
func ParseInt64(s string) (int32, string, int64) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok {
			if numErr.Err == strconv.ErrSyntax {
				return 1, "invalid syntax", 0
			} else if numErr.Err == strconv.ErrRange {
				return 1, "the target value out of range", 0
			}
		} else {
			return 1, "unknown error", 0
		}
	}

	return 0, "", i
}
