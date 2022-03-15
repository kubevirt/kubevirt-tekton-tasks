package template

import (
	"fmt"

	v1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

func GetVM(template *v1.Template) *kubevirtv1.VirtualMachine {
	for _, obj := range template.Objects {
		decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(kubevirtv1.GroupVersion)
		decoded, err := runtime.Decode(decoder, obj.Raw)
		if err != nil {
			panic(err)
		}
		vm, ok := decoded.(*kubevirtv1.VirtualMachine)
		if ok {
			return vm
		}
	}
	panic("no VM found")
}

func TemplateParam(key, value string) string {
	return fmt.Sprintf("%v:%v", key, value)
}
