package semantic

import "testing"

func m17Ref() Reference { return Reference{Kind: CompilerRecord, Citation: "M17 test"} }

func validM17Family() UnsaturatedOneRadiusCompletion {
	return NewUnsaturatedOneRadiusCompletion(ExactRational{2297, 2000}, ExactRational{1, 1}, ExactRational{171, 1000}, m17Ref())
}

func validM17PD() OneRadiusWholeLineCertificate {
	return OneRadiusWholeLineCertificate{
		ID: "test", FourierVariable: "t=2*pi*xi", FourierDensity: OneRadiusFourierDensity,
		CompactInterval: "0<=|t|<=4", GridStep: ExactRational{1, 5000}, TaylorDegree: 40,
		LipschitzBound: ExactRational{4943, 6000}, TailInterval: "|t|>=4",
		TailLowerBound: "hat(P)>=1-2*w-4/t^2-2*(c-1)/|t|>0 for |t|>=4", TailAnchor: ExactRational{267, 800},
		OmittedDirections: "exact Taylor enclosures plus a Lipschitz cover on |t|<=4 and a radius-uniform analytic bound on |t|>=4 cover every real t", WholeLine: true,
	}
}

func TestM17FamilyLegalityAndExactTaylor(t *testing.T) {
	f := validM17Family()
	if err := f.Validate(); err != nil {
		t.Fatal(err)
	}
	want := OneRadiusTaylor{ExactRational{9, 200}, ExactRational{-229, 6000}, ExactRational{3239, 360000}, ExactRational{-1847, 5040000}}
	if f.Taylor != want {
		t.Fatalf("wrong Taylor coefficients: %+v", f.Taylor)
	}
	for _, mutate := range []func(*UnsaturatedOneRadiusCompletion){
		func(x *UnsaturatedOneRadiusCompletion) {
			x.Constant = ExactRational{0, 1}
			x.ExteriorDensity = x.Constant
			x.Taylor = deriveOneRadiusTaylor(x.Constant, x.Radius, x.Weight)
		},
		func(x *UnsaturatedOneRadiusCompletion) {
			x.Radius = ExactRational{9, 10}
			x.Taylor = deriveOneRadiusTaylor(x.Constant, x.Radius, x.Weight)
		},
		func(x *UnsaturatedOneRadiusCompletion) {
			x.Weight = ExactRational{-1, 100}
			x.Taylor = deriveOneRadiusTaylor(x.Constant, x.Radius, x.Weight)
		},
	} {
		bad := f
		mutate(&bad)
		if bad.Validate() == nil {
			t.Fatal("illegal parameter accepted")
		}
	}
}

func TestM17SaturationIsCheckedStrictSubtypeAndRecoversM15(t *testing.T) {
	unsaturated := validM17Family()
	if _, err := AsSaturatedOneRadius(unsaturated); err == nil {
		t.Fatal("positive-slack family was marked saturated")
	}
	m15 := NewUnsaturatedOneRadiusCompletion(ExactRational{9, 8}, ExactRational{1, 1}, ExactRational{1, 8}, m17Ref())
	s, err := AsSaturatedOneRadius(m15)
	if err != nil || SaturatedOneRadiusCeiling(s) != nil {
		t.Fatalf("M15 regression failed: %v", err)
	}
	above := NewUnsaturatedOneRadiusCompletion(ExactRational{113, 100}, ExactRational{1, 1}, ExactRational{13, 100}, m17Ref())
	s, err = AsSaturatedOneRadius(above)
	if err != nil || SaturatedOneRadiusCeiling(s) == nil || VerifyOneRadiusLocalLegality(above) == nil {
		t.Fatal("saturated c>9/8 was not rejected")
	}
}

func TestM17PositiveA0AllowsNegativeA2(t *testing.T) {
	f := validM17Family()
	if rat(f.Taylor.Constant).Sign() <= 0 || rat(f.Taylor.Quadratic).Sign() >= 0 {
		t.Fatal("fixture does not exercise unsaturated local semantics")
	}
	if err := VerifyOneRadiusLocalLegality(f); err != nil {
		t.Fatal(err)
	}
}

func TestM17TailAndCompactProofParametersAreRecomputed(t *testing.T) {
	f, p := validM17Family(), validM17PD()
	p.TailAnchor = ExactRational{1, 3}
	if VerifyOneRadiusCertificate(f, p) == nil {
		t.Fatal("incorrect tail anchor accepted")
	}
	p = validM17PD()
	p.LipschitzBound = ExactRational{4, 5}
	if VerifyOneRadiusCertificate(f, p) == nil {
		t.Fatal("incorrect compact derivative bound accepted")
	}
	p = validM17PD()
	p.GridStep = ExactRational{3, 10000}
	if VerifyOneRadiusCertificate(f, p) == nil {
		t.Fatal("grid that leaves a region gap accepted")
	}
}

func TestM17M16OneRadiusRegressionWithoutTwoRadiusObject(t *testing.T) {
	f := NewUnsaturatedOneRadiusCompletion(ExactRational{573, 500}, ExactRational{1, 1}, ExactRational{21, 125}, m17Ref())
	if err := VerifyOneRadiusLocalLegality(f); err != nil {
		t.Fatal(err)
	}
	if f.Taylor.Constant != (ExactRational{11, 250}) || f.Taylor.Quadratic != (ExactRational{-9, 250}) {
		t.Fatalf("wrong M16 one-radius regression coefficients: %+v", f.Taylor)
	}
}

func TestM17FamilyCeilingRejectsAdversarialCandidate(t *testing.T) {
	ceiling := DeriveOneRadiusFamilyCeiling()
	if ceiling.RationalUpper != (ExactRational{1149, 1000}) || !ceiling.StrictForRadiusAboveOne {
		t.Fatal("ceiling theorem changed")
	}
	bad := NewUnsaturatedOneRadiusCompletion(ExactRational{23, 20}, ExactRational{1, 1}, ExactRational{1, 5}, m17Ref())
	if RejectAboveOneRadiusRationalCeiling(bad) == nil {
		t.Fatal("candidate above proved one-radius ceiling accepted")
	}
}
