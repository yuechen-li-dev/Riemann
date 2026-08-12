package compiler

import (
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	M11TraceIdentityTheoremID      semantic.TheoremID = "hermitian-trace-spectral-identity"
	M11FrobeniusIdentityTheoremID  semantic.TheoremID = "hermitian-frobenius-spectral-identity"
	M11MomentCountTheoremID        semantic.TheoremID = "thresholded-cauchy-schwarz-moment-count"
	M11PrimeMomentTheoremID        semantic.TheoremID = "anthropic-prime-side-first-second-moments"
	M11EventuallyBoundTheoremID    semantic.TheoremID = "little-o-to-eventual-finite-bound"
	M11RiemannVonMangoldtTheoremID semantic.TheoremID = "riemann-von-mangoldt-window-asymptotic"
)

var anthropicM11Reference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "Claude, More Than Two Thirds of the Zeta Zeros on the Critical Line (August 10, 2026), Lemma 3.3, Proposition 4.5, Theorem 5.8 (5.11)-(5.13), and Appendix C.5",
	URI:      "https://www-cdn.anthropic.com/564f962e60643842f5fcb4a17c9dbc8f608f1c37.pdf",
}

var m11LinearAlgebraReference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "Horn and Johnson, Matrix Analysis, 2nd ed. (2013), spectral theorem and Cauchy-Schwarz; specialized as Claude (2026), Lemma 3.3",
	URI:      "https://doi.org/10.1017/CBO9781139020411",
}

var m11DerivedReference = semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Riemann M11 exact-rational reconstruction of thresholded Cauchy-Schwarz and composition with the existing M10 Proposition 4.5 bridge"}

type M11ThresholdScaling struct {
	MomentError         string `json:"moment_error_scale"`
	DimensionAsymptotic string `json:"dimension_asymptotic"`
	Threshold           string `json:"threshold"`
	Penalty             string `json:"dimension_threshold_penalty"`
	RelativePenalty     string `json:"relative_penalty"`
	Conclusion          string `json:"conclusion"`
}

type M11AsymptoticCount struct {
	FiniteEpsilonBound     string `json:"finite_epsilon_bound"`
	NormalizedLowerBound   string `json:"normalized_lower_bound"`
	BandwidthFunction      string `json:"bandwidth_function"`
	EndpointIndexBound     string `json:"endpoint_index_bound"`
	M10SimpleComposition   string `json:"m10_simple_composition"`
	M10DistinctComposition string `json:"m10_distinct_composition"`
	Fringe                 string `json:"fringe"`
	SimpleCriticalLiminf   string `json:"simple_critical_liminf"`
	DistinctLiminf         string `json:"distinct_liminf"`
	HalfTypeReproduced     bool   `json:"half_type_reproduced"`
	ExactSimpleConstant    string `json:"exact_simple_constant"`
}

type M11Counterexample struct {
	RejectedCandidate string `json:"rejected_candidate"`
	ExactSpectrum     string `json:"exact_spectrum"`
	Reason            string `json:"reason"`
}

type M11Experiment struct {
	Path                   string   `json:"path"`
	CheckCommand           string   `json:"check_command"`
	RunCommand             string   `json:"run_command"`
	Setup                  string   `json:"synthetic_spectrum_setup"`
	Candidates             []string `json:"candidate_theorems"`
	Trials                 int      `json:"trials"`
	Findings               []string `json:"findings"`
	Execution              string   `json:"execution"`
	CompilerIdentity       string   `json:"compiler_identity"`
	TimingAndLimits        string   `json:"timing_and_limits"`
	EvidenceClassification string   `json:"evidence_classification"`
}

type M11Result struct {
	M10                  M10Result                            `json:"-"`
	Moments              []semantic.SpectralMomentClaim       `json:"moment_inputs"`
	AsymptoticMoments    []semantic.AsymptoticMomentStatement `json:"asymptotic_moments"`
	EventuallyBounds     []semantic.EventuallyBound           `json:"eventually_finite_bounds"`
	FiniteTheorem        semantic.FiniteMomentCountTheorem    `json:"finite_spectral_theorem"`
	ExactSanity          semantic.MomentCountResult           `json:"exact_sanity"`
	ThresholdScaling     M11ThresholdScaling                  `json:"threshold_scaling"`
	AsymptoticCount      M11AsymptoticCount                   `json:"asymptotic_count"`
	M10ReuseSanity       WindowCountBounds                    `json:"m10_reuse_sanity"`
	Counterexamples      []M11Counterexample                  `json:"counterexamples"`
	Experiment           M11Experiment                        `json:"oct_experiment"`
	UtilitySchedulerUsed bool                                 `json:"when_utility_used"`
	Fusion               []string                             `json:"representation_fusion"`
	ImportedMathematics  []string                             `json:"imported_from_literature"`
	DerivedMathematics   []string                             `json:"newly_derived_by_compiler_research_loop"`
	Sources              []semantic.Reference                 `json:"sources"`
}

