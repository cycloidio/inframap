package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cycloidio/inframap/prune"
)

var (
	canonicals bool
	pruneCmd   = &cobra.Command{
		Use:     "prune [FILE]",
		Short:   "Prunes the file",
		Long:    "Prunes the TFState or HCL file",
		Example: "inframap prune state.tfstate",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: preRunFile,
		RunE: func(cmd *cobra.Command, args []string) error {
			if tfstate {
				s, err := prune.Prune(file, canonicals)
				if err != nil {
					return err
				}

				fmt.Println(string(s))
			} else {
				return errors.New("prune does not support HCL yet")
			}

			return nil
		},
	}
)

func init() {
	pruneCmd.Flags().BoolVar(&canonicals, "canonicals", false, "If the prune command will also assign random names to the resources, EX: aws_lb.front => aws_lb.123")
}
