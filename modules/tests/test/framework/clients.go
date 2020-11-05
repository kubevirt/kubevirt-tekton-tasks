package framework

import (
	"fmt"
	templatev1 "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	tknclientversioned "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	tknclientv1beta1 "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1beta1"
	"k8s.io/client-go/kubernetes"
	kubeclient "k8s.io/client-go/kubernetes"
	kubeclientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cdicliv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"

	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
)

type Clients struct {
	RestConfig *rest.Config

	K8sClient      *kubernetes.Clientset
	CoreV1Client   kubeclientcorev1.CoreV1Interface
	TknClient      tknclientv1beta1.TektonV1beta1Interface
	TemplateClient *templatev1.TemplateV1Client
	CdiClient      *cdicliv1beta1.CdiV1beta1Client
	KubevirtClient kubevirtcliv1.KubevirtClient
}

func InitClients(clients *Clients, testOptions *testoptions.TestOptions) error {
	restConf, err := rest.InClusterConfig()
	if err != nil {
		restConf, err = clientcmd.BuildConfigFromFlags("", testOptions.KubeConfigPath)
	}
	if err != nil {
		return fmt.Errorf("could not load KUBECONFIG: %v", err)
	}

	k8sClient, err := kubeclient.NewForConfig(restConf)
	if err != nil {
		return fmt.Errorf("could not load K8sClient: %v", err)
	}

	tknClientset, err := tknclientversioned.NewForConfig(restConf)
	if err != nil {
		return fmt.Errorf("could not create TknClient: %v", err)
	}

	templateClient, err := templatev1.NewForConfig(restConf)

	if err != nil {
		return fmt.Errorf("could not create TemplateClient: %v", err)
	}

	cdiClient, err := cdicliv1beta1.NewForConfig(restConf)
	if err != nil {
		return fmt.Errorf("could not create CdiClient: %v", err)
	}

	kubevirtClient, err := kubevirtcliv1.GetKubevirtClientFromRESTConfig(restConf)
	if err != nil {
		return fmt.Errorf("could not create KubevirtClient: %v", err)
	}

	clients.RestConfig = restConf
	clients.K8sClient = k8sClient
	clients.CoreV1Client = k8sClient.CoreV1()
	clients.TknClient = tknClientset.TektonV1beta1()
	clients.TemplateClient = templateClient
	clients.CdiClient = cdiClient
	clients.KubevirtClient = kubevirtClient

	return nil
}
