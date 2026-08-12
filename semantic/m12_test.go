package semantic

import "testing"

func er(n int64) ExactRational { return ExactRational{Numerator: n, Denominator: 1} }

func TestRankTraceExactEqualityAndNegativeTrace(t *testing.T) {
	eq, err := CheckDiagonalRankTrace([]ExactRational{er(1), er(0)}, []ExactRational{er(0), er(2)}, 1, 1, er(2))
	if err != nil || eq.Slack.Numerator != 0 || eq.FrobeniusSquared.Numerator != 5 || eq.RightHandSide.Numerator != 5 {
		t.Fatalf("equality fixture failed: %+v %v", eq, err)
	}
	neg, err := CheckDiagonalRankTrace([]ExactRational{er(1)}, []ExactRational{er(-2)}, 1, 0, er(2))
	if err != nil || neg.Slack.Numerator != 8 {
		t.Fatalf("negative trace(Q) fixture failed: %+v %v", neg, err)
	}
}

func TestRankTracePremisesRejectInvalidInputs(t *testing.T) {
	if _, err := CheckDiagonalRankTrace([]ExactRational{er(-10)}, []ExactRational{er(12)}, 1, 1, er(2)); err == nil {
		t.Fatal("non-PSD P accepted")
	}
	if _, err := CheckDiagonalRankTrace([]ExactRational{er(1)}, []ExactRational{er(0)}, 0, 0, er(2)); err == nil {
		t.Fatal("too-small rank budget accepted")
	}
	if _, err := CheckDiagonalRankTrace([]ExactRational{er(0)}, []ExactRational{er(2)}, 0, 0, er(2)); err == nil {
		t.Fatal("too-small positive-index budget accepted")
	}
	if _, err := CheckDiagonalRankTrace([]ExactRational{er(0)}, []ExactRational{er(0)}, 0, 0, er(0)); err == nil {
		t.Fatal("c=0 accepted")
	}
}

func TestRankTraceEndpointBudgets(t *testing.T) {
	if got, err := CheckDiagonalRankTrace([]ExactRational{er(0)}, []ExactRational{er(-2)}, 0, 0, er(2)); err != nil || got.Slack.Numerator < 0 {
		t.Fatalf("r=b=0 endpoint failed: %+v %v", got, err)
	}
	if got, err := CheckDiagonalRankTrace([]ExactRational{er(0), er(0)}, []ExactRational{er(2), er(0)}, 0, 1, er(2)); err != nil || got.Slack.Numerator != 0 {
		t.Fatalf("dependent semidefinite block failed: %+v %v", got, err)
	}
}

func TestApproximateDecompositionIsNotTheoremEvidence(t *testing.T) {
	d := HermitianComponentDecomposition{
		TotalMatrixID: "G", PSDMatrixID: "P", IndexMatrixID: "Q", Identity: "G=P+Q",
		SameDimension: true, TotalHermitian: true, PSDHermitian: true, QHermitian: true, PSD: true,
		Evidence: ApproximateEigenEvidence, IdentityTheorem: "candidate", Provenance: Reference{Kind: CompilerRecord, Citation: "M12 test"},
	}
	if err := d.Validate(); err == nil {
		t.Fatal("approximate decomposition promoted to theorem premise")
	}
}
