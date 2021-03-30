package template

import (
	v1 "github.com/openshift/api/template/v1"
	"sigs.k8s.io/yaml"
)

const fedoraServerTinyTemplateYAML = `
apiVersion: template.openshift.io/v1
kind: Template
metadata:
  annotations:
    datavolume.template.kubevirt.io/fedora29: test-base
    datavolume.template.kubevirt.io/namespace: openshift-cnv-base-images
    defaults.template.kubevirt.io/disk: rootdisk
    description: |-
      This template can be used to create a VM suitable for Fedora 23 and newer. The template assumes that a PVC is available which is providing the necessary Fedora disk image.
      Recommended disk image (needs to be converted to raw) https://download.fedoraproject.org/pub/fedora/linux/releases/30/Cloud/x86_64/images/Fedora-Cloud-Base-30-1.2.x86_64.qcow2
    iconClass: icon-fedora
    name.os.template.kubevirt.io/fedora27: Fedora 27 or higher
    name.os.template.kubevirt.io/fedora28: Fedora 27 or higher
    name.os.template.kubevirt.io/fedora29: Fedora 27 or higher
    name.os.template.kubevirt.io/silverblue28: Fedora 27 or higher
    name.os.template.kubevirt.io/silverblue29: Fedora 27 or higher
    openshift.io/display-name: Fedora 23+ VM
    openshift.io/documentation-url: https://github.com/kubevirt/common-templates
    openshift.io/provider-display-name: KubeVirt
    openshift.io/support-url: https://github.com/kubevirt/common-templates/issues
    tags: hidden,kubevirt,virtualmachine,fedora,rhel
    template.kubevirt.io/editable: |
      /objects[0].spec.template.spec.domain.cpu.sockets
      /objects[0].spec.template.spec.domain.cpu.cores
      /objects[0].spec.template.spec.domain.cpu.threads
      /objects[0].spec.template.spec.domain.resources.requests.memory
      /objects[0].spec.template.spec.domain.devices.disks
      /objects[0].spec.template.spec.volumes
      /objects[0].spec.template.spec.networks
    template.kubevirt.io/version: v1alpha1
    template.openshift.io/bindable: "false"
    validations: |
      [
        {
          "name": "minimal-required-memory",
          "path": "jsonpath::.spec.domain.resources.requests.memory",
          "rule": "integer",
          "message": "This VM requires more memory.",
          "min": 1073741824
        }
      ]
  creationTimestamp: "2020-04-24T07:18:26Z"
  labels:
    flavor.template.kubevirt.io/tiny: "true"
    os.template.kubevirt.io/fedora27: "true"
    os.template.kubevirt.io/fedora28: "true"
    os.template.kubevirt.io/fedora29: "true"
    os.template.kubevirt.io/silverblue28: "true"
    os.template.kubevirt.io/silverblue29: "true"
    template.kubevirt.io/type: base
    template.kubevirt.io/version: 0.3.2
    workload.template.kubevirt.io/server: "true"
  name: fedora-server-tiny-v0.7.0
  namespace: openshift
  resourceVersion: "12463733"
  selfLink: /apis/template.openshift.io/v1/namespaces/openshift/templates/fedora-server-tiny-v0.7.0
  uid: fdab5553-7ba8-47af-9b01-9320927406e6
objects:
- apiVersion: kubevirt.io/v1
  kind: VirtualMachine
  metadata:
    labels:
      app: ${NAME}
      vm.kubevirt.io/template: fedora-server-tiny
      vm.kubevirt.io/template.revision: "147"
      vm.kubevirt.io/template.version: 0.3.2
    name: ${NAME}
  spec:
    running: false
    template:
      metadata:
        labels:
          kubevirt.io/domain: ${NAME}
          kubevirt.io/size: tiny
      spec:
        domain:
          cpu:
            cores: 1
            sockets: 1
            threads: 1
          devices:
            disks:
            - disk:
                bus: virtio
              name: rootdisk
            - disk:
                bus: virtio
              name: cloudinitdisk
            interfaces:
            - bridge: {}
              name: default
            networkInterfaceMultiqueue: true
            rng: {}
          resources:
            requests:
              memory: 1Gi
        networks:
        - name: default
          pod: {}
        terminationGracePeriodSeconds: 0
        volumes:
        - name: rootdisk
          persistentVolumeClaim:
            claimName: ${PVCNAME}
        - cloudInitNoCloud:
            userData: |-
              #cloud-config
              Password: fedora
              Chpasswd: { Expire: False }
          name: cloudinitdisk
parameters:
- description: VM name
  from: fedora-[a-z0-9]{16}
  generate: expression
  name: NAME
- description: Name of the PVC with the disk image
  name: PVCNAME
  required: true
`

func NewFedoraServerTinyTemplate() *TestTemplate {
	var template v1.Template
	err := yaml.Unmarshal([]byte(fedoraServerTinyTemplateYAML), &template)

	if err != nil {
		panic(err)
	}

	return &TestTemplate{
		&template,
	}
}
