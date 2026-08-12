package compiler

import (
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	M10WindowSplitTheoremID    semantic.TheoremID = "height-window-zero-side-split"
	M10FarZeroBoundTheoremID   semantic.TheoremID = "anthropic-far-zero-operator-bound"
	M10WeylThresholdTheoremID  semantic.TheoremID = "weyl-thresholded-positive-index-transfer"
	M10WindowCountingTheoremID semantic.TheoremID = "finite-height-window-counting-bridge"
	M10OrbitPullbackTheoremID  semantic.TheoremID = "windowed-orbit-inertia-pullback"
	M10LocalZeroCountTheoremID semantic.TheoremID = "titchmarsh-local-zero-count"
)

var anthropicM10Reference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "Claude, More Than Two Thirds of the Zeta Zeros on the Critical Line (August 10, 2026), §§2.2, 4.1–4.3, especially Propositions 4.1, 4.2 and 4.5",
	URI:      "https://www-cdn.anthropic.com/564f962e60643842f5fcb4a17c9dbc8f608f1c37.pdf",
}

var titchmarshZeroCountReference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "E. C. Titchmarsh, The Theory of the Riemann Zeta-Function, 2nd ed. (1986), Theorem 9.2; local zero-count bound N(t+1)-N(t)=O(log(t+3))",
	URI:      "https://doi.org/10.1093/oso/9780198533696.001.0001",
}

var m10DerivedReference = semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Riemann M10 structural reconstruction of the finite Proposition 4.5 counting bridge, reusing M8 orbit blocks and M9 positive-index accounting"}

type WindowCountBounds struct {
	SimpleCriticalLowerBound int `json:"simple_critical_lower_bound"`
	DistinctAllLowerBound    int `json:"distinct_all_lower_bound"`
}

// FiniteWindowCountLowerBounds is Proposition 4.5's exact natural-number
// conversion after transferring the thresholded count from G to the enlarged
// near window.  It intentionally does not infer independence or multiplicity.
func FiniteWindowCountLowerBounds(thresholdedPositiveCount, targetTotalMultiplicity, fringeMultiplicity int) (WindowCountBounds, error) {
	if thresholdedPositiveCount < 0 || targetTotalMultiplicity < 0 || fringeMultiplicity < 0 {
		return WindowCountBounds{}, fmt.Errorf("window counts must be nonnegative")
	}
	simple := 2*thresholdedPositiveCount - targetTotalMultiplicity - 2*fringeMultiplicity
	distinct := thresholdedPositiveCount - fringeMultiplicity
	if simple < 0 {
		simple = 0
	}
	if distinct < 0 {
		distinct = 0
	}
	return WindowCountBounds{SimpleCriticalLowerBound: simple, DistinctAllLowerBound: distinct}, nil
}

// ThresholdedCriticalRankLowerBound is the M10 adapter around M9's exact
// resource accounting.  The perturbation theorem supplies the observed near
// positive count; M9 subtracts the off-critical positive-index budget.
func ThresholdedCriticalRankLowerBound(thresholdedPositiveCount, offCriticalPositiveBudget int) (int, error) {
	return CriticalRankLowerBound(thresholdedPositiveCount, offCriticalPositiveBudget)
}

type M10Counterexample struct {
	RejectedCandidate string `json:"rejected_candidate"`
	ExactFixture      string `json:"exact_fixture"`
	Reason            string `json:"reason"`
}

type M10Experiment struct {
	Path                   string   `json:"path"`
	Command                string   `json:"command"`
	Setup                  string   `json:"matrix_generation_setup"`
	Threshold              string   `json:"threshold"`
	PerturbationBound      string   `json:"perturbation_bound"`
	Candidates             []string `json:"candidate_inequalities"`
	Trials                 int      `json:"trials"`
	CounterexamplesFound   []string `json:"counterexamples_found"`
	Execution              string   `json:"execution"`
	EvidenceClassification string   `json:"evidence_classification"`
}

type ProofArchitectureMap struct {
	WindowLocalization       string `json:"window_localization"`
	FarZeroControl           string `json:"far_zero_control"`
	ThresholdedSpectralCount string `json:"thresholded_spectral_count"`
	OrbitAccounting          string `json:"critical_off_critical_accounting"`
	AsymptoticNormalization  string `json:"asymptotic_normalization"`
	CompilerFit              string `json:"compiler_fit"`
}

