package compiler

import (
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	RankSubadditivityTheoremID          semantic.TheoremID = "matrix-rank-subadditivity"
	PositiveIndexSubadditivityTheoremID semantic.TheoremID = "hermitian-positive-index-subadditivity"
	NegativeIndexSubadditivityTheoremID semantic.TheoremID = "hermitian-negative-index-subadditivity"
	PositiveIndexRankBoundTheoremID     semantic.TheoremID = "hermitian-positive-index-at-most-rank"
	M9CriticalRankBoundTheoremID        semantic.TheoremID = "finite-critical-rank-from-positive-index-budget"
)

var hornJohnsonRankReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "Horn and Johnson, Matrix Analysis, 2nd ed. (Cambridge University Press, 2013), §0.4.5, rank inequalities", URI: "https://doi.org/10.1017/CBO9781139020411"}
var anthropicInertiaReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "Claude, More Than Two Thirds of the Zeta Zeros on the Critical Line (2026), §3, Lemma 3.1 (Sylvester-law pull-back and positive-index subadditivity)", URI: "https://www-cdn.anthropic.com/564f962e60643842f5fcb4a17c9dbc8f608f1c37.pdf"}
var m9DerivedReference = semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Riemann M9 exact deduction by combining M8 orbit budgets with the M7/explicit-formula spectral view of the same G"}

type M9Counterexample struct {
	RejectedCandidate string `json:"rejected_candidate"`
	ExactFixture      string `json:"exact_fixture"`
	Reason            string `json:"reason"`
}

type M9Experiment struct {
	Path                   string   `json:"path"`
	Command                string   `json:"command"`
	Setup                  string   `json:"matrix_generation_setup"`
	Candidates             []string `json:"candidate_inequalities"`
	Trials                 int      `json:"trials"`
	CounterexamplesFound   []string `json:"counterexamples_found"`
	EvidenceClassification string   `json:"evidence_classification"`
}

type ClaudeInertiaComparison struct {
	FiniteRoute   string   `json:"finite_route"`
	Reproduced    bool     `json:"one_half_bound_reproduced"`
	MissingInputs []string `json:"missing_inputs"`
	PaperLocation string   `json:"paper_location"`
	RankTraceUsed bool     `json:"rank_trace_used"`
}

type M9Result struct {
	M8                M8Result                                   `json:"-"`
	Compression       semantic.FiniteCompression                 `json:"finite_compression"`
	Invariants        []semantic.SpectralInvariantClaim          `json:"spectral_invariants"`
	Contracts         []semantic.SpectralTheoremContract         `json:"spectral_theorem_contracts"`
	CriticalBudget    semantic.CriticalAggregateBudget           `json:"critical_aggregate_budget"`
	OffCriticalBudget semantic.OffCriticalAggregateBudget        `json:"off_critical_aggregate_budget"`
	Fusion            semantic.RepresentationFusion              `json:"representation_fusion"`
	DerivedTheorem    semantic.FiniteCriticalContributionTheorem `json:"derived_theorem"`
	Counterexamples   []M9Counterexample                         `json:"counterexamples"`
	Experiment        M9Experiment                               `json:"oct_experiment"`
	ClaudeComparison  ClaudeInertiaComparison                    `json:"claude_inertia_comparison"`
	FinitePSDReverse  ProofAttempt                               `json:"finite_psd_reverse_inference"`
}

func CompileM9() (M9Result, error) {
	m8, err := CompileM8()
	if err != nil {
		return M9Result{}, err
	}
	return compileM9FromM8(m8)
}

