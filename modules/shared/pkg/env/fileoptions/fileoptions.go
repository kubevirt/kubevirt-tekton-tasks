package fileoptions

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
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

	*output = strings.ToLower(tmp) == "true"

	return nil
}
