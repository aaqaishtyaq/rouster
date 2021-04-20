package dockerutils

import (
	"fmt"
)

type DockerBuildOpts struct {
	Directory string
	ImageName string
	ImageTag string
}

func Build(dir string, dockerBuildOptions DockerBuildOpts) {
	// To build a docker image from local files is to compress those
	// files into tar archive first

	// tar, err := archive.TarWithOptions(dir, &archive.TarOptions{})
	// if err != nil {
	// 	return err
	// }

	// opts := types.ImageBuildOptions{
	// 	Dockerfile: "Dockerfile",
	// 	Tags: []string{dockerRegistryUserID}
	// }
	fmt.Println(dir, dockerBuildOptions)
}
