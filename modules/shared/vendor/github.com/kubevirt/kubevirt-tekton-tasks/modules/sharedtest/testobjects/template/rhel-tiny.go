package template

import (
	v1 "github.com/openshift/api/template/v1"
	"sigs.k8s.io/yaml"
)

const (
	RhelTemplateName = "rhel8-desktop-tiny"
)
const rhelDesktopTinyTemplateYAML = `
apiVersion: template.openshift.io/v1
kind: Template
metadata:
  name: rhel8-desktop-tiny
  annotations:
    openshift.io/display-name: "Red Hat Enterprise Linux 8.0+ VM"
    description: >-
      Template for Red Hat Enterprise Linux 8 VM or newer.
      A PVC with the RHEL disk image must be available.
    tags: "hidden,kubevirt,virtualmachine,linux,rhel"
    iconClass: "icon-rhel"
    openshift.io/provider-display-name: "KubeVirt"
    openshift.io/documentation-url: "https://github.com/kubevirt/common-templates"
    openshift.io/support-url: "https://github.com/kubevirt/common-templates/issues"
    template.openshift.io/bindable: "false"
    template.kubevirt.io/version: v1alpha1
    defaults.template.kubevirt.io/disk: rootdisk
    template.kubevirt.io/containerdisks: |
      registry.redhat.io/rhel8/rhel-guest-image
    template.kubevirt.io/editable: |
      /objects[0].spec.template.spec.domain.cpu.sockets
      /objects[0].spec.template.spec.domain.cpu.cores
      /objects[0].spec.template.spec.domain.cpu.threads
      /objects[0].spec.template.spec.domain.resources.requests.memory
      /objects[0].spec.template.spec.domain.devices.disks
      /objects[0].spec.template.spec.volumes
      /objects[0].spec.template.spec.networks
    name.os.template.kubevirt.io/rhel8.0: Red Hat Enterprise Linux 8.0 or higher
    name.os.template.kubevirt.io/rhel8.1: Red Hat Enterprise Linux 8.0 or higher
    name.os.template.kubevirt.io/rhel8.2: Red Hat Enterprise Linux 8.0 or higher
    name.os.template.kubevirt.io/rhel8.3: Red Hat Enterprise Linux 8.0 or higher
    name.os.template.kubevirt.io/rhel8.4: Red Hat Enterprise Linux 8.0 or higher
    name.os.template.kubevirt.io/rhel8.5: Red Hat Enterprise Linux 8.0 or higher
  labels:
    os.template.kubevirt.io/rhel8.0: "true"
    os.template.kubevirt.io/rhel8.1: "true"
    os.template.kubevirt.io/rhel8.2: "true"
    os.template.kubevirt.io/rhel8.3: "true"
    os.template.kubevirt.io/rhel8.4: "true"
    os.template.kubevirt.io/rhel8.5: "true"
    workload.template.kubevirt.io/desktop: "true"
    flavor.template.kubevirt.io/tiny: "true"
    template.kubevirt.io/type: "base"
    template.kubevirt.io/version: "v0.19.3"
objects:
- apiVersion: kubevirt.io/v1
  kind: VirtualMachine
  metadata:
    name: ${NAME}
    labels:
      vm.kubevirt.io/template: rhel8-desktop-tiny
      vm.kubevirt.io/template.version: "v0.19.3"
      vm.kubevirt.io/template.revision: "20"
      app: ${NAME}
    annotations:
      vm.kubevirt.io/validations: |
        [
          {
            "name": "minimal-required-memory",
            "path": "jsonpath::.spec.domain.resources.requests.memory",
            "rule": "integer",
            "message": "This VM requires more memory.",
            "min": 1610612736
          }
        ]
  spec:
    dataVolumeTemplates:
    - apiVersion: cdi.kubevirt.io/v1beta1
      kind: DataVolume
      metadata:
        name: ${NAME}
      spec:
        storage:
          resources:
            requests:
              storage: 30Gi
        sourceRef:
          kind: DataSource
          name: ${DATA_SOURCE_NAME}
          namespace: ${DATA_SOURCE_NAMESPACE}
    running: false
    template:
      metadata:
        annotations:
          vm.kubevirt.io/os: "rhel8"
          vm.kubevirt.io/workload: "desktop"
          vm.kubevirt.io/flavor: "tiny"
        labels:
          kubevirt.io/domain: ${NAME}
          kubevirt.io/size: tiny
      spec:
        domain:
          cpu:
            sockets: 1
            cores: 1
            threads: 1
          resources:
            requests:
              memory: 1.5Gi
          devices:
            rng: {}
            networkInterfaceMultiqueue: true
            inputs:
              - type: tablet
                bus: virtio
                name: tablet
            disks:
            - disk:
                bus: virtio
              name: rootdisk
            - disk:
                bus: virtio
              name: cloudinitdisk
            interfaces:
            - masquerade: {}
              name: default
        terminationGracePeriodSeconds: 180
        networks:
        - name: default
          pod: {}
        volumes:
        - dataVolume:
            name: ${NAME}
          name: rootdisk
        - cloudInitNoCloud:
            userData: |-
              #cloud-config
              user: cloud-user
              password: ${CLOUD_USER_PASSWORD}
              chpasswd: { expire: False }
          name: cloudinitdisk
parameters:
- description: VM name
  from: 'rhel8-[a-z0-9]{16}'
  generate: expression
  name: NAME
- name: DATA_SOURCE_NAME
  description: Name of the DataSource to clone
  value: 'rhel8'
- name: DATA_SOURCE_NAMESPACE
  description: Namespace of the DataSource
  value: kubevirt-os-images
- description: Randomized password for the cloud-init user cloud-user
  from: '[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}'
  generate: expression
  name: CLOUD_USER_PASSWORD

`

func NewRhelDesktopTinyTemplate() *TestTemplate {
	var template v1.Template
	err := yaml.Unmarshal([]byte(rhelDesktopTinyTemplateYAML), &template)

	if err != nil {
		panic(err)
	}

	return &TestTemplate{
		&template,
	}
}
