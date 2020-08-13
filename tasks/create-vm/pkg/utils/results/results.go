package results

import (
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"path/filepath"
)

func RecordResults(results map[string]string) error {
	if results == nil || len(results) == 0 {
		return nil
	}

	logger.GetLogger().Debug("recording results", zap.Reflect("results", results))

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
