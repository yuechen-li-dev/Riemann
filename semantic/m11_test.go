package semantic

import "testing"

func q(n int64) ExactRational           { return ExactRational{Numerator: n, Denominator: 1} }
func z(r, i int64) ExactComplexRational { return ExactComplexRational{Real: q(r), Imag: q(i)} }

func TestM11ExactHermitianTraceAndFrobeniusSemantics(t *testing.T) {
	// [[1,2+i],[2-i,-3]] has trace -2 and Frobenius square
	// 1+9+2*(4+1)=20, equal to tr(G^2).
	trace, frob, err := ExactHermitianMoments(2, []ExactComplexRational{z(1, 0), z(2, 1), z(2, -1), z(-3, 0)})
	if err != nil {
		t.Fatal(err)
	}
	if trace != q(-2) || frob != q(20) {
		t.Fatalf("got trace=%+v frob=%+v", trace, frob)
	}
}

func TestM11ExactHermitianMomentsRejectNonHermitianInput(t *testing.T) {
	if _, _, err := ExactHermitianMoments(2, []ExactComplexRational{z(1, 0), z(2, 1), z(2, 1), z(1, 0)}); err == nil {
		t.Fatal("non-Hermitian matrix accepted")
	}
}

func TestM11ThresholdedCountUsesResidualCauchyAndCeiling(t *testing.T) {
	r, err := ThresholdedCountFromMoments(q(6), q(14), 3, q(1))
	if err != nil {
		t.Fatal(err)
	}
	if !r.Applicable || r.TraceResidual != q(3) || r.RealLowerBound != (ExactRational{Numerator: 9, Denominator: 14}) || r.IntegerLowerBound != 1 {
		t.Fatalf("unexpected count result: %+v", r)
	}
	// Equality spectrum diag(1,1,0,0), theta=0: trace^2/F^2=2.
	eq, err := ThresholdedCountFromMoments(q(2), q(2), 4, q(0))
	if err != nil || eq.IntegerLowerBound != 2 || eq.RealLowerBound != q(2) {
		t.Fatalf("equality case: %+v %v", eq, err)
	}
}

func TestM11NegativeEigenvalueCaseAndNonpositiveResidual(t *testing.T) {
	// diag(4,-1,-1): trace=2, F^2=18, exactly one eigenvalue above 0.
	r, err := ThresholdedCountFromMoments(q(2), q(18), 3, q(0))
	if err != nil || r.IntegerLowerBound != 1 {
		t.Fatalf("negative spectrum: %+v %v", r, err)
	}
	vacuous, err := ThresholdedCountFromMoments(q(3), q(3), 3, q(1))
	if err != nil || vacuous.Applicable || vacuous.IntegerLowerBound != 0 {
		t.Fatalf("residual must be vacuous: %+v %v", vacuous, err)
	}
}

func TestM11TheoremRejectsNegativeThresholdAndInfeasiblePremises(t *testing.T) {
	if _, err := ThresholdedCountFromMoments(q(1), q(1), 1, q(-1)); err == nil {
		t.Fatal("negative theta accepted")
	}
	if _, err := ThresholdedCountFromMoments(q(10), q(1), 2, q(0)); err == nil {
		t.Fatal("infeasible moment premises accepted")
	}
}

func TestM11ApproximateAndAsymptoticMomentsCannotEnterFiniteTheorem(t *testing.T) {
	ref := Reference{Kind: ExperimentRecord, Citation: "fixture"}
	value := q(1)
	approx := SpectralMomentClaim{ID: "a", MatrixID: "G", Kind: Trace, Relation: AtLeastBound, Expression: "1.0", ExactValue: &value, Evidence: ApproximateMomentEvidence, Provenance: ref}
	if err := approx.Validate(); err == nil {
		t.Fatal("approximate value laundered as exact")
	}
	asymptotic := SpectralMomentClaim{ID: "b", MatrixID: "G", Kind: Trace, Relation: AtLeastBound, Expression: "A+o(A)", Evidence: AsymptoticMomentEvidence, Theorems: []TheoremID{"source"}, Provenance: ref}
	if err := asymptotic.Validate(); err != nil {
		t.Fatal(err)
	}
	if asymptotic.FiniteTheoremPremise() {
		t.Fatal("raw asymptotic entered finite theorem")
	}
}
