package templates

import (
	"strconv"
	"strings"
	"unicode"
)

type textIDs []string

func (r textIDs) Len() int { return len(r) }
func (r textIDs) Less(i, j int) bool {
	verI := splitID(r[i])
	verJ := splitID(r[j])

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
	return true

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
