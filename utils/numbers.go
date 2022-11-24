package utils

import "strconv"

func FloatToString(precision int, f float64) string {
	return strconv.FormatFloat(f, 'f', precision, 64)
}

func IntToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StringToUint64(n string) (uint64, error) {
	return strconv.ParseUint(n, 10, 64)
}