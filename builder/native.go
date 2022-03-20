package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

// NativeDockerBuildOpts uses docker golang client to build docker images
// subject to experimental usage
type NativeDockerBuildOpts struct {
	Directory string
	ImageName string
	ImageTag  string
}

// Build docker image
func (d NativeDockerBuildOpts) Build(ctx *context.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	d.buildContainerImage(ctx, cli)
}

func (d NativeDockerBuildOpts) buildContainerImage(ctx *context.Context, client *client.Client) {
	// To build a docker image from local files is to compress those
	// files into tar archive first
	l := log.Default()
	l.Println(d.Directory)
	tar, err := archive.TarWithOptions(d.Directory, &archive.TarOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{d.ImageName},
		Remove:     true,
	}

	res, err := client.ImageBuild(*ctx, tar, opts)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer res.Body.Close()

	printResponse(res.Body)
}

func printResponse(ioReader io.Reader) {
	var lastLine string

	scanner := bufio.NewScanner(ioReader)

	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(lastLine)
	}

	errLine := &ErrorLine{}

	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		fmt.Println(errLine.ErrorDetails.Message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err.Error())
	}
}
