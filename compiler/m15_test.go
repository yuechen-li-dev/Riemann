package compiler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM15CertifiesFirstNontrivialDualAndExactFamilyCeiling(t *testing.T) {
	r, err := testM15()
	if err != nil {
		t.Fatal(err)
	}
	if !r.Witness.CertifiesInfiniteClass() || r.Bound.MultiplicityLower != (semantic.ExactRational{Numerator: 9, Denominator: 8}) || r.Bound.SimpleUpper != (semantic.ExactRational{Numerator: 7, Denominator: 8}) {
		t.Fatalf("wrong certified witness: %+v %+v", r.Witness, r.Bound)
	}
	if r.FamilyCeiling != (semantic.ExactRational{Numerator: 9, Denominator: 8}) {
		t.Fatal("family ceiling missing")
	}
}

func TestM15PrimalDualDirectionsAndSafeDecimalComparison(t *testing.T) {
	r, err := testM15()
	if err != nil {
		t.Fatal(err)
	}
	if r.M14A.KnownPrimals[1].MultiplicityUpper == nil || *r.M14A.KnownPrimals[1].MultiplicityUpper != (semantic.ExactRational{Numerator: 1651, Denominator: 1250}) {
		t.Fatal("CGdL primal direction changed")
	}
	if r.CertifiedCBracket != "9/8 <= c_* < 1651/1250 (the upper endpoint is strict)" || r.CertifiedJBracket != "849/1250 < J_* <= 7/8 (the lower endpoint is strict)" {
		t.Fatal("bound direction changed")
	}
	if !strings.Contains(r.AnthropicComparison, "inside") || !strings.Contains(r.AnthropicComparison, "neither confirmed nor contradicted") {
		t.Fatal("0.68185 was overclaimed")
	}
}

func TestM15WitnessArtifactIsConsumedAndReportDeterministic(t *testing.T) {
	r, err := testM15()
	if err != nil {
		t.Fatal(err)
	}
	if err := semantic.VerifyBoundaryAtomCertificate(r.Family, r.PDCertificate); err != nil {
		t.Fatal(err)
	}
	a, err := M15JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M15JSONReport(r)
	if err != nil || string(a) != string(b) {
		t.Fatal("M15 JSON report is nondeterministic")
	}
	var decoded map[string]any
	if err := json.Unmarshal(a, &decoded); err != nil {
		t.Fatal(err)
	}
	h := M15HumanReport(r)
	for _, want := range []string{"Success A", "9/8", "7/8", "WHOLE-LINE", "OCT TOOLING", "still unresolved", "RH\n  unresolved", "ONE NEXT MILESTONE"} {
		if !strings.Contains(h, want) {
			t.Fatalf("report missing %q", want)
		}
	}
}

func TestM15ObjectiveNotationConflictIsExplicit(t *testing.T) {
	r, err := testM15()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.ObjectiveNotationAudit, "plus-sign") && !strings.Contains(r.ObjectiveNotationAudit, "+ integral") {
		t.Fatal("objective sign audit missing")
	}
}
