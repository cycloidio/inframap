package cmd

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/cycloidio/terracognita/log"
	"github.com/spf13/cobra"
)

var (
	hcl     bool
	tfstate bool
	file    []byte
	path    string
	debug   bool
	logsOut io.Writer = ioutil.Discard

	rootCmd = &cobra.Command{
		Use:   "inframap",
		Short: "Reads the TFState or HCL to generate a Graphical view",
		Long:  "Reads the TFState or HCL to generate a Graphical view with Nodes and Edges.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				logsOut = os.Stdout
			}
			log.Init(logsOut, debug)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// preRunFile checks where the input is, ARGS or STDIN
// also checks if --hcl or --tfstate are setted as one of
// them is required
func preRunFile(cmd *cobra.Command, args []string) error {
	var err error
	if !hcl && !tfstate {
		return errors.New("either --hcl or --tfstate have to be defined")
	}
	if len(args) == 1 {
		path = args[0]

		fi, err := os.Stat(path)
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			file, err = ioutil.ReadFile(path)
		}
	} else {
		fi, err := os.Stdin.Stat()
		if err != nil {
			return err
		}

		if fi.Mode()&os.ModeNamedPipe == 0 {
			return errors.New("empty STDIN")
		}
		file, err = ioutil.ReadAll(os.Stdin)
	}

	return err
}

func init() {
	rootCmd.AddCommand(
		generateCmd,
		pruneCmd,
		versionCmd,
	)

	rootCmd.PersistentFlags().BoolVar(&hcl, "hcl", false, "HCL file/dir to read from")
	rootCmd.PersistentFlags().BoolVar(&tfstate, "tfstate", false, "Terraform State to read from")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Activate the debug mode wich includes TF logs via TF_LOG=TRACE|DEBUG|INFO|WARN|ERROR configuration https://www.terraform.io/docs/internals/debugging.html")
}
