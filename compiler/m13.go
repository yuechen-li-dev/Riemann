package compiler

import (
	"fmt"
	"math"
	"math/big"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	M13WindowTheoremID       semantic.TheoremID = "montgomery-taylor-window-functional"
	M13MomentTheoremID       semantic.TheoremID = "general-window-prime-moment-asymptotic"
	M13ExtremalTheoremID     semantic.TheoremID = "montgomery-taylor-cauchy-schwarz-extremal"
	M13OptimizationTheoremID semantic.TheoremID = "m13-bandwidth-endpoint-optimization"
)

var anthropicM13Reference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "Claude, More Than Two Thirds of the Zeros of the Riemann Zeta Function Lie on the Critical Line (August 10, 2026), Section 7.1, equations (7.1)-(7.4) and Theorem D",
	URI:      "https://www-cdn.anthropic.com/564f962e60643842f5fcb4a17c9dbc8f608f1c37.pdf",
}

var montgomeryM13Reference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "H. L. Montgomery, Distribution of the Zeros of the Riemann Zeta Function, Proc. ICM Vancouver 1974, Vol. 1 (1975), 379-381",
}

var carneiroM13Reference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "E. Carneiro, V. Chandee, F. Littmann, M. B. Milinovich, Hilbert spaces and the pair correlation of zeros of the Riemann zeta-function, Crelle 725 (2017), Corollary 14",
	URI:      "https://doi.org/10.1515/crelle-2014-0078",
}

var m13DerivedReference = semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Riemann M13 typed window optimization and exact rational Taylor enclosure"}

type M13Optimization struct {
	ClosedObjective        semantic.ScalarExpression  `json:"closed_objective"`
	Derivative             semantic.ScalarExpression  `json:"derivative"`
	Optimizer              semantic.ExactRational     `json:"optimizer"`
	OptimizerLocation      string                     `json:"optimizer_location"`
	BoundaryBehavior       string                     `json:"boundary_behavior"`
	DerivativeSignProof    string                     `json:"derivative_sign_proof"`
	GlobalProof            string                     `json:"global_maximum_proof"`
	TaylorCertificate      semantic.TaylorCertificate `json:"taylor_certificate"`
	CertifiedSimpleLower   semantic.ExactRational     `json:"certified_simple_lower_bound"`
	CertifiedDistinctLower semantic.ExactRational     `json:"certified_distinct_lower_bound"`
	ExactSimpleValue       string                     `json:"exact_simple_value"`
	ExactDistinctValue     string                     `json:"exact_distinct_value"`
	DisplaySimplePercent   string                     `json:"display_simple_percentage"`
	DisplayDistinctPercent string                     `json:"display_distinct_percentage"`
}

type M13AsymptoticCount struct {
	M12Instantiation       string `json:"m12_instantiation"`
	TraceAsymptotic        string `json:"trace_asymptotic"`
	FrobeniusAsymptotic    string `json:"frobenius_asymptotic"`
	EventuallyComposition  string `json:"eventually_composition"`
	SimpleCriticalLiminf   string `json:"simple_critical_liminf"`
	CriticalDistinctLiminf string `json:"critical_distinct_liminf"`
	AllDistinctLiminf      string `json:"all_distinct_liminf"`
	RHStatus               string `json:"rh_status"`
}

type M13Experiment struct {
	Path                   string   `json:"path"`
	CheckCommand           string   `json:"check_command"`
	RunCommand             string   `json:"run_command"`
	ParameterRange         string   `json:"parameter_range"`
	SampleCount            int      `json:"sample_count"`
	CandidateOptimizer     string   `json:"candidate_optimizer"`
	CandidateObjective     string   `json:"candidate_objective"`
	NearbyValues           []string `json:"nearby_values"`
	Precision              string   `json:"precision"`
	PlotPath               string   `json:"plot_path"`
	Execution              string   `json:"execution"`
	EvidenceClassification string   `json:"evidence_classification"`
}

