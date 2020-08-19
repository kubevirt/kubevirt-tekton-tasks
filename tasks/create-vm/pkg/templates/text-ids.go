package templates

import (
	"strconv"
	"strings"
	"unicode"
)

type textIDs []string

func (r textIDs) Len() int { return len(r) }
func (r textIDs) Less(i, j int) bool {
	valI := r[i]
	valJ := r[j]
	verI := splitID(valI)
	verJ := splitID(valJ)

	for k := 0; k < len(verI) || k < len(verJ); k++ {
		subVerI := 0
		subVerJ := 0
		if k < len(verI) {
			subVerI = verI[k]
		}
		if k < len(verJ) {
			subVerJ = verJ[k]
		}
		if subVerI > subVerJ {
			return false
		}
		if subVerJ > subVerI {
			return true
		}
	}

	isILesser := isLesser(valI)
	isJLesser := isLesser(valJ)

	if isILesser != isJLesser {
		return isILesser
	}

	return strings.Compare(valI, valJ) == -1

}
func (r textIDs) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func splitID(value string) []int {
	idParts := strings.FieldsFunc(value, func(r rune) bool {
		return !unicode.IsDigit(r)
	})
	var result []int
	for _, val := range idParts {
		if val != "" {
			if res, err := strconv.Atoi(val); err == nil {
				result = append(result, res)
			}
		}
	}

	return result
}

var lesserKeys = []string{"silverblue"}

func isLesser(key string) bool {
	for _, lesserPrefix := range lesserKeys {
		if strings.HasPrefix(key, lesserPrefix) {
			return true
		}
	}
	return false
}
