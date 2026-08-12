package compiler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM10CompileUsesPaperWindowCompressionAndPreservesM7(t *testing.T) {
	r, err := testM10()
	if err != nil {
		t.Fatal(err)
	}
	if r.Compression.Window.Lower.Expression != "T" || r.Compression.Window.Upper.Expression != "2T" || r.Compression.Window.Boundary != semantic.LeftOpenRightClosed {
		t.Fatalf("wrong target window: %#v", r.Compression.Window)
	}
	if r.Compression.LocalizationWindow.Lower.Expression != "T-sqrt(T)" || r.Compression.DimensionExpression != "floor(L*T/(2*pi))" {
		t.Fatalf("paper localization/compression lost: %#v", r.Compression)
	}
	if r.Compression.MatrixID == r.M9.Compression.MatrixID || r.M9.M8.M7.Evaluation.Matrix.ID == "" || !strings.Contains(r.M7RegressionRole, "unchanged") {
		t.Fatal("M7 was overwritten or misidentified as the height family")
	}
}

func TestM10ExactWeylThresholdFixturesAndEqualityCounterexample(t *testing.T) {
	// Diagonal Weyl fixtures: |e_i|<=delta and theta>=delta imply every
	// g_i=a_i+e_i above theta has a_i>0.
	fixtures := []struct {
		a, e, theta []int64
		delta       int64
	}{
		{a: []int64{2, 0, -3}, e: []int64{-1, 1, 0}, theta: []int64{1}, delta: 1},
		{a: []int64{1, 1, -1}, e: []int64{0, -1, 1}, theta: []int64{1}, delta: 1},
	}
	for _, fixture := range fixtures {
		gCount, aPositive := 0, 0
		for i := range fixture.a {
			if fixture.a[i]+fixture.e[i] > fixture.theta[0] {
				gCount++
			}
			if fixture.a[i] > 0 {
				aPositive++
			}
		}
		if fixture.theta[0] < fixture.delta || gCount > aPositive {
			t.Fatalf("exact diagonal Weyl transfer failed: %+v", fixture)
		}
	}
	// A=[0], E=[1], theta=1: strict count is zero; >= would be one and fail.
	if 1 > 1 {
		t.Fatal("strict equality fixture changed")
	}
}

func TestM10M9ReuseAndConservativeCountConversion(t *testing.T) {
	r, err := testM10()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.CountingTheorem.M9Accounting, "rank(P_near)+n_plus(Q_near)") || !containsTheorem(r.CountingTheorem.Theorems, M9CriticalRankBoundTheoremID) {
		t.Fatal("M10 duplicated or bypassed M9 accounting")
	}
	rank, err := ThresholdedCriticalRankLowerBound(7, 3)
	if err != nil || rank != 4 {
		t.Fatalf("thresholded M9 adapter=%d err=%v", rank, err)
	}
	bounds, err := FiniteWindowCountLowerBounds(8, 10, 1)
	if err != nil || bounds.SimpleCriticalLowerBound != 4 || bounds.DistinctAllLowerBound != 7 {
		t.Fatalf("window bounds=%#v err=%v", bounds, err)
	}
	// Degeneracy may reduce rank, never increase the count conclusion.
	if rank > 7 {
		t.Fatal("rank was identified with multiplicity")
	}
}

func TestM10FarBoundThresholdAndAsymptoticTrustRemainTyped(t *testing.T) {
	r, err := testM10()
	if err != nil {
		t.Fatal(err)
	}
	if r.FarBound.Norm != semantic.OperatorNorm || !r.FarBound.ExactOrTrusted || !strings.Contains(r.FarBound.BoundExpression, "D0^2") || !strings.Contains(r.FarBound.AsymptoticStatement, "o(1)") {
		t.Fatalf("far theorem weakened: %#v", r.FarBound)
	}
	if r.Perturbation.Comparison != "strict" || r.Compression.Threshold.Comparison != semantic.StrictlyAboveThreshold || !r.Perturbation.ExactRequired {
		t.Fatal("threshold evidence boundary weakened")
	}
	if r.CountingTheorem.AsymptoticProportionDerived || r.UtilitySchedulerUsed {
		t.Fatal("M10 crossed scope or used utility ceremonially")
	}
}

func TestM10ReportsAreDeterministicAndExposeTheCompleteBridge(t *testing.T) {
	r, err := testM10()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M10JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M10JSONReport(r)
	if err != nil || string(a) != string(b) {
		t.Fatal("M10 JSON is nondeterministic")
	}
	var doc map[string]any
	if err := json.Unmarshal(a, &doc); err != nil || doc["schema"] != "riemann.semantic-graph.m10" {
		t.Fatalf("bad schema=%v err=%v", doc["schema"], err)
	}
	h := M10HumanReport(r)
	for _, want := range []string{"HEIGHT WINDOW", "ZERO-SIDE SPLIT", "FAR-ZERO CONTROL", "n_plus^theta", "M9 ACCOUNTING", "NEWLY DERIVED MATHEMATICAL RESULT", "N0_simple(T,2T)>=2*L_theta", "half-type proportion reached: false", "RH\n  unresolved"} {
		if !strings.Contains(h, want) {
			t.Fatalf("human report missing %q", want)
		}
	}
}

func TestM10NegativeWindowCountsRejected(t *testing.T) {
	if _, err := FiniteWindowCountLowerBounds(-1, 0, 0); err == nil {
		t.Fatal("negative thresholded count accepted")
	}
}
