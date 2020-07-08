package cmd

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var (
	hcl     bool
	tfstate bool
	file    []byte

	rootCmd = &cobra.Command{
		Use:   "infraview",
		Short: "Reads the TFState or HCL to generate a Graphical view",
		Long:  "Reads the TFState or HCL to generate a Graphical view with Nodes and Edges.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if !hcl && !tfstate {
				return errors.New("either --hcl or --tfstate have to be defined")
			}
			if len(args) == 1 {
				file, err = ioutil.ReadFile(args[0])
			} else {
				file, err = ioutil.ReadAll(os.Stdin)
			}

			return err
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(
		generateCmd,
		pruneCmd,
	)
	rootCmd.PersistentFlags().BoolVar(&hcl, "hcl", false, "HCL file/dir to read from")
	rootCmd.PersistentFlags().BoolVar(&tfstate, "tfstate", false, "Terraform State to read from")
}
