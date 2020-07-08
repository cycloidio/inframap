package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/cycloidio/infraview/infraview"
	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:     "generate",
		Short:   "Generates the Graph",
		Long:    "Generates the Graph from TFState or HCL",
		Args:    cobra.ExactValidArgs(1),
		Example: "infraview generate state.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			if tfstate {
				b, err := ioutil.ReadFile(args[0])
				if err != nil {
					return err
				}

				g, _, err := infraview.FromState(b)
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
