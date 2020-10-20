package utils

import (
	"errors"
	"flag"
	"fmt"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"strings"
)

type TestScope string

const (
	ClusterScope   TestScope = "cluster"
	NamespaceScope TestScope = "namespace"
)

var deployNamespace string
var testNamespace string
var storageClass string
var kubeConfigPath string
var scope string
var debug string

type TestOptions struct {
	DeployNamespace string
	TestNamespace   string
	StorageClass    string
	KubeConfigPath  string
	Scope           TestScope
	Debug           bool
}

func init() {
	flag.StringVar(&deployNamespace, "deploy-namespace", "", "Namespace where to deploy the tasks and taskrun")
	flag.StringVar(&testNamespace, "test-namespace", "", "Namespace where to create the vm/dv resources")
	flag.StringVar(&storageClass, "storage-class", "", "Storage class to be used for creating test DVs/PVCs")
	flag.StringVar(&kubeConfigPath, "kubeconfig-path", "", "Path to the kubeconfig")
	flag.StringVar(&scope, "scope", "", "Scope of the tests. One of: cluster|namespace")
	flag.StringVar(&debug, "debug", "", "Debug keeps all the resources alive after the tests complete. One of: true|false")
}

func GetTestOptions() (*TestOptions, error) {
	flag.Parse()

	if deployNamespace == "" {
		return nil, errors.New("--deploy-namespace must be specified")
	}

	if testNamespace == "" {
		return nil, errors.New("--test-namespace must be specified")
	}

	var result = TestOptions{
		DeployNamespace: deployNamespace,
		TestNamespace:   testNamespace,
		StorageClass:    storageClass,
		Debug:           strings.ToLower(debug) == "true",
	}

	if kubeConfigPath != "" {
		result.KubeConfigPath = kubeConfigPath
	} else {
		kubeConfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
		if file, err := os.Stat(kubeConfigPath); err == nil && file.Mode().IsRegular() {
			result.KubeConfigPath = kubeConfigPath
		}
	}

	if scope == "" {
		result.Scope = NamespaceScope
	} else if TestScope(scope) == NamespaceScope || TestScope(scope) == ClusterScope {
		result.Scope = TestScope(scope)
	} else {
		return nil, fmt.Errorf("invalid scope, only %v or %v is allowed", ClusterScope, NamespaceScope)
	}

	return &result, nil
}
