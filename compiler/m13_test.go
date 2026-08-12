package compiler

import (
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM13FamilyAndParameterAreStructuralAndAdmissible(t *testing.T) {
	r, err := testM13()
	if err != nil {
		t.Fatal(err)
	}
	if r.Family.ID != "montgomery-taylor-cosine-profile" || r.Family.SquaredProfileID == "" || !r.Family.Admissibility.Even || !r.Family.Admissibility.NonnegativeProfile {
		t.Fatalf("window semantics lost: %+v", r.Family)
	}
	if r.Family.Parameter.Domain.ContainsFloat(0) || !r.Family.Parameter.Domain.ContainsFloat(1) || r.Family.Parameter.Domain.ContainsFloat(1.01) {
		t.Fatal("illegal support parameter accepted")
	}
}

func TestM13ObjectiveIsDerivedAndFlatBasePointIsExact(t *testing.T) {
	r, err := testM13()
	if err != nil {
		t.Fatal(err)
	}
	base, err := r.FlatObjective.ExactExpression.EvalExact(map[string]semantic.ExactRational{"lambda": {Numerator: 1, Denominator: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if base != (semantic.ExactRational{Numerator: 2, Denominator: 3}) {
		t.Fatalf("flat base point is not 2/3: %+v", base)
	}
	if len(r.Objective.DerivedFrom) < 2 || !strings.Contains(strings.Join(r.Objective.DerivedFrom, " "), "M12 c=2") {
		t.Fatal("objective was not derived from M12")
	}
	derived, err := r.Objective.ExactExpression.EvalFloat(map[string]float64{"lambda": 1})
	if err != nil {
		t.Fatal(err)
	}
	if derived <= 2.0/3.0 {
		t.Fatalf("optimized family did not beat base: %.12f", derived)
	}
}

func TestM13ClosedObjectiveMatchesMomentCoefficientsAtSamples(t *testing.T) {
	r, err := testM13()
	if err != nil {
		t.Fatal(err)
	}
	for _, lambda := range []float64{0.1, 0.55, 0.9, 0.999, 1.0} {
		a, err := r.Objective.ExactExpression.EvalFloat(map[string]float64{"lambda": lambda})
		if err != nil {
			t.Fatal(err)
		}
		b, err := r.Optimization.ClosedObjective.EvalFloat(map[string]float64{"lambda": lambda})
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(a-b) > 1e-11 {
			t.Fatalf("derived/closed mismatch at %g: %.16g %.16g", lambda, a, b)
		}
	}
}

func TestM13GlobalBoundaryOptimumAndCertifiedRounding(t *testing.T) {
	r, err := testM13()
	if err != nil {
		t.Fatal(err)
	}
	if r.Optimization.Optimizer != (semantic.ExactRational{Numerator: 1, Denominator: 1}) || !strings.Contains(r.Optimization.GlobalProof, "strict increase") {
		t.Fatalf("global proof missing: %+v", r.Optimization)
	}
	for _, lambda := range []float64{0.01, 0.1, 0.5, 0.9, 0.9999} {
		d, err := r.Optimization.Derivative.EvalFloat(map[string]float64{"lambda": lambda})
		if err != nil || d <= 0 {
			t.Fatalf("derivative sign failed at %g: %g %v", lambda, d, err)
		}
	}
	lower := new(big.Rat).SetFrac(big.NewInt(r.Optimization.TaylorCertificate.ObjectiveInterval.Lower.Numerator), big.NewInt(r.Optimization.TaylorCertificate.ObjectiveInterval.Lower.Denominator))
	if lower.Cmp(big.NewRat(269, 400)) <= 0 || r.Optimization.DisplaySimplePercent != "67.25%" {
		t.Fatalf("unsafe display rounding: %s %s", lower.RatString(), r.Optimization.DisplaySimplePercent)
	}
}

func TestM13ScaleNormalizationAndWrongVariants(t *testing.T) {
	r, err := testM13()
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]semantic.ExactRational{"aL": {Numerator: 3, Denominator: 1}}
	trace, err := r.ScaleChange.ApplyExact(semantic.Trace, semantic.ExactRational{Numerator: 9, Denominator: 1}, params)
	if err != nil {
		t.Fatal(err)
	}
	frob, err := r.ScaleChange.ApplyExact(semantic.FrobeniusNormSquared, semantic.ExactRational{Numerator: 9, Denominator: 1}, params)
	if err != nil {
		t.Fatal(err)
	}
	if trace != (semantic.ExactRational{Numerator: 3, Denominator: 1}) || frob != (semantic.ExactRational{Numerator: 1, Denominator: 1}) {
		t.Fatalf("scale exponents collapsed: trace=%+v frob=%+v", trace, frob)
	}
	if len(r.Counterexamples) < 5 || !strings.Contains(r.Counterexamples[0].RejectedCandidate, "linearly") {
		t.Fatal("normalization counterexample missing")
	}
}

func TestM13AsymptoticCompositionCountsAndM12Immutability(t *testing.T) {
	r, err := testM13()
	if err != nil {
		t.Fatal(err)
	}
	if r.M12.FiniteTheorem.Conclusion != "||P+Q||_F^2 >= c*tr(P)-(c^2/4)*r+2*c*tr(Q)-c^2*b" {
		t.Fatal("M12 finite theorem changed")
	}
	joined := r.AsymptoticCount.SimpleCriticalLiminf + " " + r.AsymptoticCount.CriticalDistinctLiminf + " " + r.AsymptoticCount.AllDistinctLiminf
	for _, want := range []string{"N0_simple", "N0_distinct", "N_distinct", "269/400", "669/800"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing distinct count semantic %q", want)
		}
	}
	if !strings.Contains(r.AsymptoticCount.EventuallyComposition, "M10") || !strings.Contains(r.AsymptoticCount.EventuallyComposition, "M11") {
		t.Fatal("asymptotic logic duplicated or provenance lost")
	}
}

func TestM13ReportsAreDeterministicAndComplete(t *testing.T) {
	r, err := testM13()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M13JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M13JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatal("M13 JSON nondeterministic")
	}
	h := M13HumanReport(r)
	for _, want := range []string{"TEST WINDOW", "TYPED SCALE CHANGE", "M12 COEFFICIENTS", "OPTIMIZATION", "NEWLY RECONSTRUCTED MATHEMATICAL RESULT", "67.25%", "83.625%", "M11: total first/second moments only -> 1/2", "M12 finite theorem changed: no", "ONE NEXT MILESTONE", "RH\n  unresolved"} {
		if !strings.Contains(h, want) {
			t.Fatalf("report missing %q", want)
		}
	}
}
