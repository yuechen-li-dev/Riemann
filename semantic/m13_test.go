package semantic

import (
	"math/big"
	"testing"
)

func TestM13TypedScaleChangeSeparatesTraceAndFrobeniusSquare(t *testing.T) {
	s := ScaleChange{
		SourceObject: "G_tilde", TargetObject: "G_hat",
		Factor:                ScalarBinary(ScalarDivide, ScalarRat(1, 1), ScalarParam("aL")),
		ParameterDependencies: []string{"a", "L"},
		Provenance:            Reference{Kind: CompilerRecord, Citation: "M13 scale test"},
	}
	parameters := map[string]ExactRational{"aL": {Numerator: 3, Denominator: 1}}
	trace, err := s.ApplyExact(Trace, ExactRational{Numerator: 6, Denominator: 1}, parameters)
	if err != nil || trace != (ExactRational{Numerator: 2, Denominator: 1}) {
		t.Fatalf("trace did not scale linearly: %+v %v", trace, err)
	}
	frob, err := s.ApplyExact(FrobeniusNormSquared, ExactRational{Numerator: 9, Denominator: 1}, parameters)
	if err != nil || frob != (ExactRational{Numerator: 1, Denominator: 1}) {
		t.Fatalf("Frobenius square did not scale quadratically: %+v %v", frob, err)
	}
	wrongLinear, err := s.ApplyExact(Trace, ExactRational{Numerator: 9, Denominator: 1}, parameters)
	if err != nil || wrongLinear == frob {
		t.Fatal("linear scale was accidentally accepted for Frobenius square")
	}
}

func TestM13ScalarExpressionIsTypedAndBounded(t *testing.T) {
	e := ScalarBinary(ScalarDivide,
		ScalarBinary(ScalarAdd, ScalarParam("lambda"), ScalarRat(1, 3)),
		ScalarPow(ScalarParam("lambda"), 2),
	)
	v, err := e.EvalExact(map[string]ExactRational{"lambda": {Numerator: 1, Denominator: 1}})
	if err != nil || v != (ExactRational{Numerator: 4, Denominator: 3}) {
		t.Fatalf("typed exact arithmetic failed: %+v %v", v, err)
	}
	if err := ScalarPow(ScalarParam("lambda"), 3).Validate(); err == nil {
		t.Fatal("unneeded general powers entered the bounded IR")
	}
	if _, err := ScalarUnary(ScalarCot, ScalarParam("lambda")).EvalExact(map[string]ExactRational{"lambda": {Numerator: 1, Denominator: 1}}); err == nil {
		t.Fatal("transcendental Float evaluation was promoted to exact arithmetic")
	}
}

func TestM13ParameterDomainRejectsIllegalSupport(t *testing.T) {
	d := ScalarDomain{Lower: ExactRational{Numerator: 0, Denominator: 1}, Upper: ExactRational{Numerator: 1, Denominator: 1}, UpperIncluded: true}
	if err := d.Validate(); err != nil {
		t.Fatal(err)
	}
	if d.ContainsFloat(0) || !d.ContainsFloat(0.5) || !d.ContainsFloat(1) || d.ContainsFloat(1.000001) {
		t.Fatal("0<lambda<=1 domain semantics changed")
	}
}

func TestM13TaylorCertificateStrictlyProvesDisplayedBound(t *testing.T) {
	c, err := MontgomeryTaylorTaylorCertificate()
	if err != nil {
		t.Fatal(err)
	}
	if err := c.ObjectiveInterval.Validate(); err != nil {
		t.Fatal(err)
	}
	lower := new(big.Rat).SetFrac(big.NewInt(c.ObjectiveInterval.Lower.Numerator), big.NewInt(c.ObjectiveInterval.Lower.Denominator))
	upper := new(big.Rat).SetFrac(big.NewInt(c.ObjectiveInterval.Upper.Numerator), big.NewInt(c.ObjectiveInterval.Upper.Denominator))
	if lower.Cmp(big.NewRat(269, 400)) <= 0 {
		t.Fatalf("certificate does not prove >0.6725: %s", lower.RatString())
	}
	if new(big.Rat).Sub(upper, lower).Cmp(big.NewRat(1, 10_000_000)) >= 0 {
		t.Fatalf("certificate is too wide: [%s,%s]", lower.RatString(), upper.RatString())
	}
}
