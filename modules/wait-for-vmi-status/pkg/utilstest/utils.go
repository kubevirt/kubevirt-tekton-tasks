package utilstest

import (
	log2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utils/log"
	"github.com/onsi/gomega"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func SetupTestSuite() {
	log2.InitLogger(zap.InfoLevel)
}

func GetRequirement(key string, op selection.Operator, vals []string) labels.Requirement {
	requirement, err := labels.NewRequirement(key, op, vals)
	gomega.Expect(err).Should(gomega.Succeed())
	return *requirement
}
