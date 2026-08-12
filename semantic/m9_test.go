package semantic

import "testing"

var m9TestReference = Reference{Kind: CompilerRecord, Citation: "M9 test fixture"}

func TestM9SpectralInvariantVocabularyAndExactBoundary(t *testing.T) {
	exact := SpectralInvariantClaim{MatrixID: "G", Dimension: 2, Invariant: PositiveIndexInvariant, Relation: EqualBound, Bound: 2, Evidence: CertifiedMinorEvidence, Theorems: []TheoremID{"minor-criterion"}, Provenance: m9TestReference}
	if err := exact.Validate(); err != nil || !exact.ExactTheoremPremise() {
		t.Fatalf("certified invariant rejected: %v", err)
	}
	approximate := exact
	approximate.Evidence = ApproximateEigenEvidence
	approximate.Theorems = nil
	if err := approximate.Validate(); err != nil {
		t.Fatalf("experimental report claim rejected: %v", err)
	}
	if approximate.ExactTheoremPremise() {
		t.Fatal("approximate eigenvalues crossed the exact theorem boundary")
	}
	if (SpectralInvariant("zero_index")).Valid() {
		t.Fatal("M9 introduced an unused spectral invariant")
	}
}

func TestM9MultiplicityAccountingRequiresDirectionAndCountSeparation(t *testing.T) {
	b := CriticalAggregateBudget{MatrixID: "P", PositiveSemidefinite: true, LocalRankUpperBound: 1, RankUpperBound: "rank(P)<=C_nz<=C_mult", NonzeroVectorCountSymbol: "C_nz", MultiplicityCountSymbol: "C_mult", MultiplicityEffect: "weight only", Theorems: []TheoremID{"rank-subadditivity"}, Provenance: m9TestReference}
	if err := b.Validate(); err != nil {
		t.Fatal(err)
	}
	bad := b
	bad.MultiplicityEffect = ""
	if bad.Validate() == nil {
		t.Fatal("critical budget silently discarded multiplicity semantics")
	}
}
