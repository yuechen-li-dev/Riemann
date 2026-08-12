package semantic

import (
	"fmt"
	"strings"
)

// SpectralInvariant is deliberately limited to the invariants consumed by M9.
// Zero index and a general eigenvalue representation are not needed.
type SpectralInvariant string

const (
	RankInvariant          SpectralInvariant = "rank"
	PositiveIndexInvariant SpectralInvariant = "positive_index"
	NegativeIndexInvariant SpectralInvariant = "negative_index"
)

func (i SpectralInvariant) Valid() bool {
	return i == RankInvariant || i == PositiveIndexInvariant || i == NegativeIndexInvariant
}

type BoundRelation string

const (
	EqualBound   BoundRelation = "equal"
	AtMostBound  BoundRelation = "at_most"
	AtLeastBound BoundRelation = "at_least"
)

func (r BoundRelation) Valid() bool { return r == EqualBound || r == AtMostBound || r == AtLeastBound }

type SpectralEvidenceKind string

const (
	ExactTheoremEvidence     SpectralEvidenceKind = "exact_theorem"
	CertifiedMinorEvidence   SpectralEvidenceKind = "certified_principal_minors"
	ApproximateEigenEvidence SpectralEvidenceKind = "approximate_eigensolver"
)

// SpectralInvariantClaim keeps exact theorem facts separate from experimental
// eigensolver output. ApproximateEigenEvidence is valid report data, but it is
// never admissible to an exact theorem application.
type SpectralInvariantClaim struct {
	MatrixID   string               `json:"matrix_id"`
	Dimension  int                  `json:"dimension"`
	Invariant  SpectralInvariant    `json:"invariant"`
	Relation   BoundRelation        `json:"relation"`
	Bound      int                  `json:"bound"`
	Evidence   SpectralEvidenceKind `json:"evidence"`
	Theorems   []TheoremID          `json:"theorems"`
	Provenance Reference            `json:"provenance"`
}

func (c SpectralInvariantClaim) Validate() error {
	if strings.TrimSpace(c.MatrixID) == "" || c.Dimension < 1 || !c.Invariant.Valid() || !c.Relation.Valid() || c.Bound < 0 || c.Bound > c.Dimension || c.Provenance.Citation == "" {
		return fmt.Errorf("invalid spectral invariant claim")
	}
	if c.Evidence != ExactTheoremEvidence && c.Evidence != CertifiedMinorEvidence && c.Evidence != ApproximateEigenEvidence {
		return fmt.Errorf("invalid spectral evidence kind")
	}
	if c.Evidence != ApproximateEigenEvidence && len(c.Theorems) == 0 {
		return fmt.Errorf("exact spectral claim lacks theorem provenance")
	}
	return nil
}

func (c SpectralInvariantClaim) ExactTheoremPremise() bool {
	return c.Evidence == ExactTheoremEvidence || c.Evidence == CertifiedMinorEvidence
}

type FiniteCompression struct {
	MatrixID                 string `json:"matrix_id"`
	BasisID                  string `json:"basis_id"`
	Dimension                int    `json:"dimension"`
	FunctionSpaceRestriction bool   `json:"function_space_restriction"`
}

func (c FiniteCompression) Validate() error {
	if c.MatrixID == "" || c.BasisID == "" || c.Dimension < 1 || !c.FunctionSpaceRestriction {
		return fmt.Errorf("invalid finite compression")
	}
	return nil
}

type SpectralFactKind string

const (
	HermitianSumFact             SpectralFactKind = "hermitian_sum"
	LocalRankBudgetFact          SpectralFactKind = "local_rank_budget"
	LocalPositiveIndexBudgetFact SpectralFactKind = "local_positive_index_budget"
	LocalNegativeIndexBudgetFact SpectralFactKind = "local_negative_index_budget"
	PositiveIndexClaimFact       SpectralFactKind = "positive_index_claim"
	AggregateBudgetFact          SpectralFactKind = "aggregate_budget"
)

func (k SpectralFactKind) Valid() bool {
	switch k {
	case HermitianSumFact, LocalRankBudgetFact, LocalPositiveIndexBudgetFact, LocalNegativeIndexBudgetFact, PositiveIndexClaimFact, AggregateBudgetFact:
		return true
	default:
		return false
	}
}

// SpectralTheoremContract is a small typed contract vocabulary rather than a
// symbolic eigensolver. Its premise/conclusion kinds describe exactly the
// finite-dimensional facts that each imported theorem transports.
type SpectralTheoremContract struct {
	ID                TheoremID          `json:"id"`
	Premises          []SpectralFactKind `json:"premises"`
	Conclusion        SpectralFactKind   `json:"conclusion"`
	HermitianRequired bool               `json:"hermitian_required"`
	SameDimension     bool               `json:"same_dimension_required"`
	Statement         string             `json:"statement"`
	Trust             SpectralTrust      `json:"trust"`
	Provenance        Reference          `json:"provenance"`
}

type SpectralTrust string

const (
	TrustedSpectralTheorem      SpectralTrust = "trusted_external_theorem"
	CompilerDerivedSpectralRule SpectralTrust = "compiler_derived_rule"
)

