package parse

import (
	"bytes"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-data-object/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

const (
	dataObjectManifestOptionName = "data-object-manifest"
	objectKindOptionName         = "delete-object-kind"
	nameOptionName               = "delete-object-name"
)

type CLIOptions struct {
	DataObjectManifest  string            `arg:"--data-object-manifest,env:DATA_OBJECT_MANIFEST" placeholder:"MANIFEST" help:"YAML manifest of a data object to be created (can be set by DATA_OBJECT_MANIFEST env variable)."`
	DataObjectNamespace string            `arg:"--data-object-namespace,env:DATA_OBJECT_NAMESPACE" placeholder:"NAMESPACE" help:"Namespace where to create the data object (can be set by DATA_OBJECT_NAMESPACE env variable)."`
	WaitForSuccess      string            `arg:"--wait-for-success,env:WAIT_FOR_SUCCESS" help:"Set to \"true\" or \"false\" if container should wait for Ready condition of a DataVolume (can be set by WAIT_FOR_SUCCESS env variable)."`
	DeleteObjectName    string            `arg:"--delete-object-name,env:DELETE_OBJECT_NAME" help:"Name of the data object to delete. This parameter is used only for Delete operation."`
	DeleteObject        string            `arg:"--delete-object,env:DELETE_OBJECT" help:"Delete data object with given name. Parameters name, object-kind have to be defined."`
	DeleteObjectKind    string            `arg:"--delete-object-kind,env:DELETE_OBJECT_KIND" help:"Kind of the data object to delete. This parameter is used only for Delete operation."`
	AllowReplace        string            `arg:"--allow-replace,env:ALLOW_REPLACE" placeholder:"false" help:"Allow replacing an already existing data object (same combination name/namespace). Allowed values true/false (can be set by ALLOW_REPLACE env variable)."`
	SetOwnerReference   string            `arg:"--set-owner-reference,env:SET_OWNER_REFERENCE" placeholder:"false" help:"Set owner reference to the new object created by the task run pod. Allowed values true/false"`
	Output              output.OutputType `arg:"-o" placeholder:"FORMAT" help:"Output format. One of: yaml|json"`
	Debug               bool              `arg:"--debug" help:"Sets DEBUG log level"`

	unstructuredDataObject unstructured.Unstructured
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) GetDataObjectManifest() string {
	return c.DataObjectManifest
}

func (c *CLIOptions) GetDataObjectNamespace() string {
	return c.DataObjectNamespace
}

func (c *CLIOptions) GetWaitForSuccess() bool {
	return c.WaitForSuccess == "true"
}

func (c *CLIOptions) GetAllowReplace() bool {
	return c.AllowReplace == "true"
}

func (c *CLIOptions) GetSetOwnerReferenceValue() bool {
	return c.SetOwnerReference == "true"
}

func (c *CLIOptions) GetDeleteObject() bool {
	return c.DeleteObject == "true"
}

func (c *CLIOptions) GetObjectKind() string {
	return c.DeleteObjectKind
}

func (c *CLIOptions) GetName() string {
	return c.DeleteObjectName
}

func (c *CLIOptions) GetUnstructuredDataObject() unstructured.Unstructured {
	return c.unstructuredDataObject
}

func (c *CLIOptions) Init() error {
	c.trimSpaces()

	if err := c.assertValidParams(); err != nil {
		return err
	}

	if !c.GetDeleteObject() {
		if err := c.assertValidTypes(); err != nil {
			return err
		}
	}

	if err := c.setValues(); err != nil {
		return err
	}

	return nil
}

func (c *CLIOptions) setValues() error {
	if c.GetDataObjectNamespace() == "" {
		unstructuredDataObject := c.GetUnstructuredDataObject()
		if unstructuredDataObject.GetNamespace() != "" {
			c.DataObjectNamespace = unstructuredDataObject.GetNamespace()
		} else {
			activeNamespace, err := env.GetActiveNamespace()
			if err != nil {
				return zerrors.NewMissingRequiredError("can't get active namespace: %v", err.Error())
			}
			c.DataObjectNamespace = activeNamespace
		}
	}

	return nil
}

func (c *CLIOptions) trimSpaces() {
	for _, strVariablePtr := range []*string{&c.DataObjectManifest, &c.DataObjectNamespace, &c.WaitForSuccess} {
		*strVariablePtr = strings.TrimSpace(*strVariablePtr)
	}
}

func (c *CLIOptions) assertValidParams() error {
	if c.DeleteObject == "true" {
		if c.DeleteObjectKind == "" {
			return zerrors.NewMissingRequiredError("%s param has to be specified", objectKindOptionName)
		}

		if c.DeleteObjectName == "" {
			return zerrors.NewMissingRequiredError("%s param has to be specified", nameOptionName)
		}

		if c.DeleteObjectKind != constants.DataVolumeKind && c.DeleteObjectKind != constants.DataSourceKind && c.DeleteObjectKind != constants.PVCKind {
			return zerrors.NewMissingRequiredError("%s param has to have values %s or %s", objectKindOptionName, constants.DataVolumeKind, constants.DataSourceKind)
		}
		return nil
	}

	if c.DataObjectManifest == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", dataObjectManifestOptionName)
	}

	return nil
}

func (c *CLIOptions) assertValidTypes() error {
	if err := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(c.DataObjectManifest)), 1024).Decode(&c.unstructuredDataObject); err != nil {
		return zerrors.NewMissingRequiredError("could not read data object manifest: %v", err.Error())
	}

	if c.unstructuredDataObject.GroupVersionKind().Group != cdiv1beta1.SchemeGroupVersion.Group ||
		(c.unstructuredDataObject.GetKind() != constants.DataVolumeKind && c.unstructuredDataObject.GetKind() != constants.DataSourceKind) {
		return zerrors.NewMissingRequiredError("could not identify data object, wrong group or kind")
	}

	if !output.IsOutputType(string(c.Output)) {
		return zerrors.NewMissingRequiredError("%v is not a valid output type", c.Output)
	}

	return nil
}
