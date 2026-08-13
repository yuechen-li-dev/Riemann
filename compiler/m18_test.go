package compiler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestM18CompilesExactAttainment(t *testing.T) {
	r, err := testM18()
	if err != nil {
		t.Fatal(err)
	}
	if !r.Optimum.EqualityDerived || r.Optimum.ExactValue.Expression != "1+pi/4-2/pi" || !strings.Contains(r.Outcome, "Success A") {
		t.Fatal("M18 did not derive exact one-radius attainment")
	}
}

func TestM18ArtifactAndReportsAreDeterministic(t *testing.T) {
	r, err := testM18()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M18JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M18JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("M18 JSON report is nondeterministic")
	}
	h := M18HumanReport(r)
	for _, want := range []string{"Success A", "1+pi/4-2/pi", "1/pi-4/pi^2+8/pi^3", "107/2000", "WHOLE-LINE", "ANTHROPIC 0.68185", "RH\n  unresolved", "ONE PROPOSED NEXT MILESTONE"} {
		if !strings.Contains(h, want) {
			t.Fatalf("human report missing %q", want)
		}
	}
}

func TestM18WitnessArtifactIsDeterministic(t *testing.T) {
	b, err := os.ReadFile("artifacts/m18_exact_one_radius_extremal.json")
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprintf("%x", sha256.Sum256(b)); got != "d0a077dfebdf314e8aac7d9bf4b46a3053f62177b46dae075d8a9460a7185ad9" {
		t.Fatalf("witness artifact changed: %s", got)
	}
}

func TestM18TangencyPlotIsDeterministic(t *testing.T) {
	b, err := os.ReadFile("../artifacts/m18/m18_exact_tangency.png")
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprintf("%x", sha256.Sum256(b)); got != "f37dd4e48fd82b7c4b4a2fe3097decc7338b11b8616a7ff03b810b51c3e19b3d" {
		t.Fatalf("tangency plot changed: %s", got)
	}
}
