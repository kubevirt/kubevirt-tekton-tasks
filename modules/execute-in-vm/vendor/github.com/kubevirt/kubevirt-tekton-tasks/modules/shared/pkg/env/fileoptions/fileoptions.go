package fileoptions

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"io/ioutil"
	"os"
	"path"
)

func ReadFileOption(output *string, optionPath string) error {
	result, err := ioutil.ReadFile(optionPath)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	path.Join()

	*output = string(result)

	return nil
}

func ReadFileOptionBool(output *bool, optionPath string) error {
	var tmp string
	if err := ReadFileOption(&tmp, optionPath); err != nil {
		return err
	}

	*output = zutils.IsTrue(tmp)

	return nil
}
