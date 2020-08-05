package utils

func GetLast(values []string) string {
	l := len(values)
	if l <= 0 {
		return ""
	}
	return values[l-1]
}

func ConcatStringSlices(a []string, b []string) []string {
	return append(append([]string{}, a...), b...)
}
