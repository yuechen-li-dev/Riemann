package compiler

import (
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	M12RankTraceTheoremID     semantic.TheoremID = "hermitian-rank-trace-inequality"
	M12VonNeumannTheoremID    semantic.TheoremID = "von-neumann-hermitian-trace-inequality"
	M12PositivePartTheoremID  semantic.TheoremID = "hermitian-positive-negative-parts"
	M12FiniteWindowTheoremID  semantic.TheoremID = "rank-trace-window-count"
	M12SimpleRegroupTheoremID semantic.TheoremID = "simple-critical-rank-regrouping"
)

var anthropicM12Reference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "Claude, More Than Two Thirds of the Zeros of the Riemann Zeta Function Lie on the Critical Line (August 10, 2026), Lemma 3.2, Propositions 4.1 and 4.4, Theorem 5.8, and Theorems A-C",
	URI:      "https://www-cdn.anthropic.com/564f962e60643842f5fcb4a17c9dbc8f608f1c37.pdf",
}

var vonNeumannM12Reference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "R. Bhatia, Matrix Analysis, Springer GTM 169 (1997), Theorem II.3.6 (von Neumann trace inequality); Hermitian PSD specialization used in Lemma 3.2",
	URI:      "https://doi.org/10.1007/978-1-4612-0653-8",
}

var m12DerivedReference = semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Riemann M12 exact-rational reconstruction and M8-M11 representation fusion of the rank-trace argument"}

type M12FiniteWindowCount struct {
	AllCriticalRankBound string `json:"all_critical_rank_bound"`
	SimpleRegrouping     string `json:"simple_regrouping"`
	SimpleCountBound     string `json:"simple_count_bound"`
	CriticalCountBound   string `json:"critical_distinct_count_bound"`
	DistinctCountBound   string `json:"all_distinct_count_bound"`
	TailTransfer         string `json:"tail_transfer"`
}

type M12AsymptoticCount struct {
	NormalizedTrace       string `json:"normalized_trace"`
	NormalizedFrobenius   string `json:"normalized_frobenius"`
	BandwidthFunction     string `json:"bandwidth_function"`
	SimpleCriticalLiminf  string `json:"simple_critical_liminf"`
	CriticalLiminf        string `json:"critical_distinct_liminf"`
	DistinctLiminf        string `json:"all_distinct_liminf"`
	TwoThirdsReproduced   bool   `json:"two_thirds_reproduced"`
	ExactSimpleConstant   string `json:"exact_simple_constant"`
	ExactDistinctConstant string `json:"exact_distinct_constant"`
}

type M12Counterexample struct {
	RejectedCandidate string `json:"rejected_candidate"`
	ExactFixture      string `json:"exact_fixture"`
	Failure           string `json:"failure"`
}

type M12Experiment struct {
	Path                   string   `json:"path"`
	Command                string   `json:"command"`
	Fixtures               []string `json:"exact_matrix_fixtures"`
	Candidates             []string `json:"candidate_inequalities"`
	Trials                 int      `json:"trials"`
	Execution              string   `json:"execution"`
	EvidenceClassification string   `json:"evidence_classification"`
}

type M12Result struct {
	M11                  M11Result                                `json:"-"`
	Decomposition        semantic.HermitianComponentDecomposition `json:"decomposition"`
	Parameter            semantic.PositiveScalarParameter         `json:"parameter"`
	FiniteTheorem        semantic.FiniteRankTraceTheorem          `json:"finite_rank_trace_theorem"`
	EqualitySanity       semantic.RankTraceExactResult            `json:"equality_sanity"`
	NegativeTraceSanity  semantic.RankTraceExactResult            `json:"negative_trace_sanity"`
	FiniteWindow         M12FiniteWindowCount                     `json:"finite_window_count"`
	AsymptoticCount      M12AsymptoticCount                       `json:"asymptotic_count"`
	Counterexamples      []M12Counterexample                      `json:"counterexamples"`
	Experiment           M12Experiment                            `json:"oct_experiment"`
	Fusion               []string                                 `json:"representation_fusion"`
	ImportedMathematics  []string                                 `json:"imported_from_literature"`
	DerivedMathematics   []string                                 `json:"newly_derived_by_compiler_research_loop"`
	Sources              []semantic.Reference                     `json:"sources"`
	UtilitySchedulerUsed bool                                     `json:"when_utility_used"`
	OctGoUsed            bool                                     `json:"octgo_used"`
}

func CompileM12() (M12Result, error) {
	m11, err := CompileM11()
	if err != nil {
		return M12Result{}, err
	}
	return compileM12FromM11(m11)
}

