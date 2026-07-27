package image

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"go.uber.org/zap"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/config"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"

	tar "kubevirt.io/containerdisks/pkg/build"
)

var labelEnvPairs = map[string]string{
	"instancetype.kubevirt.io/default-instancetype":      "INSTANCETYPE_KUBEVIRT_IO_DEFAULT_INSTANCETYPE",
	"instancetype.kubevirt.io/default-instancetype-kind": "INSTANCETYPE_KUBEVIRT_IO_DEFAULT_INSTANCETYPE_KIND",
	"instancetype.kubevirt.io/default-preference":        "INSTANCETYPE_KUBEVIRT_IO_DEFAULT_PREFERENCE",
	"instancetype.kubevirt.io/default-preference-kind":   "INSTANCETYPE_KUBEVIRT_IO_DEFAULT_PREFERENCE_KIND",
}

func DefaultConfig(labels map[string]string) v1.Config {
	var env []string
	for label, envVar := range labelEnvPairs {
		if value, exists := labels[label]; exists {
			env = append(env, fmt.Sprintf("%s=%s", envVar, value))
		}
	}
	return v1.Config{Env: env}
}

func Build(diskPath string, config v1.Config) (v1.Image, error) {
	layer, err := tarball.LayerFromOpener(tar.StreamLayerOpener(diskPath))
	if err != nil {
		return nil, fmt.Errorf("error creating layer from file: %v", err)
	}

	image, err := mutate.AppendLayers(empty.Image, layer)
	if err != nil {
		return nil, fmt.Errorf("error appending layer: %v", err)
	}

	configFile, err := image.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("error getting the image config file: %v", err)
	}
	configFile.Config = config

	image, err = mutate.ConfigFile(image, configFile)
	if err != nil {
		return nil, fmt.Errorf("error setting the image config file: %v", err)
	}
	return image, nil
}

func Push(image v1.Image, imageDestination string, pushTimeout int, auth *authn.Basic) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(pushTimeout))
	defer cancel()

	ref, err := name.ParseReference(imageDestination)
	if err != nil {
		return fmt.Errorf("error parsing image reference: %v", err)
	}

	totalSize, err := calculateImageSize(image)
	if err != nil {
		return err
	}

	// Buffer enough updates so remote.Write doesn't block waiting for our
	// progress-tracking goroutine. 100 provides generous headroom since
	// remote.Write sends one update per HTTP chunk during fast uploads.
	progressChan := make(chan v1.Update, 100)
	done := make(chan struct{})
	go func() {
		defer close(done)
		trackUploadProgress(progressChan, totalSize)
	}()
	err = remote.Write(ref, image,
		remote.WithAuth(auth),
		remote.WithContext(ctx),
		remote.WithProgress(progressChan),
		remote.WithRetryBackoff(remote.Backoff{
			Duration: 2 * time.Second,
			Factor:   2.0,
			Jitter:   0.1,
			Steps:    5,
			Cap:      60 * time.Second,
		}),
	)

	<-done

	if err != nil {
		return fmt.Errorf("error pushing image: %v", err)
	}
	return nil
}

func calculateImageSize(image v1.Image) (int64, error) {
	layers, err := image.Layers()
	if err != nil {
		return 0, fmt.Errorf("error getting image layers: %v", err)
	}

	var totalSize int64
	for _, layer := range layers {
		size, err := layer.Size()
		if err != nil {
			return 0, fmt.Errorf("error getting layer size: %v", err)
		}
		totalSize += size
	}

	log.Logger().Info("Image size calculated", zap.Int64("total_bytes", totalSize), zap.Int("layer_count", len(layers)))
	return totalSize, nil
}

func RegistryAuth() (*authn.Basic, error) {
	username := os.Getenv("ACCESS_KEY_ID")
	if username == "" {
		return nil, fmt.Errorf("ACCESS_KEY_ID environment variable is not set")
	}

	password := os.Getenv("SECRET_KEY")
	if password == "" {
		return nil, fmt.Errorf("SECRET_KEY environment variable is not set")
	}

	return &authn.Basic{
		Username: username,
		Password: password,
	}, nil
}

func trackUploadProgress(progressChan <-chan v1.Update, totalSize int64) {
	threshold := config.ProgressThreshold()
	lastLogged := -threshold
	startTime := time.Now()

	for update := range progressChan {
		if update.Error != nil {
			log.Logger().Error("Upload error", zap.Error(update.Error))
			continue
		}

		if update.Complete > 0 && totalSize > 0 {
			percentage := float64(update.Complete) / float64(totalSize) * 100
			if percentage-lastLogged >= threshold || percentage >= 100.0 {
				elapsed := math.Max(time.Since(startTime).Seconds(), 0.001)
				speedBps := float64(update.Complete) / elapsed
				speedMBps := speedBps / (1024 * 1024)
				remainingBytes := max(float64(totalSize)-float64(update.Complete), 0)
				timeRemaining := time.Duration(remainingBytes / speedBps * float64(time.Second))

				log.Logger().Info("Uploading image",
					zap.Float64("percentage", math.Round(percentage*10)/10),
					zap.Int64("bytes_uploaded", update.Complete),
					zap.Int64("total_bytes", totalSize),
					zap.Float64("speed_mbps", math.Round(speedMBps*100)/100),
					zap.String("time_remaining", timeRemaining.Truncate(time.Second).String()))

				lastLogged = percentage
			}
		} else if update.Complete > 0 {
			log.Logger().Info("Uploading image", zap.Int64("bytes_uploaded", update.Complete))
		}
	}
}
