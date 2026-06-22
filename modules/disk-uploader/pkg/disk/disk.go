package disk

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/config"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"go.uber.org/zap"
)

// Matches qemu-img progress output, e.g. "(2.43/100%)"
var progressRe = regexp.MustCompile(`\((\d+(?:\.\d+)?)/100%\)`)

type progressThrottler struct {
	pw         *io.PipeWriter
	lastLogged float64
	threshold  float64
	wg         sync.WaitGroup
}

func newProgressThrottler(threshold float64) *progressThrottler {
	pr, pw := io.Pipe()
	pt := &progressThrottler{
		pw:         pw,
		lastLogged: -threshold,
		threshold:  threshold,
	}

	pt.wg.Add(1)
	go pt.scanLines(pr)

	return pt
}

func (pt *progressThrottler) scanLines(r io.Reader) {
	defer pt.wg.Done()

	scanner := bufio.NewScanner(r)
	// qemu-img uses \r to overwrite progress updates in-place on the same line.
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if i := bytes.IndexByte(data, '\r'); i >= 0 {
			return i + 1, data[:i], nil
		}
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	for scanner.Scan() {
		line := scanner.Text()
		match := progressRe.FindStringSubmatch(line)
		if len(match) > 1 {
			if pct, err := strconv.ParseFloat(match[1], 64); err == nil {
				if pct-pt.lastLogged >= pt.threshold || pct >= 100.0 {
					log.Logger().Info("Downloading disk image", zap.Float64("percentage", pct))
					pt.lastLogged = pct
				}
			}
		}
	}
}

func (pt *progressThrottler) Write(p []byte) (n int, err error) {
	return pt.pw.Write(p)
}

func (pt *progressThrottler) Close() {
	pt.pw.Close()
	pt.wg.Wait()
}

func DownloadDiskImageFromURL(rawDiskUrl, headerKey, headerValue, certificatePath, diskPath string) error {
	pt := newProgressThrottler(config.ProgressThreshold())
	defer pt.Close()

	cmd := exec.Command(
		"nbdkit",
		"-r",
		"curl",
		rawDiskUrl,
		fmt.Sprintf("header=%s: %s", headerKey, headerValue),
		fmt.Sprintf("cainfo=%s", certificatePath),
		"--run",
		fmt.Sprintf("qemu-img convert \"$uri\" --progress --target-format qcow2 %s", diskPath),
	)
	cmd.Stdout = pt
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if fileInfo, err := os.Stat(diskPath); err != nil || fileInfo.Size() == 0 {
		return fmt.Errorf("disk image file does not exist or is empty")
	}
	return nil
}
