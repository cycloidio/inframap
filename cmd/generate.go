package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cycloidio/inframap/generate"
	"github.com/cycloidio/inframap/graph"
	"github.com/cycloidio/inframap/printer"
	"github.com/cycloidio/inframap/printer/factory"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	printerType string
	raw         bool
	clean       bool
	connections bool
	showIcons   bool

	generateCmd = &cobra.Command{
		Use:     "generate [FILE]",
		Short:   "Generates the Graph",
		Long:    "Generates the Graph from TFState or HCL",
		Example: "inframap generate --tfstate state.tfstate",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: preRunFile,
		RunE: func(cmd *cobra.Command, args []string) error {
			opt := generate.Options{
				Raw:         raw,
				Clean:       clean,
				Connections: connections,
			}

			var (
				g   *graph.Graph
				err error
			)

			if tfstate {
				g, _, err = generate.FromState(file, opt)
			} else {
				if len(file) == 0 {
					g, err = generate.FromHCL(afero.NewOsFs(), path, opt)
				} else {
					fs := afero.NewMemMapFs()
					path = "module.tf"

					f, err := fs.Create(path)
					if err != nil {
						return err
					}

					_, err = f.Write(file)
					if err != nil {
						return err
					}

					err = f.Sync()
					if err != nil {
						return err
					}

					g, err = generate.FromHCL(fs, path, opt)
				}
			}

			if err != nil {
				return err
			}

			p, err := factory.Get(printerType)
			if err != nil {
				return err
			}

			popt := printer.Options{
				ShowIcons: showIcons,
			}
			err = p.Print(g, popt, os.Stdout)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	generateCmd.Flags().StringVar(&printerType, "printer", "dot", fmt.Sprintf("Type of printer to use for the output. Supported ones are: %s", strings.Join(printer.TypeStrings(), ",")))
	generateCmd.Flags().BoolVar(&raw, "raw", false, "Raw will not use any specific logic from the provider, will just display the connections between elements. It's used by default if none of the Providers is known")
	generateCmd.Flags().BoolVar(&clean, "clean", true, "Clean will the generated graph will not have any Node that does not have a connection/edge")
	generateCmd.Flags().BoolVar(&connections, "connections", true, "Connections will apply the logic of the provider to remove resources that are not nodes")
	generateCmd.Flags().BoolVar(&showIcons, "show-icons", true, "Toggle the icons on the printed graph")
}
