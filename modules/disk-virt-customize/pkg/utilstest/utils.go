package utilstest

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utils/log"
	"go.uber.org/zap"
)

func SetupTestSuite() {
	log.InitLogger(zap.InfoLevel)
}
