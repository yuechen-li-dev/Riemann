package compiler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestM20SelectsExactlyOneContactAdaptedFamily(t *testing.T) {
	r, err := testM20()
	if err != nil {
		t.Fatal(err)
	}
	if r.SelectedFamily != "D: contact-adapted r=4 atomic ray" || len(r.Candidates) != 4 {
		t.Fatalf("wrong selection: %s", r.SelectedFamily)
	}
	if !strings.Contains(r.CertifiedImprovement, "None") || r.M18.Optimum.ExactValue.Expression != "1+pi/4-2/pi" {
		t.Fatal("M20 altered theorem semantics")
	}
}

func TestM20ReportsAreDeterministicAndClassified(t *testing.T) {
	r, err := testM20()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M20JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M20JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("M20 JSON report is nondeterministic")
	}
	h := M20HumanReport(r)
	for _, want := range []string{"EXACT M18 BASELINE", "FIRST VARIATION", "  A\n", "D [SELECTED]", "WHEN UTILITY", "Numerical", "None", "ANTHROPIC 0.68185", "ONE PROPOSED NEXT MILESTONE"} {
		if !strings.Contains(h, want) {
			t.Fatalf("human report missing %q", want)
		}
	}
}

func TestM20PressureArtifactsAreDeterministic(t *testing.T) {
	wants := map[string]string{
		"../artifacts/m20/m20_baseline_pressure_map.png":        "5014c88583683c84b8e41f859dea87ba8d3306b7af59c561f0066aa9e7b1099a",
		"../artifacts/m20/m20_selected_perturbation_effect.png": "04e35e80ecea7cf0c6d745f8cdbf0f5bc0e9b0e4b92a9756c80ff376153b5a71",
	}
	for path, want := range wants {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != want {
			t.Fatalf("artifact %s changed: %s", path, got)
		}
	}
}
