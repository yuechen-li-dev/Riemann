package semantic

import "testing"

func TestQuantifierKindsAreDistinctAndTyped(t *testing.T) {
	seen := map[QuantifierKind]bool{}
	for _, quantifier := range []QuantifierKind{ForAll, Exists, DensityOne} {
		if !quantifier.Valid() || seen[quantifier] {
			t.Fatalf("invalid or duplicate quantifier %q", quantifier)
		}
		seen[quantifier] = true
	}
}

func TestSelectedDomainInclusions(t *testing.T) {
	bounded := ZerosBelowHeight(RiemannZeta, 100)
	all := NontrivialZeros(RiemannZeta)
	if !IsSubset(bounded, all) || IsSubset(all, bounded) {
		t.Fatal("bounded-zero inclusion has the wrong direction")
	}
	if !IsSubset(CriticalLine(), CriticalStrip()) || !IsSubset(CriticalStrip(), ComplexPlane()) {
		t.Fatal("selected geometric inclusions are missing")
	}
	if IsSubset(HalfPlaneReGreaterThanOne(), CriticalStrip()) {
		t.Fatal("right half-plane was incorrectly included in the critical strip")
	}
}

func TestClaimValidationRejectsMalformedStructuralSemantics(t *testing.T) {
	claim := Claim{
		ID: "bad-bound",
		Proposition: QuantifiedStatement{
			Quantifier: ForAll,
			Domain:     ZerosBelowHeight(RiemannZeta, 0),
			Predicate:  Predicate{Kind: RealPartEqualsHalfPredicate, Function: RiemannZeta},
		},
		Exactness: Exact,
	}
	if err := claim.Validate(); err == nil {
		t.Fatal("zero-height bounded domain was accepted")
	}
}