func compileM9FromM8(m8 M8Result) (M9Result, error) {
	dimension := len(m8.ZeroSide.Basis.Members)
	compression := semantic.FiniteCompression{MatrixID: m8.Dual.SemanticMatrixID, BasisID: m8.ZeroSide.Basis.ID, Dimension: dimension, FunctionSpaceRestriction: true}
	if err := compression.Validate(); err != nil {
		return M9Result{}, err
	}

	// The strict M7 principal-minor certificate makes the 2x2 matrix positive
	// definite, so its positive index is exactly its structural dimension. No
	// floating-point eigenvalue is used.
	gPositive := semantic.SpectralInvariantClaim{MatrixID: compression.MatrixID, Dimension: dimension, Invariant: semantic.PositiveIndexInvariant, Relation: semantic.EqualBound, Bound: dimension, Evidence: semantic.CertifiedMinorEvidence, Theorems: []semantic.TheoremID{M7TwoByTwoPSDTheoremID}, Provenance: m8.Dual.Provenance}
	if err := gPositive.Validate(); err != nil {
		return M9Result{}, err
	}

	contracts := m9SpectralContracts()
	for _, contract := range contracts {
		if err := contract.Validate(); err != nil {
			return M9Result{}, err
		}
	}
	critical := semantic.CriticalAggregateBudget{
		MatrixID: "P", PositiveSemidefinite: true, LocalRankUpperBound: 1,
		RankUpperBound:           "rank(P) <= C_nz <= C_mult",
		NonzeroVectorCountSymbol: "C_nz", MultiplicityCountSymbol: "C_mult",
		MultiplicityEffect: "multiplicity scales u*u* by a positive coefficient; it does not create an independent direction",
		Theorems:           []semantic.TheoremID{OuterProductPSDTheoremID, RankSubadditivityTheoremID}, Provenance: m8.CriticalTemplate.Orbit.Provenance,
	}
	if err := critical.Validate(); err != nil {
		return M9Result{}, err
	}
	off := semantic.OffCriticalAggregateBudget{
		MatrixID: "Q", Hermitian: true, LocalRankUpperBound: 2, LocalPositiveIndexUpperBound: 1, LocalNegativeIndexUpperBound: 1,
		AggregateRankUpperBound:     "rank(Q) <= sum_j rank(Q_j) <= 2*B_pairs",
		AggregatePositiveIndexBound: "n_plus(Q) <= sum_j n_plus(Q_j) <= B_pairs",
		AggregateNegativeIndexBound: "n_minus(Q) <= sum_j n_minus(Q_j) <= B_pairs",
		BudgetSymbol:                "B_off", AdditivityDisclaimed: true,
		Theorems: []semantic.TheoremID{PairRankSignatureTheoremID, RankSubadditivityTheoremID, PositiveIndexSubadditivityTheoremID, NegativeIndexSubadditivityTheoremID}, Provenance: m8.OffCriticalTemplate.Orbit.Provenance,
	}
	if err := off.Validate(); err != nil {
		return M9Result{}, err
	}
	fusion := semantic.RepresentationFusion{
		SemanticMatrixID:     compression.MatrixID,
		ZeroSideFacts:        []string{"G=P+Q", "P is a sum of PSD rank<=1 critical blocks", "n_plus(Q)<=B_off and n_minus(Q)<=B_off"},
		ExplicitFormulaFacts: []string{"n_plus(G)=2 from the certified strict 2x2 principal-minor result"},
		IdentityTheorem:      m8.Dual.IdentityTheorem,
		Theorems:             []semantic.TheoremID{m8.Dual.IdentityTheorem, M7TwoByTwoPSDTheoremID, PositiveIndexSubadditivityTheoremID, PositiveIndexRankBoundTheoremID},
		Provenance:           m8.Dual.Provenance,
	}
	if err := fusion.Validate(); err != nil {
		return M9Result{}, err
	}
	derived := semantic.FiniteCriticalContributionTheorem{
		Name:                     "finite critical-rank lower bound from aggregate positive-index budget",
		Assumptions:              []string{"G=P+Q are finite Hermitian matrices of the same dimension", "n_plus(G)=g is exact or certified", "n_plus(Q)<=B_off", "P is the critical PSD aggregate", "rank(P)<=C_nz<=C_mult"},
		PositiveIndexInequality:  "n_plus(G) <= n_plus(P)+n_plus(Q) <= rank(P)+B_off",
		CriticalRankLowerBound:   "rank(P) >= max(0, n_plus(G)-B_off)",
		CriticalCountConsequence: "C_mult >= C_nz >= rank(P) >= max(0, n_plus(G)-B_off)",
		M7SanityInstance:         "rank(P) >= max(0, 2-B_off); without a finite B_off input this excludes no off-critical block",
		Theorems:                 []semantic.TheoremID{PositiveIndexSubadditivityTheoremID, PositiveIndexRankBoundTheoremID, RankSubadditivityTheoremID, M9CriticalRankBoundTheoremID},
		NewlyDerived:             true, AsymptoticConsequenceDerived: false, Provenance: m9DerivedReference,
	}
	if err := derived.Validate(); err != nil {
		return M9Result{}, err
	}

	return M9Result{
		M8: m8, Compression: compression, Invariants: []semantic.SpectralInvariantClaim{gPositive}, Contracts: contracts,
		CriticalBudget: critical, OffCriticalBudget: off, Fusion: fusion, DerivedTheorem: derived,
		Counterexamples: []M9Counterexample{
			{RejectedCandidate: "aggregate inertia equals the sum of local inertias", ExactFixture: "diag(1,-1)+diag(-1,1)=0", Reason: "each summand has inertia (1,1), while the aggregate has (0,0); only subadditive bounds survive"},
			{RejectedCandidate: "finite positive definiteness of G forces Q=0", ExactFixture: "G=I_2=P+Q with Q=diag(1,-1), P=diag(0,2)", Reason: "P is PSD and Q is an allowed (1,1) block, yet G is positive definite"},
			{RejectedCandidate: "multiplicity m supplies m independent critical rank directions", ExactFixture: "m*u*u* has rank one for every m>0 and u!=0", Reason: "multiplicity changes weight, not vector direction"},
		},
		Experiment:       M9Experiment{Path: "experiments/m9_spectral_budget.octest", Command: "go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m9_spectral_budget.octest --execution interpreted", Setup: "deterministic exact 2x2 critical outer products and reflection-pair blocks, including repeated, cancelling, zero, and dependent directions", Candidates: []string{"n_plus(G)<=rank(P)+B_off", "n_minus(Q)<=B_off", "rank(P)<=number of nonzero critical vectors", "local inertia is additive (rejected)"}, Trials: 8, CounterexamplesFound: []string{"local inertia additivity fails by exact cancellation", "finite PSD does not imply Q=0", "multiplicity does not add rank"}, EvidenceClassification: "bounded Oct experiment; counterexample search only, never theorem evidence"},
		ClaudeComparison: ClaudeInertiaComparison{FiniteRoute: "paper Proposition 4.5: s_1 >= 2*n_plus^theta(G)-N_window; Cauchy-Schwarz first/second moments give the earlier 1/2 stage", Reproduced: false, MissingInputs: []string{"a height-window total-zero accounting object N_window tied to this compression", "an exact/imported asymptotic first-and-second-moment theorem giving n_plus^theta(G) >= (3/4-o(1))*N_window", "far-zero operator-norm control needed for the threshold theta"}, PaperLocation: "Claude (2026), Proposition 4.5 and Appendix C.5", RankTraceUsed: false},
		FinitePSDReverse: m8.FinitePSDReverseProof,
	}, nil
}

