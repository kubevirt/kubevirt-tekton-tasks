package results

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"io/ioutil"
	"path/filepath"
)

func RecordResults(results map[string]string) error {
	if results == nil || len(results) == 0 {
		return nil
	}

	resultsDir := env.GetTektonResultsDir()

	for resKey, resVal := range results {
		filename := filepath.Join(resultsDir, resKey)
		err := ioutil.WriteFile(filename, []byte(resVal), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