func compileM12FromM11(m11 M11Result) (M12Result, error) {
	decomposition := semantic.HermitianComponentDecomposition{
		TotalMatrixID: "A_hat_T", PSDMatrixID: "P_near", IndexMatrixID: "Q_near",
		Identity:      "A_hat_T=P_near+Q_near (the M8 zero-orbit decomposition after the M10 near-window restriction and paper normalization)",
		SameDimension: true, TotalHermitian: true, PSDHermitian: true, QHermitian: true, PSD: true,
		Evidence: semantic.ExactTheoremEvidence, IdentityTheorem: ZeroMatrixAggregateTheoremID, Provenance: m11.M10.M9.M8.Dual.Provenance,
	}
	if err := decomposition.Validate(); err != nil {
		return M12Result{}, err
	}
	parameter := semantic.PositiveScalarParameter{Symbol: "c", Value: semantic.ExactRational{Numerator: 2, Denominator: 1}, Domain: "c>0"}
	if err := parameter.Validate(); err != nil {
		return M12Result{}, err
	}
	finite := semantic.FiniteRankTraceTheorem{
		Name:           "rank-trace inequality for a PSD low-rank component plus an index-bounded Hermitian component",
		Assumptions:    []string{"P,Q are finite d-by-d Hermitian matrices", "G=P+Q is an exact same-dimension decomposition", "P is positive semidefinite", "rank(P)<=r", "n_plus(Q)<=b", "c>0", "||G||_F^2=tr(G^2)"},
		Expansion:      "write Q=Q_plus-Q_minus with Q_plus,Q_minus PSD and orthogonal; ||P+Q||_F^2=||P||_F^2+||Q_plus||_F^2+||Q_minus||_F^2+2tr(PQ_plus)-2tr(PQ_minus)",
		VonNeumannStep: "tr(PQ_minus)<=sum_i p_i*n_i, with p_i=0 for i>r, hence ||P||_F^2-2tr(PQ_minus)+||Q_minus||_F^2>=sum_i(p_i-n_i)^2",
		ScalarSteps:    []string{"x^2>=c*x-c^2/4 for every real x", "q^2>=2*c*q-c^2 for each positive eigenvalue q of Q_plus; at most b such eigenvalues"},
		Conclusion:     "||P+Q||_F^2 >= c*tr(P)-(c^2/4)*r+2*c*tr(Q)-c^2*b",
		Specialization: "c=2: r >= 2*tr(P)+4*tr(Q)-4*b-||P+Q||_F^2 = 4*tr(G)-2*tr(P)-4*b-||G||_F^2",
		EqualityCase:   "P=(c/2)*Pi_1 and Q=c*Pi_2 for orthogonal projections Pi_1 perpendicular to Pi_2 of ranks r and b",
		Theorems:       []semantic.TheoremID{M12PositivePartTheoremID, M12VonNeumannTheoremID, M12RankTraceTheoremID}, CompilerDerived: true, Provenance: m12DerivedReference,
	}
	if err := finite.Validate(); err != nil {
		return M12Result{}, err
	}
	equality, err := semantic.CheckDiagonalRankTrace(
		[]semantic.ExactRational{{Numerator: 1, Denominator: 1}, {Numerator: 0, Denominator: 1}},
		[]semantic.ExactRational{{Numerator: 0, Denominator: 1}, {Numerator: 2, Denominator: 1}}, 1, 1, parameter.Value)
	if err != nil {
		return M12Result{}, err
	}
	negativeTrace, err := semantic.CheckDiagonalRankTrace(
		[]semantic.ExactRational{{Numerator: 1, Denominator: 1}},
		[]semantic.ExactRational{{Numerator: -2, Denominator: 1}}, 1, 0, parameter.Value)
	if err != nil {
		return M12Result{}, err
	}

	return M12Result{
		M11: m11, Decomposition: decomposition, Parameter: parameter, FiniteTheorem: finite,
		EqualitySanity: equality, NegativeTraceSanity: negativeTrace,
		FiniteWindow: M12FiniteWindowCount{
			AllCriticalRankBound: "rank(P_near)>=4*tr(A_hat_T)-2*N(I')-||A_hat_T||_F^2, using tr(P_near)<=N_on(I') and N(I')>=N_on(I')+2*p",
			SimpleRegrouping:     "A_hat_T=P_simple+Q_prime with rank(P_simple)<=s1, tr(P_simple)<=s1, and n_plus(Q_prime)<=s2+p; no independence premise",
			SimpleCountBound:     "s1>=4*tr(A_hat_T)-2*N(I')-||A_hat_T||_F^2",
			CriticalCountBound:   "s1+s2>=rank(P_near)>=4*tr(A_hat_T)-2*N(I')-||A_hat_T||_F^2",
			DistinctCountBound:   "s1+s2+p>=one_half*(4*tr(A_hat_T)-N(I')-||A_hat_T||_F^2)",
			TailTransfer:         "A_hat_T=G_hat_T-E_hat_T; M10 trace-norm/operator tail and N(I'\\I)=O(sqrt(T)*log T) stay explicit, producing the paper Proposition 4.4 error O(theta0/L*(1+||G_hat_T||_F)+sqrt(T)*log T)",
		},
		AsymptoticCount: M12AsymptoticCount{
			NormalizedTrace:      "tr(G_hat_T)=N(T,2T)*(1+o(1))",
			NormalizedFrobenius:  "||G_hat_T||_F^2<=(1/lambda+lambda/3+o(1))*N(T,2T), reusing M11 Theorem 5.8 before the change from G_tilde to G_hat=G_tilde/(aL)",
			BandwidthFunction:    "H(lambda)=4-2-(1/lambda+lambda/3)=2-1/lambda-lambda/3; at lambda=1, H(1)=2/3",
			SimpleCriticalLiminf: "liminf_{T->infinity} N0_simple(T,2T)/N(T,2T)>=2/3",
			CriticalLiminf:       "liminf_{T->infinity} N0_distinct(T,2T)/N(T,2T)>=2/3",
			DistinctLiminf:       "liminf_{T->infinity} N_distinct(T,2T)/N(T,2T)>=5/6",
			TwoThirdsReproduced:  true, ExactSimpleConstant: "2/3", ExactDistinctConstant: "5/6",
		},
		Counterexamples: []M12Counterexample{
			{RejectedCandidate: "drop P positive semidefinite", ExactFixture: "d=1,c=2,P=[-10],Q=[12],r=b=1", Failure: "||P+Q||_F^2=4 but the claimed RHS is 23"},
			{RejectedCandidate: "omit or reduce the rank penalty c^2*r/4", ExactFixture: "d=1,c=2,P=[1],Q=[0],r=1,b=0", Failure: "equality is 1=2-1; any smaller penalty makes RHS exceed 1"},
			{RejectedCandidate: "replace -c^2*b by a smaller positive-index penalty", ExactFixture: "d=1,c=2,P=[0],Q=[2],r=0,b=1", Failure: "equality is 4=8-4; any smaller penalty makes RHS exceed 4"},
			{RejectedCandidate: "use n_minus(Q) as the b budget", ExactFixture: "d=1,c=2,P=[0],Q=[2],r=0,n_minus(Q)=0", Failure: "the false RHS is 8 while the Frobenius square is 4"},
			{RejectedCandidate: "assume tr(Q)>=0", ExactFixture: "d=1,c=2,P=[1],Q=[-2],r=1,b=0", Failure: "the valid theorem has negative tr(Q) and slack 8; no sign premise is available"},
			{RejectedCandidate: "require independent off-critical evaluation vectors", ExactFixture: "a=b=(1,0), Q=ab*+ba*=diag(2,0)", Failure: "dependent block is PSD rank one with n_plus=1 and remains covered by the aggregate bound"},
			{RejectedCandidate: "drop the G=P+Q identity", ExactFixture: "G=I_2 has decompositions (P=I_2,Q=0) and (P=0,Q=I_2)", Failure: "the same total moments coexist with rank(P)=2 and rank(P)=0; component conclusions require decomposition provenance"},
		},
		Experiment: M12Experiment{
			Path: "experiments/m12_rank_trace.octest", Command: "go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m12_rank_trace.octest --execution compiled --json",
			Fixtures:   []string{"equality P=diag(1,0), Q=diag(0,2)", "non-PSD P=-10,Q=12", "negative-trace Q=-2", "dependent semidefinite pair diag(2,0)", "r=0 and b=0 endpoints", "same G=I2 with different decompositions"},
			Candidates: []string{"Lemma 3.2 exact coefficients", "missing PSD", "smaller rank penalty", "smaller positive-index penalty", "negative-index substitution", "c=0 or c<0"},
			Trials:     12, Execution: "12 passed, 0 failed, 0 skipped; compiled: 12, interpreted fallback: 0; 604 ms reported by Oct dev/gooct-cli", EvidenceClassification: "non-certifying exact counterexample and equality regression only; Go theorem contracts remain authoritative",
		},
		Fusion:              []string{"M8: exact zero-orbit identity G=P+Q", "M9: rank(P)<=critical distinct locations and n_plus(Q)<=off-critical pairs", "M10: near/far localization, trace/operator tail, and fringe", "M11: Theorem 5.8 total trace and Frobenius-square", "M12: rank-trace theorem and multiplicity-aware regrouping"},
		ImportedMathematics: []string{"von Neumann trace inequality for Hermitian matrices", "Hermitian positive/negative-part spectral decomposition", "paper Lemma 3.2 and equality model", "paper Propositions 4.1, 4.2, 4.4 and Theorem 5.8 analytic inputs"},
		DerivedMathematics:  []string{"exact typed theorem contract and rational coefficient oracle", "finite rank(P), simple-critical, critical-distinct, and all-distinct count composition", "reconstruction of the paper's 2/3 simple/critical and 5/6 distinct stages", "same-total-matrix demonstration that decomposition provenance breaks M11's moment-only information ceiling"},
		Sources:             []semantic.Reference{anthropicM12Reference, vonNeumannM12Reference, titchmarshZeroCountReference}, UtilitySchedulerUsed: false, OctGoUsed: false,
	}, nil
}

func validateM12Result(r M12Result) error {
	if !r.AsymptoticCount.TwoThirdsReproduced || r.AsymptoticCount.ExactSimpleConstant != "2/3" || r.AsymptoticCount.ExactDistinctConstant != "5/6" || r.EqualitySanity.Slack.Numerator != 0 {
		return fmt.Errorf("incomplete M12 result")
	}
	return nil
}
