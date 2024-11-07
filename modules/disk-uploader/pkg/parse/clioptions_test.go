package parse_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/zap/zapcore"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/parse"
)

var _ = Describe("CLIOptions", func() {
	const (
		expectedExportSourceKind      = "vm"
		expectedExportSourceNamespace = "test-namespace"
		expectedExportSourceName      = "test-vmexport"
		expectedVolumeName            = "test-volume"
		expectedImageDestination      = "quay.io/kubevirt/example"
	)

	DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
		Expect(options.Init()).To(MatchError(expectedErrMessage))
	},
		Entry("no export-source-kind", "export-source-kind param has to be specified",
			&parse.CLIOptions{}),
		Entry("no export-source-name", "export-source-name param has to be specified",
			&parse.CLIOptions{ExportSourceKind: expectedExportSourceKind}),
		Entry("no volume-name", "volume-name param has to be specified",
			&parse.CLIOptions{ExportSourceKind: expectedExportSourceKind, ExportSourceName: expectedExportSourceName}),
		Entry("no image-destination", "image-destination param has to be specified",
			&parse.CLIOptions{ExportSourceKind: expectedExportSourceKind, ExportSourceName: expectedExportSourceName, VolumeName: expectedVolumeName}),
	)

	Context("valid cli options", func() {
		It("should succeed with yaml output", func() {
			options := &parse.CLIOptions{
				ExportSourceKind:      expectedExportSourceKind,
				ExportSourceNamespace: expectedExportSourceNamespace,
				ExportSourceName:      expectedExportSourceName,
				VolumeName:            expectedVolumeName,
				ImageDestination:      expectedImageDestination,
				PushTimeout:           60,
				Debug:                 true,
			}
			Expect(options.Init()).To(Succeed())
		})

		It("Init should trim spaces", func() {
			options := &parse.CLIOptions{
				ExportSourceKind:      " " + expectedExportSourceKind + " ",
				ExportSourceNamespace: " " + expectedExportSourceNamespace + " ",
				ExportSourceName:      " " + expectedExportSourceName + " ",
				VolumeName:            " " + expectedVolumeName + " ",
				ImageDestination:      " " + expectedImageDestination + " ",
				PushTimeout:           60,
				Debug:                 true,
			}

			Expect(options.Init()).To(Succeed())
			Expect(options.ExportSourceKind).To(Equal(expectedExportSourceKind), "ExportSourceKind should equal")
			Expect(options.ExportSourceNamespace).To(Equal(expectedExportSourceNamespace), "ExportSourceNamespace should equal")
			Expect(options.ExportSourceName).To(Equal(expectedExportSourceName), "ExportSourceName should equal")
			Expect(options.VolumeName).To(Equal(expectedVolumeName), "VolumeName should equal")
			Expect(options.ImageDestination).To(Equal(expectedImageDestination), "ImageDestination should equal")
		})

		It("GetDebugLevel should return correct debug level when Debug is true", func() {
			options := &parse.CLIOptions{Debug: true}
			Expect(options.GetDebugLevel()).To(Equal(zapcore.DebugLevel), "level should equal")
		})

		It("GetDebugLevel should return correct info level when Debug is false", func() {
			options := &parse.CLIOptions{Debug: false}
			Expect(options.GetDebugLevel()).To(Equal(zapcore.InfoLevel), "level should equal")
		})
	})
})
