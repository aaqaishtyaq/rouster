package client

import (
	"strings"

	"github.com/aaqaishtyaq/rouster/builder"
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
// NewBuildx returns an instance of Buildx Builder
func NewBuildx(opts *builder.BuildOpts) *builder.BuildOpts {
	dir, img := GenerateMetadata(opts.Image, opts.Context)
	opts.Context = dir
	opts.ImageName = img
	return opts
}
