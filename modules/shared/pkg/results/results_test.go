package results_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/env"
	results2 "github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/results"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	testSuffix = ".result.test"
	filenameA  = "file-a" + testSuffix
	filenameB  = "file-b" + testSuffix
	contentA   = "my content a"
	contentB   = "my content b\n"
)

var _ = Describe("Results", func() {
	Describe("Records results", func() {
		AfterEach(func() {
			for _, name := range []string{filenameA, filenameB} {
				_ = os.Remove(name) // allow not found
			}
		})
		table.DescribeTable("writes to correct file", func(results map[string]string) {
			Expect(results2.RecordResults(results)).Should(Succeed())
			files, err := filepath.Glob(env.GetTektonResultsDir() + "/*" + testSuffix)
			Expect(err).Should(Succeed())
			Expect(files).Should(HaveLen(len(results)))
			for filename, expectedContent := range results {
				content, err := ioutil.ReadFile(env.GetTektonResultsDir() + "/" + filename)
				Expect(err).Should(Succeed())
				Expect(string(content)).To(Equal(expectedContent))
			}
		},
			table.Entry("nil results", nil),
			table.Entry("no results", map[string]string{}),
			table.Entry("one result", map[string]string{
				filenameA: contentA,
			}),
			table.Entry("two results", map[string]string{
				filenameA: contentA,
				filenameB: contentB,
			}),
		)
	})

})
