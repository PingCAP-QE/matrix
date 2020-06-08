package utils

func IsInt(floatValue float64) bool {
	return floatValue == float64(int(floatValue))
}
