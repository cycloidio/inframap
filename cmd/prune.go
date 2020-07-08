package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/cycloidio/infraview/infraview"
	"github.com/spf13/cobra"
)

var (
	pruneCmd = &cobra.Command{
		Use:     "prune",
		Short:   "Prunes the file",
		Long:    "Prunes the TFState or HCL file",
		Example: "infraview prune state.json",
		Args:    cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if tfstate {
				b, err := ioutil.ReadFile(args[0])
				if err != nil {
					return err
				}

				s, err := infraview.Prune(b)
				if err != nil {
					return err
				}

				fmt.Println(string(s))
			} else {
				return errors.New("prune does not support --hcl yet")
			}

			return nil
		},
	}
)
