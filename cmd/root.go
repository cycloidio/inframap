package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	hcl     bool
	tfstate bool

	rootCmd = &cobra.Command{
		Use:   "infraview",
		Short: "Reads the TFState or HCL to generate a Graphical view",
		Long:  "Reads the TFState or HCL to generate a Graphical view with Nodes and Edges.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !hcl && !tfstate {
				return errors.New("either --hcl or --tfstate have to be defined")
			}
			return nil
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
