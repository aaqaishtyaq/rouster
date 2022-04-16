package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
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
		logrus.Fatalf("%v\n", err)
	}
}

// AddStringFlag is similar to cmd.Flags().String but supports aliases and env var
func AddStringFlag(cmd *cobra.Command, name string, aliases []string, value string, env, usage string) {
	if env != "" {
		usage = fmt.Sprintf("%s [$%s]", usage, env)
	}
	if envV, ok := os.LookupEnv(env); ok {
		value = envV
	}
	aliasesUsage := fmt.Sprintf("Alias of --%s", name)
	p := new(string)
	flags := cmd.Flags()
	flags.StringVar(p, name, value, usage)
	for _, a := range aliases {
		if len(a) == 1 {
			// pflag doesn't support short-only flags, so we have to register long one as well here
			flags.StringVarP(p, a, a, value, aliasesUsage)
		} else {
			flags.StringVar(p, a, value, aliasesUsage)
		}
	}
}
