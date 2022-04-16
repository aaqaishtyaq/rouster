package cmd

import (
	"errors"

	"github.com/aaqaishtyaq/rouster/client"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/moby/buildkit/util/appdefaults"
	"github.com/spf13/cobra"
)

// var (
// 	isPush   bool
// 	imageTag string
// )

func init() {
	buildxCmd.Flags().StringSlice("build-arg", []string{}, "Set build-time variables")
	buildxCmd.Flags().StringP("file", "f", "", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")
	buildxCmd.Flags().String("target", "", "Set the target build stage to build.")
	buildxCmd.Flags().Bool("no-cache", false, "Do not use cache when building the image")

	// Docker incompatible flags
	buildxCmd.Flags().StringP("image", "i", "", "Image to be built")
	buildxCmd.Flags().StringP("tag", "t", "", "optionally a tag which will form this 'name:tag' format")
	buildxCmd.Flags().String("buildkit-addr", appdefaults.Address, "buildkit daemon address")
	buildxCmd.Flags().BoolP("push", "p", false, "Push the image to container registry.")
	buildxCmd.Flags().Int64P("deadline", "d", 20, "Image build timeout.")
	rootCmd.AddCommand(buildxCmd)
}

var buildxCmd = &cobra.Command{
	Use:   "buildx",
	Short: "Buildx builder for docker image builds",
	Long: `Buildx builder for docker image builds.
	Supports image push, multi-arch image builds and more.
	`,
	RunE: buildxCommand,
	Args: cobra.MinimumNArgs(1),
}

func buildxCommand(cmd *cobra.Command, args []string) error {
	image, err := cmd.Flags().GetString("image")
	if err != nil && image == "" {
		return errors.New("No image provided")
	}

	tag, err := cmd.Flags().GetString("tag")
	if err != nil {
		return err
	}

	if tag == "" {
		tag = "latest"
	}

	canPush, err := cmd.Flags().GetBool("push")
	if err != nil {
		return err
	}

	ctx := appcontext.Context()
	client := client.NewBuildx(image, args[0])
	client.ImageTag = tag
	client.Push = canPush
	err = client.Make(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
