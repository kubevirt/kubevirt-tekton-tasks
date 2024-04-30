package template

import (
	v1 "github.com/openshift/api/template/v1"
	"sigs.k8s.io/yaml"
)

const (
	CirrosTemplateName = "cirros-vm-template"
)

const cirrosServerTinyTemplateYAML = `
kind: Template
apiVersion: template.openshift.io/v1
metadata:
  name: cirros-vm-template
  namespace: default
  annotations:
    name.os.template.kubevirt.io/centos7.0: CentOS 7.0 or higher
    description: VM template example
    validations: |
      [
        {
          "name": "minimal-required-memory",
          "path": "jsonpath::.spec.domain.resources.requests.memory",
          "rule": "integer",
          "message": "This VM requires more memory.",
          "min": 67108864
        }
      ]
  labels:
    os.template.kubevirt.io/centos7.0: 'true'
    flavor.template.kubevirt.io/tiny: 'true'
    workload.template.kubevirt.io/server: 'true'
    vm.kubevirt.io/template: centos-server-tiny-v0.7.0
    vm.kubevirt.io/template.namespace: openshift
    template.kubevirt.io/type: vm
parameters:
  - name: NAME
    description: Name for the new VM
    required: true
objects:
  - apiVersion: kubevirt.io/v1
    kind: VirtualMachine
    metadata:
      labels:
        app: '${NAME}'
        vm.kubevirt.io/template: centos-server-tiny
        vm.kubevirt.io/template.revision: '147'
        vm.kubevirt.io/template.version: 0.3.2
      name: '${NAME}'
    spec:
      running: false
      template:
        metadata:
          labels:
            kubevirt.io/domain: '${NAME}'
            kubevirt.io/size: tiny
        spec:
          domain:
            cpu:
              cores: 1
              sockets: 1
              threads: 1
            devices:
              disks:
                - name: containerdisk
                  bootOrder: 1
                  disk:
                    bus: virtio
                - disk:
                    bus: virtio
                  name: cloudinitdisk
              interfaces:
                - bridge: {}
                  name: default
                  model: virtio
              networkInterfaceMultiqueue: true
              rng: {}
            resources:
              requests:
                memory: 128Mi
          networks:
            - name: default
              pod: {}
          terminationGracePeriodSeconds: 0
          volumes:
            - name: containerdisk
              containerDisk:
                image: 'quay.io/kubevirt/cirros-container-disk-demo:20240426_ca94b81c6'
            - name: cloudinitdisk
              cloudInitNoCloud:
                userData: |
                  #cloud-config
                  password: cirros
                  chpasswd:
                    expire: false
          hostname: '${NAME}'
`

func NewCirrosServerTinyTemplate() *TestTemplate {
	var template v1.Template
	err := yaml.Unmarshal([]byte(cirrosServerTinyTemplateYAML), &template)

	if err != nil {
		panic(err)
	}
	template.Name = CirrosTemplateName
	return &TestTemplate{
		&template,
	}
}
