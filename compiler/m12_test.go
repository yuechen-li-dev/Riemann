package compiler

import (
	"strings"
	"testing"
)

func TestM12PreservesStructuralDecompositionAndPremises(t *testing.T) {
	r, err := testM12()
	if err != nil {
		t.Fatal(err)
	}
	if r.Decomposition.Identity != "A_hat_T=P_near+Q_near (the M8 zero-orbit decomposition after the M10 near-window restriction and paper normalization)" || !r.Decomposition.PSD || !r.Decomposition.SameDimension {
		t.Fatalf("decomposition identity lost: %+v", r.Decomposition)
	}
	joined := strings.Join(r.FiniteTheorem.Assumptions, " ")
	for _, want := range []string{"rank(P)<=r", "n_plus(Q)<=b", "c>0", "positive semidefinite", "G=P+Q"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing premise %q", want)
		}
	}
}

func TestM12ExactCoefficientsSpecializationAndEquality(t *testing.T) {
	r, err := testM12()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"c*tr(P)", "(c^2/4)*r", "2*c*tr(Q)", "c^2*b"} {
		if !strings.Contains(r.FiniteTheorem.Conclusion, want) {
			t.Fatalf("missing coefficient %q", want)
		}
	}
	if !strings.Contains(r.FiniteTheorem.Specialization, "r >= 2*tr(P)+4*tr(Q)-4*b-||P+Q||_F^2") || r.EqualitySanity.Slack.Numerator != 0 {
		t.Fatalf("c=2 specialization/equality broken: %+v", r.EqualitySanity)
	}
}

func TestM12FiniteCountsRemainConservativeAndReuseEarlierStages(t *testing.T) {
	r, err := testM12()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.FiniteWindow.CriticalCountBound, ">=rank(P_near)") || !strings.Contains(r.FiniteWindow.SimpleRegrouping, "rank(P_simple)<=s1") || !strings.Contains(r.FiniteWindow.TailTransfer, "N(I'\\I)") {
		t.Fatalf("count or M10 bridge semantics lost: %+v", r.FiniteWindow)
	}
	if !strings.Contains(r.AsymptoticCount.NormalizedFrobenius, "reusing M11") || !strings.Contains(r.FiniteWindow.TailTransfer, "M10") {
		t.Fatal("M10/M11 logic duplicated rather than reused")
	}
}

func TestM12ReproducesTwoThirdsAndFiveSixthsWithoutOptimization(t *testing.T) {
	r, err := testM12()
	if err != nil {
		t.Fatal(err)
	}
	if !r.AsymptoticCount.TwoThirdsReproduced || r.AsymptoticCount.ExactSimpleConstant != "2/3" || r.AsymptoticCount.ExactDistinctConstant != "5/6" {
		t.Fatalf("wrong M12 constants: %+v", r.AsymptoticCount)
	}
	all := strings.Join(append(append([]string{}, r.DerivedMathematics...), r.AsymptoticCount.BandwidthFunction), " ")
	if strings.Contains(all, "67.25") {
		t.Fatal("M12 crossed optimization boundary")
	}
}

func TestM12CounterexamplesEvidenceAndReport(t *testing.T) {
	r, err := testM12()
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Counterexamples) < 7 || r.Experiment.Trials != 12 || !strings.Contains(r.Experiment.EvidenceClassification, "non-certifying") || r.OctGoUsed || r.UtilitySchedulerUsed {
		t.Fatalf("research boundary missing: %+v", r.Experiment)
	}
	a, err := M12JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M12JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatal("M12 JSON nondeterministic")
	}
	h := M12HumanReport(r)
	for _, want := range []string{"RANK-TRACE INPUT", "GENERIC FINITE THEOREM", "2/3", "5/6", "NEWLY DERIVED / RECONSTRUCTED MATHEMATICAL RESULT", "COMPILER THEORY", "ONE NEXT MILESTONE", "RH\n  unresolved"} {
		if !strings.Contains(h, want) {
			t.Fatalf("report missing %q", want)
		}
	}
}
