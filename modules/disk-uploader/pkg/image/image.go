package image

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"

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
		log.Fatalf("Error creating layer from file: %v", err)
		return nil, err
	}

	image, err := mutate.AppendLayers(empty.Image, layer)
	if err != nil {
		log.Fatalf("Error appending layer: %v", err)
		return nil, err
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

func Push(image v1.Image, imageDestination string, pushTimeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(pushTimeout))
	defer cancel()

	auth := &authn.Basic{
		Username: os.Getenv("ACCESS_KEY_ID"),
		Password: os.Getenv("SECRET_KEY"),
	}
	err := crane.Push(image, imageDestination, crane.WithAuth(auth), crane.WithContext(ctx))
	if err != nil {
		log.Fatalf("Error pushing image: %v", err)
		return err
	}
	return nil
}
