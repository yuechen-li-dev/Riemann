package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuechen-li-dev/Riemann/byproducts"
)

func run() error {
	oct := flag.String("oct", "oct", "Oct CLI executable")
	out := flag.String("out", "artifacts/b3", "output directory")
	archive := flag.String("archive", "", "optional pinned source archive for B1 statement replay")
	flag.Parse()
	data, err := byproducts.CompareOctExport(*oct, "experiments/b3_chebotarev.octest", "chebotarev.txt")
	if err != nil {
		return err
	}
	groups, err := byproducts.ReadOctGroups(bytes.NewReader(data))
	if err != nil {
		return err
	}
	expected := []string{"C5", "D8", "S4"}
	if len(groups) != len(expected) {
		return fmt.Errorf("incomplete group inventory")
	}
	for i, name := range expected {
		if groups[i].Name != name {
			return fmt.Errorf("missing group %s", name)
		}
	}
	jsonData, err := json.MarshalIndent(struct {
		Revision string                        `json:"revision"`
		Scope    string                        `json:"scope"`
		Groups   []byproducts.GroupCertificate `json:"groups"`
	}{byproducts.Revision, "finite group theory and predicted density weights", groups}, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(*out, 0755); err != nil {
		return err
	}
	if *archive != "" {
		if err := replayStatement(*archive, *out); err != nil {
			return err
		}
	}
	return os.WriteFile(filepath.Join(*out, "certificates.json"), append(jsonData, '\n'), 0644)
}

func replayStatement(archive, out string) error {
	f, err := os.Open("byproducts/b1_candidates.json")
	if err != nil {
		return err
	}
	defer f.Close()
	catalog, err := byproducts.ReadCatalog(f)
	if err != nil {
		return err
	}
	for _, review := range catalog.Candidates {
		if review.Slug != "chebotarev" {
			continue
		}
		catalog.Candidates = []byproducts.Review{review}
		records, err := byproducts.Extract(archive, catalog)
		if err != nil {
			return err
		}
		data, err := json.MarshalIndent(records[0], "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(out, "source-statement.json"), append(data, '\n'), 0644)
	}
	return fmt.Errorf("B1 Chebotarev record missing")
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
