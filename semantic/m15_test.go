package semantic

import "testing"

func validM15Family() BoundaryAtomCompletion {
	return BoundaryAtomCompletion{
		Constant: ExactRational{Numerator: 9, Denominator: 8}, ExteriorDensity: ExactRational{Numerator: 9, Denominator: 8},
		BoundaryAtomMass: ExactRational{Numerator: 1, Denominator: 8}, DensitySupport: "|x|>1", AtomSupport: []int{-1, 1}, SupportLegal: true, MeasureNonnegative: true,
	}
}

func validM15PD() WholeLinePDCertificate {
	return WholeLinePDCertificate{
		ID: "fixture", FourierVariable: "t=2*pi*xi", FourierDensity: "hat(P)(xi)=1+2(cos(t)-1)/t^2+(cos(t)-sin(t)/t)/4, with value 0 at t=0", ClearedDensity: "G(t)=4*t^2*hat(P)(xi)=t^2(4+cos(t))-t*sin(t)-8(1-cos(t))",
		InnerInterval: "0<=|t|<=pi", InnerSeries: "G(t)=sum_{n>=3} (-1)^(n-1) 4(n-2)(n+1)t^(2n)/(2n)!", FirstPositiveCoefficient: ExactRational{Numerator: 1, Denominator: 45}, TermRatioUpper: ExactRational{Numerator: 25, Denominator: 56},
		OuterInterval: "|t|>=pi", OuterLowerBound: "G(t)>=3t^2-t-16>0 because |t|>=pi>3, q(3)=8, and q'(t)>0", OuterAnchor: ExactRational{Numerator: 8, Denominator: 1}, OmittedDirections: "inner alternating series plus outer coercive bound covers every real t", WholeLine: true,
	}
}

func TestM15BoundaryAtomSupportAndMeasurePositivityAreSeparate(t *testing.T) {
	f := validM15Family()
	if err := f.Validate(); err != nil {
		t.Fatal(err)
	}
	f.SupportLegal = false
	if err := f.Validate(); err == nil {
		t.Fatal("illegal support certified")
	}
	f = validM15Family()
	f.BoundaryAtomMass = ExactRational{Numerator: -1, Denominator: 8}
	if err := f.Validate(); err == nil {
		t.Fatal("negative sigma certified")
	}
}

func TestM15OriginConditionRejectsBeyondFamilyCeiling(t *testing.T) {
	zero, err := BoundaryAtomLocalQuadraticNumerator(ExactRational{Numerator: 9, Denominator: 8})
	if err != nil || zero.Sign() != 0 {
		t.Fatalf("endpoint coefficient: %v %v", zero, err)
	}
	negative, err := BoundaryAtomLocalQuadraticNumerator(ExactRational{Numerator: 563, Denominator: 500})
	if err != nil || negative.Sign() >= 0 {
		t.Fatalf("c=1.126 must fail locally: %v %v", negative, err)
	}
}

func TestM15WholeLineCertificateRequiresBothRegionsAndTailRatio(t *testing.T) {
	f, p := validM15Family(), validM15PD()
	if err := VerifyBoundaryAtomCertificate(f, p); err != nil {
		t.Fatal(err)
	}
	p.WholeLine = false
	if err := VerifyBoundaryAtomCertificate(f, p); err == nil {
		t.Fatal("missing outer directions certified")
	}
	p = validM15PD()
	p.TermRatioUpper = ExactRational{Numerator: 1, Denominator: 1}
	if err := VerifyBoundaryAtomCertificate(f, p); err == nil {
		t.Fatal("nondecreasing omitted series certified")
	}
}

func TestM15GridEvidenceStillCannotEnterCertifiedRoute(t *testing.T) {
	w := DualCompletionWitness{ID: "grid-nine-eighths", MultiplicityLower: ExactRational{Numerator: 9, Denominator: 8}, OutsideMeasure: "legal", CompletionDistribution: "sampled", FourierImage: "sampled", PositivityEvidence: GridPositivityOnly, ExactSupportScope: true, WholeLineControl: false, TailControl: false, Provenance: Reference{Kind: CompilerRecord, Citation: "fixture"}}
	if _, err := ApplyDualCompletion(testSupportOneClass(), w); err == nil {
		t.Fatal("finite grid became theorem evidence")
	}
}
