package semantic

import "testing"

func validM16Family() TwoRadiusCompletion {
	return TwoRadiusCompletion{
		Constant: ExactRational{573, 500}, ExteriorDensity: ExactRational{573, 500},
		Atoms:          []SymmetricExteriorAtom{{ExactRational{1, 1}, ExactRational{21, 125}}, {ExactRational{2, 1}, ExactRational{1, 1000}}},
		DensitySupport: "|x|>1", SupportLegal: true, MeasureNonnegative: true,
	}
}

func validM16PD() TwoRadiusWholeLineCertificate {
	return TwoRadiusWholeLineCertificate{
		ID: "test", FourierVariable: "t=2*pi*xi",
		FourierDensity:  "hat(P)(t/(2*pi))=1+2(cos(t)-1)/t^2+2(1-c)sin(t)/t+2*w1*cos(r1*t)+2*w2*cos(r2*t)",
		CompactInterval: "0<=|t|<=4", GridStep: ExactRational{1, 1000}, TaylorDegree: 40, LipschitzBound: ExactRational{1229, 1500},
		TailInterval: "|t|>=4", TailLowerBound: "hat(P)>=1-2(w1+w2)-4/t^2-2(c-1)/|t|>0 for |t|>=4", TailAnchor: ExactRational{339, 1000},
		OmittedDirections: "exact Taylor enclosures plus a Lipschitz cover on |t|<=4 and an analytic bound on |t|>=4 cover every real t", WholeLine: true,
	}
}

func TestM16RadiusPermutationAndCoincidentDegeneration(t *testing.T) {
	f := validM16Family()
	f.Atoms[0], f.Atoms[1] = f.Atoms[1], f.Atoms[0]
	c := f.Canonical()
	if c.Atoms[0].Radius != (ExactRational{1, 1}) {
		t.Fatal("radius permutation did not canonicalize")
	}
	c.Atoms[1].Radius = c.Atoms[0].Radius
	collapsed := c.CollapsedAtoms()
	if len(collapsed) != 1 || collapsed[0].Weight != (ExactRational{169, 1000}) {
		t.Fatalf("coincident radii did not reduce to one radius: %+v", collapsed)
	}
}

func TestM16ParameterLegalityIsIndependent(t *testing.T) {
	f := validM16Family()
	f.Atoms[0].Radius = ExactRational{9, 10}
	if f.Validate() == nil {
		t.Fatal("interior radius accepted")
	}
	f = validM16Family()
	f.Atoms[0].Weight = ExactRational{-1, 1000}
	if f.Validate() == nil {
		t.Fatal("negative weight accepted")
	}
}

func TestM16OriginExpansionAndSaturatedM15Regression(t *testing.T) {
	taylor, err := TwoRadiusOriginTaylor(validM16Family())
	if err != nil {
		t.Fatal(err)
	}
	if taylor.Constant != (ExactRational{23, 500}) || taylor.Quadratic != (ExactRational{-1, 25}) || taylor.Quartic != (ExactRational{911, 90000}) {
		t.Fatalf("wrong exact Taylor coefficients: %+v", taylor)
	}
	f := validM16Family()
	f.Constant, f.ExteriorDensity = ExactRational{9, 8}, ExactRational{9, 8}
	f.Atoms = []SymmetricExteriorAtom{{ExactRational{1, 1}, ExactRational{1, 8}}, {ExactRational{2, 1}, ExactRational{0, 1}}}
	taylor, err = TwoRadiusOriginTaylor(f)
	if err != nil || taylor.Constant.Numerator != 0 || taylor.Quadratic.Numerator != 0 || SaturatedTwoRadiusCeiling(f) != nil {
		t.Fatalf("M15 degeneration failed: %+v %v", taylor, err)
	}
	f.Constant, f.ExteriorDensity = ExactRational{113, 100}, ExactRational{113, 100}
	f.Atoms[0].Weight = ExactRational{13, 100}
	if SaturatedTwoRadiusCeiling(f) == nil || VerifyTwoRadiusLocalLegality(f) == nil {
		t.Fatal("saturated c>9/8 was not rejected locally")
	}
}

func TestM16CandidateWholeLineCertificateAndCoverage(t *testing.T) {
	if err := VerifyTwoRadiusCertificate(validM16Family(), validM16PD()); err != nil {
		t.Fatal(err)
	}
	p := validM16PD()
	p.CompactInterval = "0<=|t|<4"
	if VerifyTwoRadiusCertificate(validM16Family(), p) == nil {
		t.Fatal("region gap accepted")
	}
	p = validM16PD()
	p.TailAnchor = ExactRational{34, 100}
	if VerifyTwoRadiusCertificate(validM16Family(), p) == nil {
		t.Fatal("nonconservative tail claim accepted")
	}
}

func TestM16UnsaturatedOneRadiusDegenerationAlsoCertifies(t *testing.T) {
	f, p := validM16Family(), validM16PD()
	f.Atoms[1].Weight = ExactRational{0, 1}
	p.LipschitzBound = ExactRational{1223, 1500}
	p.TailAnchor = ExactRational{341, 1000}
	if err := VerifyTwoRadiusCertificate(f, p); err != nil {
		t.Fatalf("the exact witness does not depend on the second atom: %v", err)
	}
}
