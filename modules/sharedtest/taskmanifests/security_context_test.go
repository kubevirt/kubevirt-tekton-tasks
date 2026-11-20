package taskmanifests_test

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

const repoDir = "../../.."

// HOME of the default user of the quay.io/kubevirt/tekton-tasks image, see build/Containerfile.
// The passwd entry created for this user by build/bin/user_setup points at this directory.
const tektonTasksHome = "/home/tekton-tasks"

var releasedTasks = []string{
	"cleanup-vm",
	"create-vm-from-manifest",
	"disk-uploader",
	"disk-virt-customize",
	"disk-virt-sysprep",
	"execute-in-vm",
	"generate-ssh-keys",
	"modify-data-object",
	"modify-windows-iso-file",
	"wait-for-vmi-status",
}

func releasedTaskEntries() []TableEntry {
	entries := make([]TableEntry, 0, len(releasedTasks))
	for _, taskName := range releasedTasks {
		entries = append(entries, Entry(taskName, taskName))
	}
	return entries
}

func loadTask(taskName string) *pipev1.Task {
	content, err := os.ReadFile(filepath.Join(repoDir, "release", "tasks", taskName, taskName+".yaml"))
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	task := &pipev1.Task{}
	ExpectWithOffset(1, yaml.Unmarshal(content, task)).To(Succeed())
	ExpectWithOffset(1, task.Kind).To(Equal("Task"))
	ExpectWithOffset(1, task.Spec.Steps).ToNot(BeEmpty())
	return task
}

func findStep(task *pipev1.Task, stepName string) *pipev1.Step {
	for i := range task.Spec.Steps {
		if task.Spec.Steps[i].Name == stepName {
			return &task.Spec.Steps[i]
		}
	}
	Fail("task "+task.Name+" has no step "+stepName, 1)
	return nil
}

func findEnv(step *pipev1.Step, name string) string {
	for _, env := range step.Env {
		if env.Name == name {
			return env.Value
		}
	}
	return ""
}

func expectWritableVolumeAt(task *pipev1.Task, step *pipev1.Step, mountPath string) {
	var mount *corev1.VolumeMount
	for i := range step.VolumeMounts {
		if path.Clean(step.VolumeMounts[i].MountPath) == path.Clean(mountPath) {
			mount = &step.VolumeMounts[i]
		}
	}
	ExpectWithOffset(1, mount).ToNot(BeNil(), "step %s of task %s has no volume mounted at %s", step.Name, task.Name, mountPath)
	ExpectWithOffset(1, mount.ReadOnly).To(BeFalse(), "volume mounted at %s in step %s must be writable", mountPath, step.Name)

	var volume *corev1.Volume
	for i := range task.Spec.Volumes {
		if task.Spec.Volumes[i].Name == mount.Name {
			volume = &task.Spec.Volumes[i]
		}
	}
	ExpectWithOffset(1, volume).ToNot(BeNil(), "task %s does not define volume %s", task.Name, mount.Name)

	switch {
	case volume.EmptyDir != nil:
	case volume.PersistentVolumeClaim != nil:
		ExpectWithOffset(1, volume.PersistentVolumeClaim.ReadOnly).To(BeFalse(), "PVC volume %s must be writable", volume.Name)
	default:
		Fail("volume " + volume.Name + " mounted at " + mountPath + " must be an emptyDir or a PVC")
	}
}

