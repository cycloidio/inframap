package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cycloidio/inframap/provider"
)

var (
	dry    bool
	branch string
)

func init() {
	flag.BoolVar(&dry, "dry", false, "Makes a dry run")
	flag.StringVar(&branch, "branch", "master", "Branch from tfdocs to use")
}

func main() {
	flag.Parse()
	for _, pr := range provider.TypeValues() {
		if pr == provider.Raw {
			continue
		}
		url := fmt.Sprintf("https://raw.githubusercontent.com/cycloidio/tfdocs/%s/assets/%s/icons.json", branch, pr)

		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != 200 {
			log.Fatalf("expecting a 200 status code but got %d for provider %s", res.StatusCode, pr)
		}

		icons := make(map[string]string)
		err = json.NewDecoder(res.Body).Decode(&icons)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()

		paths := make(map[string]bool)
		for _, p := range icons {
			if p == "" {
				continue
			}
			paths[p] = false
		}

		rootDir := fmt.Sprintf("assets/icons/%s/", pr)
		err = filepath.Walk(rootDir, func(p string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fi.IsDir() {
				// We need to remove the rootDir as it's not on tfdocs
				// and also replace the .png for .svg
				np := strings.Replace(
					strings.TrimPrefix(p, rootDir),
					".png", ".svg", -1,
				)
				if _, ok := paths[np]; !ok {
					fmt.Printf("remove %s %s\n", pr, p)
					if !dry {
						os.Remove(p)
					}
				} else {
					paths[np] = true
				}
			}

			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		for p, ok := range paths {
			if !ok {
				fmt.Printf("missing %s %s\n", pr, p)
			}
		}

	}
}
