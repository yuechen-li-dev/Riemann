package semantic

import (
	"strings"
	"testing"
)

func validM18Certificate() OneRadiusTangencyWholeLineCertificate {
	return OneRadiusTangencyWholeLineCertificate{
		ID: "m18-exact-one-radius-extremal", PiLower: ExactRational{103993, 33102}, PiUpper: ExactRational{104348, 33215},
		LocalInterval: "3/2<=t<=17/10", LocalCellKind: NonnegativeCellWithCertifiedContact, LocalConvexityLower: ExactRational{107, 2000},
		CompactIntervals: []string{"0<=t<=3/2", "17/10<=t<=4"}, CompactCellKind: StrictPositiveCell,
		GridStep: ExactRational{1, 2500}, TaylorDegree: 40, LipschitzUpper: ExactRational{83, 100},
		TailInterval: "|t|>=4", TailLowerBound: ExactRational{1, 3},
		Coverage: "evenness; strict Taylor/Lipschitz cells on [0,3/2] and [17/10,4]; certified quadratic contact on [3/2,17/10]; strict analytic tail on [4,infinity)", WholeLine: true,
	}
}

func TestM18TangencyAlgebraAndDerivedWeight(t *testing.T) {
	c := DeriveOneRadiusTangencyCandidate()
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
	if c.ExactCExpression.Expression != "1+pi/4-2/pi" || c.ExactWExpression.Expression != "1/pi-4/pi^2+8/pi^3" {
		t.Fatalf("wrong exact candidate: %+v", c)
	}
	p, err := VerifyTangencyAlgebra(c)
	if err != nil {
		t.Fatal(err)
	}
	if !p.Verified || !strings.Contains(p.ValueEquation, "=0") || !strings.Contains(p.DerivativeEquation, "=0") || c.ContactMultiplicity != 2 {
		t.Fatal("exact quadratic contact was not certified")
	}
}

func TestM18ParameterLegalityOriginAndMultiplicity(t *testing.T) {
	c := DeriveOneRadiusTangencyCandidate()
	if evalPi(c.ExactWExpression).lo.Sign() <= 0 || evalPi(c.ExactCExpression).lo.Sign() <= 0 {
		t.Fatal("candidate parameter is not positive")
	}
	taylor := TangencyOriginTaylor(c)
	if evalPi(taylor.Constant).lo.Cmp(rat(ExactRational{1, 25})) <= 0 {
		t.Fatal("candidate is not safely unsaturated")
	}
	contact, err := VerifyTangencyAlgebra(c)
	if err != nil {
		t.Fatal(err)
	}
	if evalPi(contact.SecondDerivativeExpression).lo.Cmp(rat(ExactRational{1, 8})) <= 0 {
		t.Fatal("contact is not a strict quadratic minimum")
	}
}

func TestM18EqualityAwareWholeLineVerifier(t *testing.T) {
	if err := VerifyTangencyWholeLineCertificate(DeriveOneRadiusTangencyCandidate(), validM18Certificate()); err != nil {
		t.Fatal(err)
	}
}

func TestM18CertificateRejectsCoverageAndContactDowngrades(t *testing.T) {
	c := DeriveOneRadiusTangencyCandidate()
	for _, mutate := range []func(*OneRadiusTangencyWholeLineCertificate){
		func(p *OneRadiusTangencyWholeLineCertificate) { p.LocalCellKind = StrictPositiveCell },
		func(p *OneRadiusTangencyWholeLineCertificate) { p.CompactIntervals[1] = "18/10<=t<=4" },
		func(p *OneRadiusTangencyWholeLineCertificate) { p.TailInterval = "|t|>4" },
	} {
		p := validM18Certificate()
		mutate(&p)
		if VerifyTangencyWholeLineCertificate(c, p) == nil {
			t.Fatal("invalid equality-aware coverage accepted")
		}
	}
}

func TestM18M17UpperCompositionAndAboveCeilingRejection(t *testing.T) {
	o, err := ComposeExactOneRadiusOptimum(DeriveOneRadiusTangencyCandidate(), validM18Certificate(), DeriveOneRadiusFamilyCeiling())
	if err != nil || !o.EqualityDerived || o.Family != "UnsaturatedOneRadius" {
		t.Fatalf("exact composition failed: %+v %v", o, err)
	}
	if err := RejectAboveExactOneRadiusCeiling(ExactRational{1, 1000000}); err == nil || !strings.Contains(err.Error(), "F(pi/2)") {
		t.Fatal("above-ceiling adversary was not rejected at the active frequency")
	}
}
