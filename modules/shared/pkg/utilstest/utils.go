package utilstest

import (
	"os"

	. "github.com/onsi/gomega"
)

func SetEnv(key, value string) {
	Expect(os.Setenv(key, value)).Should(Succeed())
}

func UnSetEnv(key string) {
	Expect(os.Unsetenv(key)).Should(Succeed())
}
