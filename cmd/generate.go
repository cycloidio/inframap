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
	printerFlag string
	generateCmd = &cobra.Command{
		Use:     "generate [FILE]",
		Short:   "Generates the Graph",
		Long:    "Generates the Graph from TFState or HCL",
		Example: "infraview generate state.json",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if tfstate {
				g, _, err := infraview.FromState(file)
				if err != nil {
					return err
				}
				p, err := printer.Get(printerFlag)
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
	generateCmd.Flags().StringVar(&printerFlag, "printer", "dot", fmt.Sprintf("Type of printer to use for the output. Supported ones are: %s", strings.Join(printer.TypeStrings(), ",")))
}
