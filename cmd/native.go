package cmd

import (
	"context"
	"log"
	"time"

	"github.com/aaqaishtyaq/rouster/client"
	"github.com/muesli/coral"
)

func init() {
	nativeCmd.Flags().StringArrayVarP(&images, "image", "i", []string{""}, "Image to be built")
	nativeCmd.Flags().Int64P("deadline", "d", 20, "Image build timeout.")
	rootCmd.AddCommand(nativeCmd)
}

var nativeCmd = &coral.Command{
	Use:   "native",
	Short: "Native builder for docker image build",
	Long: `Native builder for docker image builds.
	Does not support image push.`,
	Run:  nativeCommand,
	Args: coral.MinimumNArgs(1),
}

func nativeCommand(cmd *coral.Command, args []string) {
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
			client := client.NewNative(i, args[0])
			client.Build(&ctx)
		}
	}
}
