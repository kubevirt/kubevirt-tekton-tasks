package taskmanifests_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTaskManifests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task Manifests Suite")
}