func (c SpectralTheoremContract) Validate() error {
	if c.ID == "" || len(c.Premises) == 0 || !c.Conclusion.Valid() || !c.HermitianRequired || strings.TrimSpace(c.Statement) == "" || (c.Trust != TrustedSpectralTheorem && c.Trust != CompilerDerivedSpectralRule) || c.Provenance.Citation == "" {
		return fmt.Errorf("invalid spectral theorem contract")
	}
	for _, premise := range c.Premises {
		if !premise.Valid() {
			return fmt.Errorf("invalid spectral theorem premise")
		}
	}
	return nil
}

type CriticalAggregateBudget struct {
	MatrixID                 string      `json:"matrix_id"`
	PositiveSemidefinite     bool        `json:"positive_semidefinite"`
	LocalRankUpperBound      int         `json:"local_rank_upper_bound"`
	RankUpperBound           string      `json:"rank_upper_bound"`
	NonzeroVectorCountSymbol string      `json:"nonzero_vector_count_symbol"`
	MultiplicityCountSymbol  string      `json:"multiplicity_count_symbol"`
	MultiplicityEffect       string      `json:"multiplicity_effect"`
	Theorems                 []TheoremID `json:"theorems"`
	Provenance               Reference   `json:"provenance"`
}

func (b CriticalAggregateBudget) Validate() error {
	if b.MatrixID == "" || !b.PositiveSemidefinite || b.LocalRankUpperBound != 1 || b.RankUpperBound == "" || b.NonzeroVectorCountSymbol == "" || b.MultiplicityCountSymbol == "" || b.MultiplicityEffect == "" || len(b.Theorems) == 0 || b.Provenance.Citation == "" {
		return fmt.Errorf("invalid critical aggregate budget")
	}
	return nil
}

type OffCriticalAggregateBudget struct {
	MatrixID                     string      `json:"matrix_id"`
	Hermitian                    bool        `json:"hermitian"`
	LocalRankUpperBound          int         `json:"local_rank_upper_bound"`
	LocalPositiveIndexUpperBound int         `json:"local_positive_index_upper_bound"`
	LocalNegativeIndexUpperBound int         `json:"local_negative_index_upper_bound"`
	AggregateRankUpperBound      string      `json:"aggregate_rank_upper_bound"`
	AggregatePositiveIndexBound  string      `json:"aggregate_positive_index_bound"`
	AggregateNegativeIndexBound  string      `json:"aggregate_negative_index_bound"`
	BudgetSymbol                 string      `json:"positive_index_budget_symbol"`
	AdditivityDisclaimed         bool        `json:"additivity_disclaimed"`
	Theorems                     []TheoremID `json:"theorems"`
	Provenance                   Reference   `json:"provenance"`
}

func (b OffCriticalAggregateBudget) Validate() error {
	if b.MatrixID == "" || !b.Hermitian || b.LocalRankUpperBound != 2 || b.LocalPositiveIndexUpperBound != 1 || b.LocalNegativeIndexUpperBound != 1 || b.AggregateRankUpperBound == "" || b.AggregatePositiveIndexBound == "" || b.AggregateNegativeIndexBound == "" || b.BudgetSymbol == "" || !b.AdditivityDisclaimed || len(b.Theorems) == 0 || b.Provenance.Citation == "" {
		return fmt.Errorf("invalid off-critical aggregate budget")
	}
	return nil
}

type RepresentationFusion struct {
	SemanticMatrixID     string      `json:"semantic_matrix_id"`
	ZeroSideFacts        []string    `json:"zero_side_facts"`
	ExplicitFormulaFacts []string    `json:"explicit_formula_facts"`
	IdentityTheorem      TheoremID   `json:"identity_theorem"`
	Theorems             []TheoremID `json:"theorems"`
	Provenance           Reference   `json:"provenance"`
}

func (f RepresentationFusion) Validate() error {
	if f.SemanticMatrixID == "" || len(f.ZeroSideFacts) == 0 || len(f.ExplicitFormulaFacts) == 0 || f.IdentityTheorem == "" || len(f.Theorems) == 0 || f.Provenance.Citation == "" {
		return fmt.Errorf("invalid representation fusion")
	}
	return nil
}

type FiniteCriticalContributionTheorem struct {
	Name                         string      `json:"name"`
	Assumptions                  []string    `json:"assumptions"`
	PositiveIndexInequality      string      `json:"positive_index_inequality"`
	CriticalRankLowerBound       string      `json:"critical_rank_lower_bound"`
	CriticalCountConsequence     string      `json:"critical_count_consequence"`
	M7SanityInstance             string      `json:"m7_sanity_instance"`
	Theorems                     []TheoremID `json:"theorems"`
	NewlyDerived                 bool        `json:"newly_derived"`
	AsymptoticConsequenceDerived bool        `json:"asymptotic_consequence_derived"`
	Provenance                   Reference   `json:"provenance"`
}

func (t FiniteCriticalContributionTheorem) Validate() error {
	if t.Name == "" || len(t.Assumptions) == 0 || t.PositiveIndexInequality == "" || t.CriticalRankLowerBound == "" || t.CriticalCountConsequence == "" || t.M7SanityInstance == "" || len(t.Theorems) == 0 || !t.NewlyDerived || t.AsymptoticConsequenceDerived || t.Provenance.Citation == "" {
		return fmt.Errorf("invalid finite critical-contribution theorem")
	}
	return nil
}
