package cmd

import (
	"encoding/json"
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
	printerType          string
	raw                  bool
	clean                bool
	connections          bool
	showIcons            bool
	externalNodes        bool
	descriptionFile      string
	alternativeNodeNames bool

	generateCmd = &cobra.Command{
		Use:     "generate [FILE]",
		Short:   "Generates the Graph",
		Long:    "Generates the Graph from TFState or HCL",
		Example: "inframap generate state.tfstate\ncat state.tfstate | inframap generate",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: preRunFile,
		RunE: func(cmd *cobra.Command, args []string) error {
			opt := generate.Options{
				Raw:           raw,
				Clean:         clean,
				Connections:   connections,
				ExternalNodes: externalNodes,
			}

			var (
				g     *graph.Graph
				gdesc map[string]interface{}
				err   error
			)

			if tfstate {
				g, gdesc, err = generate.FromState(file, opt)
			} else {
				if len(file) == 0 {
					g, gdesc, err = generate.FromHCL(afero.NewOsFs(), path, opt)
				} else {
					fs := afero.NewMemMapFs()
					path = "module.tf"

					var f afero.File
					f, err = fs.Create(path)
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

					g, gdesc, err = generate.FromHCL(fs, path, opt)
				}
			}

			if err != nil {
				return err
			}

			p, err := factory.Get(printerType)
			if err != nil {
				return err
			}

			if descriptionFile != "" && gdesc != nil {
				df, err := os.OpenFile(descriptionFile, os.O_APPEND|os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					return err
				}
				defer df.Close()

				// The gdesc has the description of all the elements of the graph, including the
				// edges so we have to remove them from the output
				for can := range gdesc {
					if _, err := g.GetNodeByCanonical(can); err != nil {
						delete(gdesc, can)
					}
				}

				b, err := json.Marshal(gdesc)
				if err != nil {
					return err
				}

				df.Write(b)
			}

			popt := printer.Options{
				ShowIcons:            showIcons,
				AlternativeNodeNames: alternativeNodeNames,
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
	generateCmd.Flags().BoolVar(&externalNodes, "external-nodes", true, "Toggle the addition of external nodes like 'im_out' (used to show ingress connections)")
	generateCmd.Flags().StringVar(&descriptionFile, "description-file", "", "On the given file (will be created or overwritten) we'll output the description of the returned graph, with the attributes of all the visible nodes")
	generateCmd.Flags().BoolVar(&alternativeNodeNames, "alternative-node-names", false, "Whether to try reading node names from tags, labels and other sources instead of using the canonical names")
}
