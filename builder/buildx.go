package builder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aaqaishtyaq/rouster/utils"
)

type BuildxBuildOpts struct {
	Directory string
	ImageName string
	ImageTag  string
	Platform  string
	Push      bool
	Log       *log.Logger
}

var (
	PlatformFile    = "platform"
	ociImageBuilder = "docker"
)

// Build docker image
func (d BuildxBuildOpts) Build(ctx *context.Context) {
	d.Platform = d.FetchPlatformMetadata()
	buildCommand := d.DockerBuildCommandArgs()
	utils.RunCommand(ociImageBuilder, buildCommand)
}

// FetchPlatformMetadata reads `platform` file in Dockerfile context to figure out multi-arch builds
func (d BuildxBuildOpts) FetchPlatformMetadata() string {
	platformFilePath := strings.Join([]string{d.Directory, PlatformFile}, "/")
	if _, err := os.Stat(platformFilePath); errors.Is(err, os.ErrNotExist) {
		return ""
	}

	dat, err := os.ReadFile(platformFilePath)
	if err != nil {
		d.Log.Println(err)
		return ""
	}

	platform := string(dat)
	if platform == "" {
		return platform
	}

	slic := strings.Split(string(dat), ",")
	for i := range slic {
		slic[i] = strings.TrimSpace(slic[i])
	}

	return strings.Join(slic, ",")
}

// DockerBuildCommand returns command to be used for docker build
func (d BuildxBuildOpts) DockerBuildCommandArgs() []string {
	var platformCmd string
	var loadFlag string

	if d.Push {
		loadFlag = "--push"
	} else {
		loadFlag = "--load"
	}

	if d.Platform != "" {
		platformCmd = fmt.Sprintf("--platform=%s", d.Platform)
	}

	ociImage := fmt.Sprintf("%s:%s", d.ImageName, d.ImageTag)

	cmd := []string{
		"buildx",
		"build",
		loadFlag,
		platformCmd,
		"--tag",
		ociImage,
		d.Directory,
	}

	return cmd
}