type M13Counterexample struct {
	RejectedCandidate string `json:"rejected_candidate"`
	Failure           string `json:"failure"`
}

type M13Result struct {
	M12                   M12Result                         `json:"-"`
	Family                semantic.TestWindowFamily         `json:"test_window_family"`
	ScaleChange           semantic.ScaleChange              `json:"scale_change"`
	FlatCoefficients      semantic.WindowMomentCoefficients `json:"flat_window_coefficients"`
	OptimizedCoefficients semantic.WindowMomentCoefficients `json:"montgomery_taylor_coefficients"`
	FlatObjective         semantic.ScalarObjective          `json:"flat_objective"`
	Objective             semantic.ScalarObjective          `json:"derived_objective"`
	Optimization          M13Optimization                   `json:"optimization"`
	AsymptoticCount       M13AsymptoticCount                `json:"asymptotic_count"`
	Experiment            M13Experiment                     `json:"oct_experiment"`
	Counterexamples       []M13Counterexample               `json:"counterexamples"`
	Fusion                []string                          `json:"representation_fusion"`
	ImportedMathematics   []string                          `json:"imported_from_literature"`
	DerivedMathematics    []string                          `json:"compiler_derived"`
	Sources               []semantic.Reference              `json:"sources"`
	UtilitySchedulerUsed  bool                              `json:"when_utility_used"`
	OctGoUsed             bool                              `json:"octgo_used"`
}

func add(a, b semantic.ScalarExpression) semantic.ScalarExpression {
	return semantic.ScalarBinary(semantic.ScalarAdd, a, b)
}
func sub(a, b semantic.ScalarExpression) semantic.ScalarExpression {
	return semantic.ScalarBinary(semantic.ScalarSubtract, a, b)
}
func mul(a, b semantic.ScalarExpression) semantic.ScalarExpression {
	return semantic.ScalarBinary(semantic.ScalarMultiply, a, b)
}
func div(a, b semantic.ScalarExpression) semantic.ScalarExpression {
	return semantic.ScalarBinary(semantic.ScalarDivide, a, b)
}

func coefficient(id, meaning string, expr semantic.ScalarExpression, normalization string) semantic.WindowMomentCoefficient {
	return semantic.WindowMomentCoefficient{ID: id, Meaning: meaning, Expression: expr, Parameter: "lambda", Normalization: normalization, Theorem: M13MomentTheoremID, ErrorTerm: "endpoint mollification changes the coefficient by O(1/L); Theorem 5.8 retains o(N) after M10 Eventually conversion", Provenance: anthropicM13Reference}
}

func flatWindowCoefficients() semantic.WindowMomentCoefficients {
	return semantic.WindowMomentCoefficients{
		MassA:       coefficient("a_flat", "integral of v on [-1/2,1/2]", semantic.ScalarRat(1, 1), "v(s)=1"),
		SquareMassB: coefficient("b_flat", "integral of v(s)^2", semantic.ScalarRat(1, 1), "v(s)=1"),
		DistanceJ:   coefficient("J_flat", "double integral |s-s'|v(s)v(s')", semantic.ScalarRat(1, 3), "v(s)=1"),
	}
}