type M10Result struct {
	M9                     M9Result                               `json:"-"`
	Compression            semantic.WindowCompression             `json:"window_compression"`
	CountVocabulary        []string                               `json:"zero_count_vocabulary"`
	Decomposition          semantic.NearFarZeroDecomposition      `json:"near_far_decomposition"`
	FarBound               semantic.FarZeroContributionBound      `json:"far_zero_bound"`
	Perturbation           semantic.ThresholdPerturbationContract `json:"threshold_perturbation_theorem"`
	CountingTheorem        semantic.FiniteWindowCountingTheorem   `json:"finite_window_counting_theorem"`
	ExactSanityObservation semantic.ThresholdedPositiveIndexClaim `json:"exact_threshold_sanity_observation"`
	Architecture           ProofArchitectureMap                   `json:"proof_architecture"`
	Counterexamples        []M10Counterexample                    `json:"counterexamples"`
	Experiment             M10Experiment                          `json:"oct_experiment"`
	UtilitySchedulerUsed   bool                                   `json:"when_utility_used"`
	M7RegressionRole       string                                 `json:"m7_regression_role"`
	Sources                []semantic.Reference                   `json:"sources"`
}

func CompileM10() (M10Result, error) {
	m9, err := CompileM9()
	if err != nil {
		return M10Result{}, err
	}
	window := semantic.HeightWindow{
		ID: "I(T)", Center: semantic.SymbolicHeight("3T/2"), HalfWidth: semantic.SymbolicHeight("T/2"),
		Lower: semantic.SymbolicHeight("T"), Upper: semantic.SymbolicHeight("2T"), Boundary: semantic.LeftOpenRightClosed, OrdinateConvention: semantic.PositiveOrdinateConvention,
	}
	localization := semantic.HeightWindow{
		ID: "I'(T)", Center: semantic.SymbolicHeight("3T/2"), HalfWidth: semantic.SymbolicHeight("T/2+sqrt(T)"),
		Lower: semantic.SymbolicHeight("T-sqrt(T)"), Upper: semantic.SymbolicHeight("2T+sqrt(T)"), Boundary: semantic.LeftOpenRightClosed, OrdinateConvention: semantic.PositiveOrdinateConvention,
	}
	threshold := semantic.SpectralThreshold{
		ID: "theta", Kind: semantic.FarOperatorBoundScale, Expression: "theta>=theta0(T,lambda,w,A0)", Comparison: semantic.StrictlyAboveThreshold,
		Dependencies: map[string]string{"T": "height", "lambda": "bandwidth L=lambda*log(T/2pi)", "w": "C3 taper ramp width", "A0": "local zero-count constant", "normalization": "G_tilde=G/L"}, Provenance: anthropicM10Reference,
	}
	compression := semantic.WindowCompression{
		ID: "gabor-weil-height-compression", Window: window, LocalizationWindow: localization,
		BasisFamily: "f_{T,k}(u)=phi_T(u) exp(-i tau_k u), tau_k=T+2pi*k/L", BasisHeightDependency: "L=lambda*log(T/2pi), d=floor(LT/2pi), 0<=k<d",
		DimensionExpression: "floor(L*T/(2*pi))", MatrixID: "G_tilde_T", NearMatrixID: "A_tilde_T", FarMatrixID: "E_tilde_T",
		ZeroSideIdentity: "G_tilde_T=A_tilde_T+E_tilde_T", ExplicitFormula: "G_kl=integral phi_hat(t-tau_k) phi_hat(t-tau_l) nu_X(t) dt, X=exp(L)",
		Normalization: "G_tilde_T=G_T/L", Threshold: threshold, Provenance: anthropicM10Reference,
	}
	if err := compression.Validate(); err != nil {
		return M10Result{}, err
	}
	decomposition := semantic.NearFarZeroDecomposition{
		MatrixID: "G_tilde_T", NearMatrixID: "A_tilde_T", CriticalNearMatrixID: "P_near", OffCriticalNearMatrix: "Q_near", FarMatrixID: "E_tilde_T",
		Identity: "G_tilde_T=P_near+Q_near+E_tilde_T", MembershipRule: "near iff Re(gamma_rho) is in I'(T)=(T-sqrt(T),2T+sqrt(T)]; reflection partners have the same ordinate",
		Theorems: []semantic.TheoremID{M10WindowSplitTheoremID, M10OrbitPullbackTheoremID, m9.M8.Dual.IdentityTheorem}, Provenance: anthropicM10Reference,
	}
	if err := decomposition.Validate(); err != nil {
		return M10Result{}, err
	}
	farBound := semantic.FarZeroContributionBound{
		MatrixID: "E_tilde_T", Norm: semantic.OperatorNorm, BoundSymbol: "theta0",
		BoundExpression:     "4*A0*C1^2*X^(1/2)*log(4T)/D0^2, D0=sqrt(T), C1=||phi''||_1=2||rho''||_1/w",
		Assumptions:         []string{"T>=T0", "0<lambda<=1", "1<=w<=L/8", "N(t+1)-N(t)<=A0 log(t+3)", "phi is the C3 taper of (2.12)"},
		AsymptoticStatement: "theta0=O(log(T)*T^(lambda/2-1))=o(1)", Uniformity: "fixed lambda in (0,1] and fixed taper profile; explicit dependence on w and A0 retained",
		ExactOrTrusted: true, Theorems: []semantic.TheoremID{M10FarZeroBoundTheoremID, M10LocalZeroCountTheoremID}, Provenance: anthropicM10Reference,
	}
	if err := farBound.Validate(); err != nil {
		return M10Result{}, err
	}
	perturbation := semantic.ThresholdPerturbationContract{
		ID: M10WeylThresholdTheoremID, Statement: "for Hermitian G=A+E, lambda_i(G)<=lambda_i(A)+||E||_op; hence n_plus^theta(G)<=n_plus(A)",
		ThresholdRule: "theta>=||E||_op", Comparison: "strict", Norm: semantic.OperatorNorm, ExactRequired: true, Provenance: anthropicM10Reference,
	}
	if err := perturbation.Validate(); err != nil {
		return M10Result{}, err
	}
	counting := semantic.FiniteWindowCountingTheorem{
		Name:                  "height-window compression to simple-critical-zero and distinct-zero counts",
		Assumptions:           []string{"I=(T,2T] and I'=(T-sqrt(T),2T+sqrt(T)]", "G_tilde=A_tilde+E_tilde are Hermitian", "||E_tilde||_op<=theta0<=theta", "n_plus^theta(G_tilde)>=L_theta is exact/certified", "near orbit blocks obey the M8/M9 critical/off-critical budgets", "N(t+1)-N(t)<=A0 log(t+3)"},
		ThresholdTransfer:     "n_plus^theta(G_tilde)<=n_plus(A_tilde)",
		M9Accounting:          "n_plus(A_tilde)<=rank(P_near)+n_plus(Q_near)<=s1+s2+p; this is M9 positive-index resource accounting with B_off=s2+p",
		EnlargedWindowBound:   "s1>=2*n_plus^theta(G_tilde)-N(I') and #Z(I')>=n_plus^theta(G_tilde)",
		TargetWindowBound:     "N0(T,2T)>=N0_simple(T,2T)>=2*L_theta-N(T,2T)-2*N(I'\\I)",
		DistinctZeroBound:     "N_distinct(T,2T)>=L_theta-N(I'\\I)",
		CountConversion:       "rank(P_near)<=critical distinct locations<=critical multiplicity; the stronger s1 formula additionally uses N(I')>=s1+2s2+2p; no evaluation-vector independence is assumed",
		RemainingInput:        "a certified/asymptotic lower bound L_theta for n_plus^theta(G_tilde_T); first/second-moment production of L_theta is reserved for M11",
		Theorems:              []semantic.TheoremID{M10WeylThresholdTheoremID, M10OrbitPullbackTheoremID, PositiveIndexSubadditivityTheoremID, PositiveIndexRankBoundTheoremID, M9CriticalRankBoundTheoremID, M10WindowCountingTheoremID},
		KnownLiteratureResult: true, StructurallyReconstructed: true, AsymptoticProportionDerived: false, Provenance: anthropicM10Reference,
	}
	if err := counting.Validate(); err != nil {
		return M10Result{}, err
	}
	exactThreshold := threshold
	exactThreshold.ID = "theta=1"
	exactThreshold.Kind = semantic.AbsoluteThresholdScale
	exactThreshold.Expression = "1"
	exactThreshold.ExactValue = &semantic.ExactRational{Numerator: 1, Denominator: 1}
	exactThreshold.Dependencies = map[string]string{"normalization": "exact diagonal sanity fixture"}
	sanity := semantic.ThresholdedPositiveIndexClaim{MatrixID: "diag(2,1,0)", Dimension: 3, Threshold: exactThreshold, Relation: semantic.EqualBound, Bound: 1, Evidence: semantic.ExactTheoremEvidence, Theorems: []semantic.TheoremID{M10WeylThresholdTheoremID}, Provenance: m10DerivedReference}
	if err := sanity.Validate(); err != nil {
		return M10Result{}, err
	}
	return M10Result{
		M9: m9, Compression: compression,
		CountVocabulary: []string{"N(W): total zeros with multiplicity", "N0(W): critical zeros with multiplicity", "N0Distinct(W): distinct critical locations", "N0Simple(W): simple critical locations", "NOffPairs(W): unordered off-critical reflection pairs", "rank(P_near): evaluation-vector directions (not a zero count)"},
		Decomposition:   decomposition, FarBound: farBound, Perturbation: perturbation, CountingTheorem: counting, ExactSanityObservation: sanity,
		Architecture: ProofArchitectureMap{
			WindowLocalization:       "typed I and I'; the smooth taper localizes the test family in Fourier support while zero ownership remains a hard ordinate window",
			FarZeroControl:           "paper Proposition 4.2 maps to a typed trusted operator-norm bound with an explicit o(1) scale",
			ThresholdedSpectralCount: "strict n_plus^theta is a separate exact-evidence object; its analytic lower bound remains open",
			OrbitAccounting:          "M8 blocks and M9 positive-index subadditivity map directly after adding window ownership",
			AsymptoticNormalization:  "only theta0=o(1) and fringe O(sqrt(T) log T) are imported; no proportion is derived",
			CompilerFit:              "clean except that M9's concrete two-dimensional compression had to be generalized by a separate T-dependent family rather than mutated",
		},
		Counterexamples: []M10Counterexample{
			{RejectedCandidate: "replace the strict count #{lambda>theta} by #{lambda>=theta}", ExactFixture: "A=[0], E=[1], theta=||E||=1 gives G=[1]", Reason: "the non-strict count is one while n_plus(A)=0; Weyl transfer fails at equality"},
			{RejectedCandidate: "permit theta<||E||_op", ExactFixture: "A=[0], E=[1], theta=1/2", Reason: "G has one eigenvalue above theta while A has no positive eigenvalue"},
			{RejectedCandidate: "a small Frobenius average per dimension controls every eigenvalue", ExactFixture: "E=diag(1,0,...,0) has ||E||_F/sqrt(d)->0 but ||E||_op=1", Reason: "one threshold crossing survives; Weyl requires an operator-norm bound (or a different theorem)"},
			{RejectedCandidate: "far means off-critical", ExactFixture: "a critical-line zero with ordinate 3T is far; an off-critical reflection pair at ordinate 3T/2 is near", Reason: "height localization and critical-line classification are independent axes"},
		},
		Experiment:           M10Experiment{Path: "experiments/m10_threshold_window.octest", Command: "go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m10_threshold_window.octest", Setup: "a refined StrictlyPositiveGap Concept checks value-threshold before deterministic diagonal near/far fixtures at, below, and above the threshold plus high-rank/low-rank norm fixtures", Threshold: "strict lambda>theta, represented by successful StrictlyPositiveGap(value-theta) Concept construction", PerturbationBound: "theta>=||E||_op", Candidates: []string{"n_plus^theta(A+E)<=n_plus(A)", "non-strict threshold variant", "theta below the operator bound", "dimension-averaged Frobenius substitute"}, Trials: 8, CounterexamplesFound: []string{"non-strict equality fails", "theta<||E|| fails", "small RMS Frobenius size does not prevent one crossing", "small operator norm may have high rank without violating Weyl"}, Execution: "8 passed, 0 failed, 0 skipped; compiled: 8, interpreted fallback: 0", EvidenceClassification: "bounded Oct counterexample filtering only; never theorem evidence"},
		UtilitySchedulerUsed: false,
		M7RegressionRole:     "the certified (f2,f3) matrix remains unchanged and is compiled transitively through M9; it is not identified with G_tilde_T",
		Sources:              []semantic.Reference{anthropicM10Reference, titchmarshZeroCountReference, hornJohnsonRankReference},
	}, nil
}
