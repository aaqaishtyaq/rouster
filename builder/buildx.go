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

type BuildxBuildOpts struct {
	Directory string
	ImageName string
	ImageTag  string
	Image     string
	Platform  string
	Push      bool
	Log       *logrus.Logger
}

var (
	PlatformFile    = "platform"
	ociImageBuilder = "docker"
)

func (d BuildxBuildOpts) Make(ctx context.Context, cmd *cobra.Command) error {
	bldkitAddr, err := cmd.Flags().GetString("buildkit-addr")
	if err != nil {
		return err
	}

	c, err := client.New(ctx, bldkitAddr, client.WithFailFast())
	if err != nil {
		return err
	}

	d.Image = fmt.Sprintf("%s:%s", d.ImageName, d.ImageTag)

	_, pipeW := io.Pipe()
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

	if err := eg.Wait(); err != nil {
		return err
	}

	logrus.Infof("Loaded the image %q to Docker.", d.Image)
	return nil
}

func (d BuildxBuildOpts) newSolveOpt(cmd *cobra.Command, w io.WriteCloser) (*client.SolveOpt, error) {
	buildCtx := d.Directory
	if buildCtx == "" {
		return nil, errors.New("please specify build context (e.g. \".\" for the current directory)")
	} else if buildCtx == "-" {
		return nil, errors.New("stdin not supported yet")
	}

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, err
	}

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

	target, err := cmd.Flags().GetString("target")
	if err != nil {
		return nil, err
	}
	if target != "" {
		frontendAttrs["target"] = target
	}

	noCache, err := cmd.Flags().GetBool("no-cache")
	if err != nil {
		return nil, err
	}
	if noCache {
		frontendAttrs["no-cache"] = ""
	}

	buildArgs, err := cmd.Flags().GetStringSlice("build-arg")
	if err != nil {
		return nil, err
	}

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
	cmd := exec.Command("docker", "load")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Build docker image
func (d BuildxBuildOpts) Build(ctx *context.Context) {
	d.Platform = d.FetchPlatformMetadata()

	// d.buildContainerImage(ctx, cli)
	// buildCommand := d.DockerBuildCommandArgs()
	// err := utils.RunCommand(ociImageBuilder, buildCommand)
}

// func (d BuildxBuildOpts) buildContainerImage(ctx *context.Context, client *client.Client) {
// 	// To build a docker image from local files is to compress those
// 	// files into tar archive first
// 	l := log.Default()
// 	l.Println(d.Directory)
// 	tar, err := archive.TarWithOptions(d.Directory, &archive.TarOptions{})
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	opts := types.ImageBuildOptions{
// 		Dockerfile: "Dockerfile",
// 		Tags:       []string{d.ImageName},
// 		Remove:     true,
// 		Platform:   d.Platform,
// 	}

// 	res, err := client.ImageBuild(*ctx, tar, opts)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	defer res.Body.Close()

// 	printResponse(res.Body)
// }

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
// func (d BuildxBuildOpts) DockerBuildCommandArgs() []string {
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
