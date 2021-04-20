package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aaqaishtyaq/rouster/dockerutils"
)

const (
	BASEDIR string = "dockerfiles"
	ORGANISATION string = "ghcr.io"
	USER string = "aaqaishtyaq"
)

func showPrompt() []string {
	dirs := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)

	i := 0
	fmt.Println("Which image to build: ")
	for {
		fmt.Printf("%d: ", i)
		scanner.Scan()
		text := scanner.Text()

		if len(text) != 0 {
			dirs = append(dirs, text)
			i += 1
		} else {
			break
		}
	}

	return dirs
}

func buildImage(suite string) {
	// directory is of format -> base-debian or shellcheck
	// We need to compute for both
	// and also split on - to /
	directory := strings.Join(strings.Split(suite, "-"), "/")
	image_name := strings.Join([]string{ORGANISATION, USER, suite}, "/")
	// docker_cmd := strings.Join(
	// 	[]string{
	// 		"docker",
	// 		"build",
	// 		"--rm",
	// 		"--force-rm",
	// 		"-t",
	// 		image_name,
	// 		directory,
	// 	}, " ",
	// )

	// fmt.Printf("Building image %s in Directory: %s \n", image_name, directory)
	// fmt.Println(docker_cmd)

	opts := dockerutils.DockerBuildOpts{
		Directory: directory,
		ImageName: image_name,
	}

	dockerutils.Build(directory, opts)
}

func main() {
	directories := showPrompt()

	if len(directories) != 0 {
		for _, d := range directories {
			buildImage(d)
		}
	} else {
		fmt.Println("Err: No image provided!")
	}
}
