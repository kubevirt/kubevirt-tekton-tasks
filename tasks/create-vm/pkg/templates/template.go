package templates

import (
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	"strings"
)

// Get label whose key starts with labelPrefix and has value true
// returns key, val if such label is found
func GetFlagLabelByPrefix(template *templatev1.Template, labelPrefix string) (string, string) {
	if labels := template.GetLabels(); labels != nil {
		for key, val := range labels {
			if strings.HasPrefix(key, labelPrefix) && val == "true" {
				return key, val
			}
		}
	}
	return "", ""
}

func DecodeVM(template *templatev1.Template) (*kubevirtv1.VirtualMachine, error) {
	var vm = &kubevirtv1.VirtualMachine{}

	for _, obj := range template.Objects {
		decoder := kubevirtv1.Codecs.UniversalDecoder(kubevirtv1.GroupVersion)
		decoded, err := runtime.Decode(decoder, obj.Raw)
		if err != nil {
			return nil, err
		}
		done, ok := decoded.(*kubevirtv1.VirtualMachine)
		if ok {
			vm = done
			break
		}
	}
	return vm, nil
}
