package utilstest

import (
	. "github.com/onsi/gomega"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/logger"
	"go.uber.org/zap"
	"os"
)

func SetupTestSuite() {
	SetEnv(constants.OutOfClusterENV, "true")
	logger.InitLogger(zap.InfoLevel)
}

func TearDownSuite() {
	UnSetEnv(constants.OutOfClusterENV)
}

// logic should equal to utils.HandleExit func
func HandleTestExit(shouldPanic bool, shouldExitWithCode int, shouldExitWithMessage string) {
	if e := recover(); e != nil {
		if exit, ok := e.(utils.Exit); ok == true && exit.Soft {
			errMsg := exit.Error()
			if len(errMsg) > 0 && errMsg[len(errMsg)-1] != '\n' {
				errMsg += "\n"
			}
			Expect(errMsg).To(Equal(shouldExitWithMessage))
			Expect(exit.Code).To(Equal(shouldExitWithCode))
			Expect(shouldPanic).To(BeFalse(), "should not panic")
			return
		}
		Expect(shouldPanic).To(BeTrue(), "should panic")
	}
}

func SetEnv(key, value string) {
	Expect(os.Setenv(key, value)).Should(Succeed())
}

func UnSetEnv(key string) {
	Expect(os.Unsetenv(key)).Should(Succeed())
}
