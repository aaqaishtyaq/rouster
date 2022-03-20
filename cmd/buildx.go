package cmd

import (
	"context"
	"log"
	"time"

	"github.com/aaqaishtyaq/rouster/client"
	"github.com/muesli/coral"
)

var (
	isPush   bool
	imageTag string
)

func init() {
	buildxCmd.Flags().StringArrayVarP(&images, "image", "i", []string{""}, "Image to be built")
	buildxCmd.Flags().BoolVarP(&isPush, "push", "p", false, "Push the image to container registry.")
	buildxCmd.Flags().StringVarP(&imageTag, "tag", "t", "latest", "Tag to be used for tagging the image.")
	buildxCmd.Flags().Int64P("deadline", "d", 20, "Image build timeout.")
	rootCmd.AddCommand(buildxCmd)
}

var buildxCmd = &coral.Command{
	Use:   "buildx",
	Short: "Buildx builder for docker image builds",
	Long: `Buildx builder for docker image builds.
	Supports image push, multi-arch image builds and more.
	`,
	Run:  buildxCommand,
	Args: coral.MinimumNArgs(1),
}

func buildxCommand(cmd *coral.Command, args []string) {
	fArgs := len(images)

	if fArgs != 0 {
		// For empty flagset the len of the slice is still 1
		if fArgs == 1 && images[0] == "" {
			log.Fatal("Err... No image provided.")
		}

		for _, i := range images {
			timeout, err := cmd.Flags().GetInt64("deadline")
			if err != nil {
				log.Fatal(err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*time.Duration(timeout)))
			defer cancel()

			client := client.NewBuildx(i, args[0])
			client.ImageTag = imageTag
			client.Push = isPush

			client.Build(&ctx)
		}
	}
}
