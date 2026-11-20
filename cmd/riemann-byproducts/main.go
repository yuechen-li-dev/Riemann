package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuechen-li-dev/Riemann/byproducts"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	archive := flag.String("archive", "", "original pinned codeload tar.gz (required)")
	catalogPath := flag.String("catalog", "byproducts/b1_candidates.json", "review catalog")
	oct := flag.String("oct", "oct", "Oct CLI executable")
	program := flag.String("ranker", "experiments/b1_discovery/main.oct", "Oct ranking program")
	out := flag.String("out", "artifacts/b1", "report directory")
	flag.Parse()
	if *archive == "" {
		return fmt.Errorf("provide -archive; see docs/b1-architecture.md")
	}
	f, err := os.Open(*catalogPath)
	if err != nil {
		return err
	}
	defer f.Close()
	catalog, err := byproducts.ReadCatalog(f)
	if err != nil {
		return err
	}
	report, err := byproducts.Discover(*archive, catalog, *oct, *program)
	if err != nil {
		return err
	}
	data, err := report.JSON()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(*out, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "inventory.json"), data, 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(*out, "inventory.md"), []byte(report.Markdown()), 0644)
}