func CompileM11() (M11Result, error) {
	m10, err := CompileM10()
	if err != nil {
		return M11Result{}, err
	}

	moments := []semantic.SpectralMomentClaim{
		{ID: "m11.trace.G_tilde", MatrixID: "G_tilde_T", Kind: semantic.Trace, Relation: semantic.EqualBound, Expression: "tr(G_tilde_T)=a*L*N(T,2T)+O(L*sqrt(X))=L*N(T,2T)*(1+O(E_T))", Evidence: semantic.AsymptoticMomentEvidence, Assumptions: []string{"0<lambda<=1", "L=lambda*log(T/2pi)", "X=exp(L)", "a=L^(-1)*integral(phi^2)", "1<=w<=L/8", "T>=T0(lambda)"}, Theorems: []semantic.TheoremID{M11PrimeMomentTheoremID}, Provenance: anthropicM11Reference},
		{ID: "m11.frobenius2.G_tilde", MatrixID: "G_tilde_T", Kind: semantic.FrobeniusNormSquared, Relation: semantic.EqualBound, Expression: "tr(G_tilde_T^2)=2*pi*b*L*integral_T^(2T)mu(t)^2dt+(T/pi)*sum_{n<=X}(Lambda(n)^2/n)*g(log n)+O(L*l*log(l)*(l^2+X))=(T*L/(2*pi))*(l1^2+L^2/3)*(1+O(E_T))", Evidence: semantic.AsymptoticMomentEvidence, Assumptions: []string{"0<lambda<=1", "L=lambda*l", "X=exp(L)", "b=L^(-1)*integral(phi^4)", "g=phi^2 convolved with phi^2", "w=1", "T>=T0(lambda)"}, Theorems: []semantic.TheoremID{M11PrimeMomentTheoremID, M11FrobeniusIdentityTheoremID}, Provenance: anthropicM11Reference},
		{ID: "m11.dimension.G_tilde", MatrixID: "G_tilde_T", Kind: semantic.Dimension, Relation: semantic.EqualBound, Expression: "d(T)=floor(L*T/(2*pi))", Evidence: semantic.TrustedMomentEvidence, Assumptions: []string{"L=lambda*log(T/2pi)"}, Theorems: []semantic.TheoremID{M11PrimeMomentTheoremID}, Provenance: anthropicM11Reference},
	}
	for _, m := range moments {
		if err := m.Validate(); err != nil {
			return M11Result{}, err
		}
	}
	asymptotic := []semantic.AsymptoticMomentStatement{
		{Moment: moments[0], MainTerm: "L*N(T,2T)", Remainder: "O(E_T*L*N(T,2T))=o(L*N(T,2T))", RemainderKind: semantic.LittleORemainder, Scale: "L*N(T,2T)", Parameter: "T->infinity", Uniformity: "fixed lambda in (0,1], fixed taper profile; w=1"},
		{Moment: moments[1], MainTerm: "lambda*l^2*N(T,2T)*(1+lambda^2/3)", Remainder: "O(E_T*lambda*l^2*N(T,2T)*(1+lambda^2/3))=o(l^2*N(T,2T))", RemainderKind: semantic.LittleORemainder, Scale: "l^2*N(T,2T)", Parameter: "T->infinity", Uniformity: "fixed lambda in (0,1], fixed taper profile; w=1"},
	}
	for _, a := range asymptotic {
		if err := a.Validate(); err != nil {
			return M11Result{}, err
		}
	}
	eventually := []semantic.EventuallyBound{
		{SourceMomentID: moments[0].ID, Epsilon: "every epsilon in (0,1)", Threshold: "T>=T_trace(epsilon,lambda,rho)", FiniteBound: "tr(G_tilde_T)>=(1-epsilon)*L*N(T,2T)", Relation: semantic.AtLeastBound, Theorem: M11EventuallyBoundTheoremID, Provenance: m11DerivedReference},
		{SourceMomentID: moments[1].ID, Epsilon: "every epsilon>0", Threshold: "T>=T_F(epsilon,lambda,rho)", FiniteBound: "||G_tilde_T||_F^2<=(1+epsilon)*lambda*l^2*N(T,2T)*(1+lambda^2/3)", Relation: semantic.AtMostBound, Theorem: M11EventuallyBoundTheoremID, Provenance: m11DerivedReference},
	}
	for _, b := range eventually {
		if err := b.Validate(); err != nil {
			return M11Result{}, err
		}
	}
	finite := semantic.FiniteMomentCountTheorem{
		Name:              "thresholded Cauchy-Schwarz count from first and second spectral moments",
		Assumptions:       []string{"G is a finite d-by-d Hermitian matrix", "theta>=0", "tr(G)>=A", "tr(G^2)=||G||_F^2<=B", "A-d*theta>0", "B>0", "moment premises are exact or trusted finite inequalities"},
		Partition:         "J={i:lambda_i>theta}, m=#J=n_plus^theta(G); equality lambda_i=theta belongs to the remainder",
		TraceResidual:     "S=sum_{i in J}lambda_i >= tr(G)-theta*#{i notin J} >= tr(G)-d*theta >= A-d*theta>0; negative remainder eigenvalues only help",
		CauchySchwarz:     "S^2<=m*sum_{i in J}lambda_i^2<=m*tr(G^2)<=m*B",
		RealConclusion:    "n_plus^theta(G)>=(A-d*theta)^2/B",
		IntegerConclusion: "n_plus^theta(G)>=ceil((A-d*theta)^2/B)",
		Theorems:          []semantic.TheoremID{M11TraceIdentityTheoremID, M11FrobeniusIdentityTheoremID, M11MomentCountTheoremID}, Provenance: m11LinearAlgebraReference,
	}
	if err := finite.Validate(); err != nil {
		return M11Result{}, err
	}
	exactSanity, err := semantic.ThresholdedCountFromMoments(semantic.ExactRational{Numerator: 6, Denominator: 1}, semantic.ExactRational{Numerator: 14, Denominator: 1}, 3, semantic.ExactRational{Numerator: 1, Denominator: 1})
	if err != nil {
		return M11Result{}, err
	}
	m10Sanity, err := FiniteWindowCountLowerBounds(3, 4, 1)
	if err != nil {
		return M11Result{}, err
	}

	return M11Result{
		M10: m10, Moments: moments, AsymptoticMoments: asymptotic, EventuallyBounds: eventually, FiniteTheorem: finite, ExactSanity: exactSanity,
		ThresholdScaling: M11ThresholdScaling{
			MomentError:         "E_T=w/L+((l^2+X)*log(l))/(T*l)+T^(lambda/2-1), X=(T/(2pi))^lambda; with w=1, E_T=O(log(l)/l)=o(1) (and O(1/l) for fixed lambda<1)",
			DimensionAsymptotic: "d=floor(LT/(2pi))~LT/(2pi)~lambda*N(T,2T)",
			Threshold:           "theta=theta0=O(l*T^(lambda/2-1)) from M10 Proposition 4.2",
			Penalty:             "d*theta0=O(L*T*l*T^(lambda/2-1))",
			RelativePenalty:     "d*theta0/tr(G_tilde_T)=O(T^(lambda/2-1))",
			Conclusion:          "for every fixed 0<lambda<=1 the threshold penalty is o(tr(G_tilde_T)); theta0=o(1) alone would not suffice",
		},
		AsymptoticCount: M11AsymptoticCount{
			FiniteEpsilonBound:     "for every epsilon>0 and sufficiently large T: n_plus^theta0(G_tilde_T)>=ceil((((1-epsilon)LN-d*theta0)^2)/((1+epsilon)lambda*l^2*N*(1+lambda^2/3)))",
			NormalizedLowerBound:   "n_plus^theta0(G_tilde_T)>=(F(lambda)-o(1))*N(T,2T)",
			BandwidthFunction:      "F(lambda)=lambda/(1+lambda^2/3)",
			EndpointIndexBound:     "lambda=1: n_plus^theta0(G_tilde_T)>=(3/4-o(1))*N(T,2T)",
			M10SimpleComposition:   "N0_simple(T,2T)>=2*n_plus^theta0(G_tilde_T)-N(T,2T)-2*N(I'\\I)",
			M10DistinctComposition: "N_distinct(T,2T)>=n_plus^theta0(G_tilde_T)-N(I'\\I)",
			Fringe:                 "N(I'\\I)=O(sqrt(T)*log T)=o(N(T,2T)) via the M10 local zero-count input and Riemann-von Mangoldt N(T,2T)~T*log(T/(2pi))/(2pi)",
			SimpleCriticalLiminf:   "liminf_{T->infinity} N0_simple(T,2T)/N(T,2T)>=2*F(1)-1=1/2",
			DistinctLiminf:         "liminf_{T->infinity} N_distinct(T,2T)/N(T,2T)>=F(1)=3/4",
			HalfTypeReproduced:     true, ExactSimpleConstant: "1/2",
		},
		M10ReuseSanity: m10Sanity,
		Counterexamples: []M11Counterexample{
			{RejectedCandidate: "omit the dimension correction", ExactSpectrum: "G=theta*I_d has tr(G)=d*theta and n_plus^theta(G)=0", Reason: "trace can be entirely carried by eigenvalues equal to the strict threshold"},
			{RejectedCandidate: "allow theta<0 in the same residual proof", ExactSpectrum: "G=diag(-1), theta=-2: J={-1}, S=-1 although tr(G)-d*theta=1", Reason: "the selected eigenvalues need not be positive, so the S^2 Cauchy step cannot use S as positive mass"},
			{RejectedCandidate: "use a lower bound on the second moment", ExactSpectrum: "G=diag(M,1,...,1), theta=0 with fixed trace lower bound and arbitrarily large M", Reason: "Cauchy-Schwarz needs an upper bound on sum lambda_i^2; a lower bound gives no count control"},
			{RejectedCandidate: "positive trace forces many positive directions", ExactSpectrum: "G=diag(M,-1,...,-1), M>d-1", Reason: "one positive eigenvalue can carry the trace while negatives cancel"},
			{RejectedCandidate: "replace strict >theta by >=theta without changing the partition", ExactSpectrum: "G=theta*I_d", Reason: "threshold equality changes ownership and invalidates the strict M10/Weyl composition"},
			{RejectedCandidate: "omit the Frobenius upper bound", ExactSpectrum: "G=diag(A+(d-1)M,-M,...,-M)", Reason: "fixed positive trace A is compatible with one arbitrarily large positive direction and many negative directions"},
			{RejectedCandidate: "replace Frobenius-square by operator-norm-square verbatim", ExactSpectrum: "G=diag(1,1) has trace 2, ||G||_op^2=1, and two positive directions", Reason: "the substituted quotient is 4 and exceeds the dimension; only the weaker derived bound ||G||_F^2<=d||G||_op^2 would be valid"},
		},
		Experiment:           M11Experiment{Path: "experiments/m11_moment_count.octest", CheckCommand: "standalone 'oct check' unavailable in this Oct CLI; 'oct test --execution compiled' compiles every fact before execution", RunCommand: "go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m11_moment_count.octest --execution compiled --json", Setup: "bounded exact-integer diagonal spectra; StrictlyPositiveGap, StrictlyPositiveResidual, and NonnegativeThreshold Concepts enforce theorem preconditions through Require; trace, sum of squares, strict counts, and ceiling cross-products use no floating point", Candidates: []string{"dimension-corrected thresholded Cauchy-Schwarz", "dimension-free residual", "negative-threshold extension", "second-moment lower-bound variant", "operator-norm substitution", "strict/non-strict threshold substitution"}, Trials: 11, Findings: []string{"accepted theorem survives negative eigenvalues, a single large positive eigenvalue, and threshold equality", "dimension-free, theta<0, wrong second-moment direction, operator-norm substitution, and non-strict variants have exact counterexamples", "theta=0 equality spectra attain the real and integer bound"}, Execution: "11 passed, 0 failed, 0 skipped; compiled: 11, interpreted fallback: 0", CompilerIdentity: "Oct dev; executionIdentity=gooct-cli; real CLI at C:/Users/yuech/source/repos/oct", TimingAndLimits: "635 ms reported by Oct; 11 deterministic facts; CPU-only integer arithmetic; no portability or GPU claim", EvidenceClassification: "non-certifying counterexample and sharpness exploration only"},
		UtilitySchedulerUsed: false,
		Fusion:               []string{"analytic/explicit-formula representation -> trusted first and second moments", "spectral representation -> generic thresholded count theorem", "window representation -> M10 far threshold and fringe", "zero-orbit representation -> M10 simple/distinct count conversion"},
		ImportedMathematics:  []string{"paper Lemma 3.3 thresholded Cauchy-Schwarz argument", "paper Theorem 5.8 first and second moments for G_tilde=G/L", "paper Proposition 4.2 threshold and Proposition 4.5 counting bridge (the latter already compiled in M10)", "Riemann-von Mangoldt and local zero-count asymptotics"},
		DerivedMathematics:   []string{"exact-rational reusable finite theorem with structural ceiling and theorem-evidence boundary", "explicit dimension-weighted threshold discharge", "four-representation composition reconstructing the paper's earlier one-half simple-critical stage; standard/reconstructed, not claimed novel"},
		Sources:              []semantic.Reference{anthropicM11Reference, m11LinearAlgebraReference, titchmarshZeroCountReference},
	}, nil
}

func validateM11Result(r M11Result) error {
	if len(r.Moments) != 3 || len(r.AsymptoticMoments) != 2 || !r.AsymptoticCount.HalfTypeReproduced || r.AsymptoticCount.ExactSimpleConstant != "1/2" {
		return fmt.Errorf("incomplete M11 result")
	}
	return nil
}
