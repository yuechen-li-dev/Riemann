package semantic

import (
	"math"
	"strings"
	"testing"
)

func TestM23QSystemCandidateAndIndependentStationarity(t *testing.T) {
	m, r, x := 0.0069950902804526323, 3.5078534077665942, 1.795313437115535
	f, err := DeriveContactPreservingTwoAtomFamily(ExactRational{6995090280452632, 1000000000000000000}, ExactRational{3507853407766594, 1000000000000000})
	if err != nil {
		t.Fatal(err)
	}
	q1 := ContactPreservingTwoAtomSpectrum(f, x)
	q2 := ContactPreservingTwoAtomSpectrumPrime(f, x)
	a := math.Pi / 2
	A, S := math.Cos(a*r), math.Sin(a*r)
	Ap := -a * S
	Bprime := -2*S - a*r*A
	h := TwoAtomResponseFloat(r, x)
	hr := -math.Pi*Ap*math.Sin(x)/x + 2*Bprime*math.Cos(x) - 2*x*math.Sin(r*x)
	q3 := Ap*h - A*hr
	if math.Abs(q1) > 1e-10 || math.Abs(q2) > 1e-9 || math.Abs(q3) > 1e-9 || m <= 0 {
		t.Fatalf("Q candidate residuals: %g %g %g", q1, q2, q3)
	}
}

func TestM23CertifiedJointRootAndLocalMaximum(t *testing.T) {
	s, b, e, err := CertifyVariableRadiusOptimizer()
	if err != nil {
		t.Fatal(err)
	}
	if !b.ExistsUnique || !e.CertifiedLocalMaximum || !strings.Contains(s.Q3, "A'(r)H") {
		t.Fatalf("incomplete M23 certificate: %+v %+v %+v", s, b, e)
	}
	if b.Mass.Lower != "0.0069950800000000" || b.Radius.Upper != "3.5078570000000000" || b.Contact.Upper != "1.7953170000000000" {
		t.Fatalf("unexpected certified box: %+v", b)
	}
}

func TestM23M22RegressionAndFixedContact(t *testing.T) {
	f, err := DeriveContactPreservingTwoAtomFamily(ExactRational{17, 2500}, ExactRational{7, 2})
	if err != nil {
		t.Fatal(err)
	}
	x := math.Pi / 2
	if math.Abs(ContactPreservingTwoAtomSpectrum(f, x)) > 2e-15 || math.Abs(ContactPreservingTwoAtomSpectrumPrime(f, x)) > 2e-15 {
		t.Fatal("M22 fixed contact regression failed")
	}
}

func TestM23ExplicitWitnessMetadataAndFastProof(t *testing.T) {
	w := M23ExplicitWitness()
	if err := VerifyM23ExplicitWitness(w, false); err != nil {
		t.Fatal(err)
	}
	if !w.ImprovesM22 || w.Witness.Family.Mass != (ExactRational{69, 10000}) {
		t.Fatalf("wrong M23 witness: %+v", w)
	}
}
