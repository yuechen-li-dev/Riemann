package compiler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestM24OctArtifactsAreDeterministic(t *testing.T) {
	wants := map[string]string{
		"../artifacts/m24/m24_global_radius_envelope.png": "2b6c36de1cd460c422828021fd991082fc26e810bd172064162216c4480f2d82",
		"../artifacts/m24/m24_active_contact_atlas.png":   "d71d3711429e567447045f381c51499c453afac76b5909013be4ebe0ada5c879",
		"../artifacts/m24/m24_competitor_mass_gap.png":    "04733e17e3c25f35a433084a0f0883ddf35afa01dc6291d30c1e26dc8cd9aaa2",
		"../semantic/m24_radius_barriers.generated.go":    "68592d9c19e79e237c1ce5f7f3f74aa45ec81b3011fe778604005f05abbbf4af",
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

func TestM24SuccessCIsPreciseAndDoesNotPromoteGlobality(t *testing.T) {
	result, err := testM24()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Outcome, "Success C") || result.Certificate.M23ControlsComponent || result.Certificate.OptimizerGloballyFeasible {
		t.Fatalf("wrong M24 theorem status: %+v", result.Certificate)
	}
	if result.Certificate.GlobalEvidence.UnresolvedRadius.Lower != "3.5030000000000000" || result.Certificate.GlobalEvidence.UnresolvedRadius.Upper != "3.5120000000000000" {
		t.Fatalf("wrong exact obstruction: %+v", result.Certificate.GlobalEvidence)
	}
	if result.BestCertifiedC == "" || !strings.Contains(result.BestGlobalWitness, "69/10000") || result.RHStatus != "unresolved" {
		t.Fatal("M23 theorem witness or RH boundary was lost")
	}
}

func TestM24ReportsAreDeterministicAndComplete(t *testing.T) {
	result, err := testM24()
	if err != nil {
		t.Fatal(err)
	}
	first, err := M24JSONReport(result)
	if err != nil {
		t.Fatal(err)
	}
	second, err := M24JSONReport(result)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("M24 JSON report is nondeterministic")
	}
	human := M24HumanReport(result)
	for _, want := range []string{"LEGAL RADIUS COMPONENT", "GLOBAL ENVELOPE REDUCTION", "BRANCH ATLAS", "COMPONENT CONCLUSION", "OCT-FIRST IMPLEMENTATION", "ONE PROPOSED NEXT MILESTONE"} {
		if !strings.Contains(human, want) {
			t.Fatalf("M24 human report misses %q", want)
		}
	}
}

func TestM24UsesOneNextMilestoneAndNoForbiddenExpansion(t *testing.T) {
	result, err := testM24()
	if err != nil {
		t.Fatal(err)
	}
	if result.NextMilestone == "" || strings.Contains(result.NextMilestone, "third atom") || strings.Contains(result.NextMilestone, "density") || strings.Contains(result.NextMilestone, "Aetheris") {
		t.Fatalf("wrong next milestone: %s", result.NextMilestone)
	}
	if !strings.Contains(result.Certificate.Component.SupportLegality, "no third atom") {
		t.Fatal("family scope is not explicit")
	}
}
