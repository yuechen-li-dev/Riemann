package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yuechen-li-dev/Riemann/byproducts"
)

func run() error {
	octFlag := flag.String("oct", "oct", "Oct executable (compiled export needs octxiliary-io)")
	out := flag.String("out", "artifacts/b2", "output directory")
	flag.Parse()
	oct, err := exec.LookPath(*octFlag)
	if err != nil {
		return err
	}
	oct, err = filepath.Abs(oct)
	if err != nil {
		return err
	}
	source, err := filepath.Abs("experiments/b2_cyclic_duality.octest")
	if err != nil {
		return err
	}
	var previous []byte
	for _, mode := range []string{"interpreted", "compiled"} {
		data, err := runOct(oct, source, mode)
		if err != nil {
			return err
		}
		if previous != nil && !bytes.Equal(previous, data) {
			return fmt.Errorf("Oct exports differ between execution modes")
		}
		previous = data
	}
	certificates, err := byproducts.ReadOctCyclic(bytes.NewReader(previous))
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

func runOct(oct, source, mode string) ([]byte, error) {
	dir, err := os.MkdirTemp("", "riemann-b2-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	cmd := exec.Command(oct, "test", source, "--execution", mode)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		return nil, fmt.Errorf("Oct %s: %w", mode, err)
	}
	return os.ReadFile(filepath.Join(dir, "certificates.txt"))
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
