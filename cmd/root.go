package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var images []string

var rootCmd = &cobra.Command{
	Use:   "rouster",
	Short: "Rather experimental docker image builder",
	Long: `Experimental docker image builder for
building cross platform docker images.
Suitable for a dockerfiles repository`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
