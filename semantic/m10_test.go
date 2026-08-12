package semantic

import "testing"

var m10TestReference = Reference{Kind: CompilerRecord, Citation: "M10 test fixture"}

func TestM10HeightWindowMembershipUsesLeftOpenRightClosedPositiveConvention(t *testing.T) {
	w, err := NewConcreteHeightWindow("I", 10, 20)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[int64]bool{-20: false, -10: false, 10: false, 11: true, 20: true, 21: false}
	for ordinate, want := range cases {
		got, err := w.ContainsExactOrdinate(ordinate)
		if err != nil || got != want {
			t.Fatalf("ordinate %d: got %t want %t err=%v", ordinate, got, want, err)
		}
	}
	if _, err := (HeightWindow{ID: "symbolic", Center: SymbolicHeight("3T/2"), HalfWidth: SymbolicHeight("T/2"), Lower: SymbolicHeight("T"), Upper: SymbolicHeight("2T"), Boundary: LeftOpenRightClosed, OrdinateConvention: PositiveOrdinateConvention}).ContainsExactOrdinate(15); err == nil {
		t.Fatal("symbolic membership was guessed numerically")
	}
}

func TestM10WindowOrbitCountingSeparatesMultiplicityPairsLocationsAndRank(t *testing.T) {
	w, _ := NewConcreteHeightWindow("I", 10, 20)
	counts, err := CountWindowZeros(w, []WindowZeroPoint{
		{LocationID: "critical-simple", Ordinate: 11, Multiplicity: 1, CriticalLine: true},
		{LocationID: "critical-multiple", Ordinate: 20, Multiplicity: 3, CriticalLine: true},
		{LocationID: "left-boundary", Ordinate: 10, Multiplicity: 9, CriticalLine: true},
		{LocationID: "off-left", Ordinate: 15, Multiplicity: 2, ReflectionPairID: "pair-1"},
		{LocationID: "off-right", Ordinate: 15, Multiplicity: 2, ReflectionPairID: "pair-1"},
		{LocationID: "conjugate-below", Ordinate: -15, Multiplicity: 2, ReflectionPairID: "pair-conjugate"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if counts.TotalMultiplicity != 8 || counts.CriticalMultiplicity != 4 || counts.CriticalDistinctLocations != 2 || counts.SimpleCriticalLocations != 1 || counts.OffCriticalMultiplicity != 4 || counts.OffCriticalReflectionPairs != 1 || counts.DistinctAllLocations != 4 {
		t.Fatalf("count notions conflated: %#v", counts)
	}
	counts.EvaluationRankDirections = 1
	if err := counts.Validate(); err != nil {
		t.Fatalf("conservative rank direction rejected: %v", err)
	}
	counts.EvaluationRankDirections = 3
	if counts.Validate() == nil {
		t.Fatal("rank directions exceeded distinct critical locations")
	}
}

func TestM10ThresholdedPositiveIndexIsStrictAndExactEvidenceOnly(t *testing.T) {
	one := ExactRational{Numerator: 1, Denominator: 1}
	count, err := ExactDiagonalThresholdedPositiveIndex([]ExactRational{{Numerator: 2, Denominator: 1}, one, {Numerator: 1, Denominator: 2}, {Numerator: -1, Denominator: 1}}, one)
	if err != nil || count != 1 {
		t.Fatalf("strict threshold count=%d err=%v", count, err)
	}
	threshold := SpectralThreshold{ID: "theta", Kind: AbsoluteThresholdScale, Expression: "1", ExactValue: &one, Comparison: StrictlyAboveThreshold, Dependencies: map[string]string{"scale": "fixture"}, Provenance: m10TestReference}
	claim := ThresholdedPositiveIndexClaim{MatrixID: "D", Dimension: 4, Threshold: threshold, Relation: EqualBound, Bound: count, Evidence: ExactTheoremEvidence, Theorems: []TheoremID{"exact-diagonal"}, Provenance: m10TestReference}
	if err := claim.Validate(); err != nil || !claim.ExactTheoremPremise() {
		t.Fatalf("exact threshold claim rejected: %v", err)
	}
	claim.Evidence = ApproximateEigenEvidence
	claim.Theorems = nil
	if err := claim.Validate(); err != nil {
		t.Fatalf("report-only approximate claim rejected: %v", err)
	}
	if claim.ExactTheoremPremise() {
		t.Fatal("approximate eigenvalues crossed the exact theorem boundary")
	}
}

func TestM10FarAndOffCriticalAreIndependentAxes(t *testing.T) {
	near, _ := NewConcreteHeightWindow("near", 5, 25)
	counts, err := CountWindowZeros(near, []WindowZeroPoint{
		{LocationID: "near-off-a", Ordinate: 15, Multiplicity: 1, ReflectionPairID: "near-off"},
		{LocationID: "near-off-b", Ordinate: 15, Multiplicity: 1, ReflectionPairID: "near-off"},
		{LocationID: "far-critical", Ordinate: 30, Multiplicity: 1, CriticalLine: true},
	})
	if err != nil || counts.OffCriticalReflectionPairs != 1 || counts.CriticalMultiplicity != 0 {
		t.Fatalf("far/off axes conflated: %#v err=%v", counts, err)
	}
}
