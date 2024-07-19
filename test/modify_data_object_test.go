package test

import (
	"context"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datasource"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datavolume"
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/dataobject"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/testconfigs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

var _ = Describe("Modify data objects", func() {
	f := framework.NewFramework()

	Describe("Create DataVolume", func() {
		DescribeTable("TaskRun fails and no DataVolume is ModifyD", func(config *testconfigs.ModifyDataObjectTestConfig) {
			f.TestSetup(config)

			dataVolume := config.TaskData.DataVolume
			f.ManageDataVolumes(dataVolume)

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(nil)

			if dataVolume != nil && dataVolume.Name != "" && dataVolume.Namespace != "" {
				_, err := f.CdiClient.DataVolumes(dataVolume.Namespace).Get(context.Background(), dataVolume.Name, metav1.GetOptions{})
				Expect(err).Should(HaveOccurred())
			}
		},
			Entry("empty dv", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "data-object-manifest param has to be specified",
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume: nil,
				},
			}),
			Entry("malformed dv", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "could not read data object manifest",
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume: datavolume.NewBlankDataVolume("malformed").WithoutTypeMeta().Build(),
				},
			}),
			Entry("missing name", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "invalid: metadata.name: Required value: name",
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume: datavolume.NewBlankDataVolume("").Build(),
				},
			}),
		)

		DescribeTable("DataVolume and PVC is modified successfully", func(config *testconfigs.ModifyDataObjectTestConfig) {
			f.TestSetup(config)

			dv := config.TaskData.DataVolume
			f.ManageDataVolumes(dv)

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					ModifyDataObjectResults.Name:      dv.Name,
					ModifyDataObjectResults.Namespace: dv.Namespace,
				})

			err := dataobject.WaitForSuccessfulDataVolume(f.KubevirtClient, dv.Namespace, dv.Name, config.GetWaitForDataObjectTimeout())
			Expect(err).ShouldNot(HaveOccurred())

			dv, err = f.CdiClient.DataVolumes(dv.Namespace).Get(context.Background(), dv.Name, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			if config.TaskData.SetOwnerReference == "true" {
				Expect(dv.OwnerReferences).To(HaveLen(1), "dv should has owner reference")
				Expect(dv.OwnerReferences[0].Kind).To(Equal("Pod"), "OwnerReference should have Kind Pod")
				Expect(dv.OwnerReferences[0].Name).To(HavePrefix("e2e-tests-taskrun-modify-data-object"), "OwnerReference should be binded to correct Pod")
			} else {
				Expect(dv.OwnerReferences).To(BeEmpty(), "dv OwnerReference should be empty")
			}

		},
			Entry("blank no wait", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.SmallDVCreation,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume:        datavolume.NewBlankDataVolume("blank").Build(),
					SetOwnerReference: "true",
				},
			}),
			Entry("blank wait", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.SmallDVCreation,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume:        datavolume.NewBlankDataVolume("blank-wait").Build(),
					WaitForSuccess:    true,
					SetOwnerReference: "false",
				},
			}),
			Entry("works also in the same namespace as deploy", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.SmallDVCreation,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume:     datavolume.NewBlankDataVolume("same-ns").Build(),
					WaitForSuccess: true,
				},
			}),
		)

		It("DataVolume and PVC is modified successfully with generateName", func() {
			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.SmallDVCreation,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume:     datavolume.NewBlankDataVolume("").WithGenerateName("blank-wait-").Build(),
					WaitForSuccess: true,
				},
			}
			f.TestSetup(config)

			results := runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				GetResults()

			dv := config.TaskData.DataVolume
			dv.Name = results[ModifyDataObjectResults.Name]
			f.ManageDataVolumes(dv)
			err := dataobject.WaitForSuccessfulDataVolume(f.KubevirtClient, dv.Namespace, dv.Name, config.GetWaitForDataObjectTimeout())
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("TaskRun fails and DataVolume is modified but does not import successfully", func() {
			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.QuickTaskRun,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume: datavolume.NewBlankDataVolume("import-failed").
						WithURLSource("https://invalid.source.my.domain.fail").Build(),
					WaitForSuccess: true,
				},
			}
			f.TestSetup(config)

			dv := config.TaskData.DataVolume
			f.ManageDataVolumes(dv)

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure()

			dv, err := f.CdiClient.DataVolumes(dv.Namespace).Get(context.Background(), dv.Name, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dv.Spec.Source.HTTP.URL).To(Equal(dv.Spec.Source.HTTP.URL))
			Expect(dataobject.HasDataVolumeFailedToImport(dv)).To(BeTrue())
		})

		It("Existing DataVolume is not replaced", func() {
			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.QuickTaskRun,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume:   datavolume.NewBlankDataVolume("existing-dv").Build(),
					AllowReplace: false,
				},
			}
			f.TestSetup(config)

			dvName := config.TaskData.DataVolume.Name
			dvNamespace := config.TaskData.DataVolume.Namespace

			dv := datavolume.NewBlankDataVolume(dvName).WithNamespace(dvNamespace).Build()
			dv.Annotations["cdi.kubevirt.io/storage.deleteAfterCompletion"] = "false"

			f.ManageDataVolumes(dv)

			dv, err := f.CdiClient.DataVolumes(dvNamespace).Create(context.Background(), dv, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs("already exists")

			dv2, err := f.CdiClient.DataVolumes(dvNamespace).Get(context.Background(), dvName, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dv2.CreationTimestamp).To(Equal(dv.CreationTimestamp))
			Expect(dv2.Spec).To(Equal(dv.Spec))
		})

		It("Existing pvc is deleted and replaced", func() {
			const (
				replacedURL = "https://invalid.url.replaced"
			)

			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.QuickTaskRun,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume: datavolume.NewBlankDataVolume("replace-me").
						WithURLSource(replacedURL).Build(),
					AllowReplace: true,
				},
			}
			f.TestSetup(config)

			pvcName := config.TaskData.DataVolume.Name
			pvcNamespace := config.TaskData.DataVolume.Namespace

			pvc := &k8sv1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      pvcName,
					Namespace: pvcNamespace,
				},
				Spec: k8sv1.PersistentVolumeClaimSpec{
					AccessModes: []k8sv1.PersistentVolumeAccessMode{k8sv1.ReadWriteOnce},
					Resources: k8sv1.VolumeResourceRequirements{
						Requests: k8sv1.ResourceList{
							k8sv1.ResourceStorage: *resource.NewScaledQuantity(100, resource.Mega),
						},
					},
				},
			}

			pvc, err := f.K8sClient.CoreV1().PersistentVolumeClaims(pvcNamespace).Create(context.Background(), pvc, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectResults(map[string]string{
					ModifyDataObjectResults.Name:      pvcName,
					ModifyDataObjectResults.Namespace: pvcNamespace,
				})

			var dv *cdiv1beta1.DataVolume
			dv, err = f.CdiClient.DataVolumes(pvcNamespace).Get(context.Background(), pvcName, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataVolumes(dv)

			Expect(dv.Spec.Source.HTTP.URL).To(Equal(replacedURL))

			pvc2, err := f.K8sClient.CoreV1().PersistentVolumeClaims(pvcNamespace).Get(context.Background(), pvcName, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(pvc2.CreationTimestamp).ToNot(Equal(pvc.CreationTimestamp))

		})

		It("Existing DataVolume is replaced", func() {
			const (
				initialURL  = "https://invalid.url.initial"
				replacedURL = "https://invalid.url.replaced"
			)

			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.QuickTaskRun,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume: datavolume.NewBlankDataVolume("replace-me").
						WithURLSource(replacedURL).Build(),
					AllowReplace: true,
				},
			}
			f.TestSetup(config)

			dvName := config.TaskData.DataVolume.Name
			dvNamespace := config.TaskData.DataVolume.Namespace

			dv := datavolume.NewBlankDataVolume(dvName).WithURLSource(initialURL).WithNamespace(dvNamespace).Build()
			f.ManageDataVolumes(dv)

			dv, err := f.CdiClient.DataVolumes(dvNamespace).Create(context.Background(), dv, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dv.Spec.Source.HTTP.URL).To(Equal(initialURL))

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectResults(map[string]string{
					ModifyDataObjectResults.Name:      dvName,
					ModifyDataObjectResults.Namespace: dvNamespace,
				})

			dv2, err := f.CdiClient.DataVolumes(dvNamespace).Get(context.Background(), dvName, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dv2.CreationTimestamp).ToNot(Equal(dv.CreationTimestamp))
			Expect(dv2.Spec.Source.HTTP.URL).To(Equal(replacedURL))
		})
	})

	Describe("Create DataSource", func() {
		DescribeTable("TaskRun fails and no DataSource is modified", func(config *testconfigs.ModifyDataObjectTestConfig) {
			f.TestSetup(config)

			ds := config.TaskData.DataSource
			f.ManageDataSources(ds)

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(nil)

			if ds != nil && ds.Name != "" && ds.Namespace != "" {
				_, err := f.CdiClient.DataSources(ds.Namespace).Get(context.Background(), ds.Name, metav1.GetOptions{})
				Expect(err).Should(HaveOccurred())
			}
		},
			Entry("empty ds", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "data-object-manifest param has to be specified",
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource: nil,
				},
			}),
			Entry("malformed ds", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "could not read data object manifest",
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource: datasource.NewDataSource("malformed").WithoutTypeMeta().Build(),
				},
			}),
			Entry("missing name", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "invalid: metadata.name: Required value: name",
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource: datasource.NewDataSource("").Build(),
				},
			}),
		)

		DescribeTable("DataSource is modified successfully", func(config *testconfigs.ModifyDataObjectTestConfig) {
			f.TestSetup(config)

			ds := config.TaskData.DataSource
			f.ManageDataSources(ds)

			dv := datavolume.NewBlankDataVolume(ds.Name).Build()
			f.ManageDataVolumes(dv)

			dv, err := f.CdiClient.DataVolumes(ds.Namespace).Create(context.Background(), dv, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			ds.Spec.Source.PVC = &cdiv1beta1.DataVolumeSourcePVC{
				Name:      dv.Name,
				Namespace: dv.Namespace,
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					ModifyDataObjectResults.Name:      ds.Name,
					ModifyDataObjectResults.Namespace: ds.Namespace,
				})

			ds, err = f.CdiClient.DataSources(ds.Namespace).Get(context.Background(), ds.Name, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			err = dataobject.WaitForSuccessfulDataSource(f.CdiClient, ds.Namespace, ds.Name, config.GetWaitForDataObjectTimeout())
			Expect(err).ShouldNot(HaveOccurred())

			if config.TaskData.SetOwnerReference == "true" {
				Expect(ds.OwnerReferences).To(HaveLen(1), "ds should has owner reference")
				Expect(ds.OwnerReferences[0].Kind).To(Equal("Pod"), "OwnerReference should have Kind Pod")
				Expect(ds.OwnerReferences[0].Name).To(HavePrefix("e2e-tests-taskrun-modify-data-object"), "OwnerReference should be binded to correct Pod")
			} else {
				Expect(ds.OwnerReferences).To(BeEmpty(), "ds OwnerReference should be empty")
			}
		},
			Entry("blank no wait", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					Timeout: Timeouts.SmallDVCreation,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource:        datasource.NewDataSource("blank").Build(),
					SetOwnerReference: "true",
				},
			}),
			Entry("blank wait", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.SmallDVCreation,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource:        datasource.NewDataSource("blank-wait").Build(),
					WaitForSuccess:    true,
					SetOwnerReference: "false",
				},
			}),
			Entry("works also in the same namespace as deploy", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.SmallDVCreation,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource:     datasource.NewDataSource("same-ns").Build(),
					WaitForSuccess: true,
				},
			}),
		)

		It("DataSource is modified successfully with generateName", func() {
			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.SmallDVCreation,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource:     datasource.NewDataSource("").WithGenerateName("blank-wait-").Build(),
					WaitForSuccess: true,
				},
			}
			f.TestSetup(config)

			ds := config.TaskData.DataSource
			dv := datavolume.NewBlankDataVolume("blank-wait").Build()
			f.ManageDataVolumes(dv)

			dv, err := f.CdiClient.DataVolumes(ds.Namespace).Create(context.Background(), dv, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			ds.Spec.Source.PVC = &cdiv1beta1.DataVolumeSourcePVC{
				Name:      dv.Name,
				Namespace: dv.Namespace,
			}

			results := runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				GetResults()

			ds.Name = results[ModifyDataObjectResults.Name]
			f.ManageDataSources(ds)

			err = dataobject.WaitForSuccessfulDataSource(f.CdiClient, ds.Namespace, ds.Name, config.GetWaitForDataObjectTimeout())
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("TaskRun fails and DataSource is modified but does not get ready", func() {
			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource:     datasource.NewDataSource("import-failed").Build(),
					WaitForSuccess: true,
				},
			}
			f.TestSetup(config)

			ds := config.TaskData.DataSource
			f.ManageDataSources(ds)

			dataVolume := datavolume.NewBlankDataVolume(ds.Name).WithURLSource("https://invalid.source.my.domain.fail").Build()
			f.ManageDataVolumes(dataVolume)

			dataVolume, err := f.CdiClient.DataVolumes(ds.Namespace).Create(context.Background(), dataVolume, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			ds.Spec.Source.PVC = &cdiv1beta1.DataVolumeSourcePVC{
				Name:      dataVolume.Name,
				Namespace: dataVolume.Namespace,
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure()

			ds, err = f.CdiClient.DataSources(ds.Namespace).Get(context.Background(), ds.Name, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dataobject.IsDataSourceReady(ds)).To(BeFalse())
		})

		It("Existing DataSource is not replaced", func() {
			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					Timeout: Timeouts.QuickTaskRun,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource:   datasource.NewDataSource("existing-ds").Build(),
					AllowReplace: false,
				},
			}
			f.TestSetup(config)

			dsName := config.TaskData.DataSource.Name
			dsNamespace := config.TaskData.DataSource.Namespace

			ds := datasource.NewDataSource(dsName).WithNamespace(dsNamespace).Build()
			f.ManageDataSources(ds)

			ds, err := f.CdiClient.DataSources(dsNamespace).Create(context.Background(), ds, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs("already exists")

			ds2, err := f.CdiClient.DataSources(dsNamespace).Get(context.Background(), dsName, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ds2.CreationTimestamp).To(Equal(ds.CreationTimestamp))
			Expect(ds2.Spec).To(Equal(ds.Spec))
		})

		It("Existing DataSource is replaced", func() {
			const (
				initialPVC  = "initialPVC"
				initialNS   = "initialNS"
				replacedPVC = "replacedPVC"
				replacedNS  = "replacedNS"
			)

			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.QuickTaskRun,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource:   datasource.NewDataSource("existing-ds").WithSourcePVC(replacedPVC, replacedNS).Build(),
					AllowReplace: true,
				},
			}
			f.TestSetup(config)

			dsName := config.TaskData.DataSource.Name
			dsNamespace := config.TaskData.DataSource.Namespace

			ds := datasource.NewDataSource(dsName).WithSourcePVC(initialPVC, initialNS).WithNamespace(dsNamespace).Build()
			f.ManageDataSources(ds)

			ds, err := f.CdiClient.DataSources(dsNamespace).Create(context.Background(), ds, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ds.Spec.Source.PVC.Name).To(Equal(initialPVC))
			Expect(ds.Spec.Source.PVC.Namespace).To(Equal(initialNS))

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectResults(map[string]string{
					ModifyDataObjectResults.Name:      dsName,
					ModifyDataObjectResults.Namespace: dsNamespace,
				})

			ds2, err := f.CdiClient.DataSources(dsNamespace).Get(context.Background(), dsName, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ds2.CreationTimestamp).ToNot(Equal(ds.CreationTimestamp))
			Expect(ds2.Spec.Source.PVC.Name).To(Equal(replacedPVC))
			Expect(ds2.Spec.Source.PVC.Namespace).To(Equal(replacedNS))
		})
	})

	Describe("Unsupported apiVersion or kind", func() {
		DescribeTable("TaskRun fails and nothing is ModifyD", func(config *testconfigs.ModifyDataObjectTestConfig) {
			f.TestSetup(config)
			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs("could not identify data object, wrong group or kind").
				ExpectResults(nil)
		},
			Entry("Unsupported group", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					RawManifest: datasource.NewDataSource("unsupported-apiversion").WithAPIVersion("unsupported").ToString(),
				},
			}),
			Entry("Unsupported kind", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					RawManifest: datasource.NewDataSource("unsupported-kind").WithKind("unsupported").ToString(),
				},
			}),
			Entry("With VirtualMachine", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					RawManifest: testobjects.NewTestAlpineVM("alpine-vm").ToString(),
				},
			}),
		)
	})

	Describe("Delete DataSource", func() {
		DescribeTable("TaskRun fails and dataSource is not deleted", func(config *testconfigs.ModifyDataObjectTestConfig) {
			f.TestSetup(config)

			dsNamespace := config.TaskData.DataSource.Namespace

			ds, err := f.CdiClient.DataSources(dsNamespace).Create(context.Background(), config.TaskData.DataSource, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataSources(ds)

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs(config.GetAllExpectedLogs()...)

			ds, err = f.CdiClient.DataSources(dsNamespace).Get(context.Background(), ds.Name, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ds).ToNot(BeNil(), "dataSource should exists")
		},
			Entry("missing kind", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout:      Timeouts.QuickTaskRun,
					ExpectedLogs: "object-kind param has to be specified",
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource:       datasource.NewDataSource("existing-ds").Build(),
					DeleteObjectName: datasource.NewDataSource("existing-ds").Build().Name,
					DeleteObject:     true,
				},
			}),
			Entry("Unsupported kind", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout:      Timeouts.QuickTaskRun,
					ExpectedLogs: "name param has to be specified",
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource:       datasource.NewDataSource("existing-ds").Build(),
					DeleteObject:     true,
					DeleteObjectKind: "DataSource",
				},
			}),
		)
		It("Existing DataSource is deleted", func() {
			ds := datasource.NewDataSource("existing-ds").Build()

			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.QuickTaskRun,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataSource:       ds,
					DeleteObjectName: ds.Name,
					DeleteObject:     true,
					DeleteObjectKind: "DataSource",
				},
			}
			f.TestSetup(config)

			dsName := config.TaskData.DataSource.Name
			dsNamespace := config.TaskData.DataSource.Namespace

			ds, err := f.CdiClient.DataSources(dsNamespace).Create(context.Background(), ds, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataSources(ds)

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess()

			Eventually(func(g Gomega) {
				if _, err := f.CdiClient.DataSources(dsNamespace).Get(context.Background(), dsName, metav1.GetOptions{}); err != nil {
					g.Expect(errors.ReasonForError(err)).To(Equal(metav1.StatusReasonNotFound))
				}
			}, Timeouts.TaskRunExtraWaitDelay.Duration, time.Second).Should(Succeed(), "DataSource should be deleted")

		})
	})

	Describe("Delete DataVolume", func() {
		It("Existing DataVolume is deleted", func() {
			dv := datavolume.NewBlankDataVolume("existing-ds").Build()
			dv.Annotations["cdi.kubevirt.io/storage.deleteAfterCompletion"] = "false"
			config := &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout: Timeouts.QuickTaskRun,
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume:       dv,
					DeleteObjectName: dv.Name,
					DeleteObject:     true,
					DeleteObjectKind: "DataVolume",
				},
			}
			f.TestSetup(config)

			dvName := config.TaskData.DataVolume.Name
			dvNamespace := config.TaskData.DataVolume.Namespace

			dv, err := f.CdiClient.DataVolumes(dvNamespace).Create(context.Background(), dv, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataVolumes(dv)

			err = dataobject.WaitForSuccessfulDataVolume(f.KubevirtClient, dv.Namespace, dv.Name, 5*time.Minute)
			Expect(err).ToNot(HaveOccurred())

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess()

			Eventually(func(g Gomega) {
				if _, err := f.CdiClient.DataVolumes(dvNamespace).Get(context.Background(), dvName, metav1.GetOptions{}); err != nil {
					g.Expect(errors.ReasonForError(err)).To(Equal(metav1.StatusReasonNotFound))
				}
			}, Timeouts.TaskRunExtraWaitDelay.Duration, time.Second).Should(Succeed(), "DataVolume should be deleted")

		})

		DescribeTable("TaskRun fails and datavolume is not deleted", func(config *testconfigs.ModifyDataObjectTestConfig) {
			f.TestSetup(config)
			config.TaskData.DataVolume.Annotations["cdi.kubevirt.io/storage.deleteAfterCompletion"] = "false"
			dvNamespace := config.TaskData.DataVolume.Namespace

			dv, err := f.CdiClient.DataVolumes(dvNamespace).Create(context.Background(), config.TaskData.DataVolume, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataVolumes(dv)

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs(config.GetAllExpectedLogs()...)

			dv, err = f.CdiClient.DataVolumes(dvNamespace).Get(context.Background(), dv.Name, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dv).ToNot(BeNil(), "dataVolume should exists")
		},
			Entry("missing kind", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout:      Timeouts.QuickTaskRun,
					ExpectedLogs: "object-kind param has to be specified",
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume:       datavolume.NewBlankDataVolume("existing-ds").Build(),
					DeleteObjectName: datavolume.NewBlankDataVolume("existing-ds").Build().Name,
					DeleteObject:     true,
				},
			}),
			Entry("Unsupported kind", &testconfigs.ModifyDataObjectTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{

					Timeout:      Timeouts.QuickTaskRun,
					ExpectedLogs: "name param has to be specified",
				},
				TaskData: testconfigs.ModifyDataObjectTaskData{
					DataVolume:       datavolume.NewBlankDataVolume("existing-ds").Build(),
					DeleteObject:     true,
					DeleteObjectKind: "DataVolume",
				},
			}),
		)
	})
})
