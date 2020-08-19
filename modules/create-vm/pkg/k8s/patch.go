package k8s

import (
	"encoding/json"
	"github.com/mattbaird/jsonpatch"
)

func CreatePatch(before interface{}, after interface{}) ([]byte, error) {
	beforeJson, err := json.Marshal(before)
	if err != nil {
		return nil, err
	}

	afterJson, err := json.Marshal(after)
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.CreatePatch(beforeJson, afterJson)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(patch, "", "  ")
}
