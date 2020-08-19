package zutils

func GetLast(values []string) string {
	l := len(values)
	if l <= 0 {
		return ""
	}
	return values[l-1]
}

func ConcatStringSlices(a []string, b []string) []string {
	if len(a) == 0 && len(b) == 0 {
		return nil
	}

	return append(append([]string{}, a...), b...)
}
