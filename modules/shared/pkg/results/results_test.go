package results_test

import (
	results2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/results"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
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
		var tempDir string

		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "test-record-results-")
			Expect(err).Should(Succeed())
		})
		AfterEach(func() {
			_ = os.Remove(tempDir) // allow not found
		})

		table.DescribeTable("writes to correct file", func(results map[string]string) {
			Expect(results2.RecordResultsIn(tempDir, results)).Should(Succeed())
			files, err := filepath.Glob(filepath.Join(tempDir, "*"+testSuffix))
			Expect(err).Should(Succeed())
			Expect(files).Should(HaveLen(len(results)))
			for filename, expectedContent := range results {
				content, err := ioutil.ReadFile(filepath.Join(tempDir, filename))
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

		It("recordResults works without destination and results", func() {
			Expect(results2.RecordResults(nil)).Should(Succeed())
		})
	})

})
