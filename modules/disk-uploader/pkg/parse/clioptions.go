package parse

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap/zapcore"
)

const (
	defaultPushTimeout = "120"
)

type CLIOptions struct {
	ExportSourceKind      string `arg:"--export-source-kind" help:"Specify the export source kind (vm, vmsnapshot, pvc)"`
	ExportSourceNamespace string `arg:"--export-source-namespace" help:"Namespace of the export source"`
	ExportSourceName      string `arg:"--export-source-name" help:"Name of the export source"`
	VolumeName            string `arg:"--volumename" help:"Name of the volume (if source kind is 'pvc', then volume name is equal to source name)"`
	ImageDestination      string `arg:"--imagedestination" help:"Destination of the image in container registry"`
	PushTimeout           string `arg:"--pushtimeout" help:"Push timeout of container disk to registry"`
	Debug                 bool   `arg:"--debug" help:"Sets DEBUG log level"`
}

func (c *CLIOptions) GetExportSourceKind() string {
	return c.ExportSourceKind
}

func (c *CLIOptions) GetExportSourceNamespace() string {
	return c.ExportSourceNamespace
}

func (c *CLIOptions) GetExportSourceName() string {
	return c.ExportSourceName
}

func (c *CLIOptions) GetVolumeName() string {
	return c.VolumeName
}

func (c *CLIOptions) GetImageDestination() string {
	return c.ImageDestination
}

func (c *CLIOptions) GetPushTimeout() int {
	timeout, err := strconv.Atoi(c.PushTimeout)
	if err != nil {
		log.Fatalf("Invalid push timeout value: %v", err)
	}

	return timeout
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) trimSpaces() {
	variables := []*string{
		&c.ExportSourceKind,
		&c.ExportSourceNamespace,
		&c.ExportSourceName,
		&c.VolumeName,
		&c.ImageDestination,
	}
	for _, strVariablePtr := range variables {
		*strVariablePtr = strings.TrimSpace(*strVariablePtr)
	}
}

func (c *CLIOptions) assertValidParams() error {
	if c.ExportSourceKind == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", "export-source-kind")
	}

	if c.ExportSourceName == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", "export-source-name")
	}

	if c.VolumeName == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", "volume-name")
	}

	if c.ImageDestination == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", "image-destination")
	}
	return nil
}

func (c *CLIOptions) setValues() error {
	namespace := os.Getenv("POD_NAMESPACE")
	if namespace != "" {
		c.ExportSourceNamespace = namespace
	}

	if c.PushTimeout == "" || c.PushTimeout == "0" {
		c.PushTimeout = defaultPushTimeout
	}
	return nil
}

func (c *CLIOptions) Init() error {
	c.trimSpaces()

	if err := c.assertValidParams(); err != nil {
		return err
	}

	c.setValues()
	return nil
}
