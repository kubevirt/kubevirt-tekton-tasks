package disk

import (
	"fmt"
	"os"
	"os/exec"
)

func DownloadDiskImageFromURL(rawDiskUrl, headerKey, headerValue, certificatePath, diskPath string) error {
	cmd := exec.Command(
		"nbdkit",
		"-r",
		"curl",
		rawDiskUrl,
		fmt.Sprintf("header=%s: %s", headerKey, headerValue),
		fmt.Sprintf("cainfo=%s", certificatePath),
		"--run",
		fmt.Sprintf("qemu-img convert \"$uri\" -O qcow2 %s", diskPath),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if fileInfo, err := os.Stat(diskPath); err != nil || fileInfo.Size() == 0 {
		return fmt.Errorf("disk image file does not exist or is empty")
	}
	return nil
}
