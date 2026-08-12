package compiler

import (
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM11CompileRepresentsDistinctMomentsAndPaperNormalization(t *testing.T) {
	r, err := testM11()
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Moments) != 3 || r.Moments[0].Kind != semantic.Trace || r.Moments[1].Kind != semantic.FrobeniusNormSquared || r.Moments[2].Kind != semantic.Dimension {
		t.Fatalf("moments collapsed or missing: %+v", r.Moments)
	}
	if r.M10.Compression.Normalization != "G_tilde_T=G_T/L" || !strings.Contains(r.Moments[2].Expression, "floor(L*T/(2*pi))") {
		t.Fatal("paper matrix normalization/dimension changed")
	}
	if r.Moments[0].FiniteTheoremPremise() || r.Moments[1].FiniteTheoremPremise() {
		t.Fatal("raw asymptotic moments treated as finite inequalities")
	}
}

func TestM11FiniteTheoremIsStrictDimensionCorrectedAndIntegerRounded(t *testing.T) {
	r, err := testM11()
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join([]string{r.FiniteTheorem.Partition, r.FiniteTheorem.TraceResidual, r.FiniteTheorem.IntegerConclusion}, " ")
	for _, want := range []string{"lambda_i>theta", "d*theta", "ceil"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing %q in %s", want, joined)
		}
	}
	if !r.ExactSanity.Applicable || r.ExactSanity.IntegerLowerBound != 1 {
		t.Fatalf("bad sanity: %+v", r.ExactSanity)
	}
}

func TestM11ThresholdDimensionScalingAndAsymptoticErrorsRemainVisible(t *testing.T) {
	r, err := testM11()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.ThresholdScaling.RelativePenalty, "T^(lambda/2-1)") || !strings.Contains(r.ThresholdScaling.Conclusion, "alone would not suffice") {
		t.Fatal("dimension-weighted threshold discharge missing")
	}
	if !strings.Contains(r.ThresholdScaling.MomentError, "(l^2+X)") || !strings.Contains(r.ThresholdScaling.MomentError, "o(1)") {
		t.Fatal("Theorem 5.8 error scale missing")
	}
	for _, m := range r.AsymptoticMoments {
		if m.RemainderKind != semantic.LittleORemainder || !strings.Contains(m.Remainder, "o(") {
			t.Fatalf("error erased: %+v", m)
		}
	}
	if len(r.EventuallyBounds) != 2 || !strings.Contains(r.EventuallyBounds[0].Threshold, "epsilon") {
		t.Fatal("eventual epsilon discharge missing")
	}
}

func TestM11ReusesM10AndKeepsFringeAndCountNotionsSeparate(t *testing.T) {
	r, err := testM11()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.AsymptoticCount.M10SimpleComposition, "2*n_plus") || !strings.Contains(r.AsymptoticCount.M10SimpleComposition, "N(I'\\I)") || !strings.Contains(r.AsymptoticCount.M10DistinctComposition, "N_distinct") {
		t.Fatal("M10 composition or fringe missing")
	}
	if r.M10ReuseSanity.SimpleCriticalLowerBound != 0 || r.M10ReuseSanity.DistinctAllLowerBound != 2 {
		t.Fatalf("existing M10 helper not reused conservatively: %+v", r.M10ReuseSanity)
	}
}

func TestM11ReproducesOnlyHalfStage(t *testing.T) {
	r, err := testM11()
	if err != nil {
		t.Fatal(err)
	}
	if !r.AsymptoticCount.HalfTypeReproduced || r.AsymptoticCount.ExactSimpleConstant != "1/2" || !strings.Contains(r.AsymptoticCount.EndpointIndexBound, "3/4") {
		t.Fatalf("half stage not reconstructed: %+v", r.AsymptoticCount)
	}
	forbidden := strings.Join([]string{r.AsymptoticCount.SimpleCriticalLiminf, r.AsymptoticCount.DistinctLiminf, strings.Join(r.DerivedMathematics, " ")}, " ")
	if strings.Contains(forbidden, "2/3") || strings.Contains(forbidden, "67.25") || strings.Contains(forbidden, "rank-trace") {
		t.Fatal("M11 crossed milestone boundary")
	}
}

func TestM11CounterexamplesAndExperimentAreBounded(t *testing.T) {
	r, err := testM11()
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Counterexamples) < 6 || r.UtilitySchedulerUsed || r.Experiment.Trials != 11 || !strings.Contains(r.Experiment.EvidenceClassification, "non-certifying") {
		t.Fatalf("counterexample discipline missing: %+v", r.Experiment)
	}
}

func TestM11ReportsAreDeterministicAndStatusIsHonest(t *testing.T) {
	r, err := testM11()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M11JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M11JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatal("M11 JSON nondeterministic")
	}
	h := M11HumanReport(r)
	for _, want := range []string{"MOMENT INPUT", "FINITE SPECTRAL THEOREM", "M10 COMPOSITION", "1/2", "RH\n  unresolved", "IMPORTED FROM LITERATURE", "NEWLY DERIVED BY COMPILER / RESEARCH LOOP", "ARCHITECTURAL AWKWARDNESS", "COMPILER THEORY", "ONE NEXT MILESTONE"} {
		if !strings.Contains(h, want) {
			t.Fatalf("report missing %q", want)
		}
	}
}
