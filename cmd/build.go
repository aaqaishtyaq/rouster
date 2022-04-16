package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aaqaishtyaq/rouster/builder"
	"github.com/aaqaishtyaq/rouster/client"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/moby/buildkit/util/appdefaults"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newBuildCommand())
}

func newBuildCommand() *cobra.Command {
	var buildCommand = &cobra.Command{
		Use:          "build",
		Short:        "Build an image from a Dockerfile. Needs buildkitd to be running.",
		RunE:         buildAction,
		SilenceUsage: true,
	}

	AddStringFlag(buildCommand, "buildkit-host", nil, appdefaults.Address, "BUILDKIT_HOST", "BuildKit address")
	buildCommand.Flags().StringSlice("build-arg", []string{}, "Set build-time variables")
	buildCommand.Flags().StringP("file", "f", "", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")
	buildCommand.Flags().String("target", "", "Set the target build stage to build.")
	buildCommand.Flags().Bool("no-cache", false, "Do not use cache when building the image")

	// Docker incompatible flags
	buildCommand.Flags().StringP("image", "i", "", "Image to be built")
	buildCommand.Flags().StringP("tag", "t", "", "optionally a tag which will form this 'name:tag' format")
	buildCommand.Flags().String("buildkit-addr", appdefaults.Address, "buildkit daemon address")
	buildCommand.Flags().BoolP("push", "p", false, "Push the image to container registry.")
	buildCommand.Flags().Int64P("deadline", "d", 20, "Image build timeout.")

	return buildCommand
}

func generateBuildActionArgs(cmd *cobra.Command, args []string) (*builder.BuildOpts, error) {
	if len(args) < 1 {
		return nil, errors.New("context needs to be specified")
	}

	buildContext := args[0]

	if buildContext == "" {
		return nil, errors.New("please specify build context (e.g. \".\" for the current directory)")
	} else if buildContext == "-" || strings.Contains(buildContext, "://") {
		return nil, fmt.Errorf("unsupported build context: %q", buildContext)
	}

	bldkitAddr, err := cmd.Flags().GetString("buildkit-addr")
	if err != nil {
		return nil, err
	}

	image, err := cmd.Flags().GetString("image")
	if err != nil && image == "" {
		return nil, errors.New("No image provided")
	}

	tag, err := cmd.Flags().GetString("tag")
	if err != nil {
		return nil, err
	}

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, err
	}

	canPush, err := cmd.Flags().GetBool("push")
	if err != nil {
		return nil, err
	}

	target, err := cmd.Flags().GetString("target")
	if err != nil {
		return nil, err
	}

	noCache, err := cmd.Flags().GetBool("no-cache")
	if err != nil {
		return nil, err
	}

	buildArgs, err := cmd.Flags().GetStringSlice("build-arg")
	if err != nil {
		return nil, err
	}

	return &builder.BuildOpts{
		BuildArgs:    buildArgs,
		BuildkitAddr: bldkitAddr,
		Context:      buildContext,
		File:         file,
		Image:        image,
		NoCache:      noCache,
		Push:         canPush,
		Tag:          tag,
		Target:       target,
	}, nil
}

func buildAction(cmd *cobra.Command, args []string) error {
	buildAttrs, err := generateBuildActionArgs(cmd, args)
	if err != nil {
		return err
	}

	ctx := appcontext.Context()
	client := client.NewBuildx(buildAttrs)
	client.ImageTag = buildAttrs.Tag
	client.Push = buildAttrs.Push
	err = client.Make(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
