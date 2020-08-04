package results

import (
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/constants"
	"io/ioutil"
	"path/filepath"
)

func WriteResults(results map[string]string) error {
	resultsDir := constants.GetTektonResultsDir()

	for resKey, resVal := range results {
		filename := filepath.Join(resultsDir, resKey)
		err := ioutil.WriteFile(filename, []byte(resVal), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
