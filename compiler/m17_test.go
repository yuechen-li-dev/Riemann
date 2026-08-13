package compiler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM17CompilesImprovedOneRadiusWitnessAndCeiling(t *testing.T) {
	r, err := testM17()
	if err != nil {
		t.Fatal(err)
	}
	if !r.Bound.MultiplicityLower.GreaterThan(semantic.ExactRational{Numerator: 9, Denominator: 8}) || !r.Bound.MultiplicityLower.GreaterThan(semantic.ExactRational{Numerator: 573, Denominator: 500}) {
		t.Fatal("M17 did not improve both regressions")
	}
	if r.Bound.SimpleUpper != (semantic.ExactRational{Numerator: 1703, Denominator: 2000}) {
		t.Fatal("wrong propagated J upper bound")
	}
	if r.FamilyCeiling.RationalUpper != (semantic.ExactRational{Numerator: 1149, Denominator: 1000}) {
		t.Fatal("family ceiling missing")
	}
}

func TestM17ArtifactAndReportsAreDeterministic(t *testing.T) {
	r, err := testM17()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M17JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M17JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("M17 JSON report is nondeterministic")
	}
	h := M17HumanReport(r)
	for _, want := range []string{"Success B", "ONE-RADIUS FAMILY", "SUBFAMILY RELATION", "2297/2000", "1149/1000", "WHOLE-LINE CERTIFICATE", "ANTHROPIC 0.68185", "RH\n  unresolved", "ONE PROPOSED NEXT MILESTONE"} {
		if !strings.Contains(h, want) {
			t.Fatalf("human report missing %q", want)
		}
	}
}

func TestM17NumericalEvidenceCannotCertify(t *testing.T) {
	r, err := testM17()
	if err != nil {
		t.Fatal(err)
	}
	w := r.Witness
	w.PositivityEvidence = semantic.GridPositivityOnly
	if _, err := semantic.ApplyDualCompletion(r.M16.M15.M14A.Class, w); err == nil {
		t.Fatal("M17 numerical scan entered theorem route")
	}
}

func TestM17PrimalUpperBoundGuard(t *testing.T) {
	r, err := testM17()
	if err != nil {
		t.Fatal(err)
	}
	r.Bound.MultiplicityLower = semantic.ExactRational{Numerator: 1651, Denominator: 1250}
	if validateM17Result(r) == nil {
		t.Fatal("dual at imported strict primal upper bound accepted")
	}
}

func TestM17ResearchPlotArtifactIsDeterministic(t *testing.T) {
	b, err := os.ReadFile("../artifacts/m17/m17_one_radius_spectrum.png")
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprintf("%x", sha256.Sum256(b)); got != "1f0dd33ec965e8be79acc0faa0dc4fb853cd1273999908b23bdec46975ab4e82" {
		t.Fatalf("plot artifact changed: %s", got)
	}
}
