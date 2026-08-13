package compiler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM14AReconstructsBroaderClassAndM13Membership(t *testing.T) {
	r, err := testM14A()
	if err != nil {
		t.Fatal(err)
	}
	if r.Class.TailSemantics != semantic.NonpositiveTail || r.M13Membership.TailSemantics != semantic.ExactFourierSupport {
		t.Fatal("M13 subclass and broader class collapsed")
	}
	if err := r.M13Membership.ValidateFor(r.Class); err != nil {
		t.Fatal(err)
	}
}

func TestM14AObjectiveAndKnownPrimalsUseCorrectDirection(t *testing.T) {
	r, err := testM14A()
	if err != nil {
		t.Fatal(err)
	}
	if r.Objective.SimpleProportion != "J(g)=2-c(g)" {
		t.Fatal("objective conversion changed")
	}
	if len(r.KnownPrimals) != 2 || r.KnownPrimals[1].SimpleLower != (semantic.ExactRational{Numerator: 849, Denominator: 1250}) {
		t.Fatal("CGdL feasible lower result missing")
	}
	if r.KnownPrimals[1].MultiplicityUpper == nil || *r.KnownPrimals[1].MultiplicityUpper != (semantic.ExactRational{Numerator: 1651, Denominator: 1250}) {
		t.Fatal("primal direction reversed")
	}
}

func TestM14AOnlyBaselineDualIsCertified(t *testing.T) {
	r, err := testM14A()
	if err != nil {
		t.Fatal(err)
	}
	if !r.BaselineWitness.CertifiesInfiniteClass() || r.BaselineBound.SimpleUpper != (semantic.ExactRational{Numerator: 1, Denominator: 1}) {
		t.Fatal("analytic baseline witness failed")
	}
	if r.NumericalDual.InfiniteClassCertified || r.CeilingCertified {
		t.Fatal("numerical candidate became a ceiling theorem")
	}
	if !strings.Contains(r.NumericalDual.DerivedLocalLimit, "9/8") {
		t.Fatal("restricted-family obstruction missing")
	}
}

func TestM14ATailAndGridObstructionAreExplicit(t *testing.T) {
	r, err := testM14A()
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(r.PreciseObstruction, " ")
	for _, want := range []string{"positive-definite completion", "primal-feasible", "whole-line positive definiteness", "tail/completeness"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing obstruction %q", want)
		}
	}
}

func TestM14AReportsAreDeterministicAndDoNotClaimRH(t *testing.T) {
	r, err := testM14A()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M14AJSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M14AJSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatal("M14A JSON is nondeterministic")
	}
	var decoded map[string]any
	if err := json.Unmarshal(a, &decoded); err != nil {
		t.Fatal(err)
	}
	h := M14AHumanReport(r)
	for _, want := range []string{"SUPPORT-ONE EXTREMAL CLASS", "PRIMAL FORMULATION", "DUAL FORMULATION", "CERTIFIED DUAL WITNESS", "PRECISE OBSTRUCTION", "Anthropic", "RH\n  unresolved", "ONE NEXT MILESTONE"} {
		if !strings.Contains(h, want) {
			t.Fatalf("report missing %q", want)
		}
	}
	if strings.Contains(h, "ceiling certified: yes") {
		t.Fatal("unsupported ceiling emitted")
	}
}
