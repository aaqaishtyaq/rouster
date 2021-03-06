package cmd

import (
	"fmt"
	"os"

	"github.com/muesli/coral"
)

var images []string

var rootCmd = &coral.Command{
	Use:   "rouster",
	Short: "Rather experimental docker image builder",
	Long: `Experimental docker image builder for
building cross platform docker images.
Suitable for a dockerfiles repository`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
