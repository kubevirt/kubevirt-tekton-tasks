package requirements

import (
	"bytes"
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/util/jsonpath"
	"strings"
)

func ObjectToLabelsLookup(obj interface{}, requirements labels.Requirements) (labels.Labels, error) {
	result := make(labels.Set, len(requirements))
	for _, requirement := range requirements {
		jsonPath := jsonpath.New("requirements")

		path := requirement.Key()

		if !strings.HasPrefix(path, "{") {
			path = "{." + path + "}"
		}

		err := jsonPath.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("cannot parse jsonpath %v: %v", path, err)
		}

		buf := new(bytes.Buffer)
		err = jsonPath.Execute(buf, obj)
		if err == nil {
			result[requirement.Key()] = buf.String()
		}
	}

	return result, nil
}
