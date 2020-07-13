package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cycloidio/infraview/infraview"
	"github.com/cycloidio/infraview/printer"
	"github.com/spf13/cobra"
)

var (
	printerType string
	raw         bool
	clean       bool

	generateCmd = &cobra.Command{
		Use:     "generate [FILE]",
		Short:   "Generates the Graph",
		Long:    "Generates the Graph from TFState or HCL",
		Example: "infraview generate --tfstate state.json",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: preRunFile,
		RunE: func(cmd *cobra.Command, args []string) error {
			if tfstate {
				opt := infraview.GenerateOptions{
					Raw:   raw,
					Clean: clean,
				}
				g, _, err := infraview.FromState(file, opt)
				if err != nil {
					return err
				}
				p, err := printer.Get(printerType)
				if err != nil {
					return err
				}
				p.Print(g, os.Stdout)
			} else {
				return errors.New("generate does not support --hcl yet")
			}

			return nil
		},
	}
)

func init() {
	generateCmd.Flags().StringVar(&printerType, "printer", "dot", fmt.Sprintf("Type of printer to use for the output. Supported ones are: %s", strings.Join(printer.TypeStrings(), ",")))
	generateCmd.Flags().BoolVar(&raw, "raw", false, "Raw means that will not use any specific logic from the provider, will just display the connections between elements")
	generateCmd.Flags().BoolVar(&clean, "clean", true, "Clean means that the generated graph will not have any Node that does not have a connection/edge")
}
