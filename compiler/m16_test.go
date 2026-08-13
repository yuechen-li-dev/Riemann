package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM16CertifiesTwoRadiusImprovementAndPropagatesBound(t *testing.T) {
	r, err := testM16()
	if err != nil {
		t.Fatal(err)
	}
	if !r.Bound.MultiplicityLower.GreaterThan(semantic.ExactRational{Numerator: 9, Denominator: 8}) {
		t.Fatal("M16 did not beat 9/8")
	}
	if r.Bound.SimpleUpper != (semantic.ExactRational{Numerator: 427, Denominator: 500}) {
		t.Fatal("wrong J upper bound")
	}
	if !strings.Contains(r.AnthropicComparison, "neither confirmed nor contradicted") {
		t.Fatal("Anthropic comparison direction changed")
	}
}

func TestM16ArtifactAndReportsAreDeterministic(t *testing.T) {
	r, err := testM16()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M16JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M16JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("M16 JSON report is nondeterministic")
	}
	h := M16HumanReport(r)
	for _, want := range []string{"Success A", "573/500", "427/500", "TWO-RADIUS FAMILY", "WHOLE-LINE CERTIFICATION", "RH\n  unresolved", "ONE NEXT MILESTONE"} {
		if !strings.Contains(h, want) {
			t.Fatalf("human report missing %q", want)
		}
	}
}

func TestM16NumericalEvidenceCannotEnterCertifiedRoute(t *testing.T) {
	r, err := testM16()
	if err != nil {
		t.Fatal(err)
	}
	w := r.Witness
	w.PositivityEvidence = semantic.GridPositivityOnly
	if _, err := semantic.ApplyDualCompletion(r.M15.M14A.Class, w); err == nil {
		t.Fatal("Oct scan entered the certified route")
	}
}
