package testoptions

import (
	"errors"
	"flag"
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"strings"
)

var deployNamespace string
var testNamespace string
var storageClass string
var kubeConfigPath string
var scope string
var debug string
var isOpenshift string

type TestOptions struct {
	DeployNamespace string
	TestNamespace   string
	StorageClass    string
	KubeConfigPath  string
	TestScope       constants.TestScope
	EnvScope        constants.EnvScope
	Debug           bool

	CommonTemplatesVersion string

	targetNamespaces map[constants.TargetNamespace]string
}

func init() {
	flag.StringVar(&deployNamespace, "deploy-namespace", "", "Namespace where to deploy the tasks and taskrun")
	flag.StringVar(&testNamespace, "test-namespace", "", "Namespace where to create the vm/dv resources")
	flag.StringVar(&storageClass, "storage-class", "", "Storage class to be used for creating test DVs/PVCs")
	flag.StringVar(&kubeConfigPath, "kubeconfig-path", "", "Path to the kubeconfig")
	flag.StringVar(&isOpenshift, "is-openshift", "", "Set to true if running on openshift. One of: true|false")
	flag.StringVar(&scope, "scope", "", "Scope of the tests. One of: cluster|namespace")
	flag.StringVar(&debug, "debug", "", "Debug keeps all the resources alive after the tests complete. One of: true|false")
}

func InitTestOptions(testOptions *TestOptions) error {
	flag.Parse()

	if deployNamespace == "" {
		return errors.New("--deploy-namespace must be specified")
	}

	if testNamespace == "" {
		return errors.New("--test-namespace must be specified")
	}

	if scope == "" {
		testOptions.TestScope = constants.NamespaceTestScope
	} else if constants.TestScope(scope) == constants.NamespaceTestScope || constants.TestScope(scope) == constants.ClusterTestScope {
		testOptions.TestScope = constants.TestScope(scope)
	} else {
		return fmt.Errorf("invalid scope, only %v or %v is allowed", constants.ClusterTestScope, constants.NamespaceTestScope)
	}

	if kubeConfigPath != "" {
		testOptions.KubeConfigPath = kubeConfigPath
	} else {
		kubeConfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
		if file, err := os.Stat(kubeConfigPath); err == nil && file.Mode().IsRegular() {
			testOptions.KubeConfigPath = kubeConfigPath
		}
	}

	testOptions.DeployNamespace = deployNamespace
	testOptions.TestNamespace = testNamespace
	testOptions.StorageClass = storageClass
	if strings.ToLower(isOpenshift) == "true" {
		testOptions.EnvScope = constants.OpenshiftEnvScope
	} else {
		testOptions.EnvScope = constants.KubernetesEnvScope
	}
	testOptions.Debug = strings.ToLower(debug) == "true"

	testOptions.targetNamespaces = testOptions.resolveNamespaces()

	return nil
}

func (f *TestOptions) resolveNamespaces() map[constants.TargetNamespace]string {
	var systemNS string
	if f.DeployNamespace == "tekton-pipelines" {
		systemNS = "default"
	} else {
		systemNS = "tekton-pipelines"
	}

	return map[constants.TargetNamespace]string{
		constants.DeployTargetNS: f.DeployNamespace,
		constants.TestTargetNS:   f.TestNamespace,
		constants.SystemTargetNS: systemNS,
	}
}

func (f *TestOptions) ResolveNamespace(namespace constants.TargetNamespace) string {
	ns := f.targetNamespaces[namespace]

	if ns != "" {
		return ns
	}

	return f.TestNamespace
}
