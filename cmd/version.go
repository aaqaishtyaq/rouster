package cmd

import (
	"fmt"

	"github.com/muesli/coral"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var version = "2.0.0"

var versionCmd = &coral.Command{
	Use:   "version",
	Short: "Version of rouster",
	Long: `Returns version number of rouster
	`,
	Run:  versionCommand,
}

func versionCommand(cmd *coral.Command, args []string) {
	fmt.Println("Rouster version: " + version)
	fmt.Println("For contribution, Please visit https://github.com/aaqaishtyaq/rouster")
}
