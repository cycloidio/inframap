package cmd

import (
	"encoding/json"
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

	if len(args) == 1 {
		path = args[0]

		fi, err := os.Stat(path)
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			file, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}
		} else {
			// Only HCL can used with dirs
			hcl = true
			tfstate = false
		}
	} else {
		fi, err := os.Stdin.Stat()
		if err != nil {
			return err
		}

		if fi.Mode()&os.ModeNamedPipe == 0 {
			return errors.New("STDIN is empty")
		}

		file, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

	}

	setGenerateType(file)

	return err
}

// setGenerateType will try to guess the file content by first parsing it in JSON
// and if it fails fallback to HCL.
// If any of the flags --hcl or --tfstate are set it'll do nothing and use those
// directly as they are setted by the user
func setGenerateType(b []byte) {
	if hcl || tfstate {
		return
	}

	var aux map[string]interface{}
	if err := json.Unmarshal(b, &aux); err != nil {
		hcl = true
		tfstate = false
	} else {
		hcl = false
		tfstate = true
	}
}

func init() {
	rootCmd.AddCommand(
		generateCmd,
		pruneCmd,
		versionCmd,
	)

	rootCmd.PersistentFlags().BoolVar(&hcl, "hcl", false, "Forces to use HCL parser")
	rootCmd.PersistentFlags().BoolVar(&tfstate, "tfstate", false, "Forces to use TFState parser")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Activate the debug mode wich includes TF logs via TF_LOG=TRACE|DEBUG|INFO|WARN|ERROR configuration https://www.terraform.io/docs/internals/debugging.html")
}
