package image_test

import (
	"os"
	"path/filepath"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/image"
)

var _ = Describe("Build", func() {
	DescribeTable("sets image platform metadata from architecture", func(architecture, expectedOS string) {
		diskPath := filepath.Join(GinkgoT().TempDir(), "disk.img")
		Expect(os.WriteFile(diskPath, []byte("disk content"), 0o600)).To(Succeed())

		config := v1.Config{Env: []string{"TEST=yes"}}
		containerImage, err := image.Build(diskPath, config, architecture)
		Expect(err).NotTo(HaveOccurred())

		configFile, err := containerImage.ConfigFile()
		Expect(err).NotTo(HaveOccurred())
		Expect(configFile.Config).To(Equal(config))
		Expect(configFile.Architecture).To(Equal(architecture))
		Expect(configFile.OS).To(Equal(expectedOS))
	},
		Entry("VM architecture is set", "arm64", "linux"),
		Entry("VM architecture is empty", "", ""),
	)
})
