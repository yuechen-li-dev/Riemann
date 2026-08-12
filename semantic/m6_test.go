package semantic

import "testing"

func valueMeta(kind string, proofKind ProofObjectKind, proof string) *EvaluationMetadata {
	return &EvaluationMetadata{SemanticTargetID: "test", Backend: "test", PrecisionBits: 256, TransformConvention: LagariasMellinConvention, Error: ErrorSemantics{Kind: kind, ProofObjectKind: proofKind, ProofObject: proof}}
}

func approximateFixture() EntryValue {
	return EntryValue{Kind: ApproximateValue, DefinitionExact: true, Approximate: &ComplexValue{Real: 1.25}, Metadata: valueMeta("roundoff", "", "")}
}

func intervalFixture() EntryValue {
	return EntryValue{Kind: CertifiedInterval, DefinitionExact: true, Interval: &ComplexInterval{RealLower: 1, RealUpper: 2, ImagLower: 0, ImagUpper: 0}, Metadata: valueMeta("rigorous_bound", TheoremBackedBound, "interval-theorem")}
}

func exactFixture() EntryValue {
	return EntryValue{Kind: ExactValue, DefinitionExact: true, Exact: &ExactComplex{Real: ExactScalar{Expression: "5/4", Numerator: 5, Denominator: 4}, Imag: ExactScalar{Expression: "0", Denominator: 1}}, Metadata: valueMeta("exact", IndependentExactArgument, "independent-exact-argument")}
}

func TestM6ValueEvidenceLegalUpgrades(t *testing.T) {
	u := NewUnevaluatedEntryValue()
	a := approximateFixture()
	i := intervalFixture()
	e := exactFixture()
	for _, pair := range [][2]EntryValue{{u, a}, {u, i}, {u, e}, {i, e}} {
		if _, err := UpgradeEntryValue(pair[0], pair[1]); err != nil {
			t.Fatalf("legal upgrade rejected: %v", err)
		}
	}
	for _, pair := range [][2]EntryValue{{a, e}, {a, i}, {i, a}, {e, a}, {e, i}} {
		if _, err := UpgradeEntryValue(pair[0], pair[1]); err == nil {
			t.Fatalf("illegal upgrade %s -> %s accepted", pair[0].Kind, pair[1].Kind)
		}
	}
}

func TestM6HighPrecisionApproximationIsNotExactOrCertified(t *testing.T) {
	a := approximateFixture()
	a.Metadata.PrecisionBits = 100000
	if _, err := UpgradeEntryValue(a, exactFixture()); err == nil {
		t.Fatal("precision laundered approximation into exactness")
	}
	if _, err := UpgradeEntryValue(a, intervalFixture()); err == nil {
		t.Fatal("precision laundered approximation into certification")
	}
}

func TestM6CertifiedAndExactValuesRequireProofObjects(t *testing.T) {
	i := intervalFixture()
	i.Metadata.Error.ProofObject = ""
	if err := i.Validate(); err == nil {
		t.Fatal("uncertified interval accepted")
	}
	e := exactFixture()
	e.Metadata.Error.ProofObject = ""
	if err := e.Validate(); err == nil {
		t.Fatal("unsupported exact value accepted")
	}
	i = intervalFixture()
	i.Metadata.Error.ProofObjectKind = IndependentExactArgument
	if err := i.Validate(); err == nil {
		t.Fatal("wrong proof-object kind certified an interval")
	}
	e = exactFixture()
	e.Metadata.Error.ProofObjectKind = TheoremBackedBound
	if err := e.Validate(); err == nil {
		t.Fatal("interval-bound proof kind certified exactness")
	}
}

func TestM6ConservativeMixedEvidenceArithmetic(t *testing.T) {
	e, i, a := exactFixture(), intervalFixture(), approximateFixture()
	if got := WeakestValueKind(e, e); got != ExactValue {
		t.Fatalf("exact terms weakened to %s", got)
	}
	if got := WeakestValueKind(e, i); got != CertifiedInterval {
		t.Fatalf("exact+interval got %s", got)
	}
	if got := WeakestValueKind(e, i, a); got != ApproximateValue {
		t.Fatalf("approximate component was laundered: %s", got)
	}
	if got := WeakestValueKind(e, NewUnevaluatedEntryValue()); got != UnevaluatedExactDefinition {
		t.Fatalf("unevaluated component disappeared: %s", got)
	}
}

func TestM6DefinitionExactnessIsIndependentOfValueKind(t *testing.T) {
	for _, v := range []EntryValue{NewUnevaluatedEntryValue(), approximateFixture(), intervalFixture(), exactFixture()} {
		if !v.DefinitionExact || v.Validate() != nil {
			t.Fatalf("invalid value state: %+v", v)
		}
	}
}
