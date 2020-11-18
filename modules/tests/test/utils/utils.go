package utils

import "strconv"

func IsBVersionHigher(a []int, b []int) bool {
	for idx, partA := range a {
		partB := 0
		if idx < len(b) {
			partB = b[idx]
		}
		if partA > partB {
			return false
		}
	}
	return true
}

func JoinIntSlice(input []int, sep string) (output string) {
	for idx, value := range input {
		num := strconv.Itoa(value)
		output += num
		if idx != len(input)-1 {
			output += sep
		}
	}

	return output
}

func ConvertStringSliceToInt(input []string) ([]int, error) {
	var output []int
	for _, value := range input {
		num, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		output = append(output, num)
	}

	return output, nil
}
