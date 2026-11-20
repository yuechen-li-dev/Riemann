package compiler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
)

func TestM25SuccessAClosesOnlyContactPreservingComponent(t *testing.T) {
	r, err := testM25()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.Outcome, "Success A") || !r.Certificate.Control.M23ControlsComponent || !r.Certificate.Control.OptimizerFeasible {
		t.Fatalf("wrong M25 theorem status: %+v", r.Certificate.Control)
	}
	if !strings.Contains(r.Certificate.Optimum.Scope, "contact-preserving") || strings.Contains(r.Certificate.Optimum.Scope, "all two-atom") {
		t.Fatalf("M25 scope overclaim: %s", r.Certificate.Optimum.Scope)
	}
}

func TestM25ReportedJUsesComponentObjectiveLowerEndpoint(t *testing.T) {
	r, err := testM25()
	if err != nil {
		t.Fatal(err)
	}
	lower := r.Certificate.Optimum.Objective.Lower
	c, ok := new(big.Rat).SetString(lower)
	if !ok {
		t.Fatalf("invalid M25 objective endpoint %q", lower)
	}
	j := new(big.Rat).Sub(big.NewRat(2, 1), c)
	if !strings.HasPrefix(r.BestCertifiedC, "c_* >= "+lower+" ") ||
		!strings.HasPrefix(r.BestCertifiedJ, "J_* <= "+j.FloatString(16)+" from J=2-c ") {
		t.Fatalf("inconsistent M25 projection of %s: %s, %s", lower, r.BestCertifiedC, r.BestCertifiedJ)
	}
}

func TestM25ReportAndManifestAreDeterministic(t *testing.T) {
	r, err := testM25()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M25JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M25JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("M25 JSON report is nondeterministic")
	}
	for _, want := range []string{"B1 CERTIFIED GRAPH", "WHOLE-LINE / COMPETITOR CONTROL", "M23/M24/M25 COMPOSITION", "ONE PROPOSED NEXT MILESTONE"} {
		if !strings.Contains(M25HumanReport(r), want) {
			t.Fatalf("M25 human report misses %q", want)
		}
	}
	wants := map[string]string{
		"../semantic/m25_continuation.generated.go": "9caa7311d514431159b52b1f3d5ed054f13e116af12c970b1692db10aa565c81",
		"../artifacts/m25/m25_b1_continuation.png":  "96cfe6d47be85cc75c73b3d30678bd79ee2dc1f1d7b1aeec202d6416870250c9",
		"../artifacts/m25/m25_b1_objective.png":     "7cb7e26da05a4f651b04a6a47301b705757abaa33f96597c2e5c4111fe85aabd",
		"../artifacts/m25/m25_competitor_gap.png":   "9427858006ab48710368fd43f5f3a5087c847d06081a95e3599e204f0883be03",
	}
	for path, want := range wants {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != want {
			t.Fatalf("M25 artifact %s changed: %s", path, got)
		}
	}
}

func TestM25OneNextMilestoneAndNoForbiddenExpansion(t *testing.T) {
	r, err := testM25()
	if err != nil {
		t.Fatal(err)
	}
	if r.NextMilestone == "" || strings.Contains(r.NextMilestone, "third atom") || strings.Contains(r.NextMilestone, "density") || strings.Contains(r.NextMilestone, "Aetheris") {
		t.Fatalf("wrong M25 next milestone: %s", r.NextMilestone)
	}
	if r.RHStatus != "unresolved" {
		t.Fatal("M25 changed RH status")
	}
}
