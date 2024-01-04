package testconstants

import (
	"strings"

	k8srand "k8s.io/apimachinery/pkg/util/rand"
)

func TestRandomName(name string) string {
	// convert Full Test description into ID
	id := k8srand.String(5)

	return strings.Join([]string{name, id}, "-")
}