var _ = Describe("Released task manifests", func() {
	It("cover every task template", func() {
		entries, err := os.ReadDir(filepath.Join(repoDir, "templates"))
		Expect(err).ToNot(HaveOccurred())

		var templates []string
		for _, entry := range entries {
			if entry.IsDir() {
				templates = append(templates, entry.Name())
			}
		}
		Expect(templates).To(ConsistOf(releasedTasks))
	})

	DescribeTable("run every step with a read-only root filesystem and the existing hardening", func(taskName string) {
		task := loadTask(taskName)

		for _, step := range task.Spec.Steps {
			securityContext := step.SecurityContext
			Expect(securityContext).ToNot(BeNil(), "step %s", step.Name)
			Expect(securityContext.ReadOnlyRootFilesystem).To(HaveValue(BeTrue()), "step %s", step.Name)
			Expect(securityContext.RunAsNonRoot).To(HaveValue(BeTrue()), "step %s", step.Name)
			Expect(securityContext.AllowPrivilegeEscalation).To(HaveValue(BeFalse()), "step %s", step.Name)
			Expect(securityContext.Privileged).To(Or(BeNil(), HaveValue(BeFalse())), "step %s", step.Name)
			Expect(securityContext.Capabilities).ToNot(BeNil(), "step %s", step.Name)
			Expect(securityContext.Capabilities.Drop).To(ContainElement(corev1.Capability("ALL")), "step %s", step.Name)
			Expect(securityContext.Capabilities.Add).To(BeEmpty(), "step %s", step.Name)
		}
	},
		releasedTaskEntries(),
	)

	DescribeTable("mount writable volumes where steps write at runtime", func(taskName, stepName string, writablePaths ...string) {
		task := loadTask(taskName)
		step := findStep(task, stepName)

		for _, writablePath := range writablePaths {
			expectWritableVolumeAt(task, step, writablePath)
		}
	},
		Entry("generate-ssh-keys: ssh-keygen output in a temporary directory",
			"generate-ssh-keys", "generate-ssh-keys", "/tmp"),
		Entry("disk-uploader: downloaded disk image and export certificate",
			"disk-uploader", "disk-uploader-step", "/tmp"),
		// libguestfs keeps its cache in /var/tmp and qemu creates the appliance snapshot overlay there
		Entry("disk-virt-customize: commands file, libguestfs and qemu temporary files and target disk",
			"disk-virt-customize", "run-virt-customize", "/tmp", "/var/tmp", "/mnt/targetpvc"),
		Entry("disk-virt-sysprep: commands file, libguestfs and qemu temporary files and target disk",
			"disk-virt-sysprep", "run-virt-sysprep", "/tmp", "/var/tmp", "/mnt/targetpvc"),
		Entry("modify-windows-iso-file: libguestfs and qemu temporary files and extracted ISO files",
			"modify-windows-iso-file", "modify-iso-file", "/tmp", "/var/tmp", "/tmp/target-pvc", "/tmp/extracted-iso-files"),
		Entry("modify-windows-iso-file: created ISO file",
			"modify-windows-iso-file", "create-iso-file", "/tmp/iso-file"),
		Entry("modify-windows-iso-file: converted target disk",
			"modify-windows-iso-file", "convert-iso-file", "/tmp/target-pvc"),
	)

	DescribeTable("point HOME of ssh based steps to a writable volume", func(taskName, stepName string) {
		task := loadTask(taskName)
		step := findStep(task, stepName)

		// ssh resolves ~/.ssh (private key, known_hosts) from the passwd entry, whose home is
		// either the image default user's HOME or the HOME passed to the container runtime.
		Expect(findEnv(step, "HOME")).To(Equal(tektonTasksHome))
		expectWritableVolumeAt(task, step, tektonTasksHome)
	},
		Entry("execute-in-vm", "execute-in-vm", "execute-in-vm"),
		Entry("cleanup-vm", "cleanup-vm", "execute-in-vm"),
	)

	It("use the HOME of the default user of the tekton-tasks image", func() {
		content, err := os.ReadFile(filepath.Join(repoDir, "build", "Containerfile"))
		Expect(err).ToNot(HaveOccurred())

		userName := regexp.MustCompile(`(?m)^\s*(?:ENV\s+)?USER_NAME=(\S+)`).FindStringSubmatch(string(content))
		Expect(userName).To(HaveLen(2), "USER_NAME is not defined in build/Containerfile")

		// HOME has to be set by its own ENV instruction to be able to reference USER_NAME
		home := regexp.MustCompile(`(?m)^ENV\s+HOME=(\S+)\s*$`).FindStringSubmatch(string(content))
		Expect(home).To(HaveLen(2), "HOME is not set by a separate ENV instruction in build/Containerfile")

		Expect(strings.ReplaceAll(home[1], "${USER_NAME}", userName[1])).To(Equal(tektonTasksHome))
	})

	It("register the default user with the configured HOME as passwd home in build/bin/user_setup", func() {
		script, err := os.ReadFile(filepath.Join(repoDir, "build", "bin", "user_setup"))
		Expect(err).ToNot(HaveOccurred())

		tmpDir, err := os.MkdirTemp("", "user-setup-")
		Expect(err).ToNot(HaveOccurred())
		DeferCleanup(os.RemoveAll, tmpDir)

		// Run a copy of the script against a fixture instead of the host /etc/passwd.
		passwd := filepath.Join(tmpDir, "passwd")
		Expect(os.WriteFile(passwd, []byte("root:x:0:0:root:/root:/bin/bash\n"), 0664)).To(Succeed())
		scriptCopy := strings.ReplaceAll(string(script), "/etc/passwd", passwd)
		Expect(scriptCopy).ToNot(ContainSubstring("/etc/passwd"))
		scriptPath := filepath.Join(tmpDir, "user_setup")
		Expect(os.WriteFile(scriptPath, []byte(scriptCopy), 0755)).To(Succeed())

		// ssh resolves ~/.ssh from the passwd home, so it has to be the HOME configured for the image
		home := filepath.Join(tmpDir, tektonTasksHome)
		cmd := exec.Command("sh", scriptPath)
		cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "USER_UID=1001", "USER_NAME=tekton-tasks", "HOME=" + home}
		output, err := cmd.CombinedOutput()
		Expect(err).ToNot(HaveOccurred(), string(output))

		content, err := os.ReadFile(passwd)
		Expect(err).ToNot(HaveOccurred())
		var entries [][]string
		for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
			if fields := strings.Split(line, ":"); len(fields) == 7 && fields[2] == "1001" {
				entries = append(entries, fields)
			}
		}
		Expect(entries).To(HaveLen(1), "passwd:\n%s", content)
		Expect(entries[0][5]).To(Equal(home), "passwd home of the default user must be the configured HOME")
	})
})
