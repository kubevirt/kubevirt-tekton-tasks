package tekton

import beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

func TaskResultsToMap(results []beta1.TaskRunResult) map[string]string {
	mappedResult := make(map[string]string, len(results))

	for _, result := range results {
		mappedResult[result.Name] = result.Value
	}

	return mappedResult
}

func PipelineResultsToMap(results []beta1.PipelineRunResult) map[string]string {
	mappedResult := make(map[string]string, len(results))

	for _, result := range results {
		mappedResult[result.Name] = result.Value
	}

	return mappedResult
}
