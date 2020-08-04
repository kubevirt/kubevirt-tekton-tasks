package output

import (
	"encoding/json"
	"fmt"
	"sigs.k8s.io/yaml"
)

type OutputType string

const (
	YamlOutput OutputType = "yaml"
	JsonOutput OutputType = "json"
)

func PrettyPrint(object interface{}, outputType OutputType) {
	switch outputType {
	case YamlOutput:
		outBytes, _ := yaml.Marshal(object)
		fmt.Print(string(outBytes))
	case JsonOutput:
		outBytes, _ := json.MarshalIndent(object, "", "    ")
		fmt.Println(string(outBytes))
	}
}
