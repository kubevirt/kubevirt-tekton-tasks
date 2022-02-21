package zutils

import (
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

func IsTrue(value string) bool {
	return strings.ToLower(value) == zconstants.True
}

func DecodeVM(template *templatev1.Template) (*kubevirtv1.VirtualMachine, int, error) {
	var vm *kubevirtv1.VirtualMachine
	vmIndex := -1
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(kubevirtv1.GroupVersion, &kubevirtv1.VirtualMachine{})
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(kubevirtv1.GroupVersion)
	for i, obj := range template.Objects {

		decoded, err := runtime.Decode(decoder, obj.Raw)
		if err != nil {
			return nil, vmIndex, err
		}
		done, ok := decoded.(*kubevirtv1.VirtualMachine)
		if ok {
			vm = done
			vmIndex = i
			break
		}
	}
	if vm == nil {
		return nil, vmIndex, zerrors.NewMissingRequiredError("no VM object found in the template")
	}
	return vm, vmIndex, nil
}
