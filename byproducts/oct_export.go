package byproducts

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CompareOctExport runs one experiment in isolated working directories.
func CompareOctExport(oct, source, filename string) ([]byte, error) {
	oct, err := exec.LookPath(oct)
	if err != nil {
		return nil, err
	}
	oct, err = filepath.Abs(oct)
	if err != nil {
		return nil, err
	}
	source, err = filepath.Abs(source)
	if err != nil {
		return nil, err
	}
	interpreted, err := runOctExport(oct, source, filename, "interpreted")
	if err != nil {
		return nil, err
	}
	compiled, err := runOctExport(oct, source, filename, "compiled")
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(interpreted, compiled) {
		return nil, fmt.Errorf("Oct exports differ between execution modes")
	}
	return interpreted, nil
}

func runOctExport(oct, source, filename, mode string) ([]byte, error) {
	dir, err := os.MkdirTemp("", "riemann-oct-")
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
	return os.ReadFile(filepath.Join(dir, filename))
}
