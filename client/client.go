package client

import (
	"strings"

	"github.com/aaqaishtyaq/rouster/builder"
	"github.com/sirupsen/logrus"
)

const (
	BASEDIR      string = "dockerfiles"
	ORGANISATION string = "ghcr.io"
	USER         string = "aaqaishtyaq"
)

// GenerateMetadata returns target directory and image name
func GenerateMetadata(suite, context string) (string, string) {
	// directory is of format -> base-debian or shellcheck
	// We need to compute for both
	// and also split on - to /
	baseDir := strings.Join(strings.Split(suite, "-"), "/")
	image_name := strings.Join([]string{ORGANISATION, USER, suite}, "/")
	directory := strings.Join([]string{context, baseDir}, "/")
	return directory, image_name
}

// NewNative returns an instance of Native builder
func NewNative(suite, context string) *builder.NativeDockerBuildOpts {
	dir, img := GenerateMetadata(suite, context)
	return &builder.NativeDockerBuildOpts{
		Directory: dir,
		ImageName: img,
	}
}

// NewBuildx returns an instance of Buildx Builder
func NewBuildx(suite, context string) *builder.BuildxBuildOpts {
	dir, img := GenerateMetadata(suite, context)
	return &builder.BuildxBuildOpts{
		Directory: dir,
		ImageName: img,
		Log:       logrus.New(),
	}
}
