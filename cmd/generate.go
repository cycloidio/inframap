package cmd

import (
	"errors"
	"fmt"

	"github.com/cycloidio/infraview/infraview"
	"github.com/spf13/cobra"
)

var (
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
				fmt.Println(g)
			} else {
				return errors.New("generate does not support --hcl yet")
			}

			return nil
		},
	}
)
