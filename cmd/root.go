package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/cycloidio/infraview/infraview"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:     "infraview",
		Short:   "Reads the TFState or HCL to generate a Graphical view",
		Long:    "Reads the TFState or HCL to generate a Graphical view with Nodes and Edges.",
		Example: "infraview state.json",
		Args:    cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}
			g, _, err := infraview.FromState(b)
			if err != nil {
				return err
			}
			fmt.Println(g)
			return nil
		},
	}
)
