package zutils_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Maps", func() {
	Describe("ExtractKeysAndValuesByLastKnownKey", func() {
		table.DescribeTable("returns error", func(input []string, expectedError string) {
			result, err := zutils.ExtractKeysAndValuesByLastKnownKey(input, ":")
			Expect(result).Should(BeNil())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal(expectedError))
		},
			table.Entry("starts with key with no separator", []string{"key", "key2:val2", "key3:val3"}, "no key found before \"key\"; pair should be in \"KEY:VAL\" format"),
			table.Entry("starts with separator", []string{":key", "key2:val2", "key3:val3"}, "no key found before \":key\"; pair should be in \"KEY:VAL\" format"),
			table.Entry("starts with empty and key with no separator", []string{"", "key", "key2:val2", "key3:val3"}, "no key found before \"key\"; pair should be in \"KEY:VAL\" format"),
			table.Entry("missing key", []string{"key:val", "key2:val2", ":val3"}, "no key found before \":val3\"; pair should be in \"KEY:VAL\" format"),
		)

		table.DescribeTable("returns error", func(input []string, expected map[string]string) {
			result, err := zutils.ExtractKeysAndValuesByLastKnownKey(input, ":")
			Expect(err).Should(Succeed())
			Expect(result).To(Equal(expected))
		},
			table.Entry("basic", []string{"key:val", "key2:val2", "key3:val3"}, map[string]string{
				"key":  "val",
				"key2": "val2",
				"key3": "val3",
			}),
			table.Entry("advanced", []string{"key:val=515", "key2:val2 two", "key3:val3 three"}, map[string]string{
				"key":  "val=515",
				"key2": "val2 two",
				"key3": "val3 three",
			}),
			table.Entry("mixed", []string{"key:val=515:5:6 7", "", "hello", "world", "", "!", "key2:val2 two", "key3:val3 three", "four", "five  six", "seven"}, map[string]string{
				"key":  "val=515:5:6 7 hello world !",
				"key2": "val2 two",
				"key3": "val3 three four five  six seven",
			}),
		)
	})
})