func montgomeryTaylorCoefficients() semantic.WindowMomentCoefficients {
	lambda := semantic.ScalarParam("lambda")
	sqrt2 := semantic.ScalarUnary(semantic.ScalarSqrt, semantic.ScalarRat(2, 1))
	theta := div(lambda, sqrt2)
	k := mul(sqrt2, lambda)
	a := div(semantic.ScalarUnary(semantic.ScalarSin, theta), theta)
	b := add(semantic.ScalarRat(1, 2), div(semantic.ScalarUnary(semantic.ScalarSin, k), mul(semantic.ScalarRat(2, 1), k)))
	// From (T v)''=2v and evenness: J=a^2/2+(2a cos(theta)-2b)/k^2.
	j := add(div(semantic.ScalarPow(a, 2), semantic.ScalarRat(2, 1)), div(sub(mul(mul(semantic.ScalarRat(2, 1), a), semantic.ScalarUnary(semantic.ScalarCos, theta)), mul(semantic.ScalarRat(2, 1), b)), semantic.ScalarPow(k, 2)))
	return semantic.WindowMomentCoefficients{
		MassA:       coefficient("a_mt", "integral of cos(sqrt(2)*lambda*s)", a, "phi(u)^2=v(u/L), a=integral v"),
		SquareMassB: coefficient("b_mt", "integral of cos(sqrt(2)*lambda*s)^2", b, "b=integral v^2"),
		DistanceJ:   coefficient("J_mt", "double integral |s-s'|v(s)v(s')", j, "J=double integral |s-s'|v(s)v(s')"),
	}
}

// deriveM12Objective instantiates, without changing, M12's c=2 count bound:
// 4*trace - 2*zero_count - FrobeniusSquare.  The optimized window changes
// only the normalized Frobenius coefficient.
func deriveM12Objective(id string, parameter semantic.WindowParameter, c semantic.WindowMomentCoefficients) semantic.ScalarObjective {
	lambda := semantic.ScalarParam(parameter.Symbol)
	denominator := add(c.SquareMassB.Expression, mul(semantic.ScalarPow(lambda, 2), c.DistanceJ.Expression))
	frobenius := div(denominator, mul(lambda, semantic.ScalarPow(c.MassA.Expression, 2)))
	expr := sub(sub(mul(semantic.ScalarRat(4, 1), semantic.ScalarRat(1, 1)), semantic.ScalarRat(2, 1)), frobenius)
	return semantic.ScalarObjective{ID: id, Parameter: parameter, ExactExpression: expr, DerivedFrom: []string{"M12 c=2 rank-trace count coefficient 4*tr-2*N-||G||_F^2", "M13 Theorem 5.8 general-window coefficients a,b,J after typed G_tilde->G_hat scaling"}, Provenance: m13DerivedReference}
}

func floorRational(value semantic.ExactRational, places int) (semantic.ExactRational, error) {
	if err := value.Validate(); err != nil {
		return semantic.ExactRational{}, err
	}
	scale := big.NewInt(1)
	for i := 0; i < places; i++ {
		scale.Mul(scale, big.NewInt(10))
	}
	n := new(big.Int).Mul(big.NewInt(value.Numerator), scale)
	q := new(big.Int).Quo(n, big.NewInt(value.Denominator))
	if !q.IsInt64() || !scale.IsInt64() {
		return semantic.ExactRational{}, fmt.Errorf("display rational overflow")
	}
	return semantic.ExactRational{Numerator: q.Int64(), Denominator: scale.Int64()}, nil
}

