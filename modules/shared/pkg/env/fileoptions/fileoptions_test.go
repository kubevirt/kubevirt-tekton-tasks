package fileoptions_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env/fileoptions"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	i18nContent = "/[;]!@#$#^^&^&*%ščřžýýá15adf\n\tげんきですか？{\n"
)

var _ = Describe("Fileoptions", func() {

	DescribeTable("should read file option", func(content string, expectedContent interface{}) {
		By("prepare initial data and defer temp file removal")
		tempFile, err := ioutil.TempFile("", "shared-file-options-test")
		Expect(err).Should(Succeed())
		defer os.Remove(tempFile.Name())

		_, err = tempFile.Write([]byte(content))
		Expect(err).Should(Succeed())

		err = tempFile.Close()
		Expect(err).Should(Succeed())

		By("test by type")
		switch expected := expectedContent.(type) {
		case bool:
			var result bool
			err := fileoptions.ReadFileOptionBool(&result, tempFile.Name())
			Expect(err).Should(Succeed())
			Expect(result).Should(Equal(expected))
		default:
			var result string
			err := fileoptions.ReadFileOption(&result, tempFile.Name())
			Expect(err).Should(Succeed())
			Expect(result).Should(Equal(expected))
		}
	},
		Entry("False", "false", false),
		Entry("Bad", "falzee", false),
		Entry("UpperCase", "FALSE", false),
		Entry("Partially UpperCase", "FAlsE", false),

		Entry("True", "true", true),
		Entry("UpperCase", "TRUE", true),
		Entry("Partially UpperCase", "True", true),

		Entry("basic content", "hello world", "hello world"),
		Entry("i18n content", i18nContent, i18nContent),
	)

	DescribeTable("should return default when file does not exist", func(expectedContent interface{}) {
		nonExistentFile := path.Join(os.TempDir(), "non-existent", fmt.Sprintf("non-existent-%v", rand.Float64()))
		switch expected := expectedContent.(type) {
		case bool:
			var result bool
			err := fileoptions.ReadFileOptionBool(&result, nonExistentFile)
			Expect(err).Should(Succeed())
			Expect(result).Should(Equal(expected))
		default:
			var result string
			err := fileoptions.ReadFileOption(&result, nonExistentFile)
			Expect(err).Should(Succeed())
			Expect(result).Should(Equal(expected))
		}
	},
		Entry("empty bool", false),
		Entry("empty string", ""),
	)
})
