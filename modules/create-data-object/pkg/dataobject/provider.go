package dataobject

import (
	"context"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-data-object/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-data-object/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	cdiclientv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"
)

type dataObjectProvider struct {
	client cdiclientv1beta1.CdiV1beta1Interface
}

type DataObjectProvider interface {
	GetDv(string, string) (*cdiv1beta1.DataVolume, error)
	GetDs(string, string) (*cdiv1beta1.DataSource, error)
	CreateDo(*unstructured.Unstructured, bool) error
}

func NewDataObjectProvider(client cdiclientv1beta1.CdiV1beta1Interface) DataObjectProvider {
	return &dataObjectProvider{
		client: client,
	}
}

func (d *dataObjectProvider) GetDv(namespace string, name string) (*cdiv1beta1.DataVolume, error) {
	return d.client.DataVolumes(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (d *dataObjectProvider) GetDs(namespace string, name string) (*cdiv1beta1.DataSource, error) {
	return d.client.DataSources(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (d *dataObjectProvider) CreateDo(obj *unstructured.Unstructured, allowReplace bool) error {
	dc := discovery.NewDiscoveryClient(d.client.RESTClient())
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	gvk := obj.GroupVersionKind()
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}
	helper := resource.NewHelper(d.client.RESTClient(), mapping)

	existing, err := helper.Get(obj.GetNamespace(), obj.GetName())
	if existing != nil && allowReplace {
		if _, err := helper.Delete(obj.GetNamespace(), obj.GetName()); err != nil {
			return err
		}
	}

	_, err = helper.Create(obj.GetNamespace(), false, obj)
	return err
}

type DataObjectCreator struct {
	cliOptions         *parse.CLIOptions
	dataObjectProvider DataObjectProvider
}

func NewDataObjectCreator(cliOptions *parse.CLIOptions) (*DataObjectCreator, error) {
	log.Logger().Debug("initialized clients and providers")

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return &DataObjectCreator{
		cliOptions:         cliOptions,
		dataObjectProvider: NewDataObjectProvider(cdiclientv1beta1.NewForConfigOrDie(config)),
	}, nil
}

func (d *DataObjectCreator) CreateDataObject() (*unstructured.Unstructured, error) {
	do := d.cliOptions.GetUnstructuredDataObject()
	do.SetNamespace(d.cliOptions.GetDataObjectNamespace())

	var waitForSuccess func() error
	switch do.GetKind() {
	case constants.DataVolumeKind:
		waitForSuccess = d.waitForSuccessDv
	case constants.DataSourceKind:
		waitForSuccess = d.waitForSuccessDs
	default:
		return nil, zerrors.NewSoftError("unsupported data object kind")
	}

	err := d.dataObjectProvider.CreateDo(&do, d.cliOptions.GetAllowReplace())
	if err != nil {
		return nil, zerrors.NewSoftError("could not create data object: %v", err.Error())
	}

	if d.cliOptions.GetWaitForSuccess() {
		log.Logger().Debug("waiting for success of data object", zap.Reflect("do", do))
		if err := waitForSuccess(); err != nil {
			return nil, zerrors.NewSoftError("Failed to wait for success of data object: %v", err.Error())
		}
	}

	return &do, nil
}

func (d *DataObjectCreator) waitForSuccessDv() error {
	return wait.PollImmediate(constants.PollInterval, constants.PollTimeout, func() (bool, error) {
		do := d.cliOptions.GetUnstructuredDataObject()
		dv, err := d.dataObjectProvider.GetDv(do.GetNamespace(), do.GetName())
		if err != nil {
			return false, err
		}

		if isDataVolumeImportStatusSuccessful(dv) {
			return true, nil
		}

		if hasDataVolumeFailedToImport(dv) {
			return false, zerrors.NewSoftError("Import of DV failed: %v", dv)
		}

		if dv.Status.Phase == cdiv1beta1.Failed {
			return false, zerrors.NewSoftError("DV is in phase failed: %v", dv)
		}

		return false, nil
	})
}

func (d *DataObjectCreator) waitForSuccessDs() error {
	return wait.PollImmediate(constants.PollInterval, constants.PollTimeout, func() (bool, error) {
		do := d.cliOptions.GetUnstructuredDataObject()
		ds, err := d.dataObjectProvider.GetDs(do.GetNamespace(), do.GetName())
		if err != nil {
			return false, err
		}

		if isDataSourceReady(ds) {
			return true, nil
		}

		return false, nil
	})
}
