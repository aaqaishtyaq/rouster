package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
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
	logrus.Infoln("Native Builder is deprecated...")
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