func m9SpectralContracts() []semantic.SpectralTheoremContract {
	return []semantic.SpectralTheoremContract{
		{ID: RankSubadditivityTheoremID, Premises: []semantic.SpectralFactKind{semantic.HermitianSumFact, semantic.LocalRankBudgetFact}, Conclusion: semantic.AggregateBudgetFact, HermitianRequired: true, SameDimension: true, Statement: "rank(sum_j A_j) <= sum_j rank(A_j)", Trust: semantic.TrustedSpectralTheorem, Provenance: hornJohnsonRankReference},
		{ID: PositiveIndexSubadditivityTheoremID, Premises: []semantic.SpectralFactKind{semantic.HermitianSumFact, semantic.LocalPositiveIndexBudgetFact}, Conclusion: semantic.AggregateBudgetFact, HermitianRequired: true, SameDimension: true, Statement: "n_plus(A+B) <= n_plus(A)+n_plus(B)", Trust: semantic.TrustedSpectralTheorem, Provenance: anthropicInertiaReference},
		{ID: NegativeIndexSubadditivityTheoremID, Premises: []semantic.SpectralFactKind{semantic.HermitianSumFact, semantic.LocalNegativeIndexBudgetFact}, Conclusion: semantic.AggregateBudgetFact, HermitianRequired: true, SameDimension: true, Statement: "n_minus(A+B) <= n_minus(A)+n_minus(B), obtained by positive-index subadditivity for -A and -B", Trust: semantic.TrustedSpectralTheorem, Provenance: anthropicInertiaReference},
		{ID: PositiveIndexRankBoundTheoremID, Premises: []semantic.SpectralFactKind{semantic.PositiveIndexClaimFact}, Conclusion: semantic.AggregateBudgetFact, HermitianRequired: true, SameDimension: false, Statement: "n_plus(A) <= rank(A)", Trust: semantic.TrustedSpectralTheorem, Provenance: anthropicInertiaReference},
		{ID: M9CriticalRankBoundTheoremID, Premises: []semantic.SpectralFactKind{semantic.HermitianSumFact, semantic.PositiveIndexClaimFact, semantic.AggregateBudgetFact}, Conclusion: semantic.AggregateBudgetFact, HermitianRequired: true, SameDimension: true, Statement: "G=P+Q and n_plus(Q)<=B imply rank(P)>=max(0,n_plus(G)-B)", Trust: semantic.CompilerDerivedSpectralRule, Provenance: m9DerivedReference},
	}
}

// CriticalRankLowerBound is the exact natural-number rearrangement used by
// the M9 theorem contract. It intentionally consumes already-certified integer
// invariants rather than eigenvalues.
func CriticalRankLowerBound(positiveIndexG, offCriticalPositiveBudget int) (int, error) {
	if positiveIndexG < 0 || offCriticalPositiveBudget < 0 {
		return 0, fmt.Errorf("spectral budgets must be nonnegative")
	}
	if positiveIndexG <= offCriticalPositiveBudget {
		return 0, nil
	}
	return positiveIndexG - offCriticalPositiveBudget, nil
}
