package dockerutils

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type DockerBuildOpts struct {
	Directory string
	ImageName string
	ImageTag string
}

type ErrorLine struct {
	Error string `json:"error"`
	ErrorDetails ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func Build(dockerBuildOptions DockerBuildOpts) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	buildContainerImage(cli, dockerBuildOptions)
}

func buildContainerImage(client *client.Client, buildOpts DockerBuildOpts) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	// To build a docker image from local files is to compress those
	// files into tar archive first
	tar, err := archive.TarWithOptions(buildOpts.Directory, &archive.TarOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags: []string{buildOpts.ImageName},
		Remove: true,
	}

	res, err := client.ImageBuild(ctx, tar, opts)
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