func CompileM13() (M13Result, error) {
	m12, err := CompileM12()
	if err != nil {
		return M13Result{}, err
	}
	domain := semantic.ScalarDomain{Lower: semantic.ExactRational{Numerator: 0, Denominator: 1}, Upper: semantic.ExactRational{Numerator: 1, Denominator: 1}, LowerIncluded: false, UpperIncluded: true}
	parameter := semantic.WindowParameter{Symbol: "lambda", Domain: domain, Meaning: "normalized support/bandwidth L/log(T/2pi)"}
	lambda := semantic.ScalarParam("lambda")
	sqrt2 := semantic.ScalarUnary(semantic.ScalarSqrt, semantic.ScalarRat(2, 1))
	theta := div(lambda, sqrt2)
	family := semantic.TestWindowFamily{
		ID: "montgomery-taylor-cosine-profile", Parameter: parameter, WindowObjectID: "phi_lambda,L", SquaredProfileID: "v_lambda(s)=phi_lambda(Ls)^2", ProfileArgument: mul(sqrt2, mul(lambda, semantic.ScalarParam("s"))), SupportScale: mul(lambda, semantic.ScalarParam("ell")),
		TransformConvention: "hat(f)(z)=integral_R f(u) exp(i z u) du; K=|hat(v)|^2", Normalization: "L=lambda*ell; v(s)=phi(Ls)^2 on |s|<=1/2; phi is sqrt(cos(sqrt(2)*u/ell)) with a fixed-width C3 endpoint ramp; G_hat=G_tilde/(aL)",
		Admissibility: semantic.WindowAdmissibility{Even: true, NonnegativeProfile: true, CompactSupport: true, MonotoneInAbs: true, TwiceDifferentiable: true, FixedWidthRamp: true, FourierSupportAtMost: semantic.ExactRational{Numerator: 1, Denominator: 1}}, Theorem: M13WindowTheoremID, Provenance: anthropicM13Reference,
	}
	if err := family.Validate(); err != nil {
		return M13Result{}, err
	}
	scale := semantic.ScaleChange{SourceObject: "G_tilde_T", TargetObject: "G_hat_T", Factor: div(semantic.ScalarRat(1, 1), semantic.ScalarParam("aL")), ParameterDependencies: []string{"a", "L", "lambda"}, Provenance: anthropicM13Reference}
	if err := scale.Validate(); err != nil {
		return M13Result{}, err
	}
	flat := flatWindowCoefficients()
	if err := flat.Validate(); err != nil {
		return M13Result{}, err
	}
	mt := montgomeryTaylorCoefficients()
	if err := mt.Validate(); err != nil {
		return M13Result{}, err
	}
	flatObjective := deriveM12Objective("m13.flat-window.m12-objective", parameter, flat)
	objective := deriveM12Objective("m13.montgomery-taylor.m12-objective", parameter, mt)
	if err := flatObjective.Validate(); err != nil {
		return M13Result{}, err
	}
	if err := objective.Validate(); err != nil {
		return M13Result{}, err
	}
	closed := sub(sub(semantic.ScalarRat(2, 1), div(lambda, semantic.ScalarRat(2, 1))), mul(div(semantic.ScalarRat(1, 1), sqrt2), semantic.ScalarUnary(semantic.ScalarCot, theta)))
	derivative := mul(semantic.ScalarRat(1, 2), semantic.ScalarPow(semantic.ScalarUnary(semantic.ScalarCot, theta), 2))
	certificate, err := semantic.MontgomeryTaylorTaylorCertificate()
	if err != nil {
		return M13Result{}, err
	}
	if new(big.Rat).SetFrac(big.NewInt(certificate.ObjectiveInterval.Lower.Numerator), big.NewInt(certificate.ObjectiveInterval.Lower.Denominator)).Cmp(big.NewRat(269, 400)) <= 0 {
		return M13Result{}, fmt.Errorf("Taylor enclosure does not certify 67.25 percent")
	}
	displaySimple, err := floorRational(certificate.ObjectiveInterval.Lower, 4)
	if err != nil {
		return M13Result{}, err
	}
	candidateValue, err := closed.EvalFloat(map[string]float64{"lambda": 1})
	if err != nil {
		return M13Result{}, err
	}
	nearbyValues := make([]string, 0, 3)
	for _, x := range []float64{0.9998, 0.9999, 1} {
		value, evalErr := closed.EvalFloat(map[string]float64{"lambda": x})
		if evalErr != nil {
			return M13Result{}, evalErr
		}
		nearbyValues = append(nearbyValues, fmt.Sprintf("J(%.4f)=%.12f", x, value))
	}
	distinctLowerRat := new(big.Rat).Quo(new(big.Rat).Add(big.NewRat(1, 1), new(big.Rat).SetFrac(big.NewInt(certificate.ObjectiveInterval.Lower.Numerator), big.NewInt(certificate.ObjectiveInterval.Lower.Denominator))), big.NewRat(2, 1))
	distinctFloorInput := semantic.ExactRational{Numerator: distinctLowerRat.Num().Int64(), Denominator: distinctLowerRat.Denom().Int64()}
	displayDistinct, err := floorRational(distinctFloorInput, 5)
	if err != nil {
		return M13Result{}, err
	}
	certifiedSimple := displaySimple
	certifiedDistinct := displayDistinct
	optimization := M13Optimization{
		ClosedObjective: closed, Derivative: derivative, Optimizer: semantic.ExactRational{Numerator: 1, Denominator: 1}, OptimizerLocation: "included upper boundary lambda=1 (no interior stationary point)", BoundaryBehavior: "J(lambda)->-infinity as lambda->0+; J(1) is finite; the derivative exists on (0,1] from the left at 1",
		DerivativeSignProof: "J'(lambda)=cot(lambda/sqrt(2))^2/2>0 because 0<lambda/sqrt(2)<=1/sqrt(2)<pi/2", GlobalProof: "strict increase on (0,1] and inclusion of lambda=1 make lambda=1 the unique global maximizer; the excluded lower boundary tends to -infinity",
		TaylorCertificate: certificate, CertifiedSimpleLower: certifiedSimple, CertifiedDistinctLower: certifiedDistinct, ExactSimpleValue: "3/2-(1/sqrt(2))*cot(1/sqrt(2))", ExactDistinctValue: "5/4-(1/(2*sqrt(2)))*cot(1/sqrt(2))", DisplaySimplePercent: fmt.Sprintf("%.2f%%", 100*float64(displaySimple.Numerator)/float64(displaySimple.Denominator)), DisplayDistinctPercent: fmt.Sprintf("%.3f%%", 100*float64(displayDistinct.Numerator)/float64(displayDistinct.Denominator)),
	}
	result := M13Result{
		M12: m12, Family: family, ScaleChange: scale, FlatCoefficients: flat, OptimizedCoefficients: mt, FlatObjective: flatObjective, Objective: objective, Optimization: optimization,
		AsymptoticCount: M13AsymptoticCount{M12Instantiation: "unchanged M12 c=2 theorem: simple rank >=4*tr(G_hat)-2*N(I')-||G_hat||_F^2", TraceAsymptotic: "tr(G_hat)=N(T,2T)*(1+o(1))", FrobeniusAsymptotic: "||G_hat||_F^2=(1/c*_lambda+o(1))*N, c*_lambda=sqrt(2)*tan(theta)/(1+theta*tan(theta)), theta=lambda/sqrt(2)", EventuallyComposition: "M10 fringe/tail and M11 o/O terms are reused; for every epsilon>0 the finite M12 bound holds beyond T0(lambda,epsilon), then epsilon tends to zero", SimpleCriticalLiminf: "liminf N0_simple(T,2T)/N(T,2T)>=3/2-(1/sqrt(2))*cot(1/sqrt(2))>269/400", CriticalDistinctLiminf: "liminf N0_distinct(T,2T)/N(T,2T)>=3/2-(1/sqrt(2))*cot(1/sqrt(2))>269/400", AllDistinctLiminf: "liminf N_distinct(T,2T)/N(T,2T)>=5/4-(1/(2*sqrt(2)))*cot(1/sqrt(2))>669/800", RHStatus: "unresolved"},
		Experiment:      M13Experiment{Path: "experiments/m13_window_optimization.octest", CheckCommand: "Oct CLI has no standalone check command; oct test performs parse/type-check before execution", RunCommand: "go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m13_window_optimization.octest --execution interpreted (from C:/Users/yuech/source/repos/oct)", ParameterRange: "lambda in [0.0001,1], constrained by the theorem to 0<lambda<=1", SampleCount: 10000, CandidateOptimizer: "lambda=1.0000 (included endpoint)", CandidateObjective: fmt.Sprintf("%.10f", candidateValue), NearbyValues: nearbyValues, Precision: "Oct Float (binary64), step 0.0001; research guidance only", PlotPath: "experiments/m13_window_objective.png (zoom lambda in [0.5,1])", Execution: "7 passed, 0 failed, 0 skipped; interpreted; 4.0 s wall time", EvidenceClassification: "non-certifying numerical scan and plot; exact derivative and rational alternating-series enclosure are theorem evidence"},
		Counterexamples: []M13Counterexample{
			{RejectedCandidate: "Frobenius square scaled linearly by 1/(aL)", Failure: "typed scale fixture gives 5/9 from source value 5 and factor 1/3, not 5/3"},
			{RejectedCandidate: "optimize flat H(lambda) outside lambda<=1", Failure: "support theorem rejects lambda>1; within (0,1] flat H is maximized at 1 and remains 2/3"},
			{RejectedCandidate: "omit b or the denominator b+lambda^2 J", Failure: "changes the derived normalized Frobenius coefficient and fails the flat-window 4/3 regression"},
			{RejectedCandidate: "use phi rather than v=phi^2 in a,b,J", Failure: "violates equation (7.1) normalization and no longer agrees with the typed G_tilde->G_hat scale"},
			{RejectedCandidate: "round numerical optimizer upward into theorem semantics", Failure: "display percentage is floored from the rigorous lower enclosure; the Float scan is research evidence only"},
		},
		Fusion:              []string{"M8: zero-orbit decomposition", "M9: rank and positive-index budgets", "M10: height window, fringe, tail and Eventually transfer", "M11: first/second prime-side moments", "M12: unchanged c=2 rank-trace theorem", "M13: admissible Montgomery-Taylor window and certified scalar optimization"},
		ImportedMathematics: []string{"paper Section 7.1 general-window moment theorem and endpoint mollification", "Montgomery-Taylor Cauchy-Schwarz extremal profile v(s)=cos(sqrt(2)*lambda*s)", "support ceiling lambda<=1 from the available prime-side theorem", "CCLM17 Corollary 14 optimality identification (documented, not reimplemented)"},
		DerivedMathematics:  []string{"typed a,b,J coefficient expressions and typed normalization scale", "compiler construction of the M12 objective 4-2-(b+lambda^2 J)/(lambda a^2)", "flat-window base regression J_flat(1)=2/3", "exact derivative J'(lambda)=cot(lambda/sqrt(2))^2/2 and boundary global-optimum proof", "rational alternating-series enclosure proving the reported 67.25% and 83.625% lower bounds"},
		Sources:             []semantic.Reference{anthropicM13Reference, montgomeryM13Reference, carneiroM13Reference, titchmarshZeroCountReference}, UtilitySchedulerUsed: false, OctGoUsed: false,
	}
	if err := validateM13Result(result); err != nil {
		return M13Result{}, err
	}
	return result, nil
}

func validateM13Result(r M13Result) error {
	if err := r.Family.Validate(); err != nil {
		return err
	}
	if err := r.ScaleChange.Validate(); err != nil {
		return err
	}
	if err := r.Objective.Validate(); err != nil {
		return err
	}
	flat, err := r.FlatObjective.ExactExpression.EvalFloat(map[string]float64{"lambda": 1})
	if err != nil {
		return err
	}
	if math.Abs(flat-2.0/3.0) > 1e-15 {
		return fmt.Errorf("M12 base point changed: %.17g", flat)
	}
	derived, err := r.Objective.ExactExpression.EvalFloat(map[string]float64{"lambda": 1})
	if err != nil {
		return err
	}
	closed, err := r.Optimization.ClosedObjective.EvalFloat(map[string]float64{"lambda": 1})
	if err != nil {
		return err
	}
	if derived-closed > 1e-12 || closed-derived > 1e-12 {
		return fmt.Errorf("derived and closed objectives disagree")
	}
	if !r.Family.Parameter.Domain.ContainsFloat(1) || r.Family.Parameter.Domain.ContainsFloat(0) || r.Family.Parameter.Domain.ContainsFloat(1.0001) {
		return fmt.Errorf("wrong window parameter domain")
	}
	return nil
}
