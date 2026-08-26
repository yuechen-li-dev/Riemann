package semantic

import (
	"math"
	"strings"
	"testing"
)

func validM22Certificate() TwoAtomWholeLineCertificate {
	return TwoAtomWholeLineCertificate{
		ID:              "m22-r7-over-2-m17-over-2500-degree100",
		OriginalContact: TwoAtomCompactPiece{Left: ExactRational{154, 100}, Right: ExactRational{160, 100}, GridStep: ExactRational{1, 10000}, TaylorDegree: 100, Kind: NonnegativeCellWithCertifiedContact},
		ContactGridStep: ExactRational{1, 10000}, ContactThirdAbsUpper: ExactRational{6, 5}, ContactConvexityLower: ExactRational{1, 1000},
		StrictPieces: []TwoAtomCompactPiece{
			{Left: ExactRational{0, 1}, Right: ExactRational{7, 5}, GridStep: ExactRational{1, 5000}, TaylorDegree: 100, Kind: StrictPositiveCell},
			{Left: ExactRational{7, 5}, Right: ExactRational{154, 100}, GridStep: ExactRational{1, 100000}, TaylorDegree: 100, Kind: StrictPositiveCell},
			{Left: ExactRational{160, 100}, Right: ExactRational{17, 10}, GridStep: ExactRational{1, 100000}, TaylorDegree: 100, Kind: StrictPositiveCell},
			{Left: ExactRational{17, 10}, Right: ExactRational{19, 10}, GridStep: ExactRational{1, 10000}, TaylorDegree: 100, Kind: StrictPositiveCell},
			{Left: ExactRational{19, 10}, Right: ExactRational{4, 1}, GridStep: ExactRational{1, 5000}, TaylorDegree: 100, Kind: StrictPositiveCell},
		},
		LipschitzUpper: ExactRational{19, 20}, TailStart: ExactRational{4, 1}, TailLower: ExactRational{1, 4},
		Coverage: "evenness; strict rational Taylor cells on [0,154/100] and [160/100,4]; certified convex original-contact cell [154/100,160/100]; radius-uniform analytic tail [4,infinity)", WholeLine: true,
	}
}

func TestM22M21StructuralRegressionAndContactPreservation(t *testing.T) {
	m := ExactRational{19, 5000}
	f, err := DeriveM21Subfamily(m)
	if err != nil {
		t.Fatal(err)
	}
	ray, err := DeriveContactAdaptedAtomicRay(m)
	if err != nil {
		t.Fatal(err)
	}
	c, w := TwoAtomParametersFloat(f)
	wantC := 1 + math.Pi/4 - 2/math.Pi + math.Pi*19/10000
	wantW := 1/math.Pi - 4/(math.Pi*math.Pi) + 8/(math.Pi*math.Pi*math.Pi) + 38/(5000*math.Pi)
	if math.Abs(c-wantC) > 1e-14 || math.Abs(w-wantW) > 1e-14 || f.Radius != ray.Radius2 || f.Mass != ray.Mass {
		t.Fatalf("r=4 does not recover M21: c=%g w=%g", c, w)
	}
	a := math.Pi / 2
	if math.Abs(ContactPreservingTwoAtomSpectrum(f, a)) > 2e-15 || math.Abs(ContactPreservingTwoAtomSpectrumPrime(f, a)) > 2e-15 {
		t.Fatal("derived M22 surface does not preserve the M18 contact")
	}
}

func TestM22LegalityAndNoSingularRadiusBranches(t *testing.T) {
	for _, tc := range []struct{ m, r ExactRational }{
		{ExactRational{-1, 100}, ExactRational{7, 2}},
		{ExactRational{1, 100}, ExactRational{1, 1}},
		{ExactRational{1, 100}, ExactRational{0, 1}},
		{ExactRational{1, 10}, ExactRational{5, 1}}, // w1<0
		{ExactRational{1, 1}, ExactRational{2, 1}},  // c<0
	} {
		if _, err := DeriveContactPreservingTwoAtomFamily(tc.m, tc.r); err == nil {
			t.Fatalf("accepted illegal point %+v", tc)
		}
	}
	if !strings.Contains(ContactPreservationSingularRadii(), "none") || !strings.Contains(ContactPreservationSingularRadii(), "8/pi") {
		t.Fatal("singular-branch classification is incomplete")
	}
}

func TestM22RadiusSensitivityMatchesFiniteDifference(t *testing.T) {
	m := ExactRational{387226734364153, 100000000000000000}
	base, _ := DeriveContactPreservingTwoAtomFamily(m, ExactRational{4, 1})
	t1, eps := 1.36474957003247, 1e-5
	plus, _ := DeriveContactPreservingTwoAtomFamily(m, ExactRational{400001, 100000})
	minus, _ := DeriveContactPreservingTwoAtomFamily(m, ExactRational{399999, 100000})
	fd := (ContactPreservingTwoAtomSpectrum(plus, t1) - ContactPreservingTwoAtomSpectrum(minus, t1)) / (2 * eps)
	analytic := DeriveM21ContactMotionEvidence(float64(m.Numerator)/float64(m.Denominator), t1).DeltaFAtOldContact
	if math.Abs(fd-analytic) > 2e-10 || math.Abs(ContactPreservingTwoAtomSpectrum(base, math.Pi/2)) > 2e-15 {
		t.Fatalf("radius sensitivity mismatch fd=%g analytic=%g", fd, analytic)
	}
}

func TestM22NumericalCandidateHasExpectedDoubleContact(t *testing.T) {
	f, err := DeriveContactPreservingTwoAtomFamily(ExactRational{6995089360323918, 1000000000000000000}, ExactRational{3507853493599018, 1000000000000000})
	if err != nil {
		t.Fatal(err)
	}
	t1 := 1.7953133218061055
	if math.Abs(ContactPreservingTwoAtomSpectrum(f, t1)) > 2e-14 || math.Abs(ContactPreservingTwoAtomSpectrumPrime(f, t1)) > 2e-8 || ContactPreservingTwoAtomSpectrumSecond(f, t1) <= 0 {
		t.Fatalf("numerical candidate lacks the expected double contact: F=%g F'=%g F''=%g", ContactPreservingTwoAtomSpectrum(f, t1), ContactPreservingTwoAtomSpectrumPrime(f, t1), ContactPreservingTwoAtomSpectrumSecond(f, t1))
	}
}

func TestM22CertifiedWitnessFastStructure(t *testing.T) {
	f, err := DeriveContactPreservingTwoAtomFamily(ExactRational{17, 2500}, ExactRational{7, 2})
	if err != nil {
		t.Fatal(err)
	}
	w := CertifiedTwoAtomWitness{Family: f, WholeLine: validM22Certificate(), CExact: "1+pi/4-2/pi+17*pi*sqrt(2)/10000", JExact: "1-pi/4+2/pi-17*pi*sqrt(2)/10000", Verified: true}
	if err := VerifyCertifiedTwoAtomWitness(w, false); err != nil {
		t.Fatal(err)
	}
	if ContactPreservingTwoAtomSpectrumSecond(f, math.Pi/2) <= 0 {
		t.Fatal("original contact lost positive curvature")
	}
	if ContactPreservingTwoAtomSpectrum(f, 1.8034) <= 0 {
		t.Fatal("released contact was not lifted")
	}
}
