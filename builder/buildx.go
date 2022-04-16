package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/containerd/console"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type BuildOpts struct {
	BuildArgs    []string
	BuildkitAddr string
	Context      string
	File         string
	Image        string
	Tag          string
	Push         bool
	NoCache      bool
	Target       string
	ImageName    string
	ImageTag     string
	Platform     string
}

var (
	PlatformFile    = "platform"
	ociImageBuilder = "docker"
)

func (d BuildOpts) Make(ctx context.Context, cmd *cobra.Command) error {
	bldkitAddr, err := cmd.Flags().GetString("buildkit-addr")
	if err != nil {
		return err
	}

	c, err := client.New(ctx, bldkitAddr, client.WithFailFast())
	if err != nil {
		return err
	}

	d.Image = fmt.Sprintf("%s:%s", d.ImageName, d.ImageTag)

	pipeR, pipeW := io.Pipe()
	solveOpt, err := d.newSolveOpt(cmd, pipeW)
	if err != nil {
		return err
	}

	ch := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		_, err = c.Solve(ctx, nil, *solveOpt, ch)
		if err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		var c console.Console
		if cn, err := console.ConsoleFromFile(os.Stderr); err == nil {
			c = cn
		}

		// not using shared context to not disrupt display but let is finish reporting errors
		_, err = progressui.DisplaySolveStatus(context.TODO(), "", c, os.Stdout, ch)
		return err
	})

	eg.Go(func() error {
		if err := loadDockerTar(pipeR); err != nil {
			return err
		}
		return pipeR.Close()
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	logrus.Infof("Loaded the image %q to Docker.", d.Image)
	return nil
}

func (d BuildOpts) newSolveOpt(cmd *cobra.Command, w io.WriteCloser) (*client.SolveOpt, error) {
	buildCtx := d.Context
	file := d.File
	if file == "" {
		file = filepath.Join(buildCtx, "Dockerfile")
	}

	localDirs := map[string]string{
		"context":    buildCtx,
		"dockerfile": filepath.Dir(file),
	}

	frontend := "dockerfile.v0"

	frontendAttrs := map[string]string{
		"filename": filepath.Base(file),
	}

	target := d.Target
	if target != "" {
		frontendAttrs["target"] = target
	}

	noCache := d.NoCache
	if noCache {
		frontendAttrs["no-cache"] = ""
	}

	buildArgs := d.BuildArgs

	for _, buildArg := range buildArgs {
		kv := strings.SplitN(buildArg, "=", 2)
		if len(kv) != 2 {
			return nil, errors.Errorf("invalid build-arg value %s", buildArg)
		}
		frontendAttrs["build-arg:"+kv[0]] = kv[1]
	}

	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type: "docker", // TODO: use containerd image store when it is integrated to Docker
				Attrs: map[string]string{
					"name": d.Image,
				},
				Output: func(_ map[string]string) (io.WriteCloser, error) {
					return w, nil
				},
			},
		},
		LocalDirs:     localDirs,
		Frontend:      frontend,
		FrontendAttrs: frontendAttrs,
	}, nil
}

func loadDockerTar(r io.Reader) error {
	// no need to use moby/moby/client here
	cmd := exec.Command("nerdctl", "load")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Build docker image
func (d BuildOpts) Build(ctx *context.Context) {
	d.Platform = d.FetchPlatformMetadata()

	// d.buildContainerImage(ctx, cli)
	// buildCommand := d.DockerBuildCommandArgs()
	// err := utils.RunCommand(ociImageBuilder, buildCommand)
}


// FetchPlatformMetadata reads `platform` file in Dockerfile context to figure out multi-arch builds
func (d BuildOpts) FetchPlatformMetadata() string {
	platformFilePath := strings.Join([]string{d.Context, PlatformFile}, "/")
	if _, err := os.Stat(platformFilePath); errors.Is(err, os.ErrNotExist) {
		return ""
	}

	dat, err := os.ReadFile(platformFilePath)
	if err != nil {
		// d.Log.Println(err)
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
// func (d BuildOpts) DockerBuildCommandArgs() []string {
// 	var platformCmd string
// 	var loadFlag string

// 	if d.Push {
// 		loadFlag = "--push"
// 	} else {
// 		loadFlag = "--load"
// 	}

// 	if d.Platform != "" {
// 		platformCmd = fmt.Sprintf("--platform=%s", d.Platform)
// 	}

// 	ociImage := fmt.Sprintf("%s:%s", d.ImageName, d.ImageTag)

// 	cmd := []string{
// 		"buildx",
// 		"build",
// 		loadFlag,
// 		platformCmd,
// 		"--tag",
// 		ociImage,
// 		d.Directory,
// 	}

// 	return cmd
// }
