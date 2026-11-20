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
	octFlag := flag.String("oct", "oct", "Oct executable (compiled export needs octxiliary-io)")
	out := flag.String("out", "artifacts/b2", "output directory")
	flag.Parse()
	data, err := byproducts.CompareOctExport(*octFlag, "experiments/b2_cyclic_duality.octest", "certificates.txt")
	if err != nil {
		return err
	}
	certificates, err := byproducts.ReadOctCyclic(bytes.NewReader(data))
	if err != nil {
		return err
	}
	if len(certificates) != 435 {
		return fmt.Errorf("expected 435 certificates; found %d", len(certificates))
	}
	g, err := byproducts.CyclicGraph(certificates[10])
	if err != nil {
		return err
	}
	report := struct {
		Revision     string                         `json:"revision"`
		Certificates []byproducts.CyclicCertificate `json:"certificates"`
	}{byproducts.Revision, certificates}
	graph := map[string]any{"claims": g.Claims(), "transformations": g.Transformations()}
	if err := os.MkdirAll(*out, 0755); err != nil {
		return err
	}
	for name, object := range map[string]any{"certificates.json": report, "proof-graph.json": graph} {
		data, err := json.MarshalIndent(object, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(*out, name), append(data, '\n'), 0644); err != nil {
			return err
		}
	}
	fmt.Printf("Replayed %d exact certificates from byte-identical Oct exports.\n", len(certificates))
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
