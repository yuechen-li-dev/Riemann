package semantic

import "testing"

func validM21Certificate() ContactAdaptedRayWholeLineCertificate {
	return ContactAdaptedRayWholeLineCertificate{
		ID:              "m21-contact-adapted-ray-19-over-5000",
		ContactPiece:    RayCompactPiece{Left: ExactRational{151, 100}, Right: ExactRational{17, 10}, GridStep: ExactRational{1, 10000}, TaylorDegree: 100, Kind: NonnegativeCellWithCertifiedContact},
		ContactGridStep: ExactRational{1, 10000}, ContactThirdAbsUpper: ExactRational{11, 10}, ContactConvexityLower: ExactRational{1, 10000},
		StrictPieces: []RayCompactPiece{
			{Left: ExactRational{0, 1}, Right: ExactRational{13, 10}, GridStep: ExactRational{1, 10000}, TaylorDegree: 100, Kind: StrictPositiveCell},
			{Left: ExactRational{13, 10}, Right: ExactRational{151, 100}, GridStep: ExactRational{1, 200000}, TaylorDegree: 100, Kind: StrictPositiveCell},
			{Left: ExactRational{17, 10}, Right: ExactRational{4, 1}, GridStep: ExactRational{1, 5000}, TaylorDegree: 100, Kind: StrictPositiveCell},
		},
		LipschitzUpper: ExactRational{9, 10}, TailStart: ExactRational{4, 1}, TailLower: ExactRational{3, 10},
		Coverage: "evenness; strict Taylor/Lipschitz cells on [0,13/10], [13/10,151/100], and [17/10,4]; certified convex contact cell [151/100,17/10]; analytic atomic tail [4,infinity)", WholeLine: true,
	}
}

func TestM21ExactRayAndContactAlgebra(t *testing.T) {
	r, err := DeriveContactAdaptedAtomicRay(ExactRational{19, 5000})
	if err != nil {
		t.Fatal(err)
	}
	a, err := VerifyContactAdaptedRayAlgebra(r)
	if err != nil {
		t.Fatal(err)
	}
	if !a.Verified || a.ResponseAtContact != "-2+2=0" || a.ResponsePrimeAtContact != "4/pi-4/pi=0" {
		t.Fatalf("incomplete exact contact algebra: %+v", a)
	}
	base, err := DeriveContactAdaptedAtomicRay(ExactRational{0, 1})
	if err != nil || !samePiValue(base.C, DeriveOneRadiusTangencyCandidate().ExactCExpression) || !samePiValue(base.Weight1, DeriveOneRadiusTangencyCandidate().ExactWExpression) {
		t.Fatalf("m=0 does not recover M18: %v", err)
	}
}

func TestM21WholeLineCertificateStructure(t *testing.T) {
	r, err := DeriveContactAdaptedAtomicRay(ExactRational{19, 5000})
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyContactAdaptedRayWholeLine(r, validM21Certificate(), false); err != nil {
		t.Fatal(err)
	}
}

func TestM21CertifiedRayCeilingUpper(t *testing.T) {
	b := CertifiedRayCeilingBracket{LowerWitness: ExactRational{19, 5000}, UpperBound: ExactRational{3873, 1000000}, UpperWitnessFrequency: ExactRational{5459, 4000}}
	if err := VerifyRayCeilingUpper(b, 100); err != nil {
		t.Fatal(err)
	}
}

func TestM21NumericalEvidenceCannotClaimTheorem(t *testing.T) {
	if err := VerifyRayCeilingUpper(CertifiedRayCeilingBracket{LowerWitness: ExactRational{19, 5000}, UpperBound: ExactRational{3873, 1000000}, UpperWitnessFrequency: ExactRational{0, 1}}, 100); err == nil {
		t.Fatal("a numerical ceiling without a certified failing frequency became theorem evidence")
	}
}
