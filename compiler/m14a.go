package compiler

import (
	"fmt"
	"math/big"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

var cgdlM14AReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "A. Chirre, F. Goncalves, D. de Laat, Pair correlation estimates for the zeros of the zeta function via semidefinite programming, Adv. Math. 361 (2020), 106926, Sections 3-4 and Theorem 1", URI: "https://arxiv.org/abs/1810.08843"}
var ramosM14AReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "A. Ramos, Fourier Optimization and Pair Correlation Problems, PhD thesis (2026), EP3.1, equations (3.10)-(3.12), Proposition 3.14", URI: "https://iris.sissa.it/retrieve/10bb6625-4425-47a4-ab11-d3a72a4267d6/thesis_last_corrections_antonioramos.pdf"}
var anthropicM14AReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "Claude, More Than Two Thirds of the Zeros of the Riemann Zeta Function Lie on the Critical Line (August 10, 2026), Remark 1.1", URI: "https://www-cdn.anthropic.com/564f962e60643842f5fcb4a17c9dbc8f608f1c37.pdf"}
var m14ADerivedReference = semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Riemann M14A weak-dual positive-definite completion derivation"}

type M14AKnownPrimal struct {
	ID                string                  `json:"id"`
	ClassMembership   string                  `json:"class_membership"`
	MultiplicityUpper *semantic.ExactRational `json:"multiplicity_upper,omitempty"`
	SimpleLower       semantic.ExactRational  `json:"simple_lower"`
	Strict            bool                    `json:"strict"`
	Evidence          string                  `json:"evidence"`
	Provenance        semantic.Reference      `json:"provenance"`
}

type M14ANumericalDual struct {
	Family                 string `json:"family"`
	DerivedLocalLimit      string `json:"derived_local_limit"`
	BestGridCandidate      string `json:"best_grid_candidate"`
	WholeLineCertified     bool   `json:"whole_line_certified"`
	InfiniteClassCertified bool   `json:"infinite_class_certified"`
	Failure                string `json:"failure"`
}

type M14AExperiment struct {
	Path                   string   `json:"path"`
	CheckCommand           string   `json:"check_command"`
	RunCommand             string   `json:"run_command"`
	Parameterization       string   `json:"parameterization"`
	CandidateValues        []string `json:"candidate_values"`
	Convergence            string   `json:"convergence"`
	Counterexamples        []string `json:"counterexamples"`
	Execution              string   `json:"execution"`
	EvidenceClassification string   `json:"evidence_classification"`
}

type M14AResult struct {
	M13                  M13Result                            `json:"m13"`
	Class                semantic.SupportOneExtremalClass     `json:"support_one_class"`
	Objective            semantic.ExtremalObjective           `json:"objective"`
	M13Membership        semantic.ExtremalCandidateMembership `json:"m13_membership"`
	KnownPrimals         []M14AKnownPrimal                    `json:"known_primal_bounds"`
	DualFormulation      string                               `json:"dual_formulation"`
	WeakDuality          string                               `json:"weak_duality"`
	BaselineWitness      semantic.DualCompletionWitness       `json:"baseline_dual_witness"`
	BaselineBound        semantic.CertifiedExtremalBound      `json:"baseline_certified_bound"`
	NumericalDual        M14ANumericalDual                    `json:"numerical_dual_candidate"`
	CertifiedBracket     string                               `json:"certified_bracket"`
	AnthropicComparison  string                               `json:"anthropic_comparison"`
	CeilingCertified     bool                                 `json:"claimed_ceiling_certified"`
	PreciseObstruction   []string                             `json:"precise_obstruction"`
	FailedWitnesses      []string                             `json:"failed_witnesses"`
	Experiment           M14AExperiment                       `json:"oct_experiment"`
	ImportedMathematics  []string                             `json:"imported_mathematics"`
	DerivedMathematics   []string                             `json:"derived_mathematics"`
	Sources              []semantic.Reference                 `json:"sources"`
	UtilitySchedulerUsed bool                                 `json:"when_utility_used"`
	RHStatus             string                               `json:"rh_status"`
}

func CompileM14A() (M14AResult, error) {
	m13, err := CompileM13()
	if err != nil {
		return M14AResult{}, err
	}
	return compileM14AFromM13(m13)
}

func compileM14AFromM13(m13 M13Result) (M14AResult, error) {
	class := semantic.SupportOneExtremalClass{
		ID: "zeta-support-one-cohn-elkies-ep3.1", DataRadius: semantic.ExactRational{Numerator: 1, Denominator: 1}, TailSemantics: semantic.NonpositiveTail,
		TransformConvention: "hat g(alpha)=integral_R g(x) exp(-2*pi*i*x*alpha) dx", Even: true, Continuous: true, FunctionL1: true, TransformL1: true, FunctionNonnegative: true,
		Normalization: "g(0)=1 (homogeneous normalization)", ObjectiveID: "support-one-multiplicity-ratio", Provenance: []semantic.Reference{cgdlM14AReference, ramosM14AReference},
	}
	objective := semantic.ExtremalObjective{ID: class.ObjectiveID, MultiplicityRatio: "c(g)=[hat g(0)+integral_{-1}^{1}|alpha| hat g(alpha)dalpha]/g(0)", SimpleProportion: "J(g)=2-c(g)", Homogeneous: true, DataMeasure: "nu=delta_0+|alpha| 1_{[-1,1]}(alpha)dalpha", PipelineContract: "M8-M12 fixed rank-trace machinery maps c(g) to simple-critical lower bound 2-c(g)", Provenance: []semantic.Reference{cgdlM14AReference, ramosM14AReference, anthropicM14AReference}}
	membership := semantic.ExtremalCandidateMembership{CandidateID: m13.Family.ID, Even: true, Continuous: true, FunctionL1: true, TransformL1: true, FunctionNonnegative: true, TailSemantics: semantic.ExactFourierSupport, TailStartsAt: semantic.ExactRational{Numerator: 1, Denominator: 1}, ValueAtZero: semantic.ExactRational{Numerator: 1, Denominator: 1}, Normalization: class.Normalization, Proof: "M13/CCLM optimizer is nonnegative and bandlimited to [-1,1]; exact support implies the broader nonpositive-tail condition, after homogeneous normalization"}
	if err := membership.ValidateFor(class); err != nil {
		return M14AResult{}, err
	}

	mtLower := m13.Optimization.TaylorCertificate.ObjectiveInterval.Lower
	mtMultUpperRat := new(big.Rat).Sub(big.NewRat(2, 1), new(big.Rat).SetFrac(big.NewInt(mtLower.Numerator), big.NewInt(mtLower.Denominator)))
	mtMultUpper, err := exactCompilerRat(mtMultUpperRat)
	if err != nil {
		return M14AResult{}, err
	}
	baseline := semantic.DualCompletionWitness{ID: "sine-process-baseline-completion", MultiplicityLower: semantic.ExactRational{Numerator: 1, Denominator: 1}, OutsideMeasure: "sigma=1_{|x|>1} dx", CompletionDistribution: "delta_0+|x|1_{[-1,1]}dx-dx+sigma = delta_0-(1-|x|)_+ dx", FourierImage: "1-(sin(pi*xi)/(pi*xi))^2 >= 0 on R", PositivityEvidence: semantic.GlobalAnalyticPD, ExactSupportScope: true, WholeLineControl: true, TailControl: true, Provenance: m14ADerivedReference}
	baselineBound, err := semantic.ApplyDualCompletion(class, baseline)
	if err != nil {
		return M14AResult{}, err
	}

	result := M14AResult{
		M13: m13, Class: class, Objective: objective, M13Membership: membership,
		KnownPrimals: []M14AKnownPrimal{
			{ID: "Montgomery-Taylor/CCLM", ClassMembership: "certified exact-support subclass member", MultiplicityUpper: &mtMultUpper, SimpleLower: m13.Optimization.TaylorCertificate.ObjectiveInterval.Lower, Strict: true, Evidence: "compiler exact rational Taylor enclosure", Provenance: carneiroM13Reference},
			{ID: "CGdL degree-40 feasible certificate", ClassMembership: "authoritative EP3.1/ALP feasible result; explicit normalization equivalence recorded, coefficient artifact not imported", MultiplicityUpper: &semantic.ExactRational{Numerator: 1651, Denominator: 1250}, SimpleLower: semantic.ExactRational{Numerator: 849, Denominator: 1250}, Strict: true, Evidence: "published Arb/interval-verified feasible upper c<1.3208, hence J>0.6792; imported theorem evidence", Provenance: cgdlM14AReference},
		},
		DualFormulation: "maximize c over nonnegative tempered measures sigma supported on R\\(-1,1) such that P=nu-c*Lebesgue+sigma is a positive-definite tempered distribution",
		WeakDuality:     "if P is positive definite, then <P,hat g>>=0; since hat g<=0 on supp(sigma), Phi_nu(g)-c*g(0)>=-<sigma,hat g>>=0",
		BaselineWitness: baseline, BaselineBound: baselineBound,
		NumericalDual:       M14ANumericalDual{Family: "sigma=c*1_{|x|>1}dx+(c-1)(delta_{-1}+delta_1)", DerivedLocalLimit: "Fourier density has quadratic coefficient (9-8c)/12, so this restricted family requires c<=9/8", BestGridCandidate: "c=9/8 (simple ceiling 7/8 in this diagnostic family)", WholeLineCertified: false, InfiniteClassCertified: false, Failure: "frequency-grid nonnegativity is not global positivity, and the restricted measure ansatz cannot bound the full dual class"},
		CertifiedBracket:    "1 <= c_* < 1.3208, equivalently 0.6792 < J_* <= 1; the upper endpoint on J is only the exact baseline dual and is not the claimed 0.68185 ceiling",
		AnthropicComparison: "Remark 1.1's approximately 0.68185 is consistent with, but not implied by, the wide certified bracket; no released extremal law or witness identifies it",
		CeilingCertified:    false,
		PreciseObstruction: []string{
			"Anthropic does not define its configuration-law class beyond bandwidth-one data and configurationwise validity, so equality with the maximal EP3.1 class is an informed reconstruction rather than a cited theorem",
			"a useful ceiling requires a positive-definite completion with c near 1.31815; no concrete outside measure or Fourier-positive distribution is supplied by the audited sources",
			"CGdL's SDP is primal-feasible and proves c_* is at most 1.3208; reversing that inequality requires an infinite-dimensional dual witness",
			"finite basis or frequency-grid duals do not cover omitted directions; the missing promotion theorem is whole-line positive definiteness plus tail/completeness control",
		},
		FailedWitnesses: []string{"binary64/grid positivity: rejected because a missed frequency is a legal counterexample", "finite polynomial/SDP truncation: rejected as a class ceiling without an omitted-basis bound", "one-atom completion: local expansion restricts it to c<=9/8 and it is far too weak even if globally certified", "calling exact Fourier support the broader class: rejected because EP3.1 permits a nonpositive Fourier tail"},
		Experiment: M14AExperiment{
			Path:                   "experiments/m14a_dual_completion.octest; experiments/m14a_dual_completion_compiled.octest; experiments/m14a_plot_artifact.octest",
			CheckCommand:           "Oct has no standalone check command; oct test/artifact parse and type-check before execution",
			RunCommand:             "from the Oct repository: go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m14a_dual_completion.octest --json; go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m14a_dual_completion_compiled.octest --execution compiled --json; go run ./cmd/oct artifact C:/Users/yuech/source/repos/Riemann/experiments/m14a_plot_artifact.octest --output-root C:/Users/yuech/source/repos/Riemann --json",
			Parameterization:       "one-atom positive-definite completion; c scan and adversarial frequency scan, exact support radius fixed at one",
			CandidateValues:        []string{"baseline c=1: analytic positive-definite witness", "restricted numerical candidate c=9/8", "c>9/8: killed by negative local quadratic coefficient"},
			Convergence:            "bounded scans refine toward 9/8 inside the diagnostic family; no inference to the infinite dual",
			Counterexamples:        []string{"support radius >1 rejected", "grid-only positivity rejected", "tail-free finite witness rejected", "c=1.126 violates the local necessary condition"},
			Execution:              "interpreted research suite 6/6 in 4.811 s; independent compiled checker 3/3 in 0.633 s with zero fallback; [Artifact] plot probe 1/1 in 0.395 s but reported zero attributed outputs; compiled PlotLine remains unsupported",
			EvidenceClassification: "Oct numerical research evidence; the artifact lane is a plotting workaround, not theorem evidence; only the analytic c=1 completion enters theorem state",
		},
		ImportedMathematics: []string{"CGdL definition of ALP, functional Z, zeta mapping, and interval-certified feasible c<1.3208", "Ramos EP3.1 class A_1, measure nu=delta_0+|alpha|dalpha, and normalization map to ALP", "CCLM/M13 exact-support Montgomery-Taylor member", "Bochner-Schwartz positive-definite distribution pairing principle (weak-dual direction only)"},
		DerivedMathematics:  []string{"typed distinction between exact Fourier support and a nonpositive tail outside the data interval", "dual positive-definite completion sufficient condition and weak-duality inequality", "analytic baseline completion P=delta_0-(1-|x|)_+ with nonnegative Fourier transform", "one-atom diagnostic local obstruction c<=9/8", "precise identification of global positive-definiteness and tail/completeness as the missing certificate"},
		Sources:             []semantic.Reference{anthropicM14AReference, cgdlM14AReference, ramosM14AReference, carneiroM13Reference}, UtilitySchedulerUsed: false, RHStatus: "unresolved",
	}
	if err := validateM14AResult(result); err != nil {
		return M14AResult{}, err
	}
	return result, nil
}

func exactCompilerRat(v *big.Rat) (semantic.ExactRational, error) {
	if !v.Num().IsInt64() || !v.Denom().IsInt64() {
		return semantic.ExactRational{}, fmt.Errorf("rational does not fit M14A scalar")
	}
	return semantic.ExactRational{Numerator: v.Num().Int64(), Denominator: v.Denom().Int64()}, nil
}

func validateM14AResult(r M14AResult) error {
	if err := r.Class.Validate(); err != nil {
		return err
	}
	if err := r.Objective.Validate(); err != nil {
		return err
	}
	if err := r.M13Membership.ValidateFor(r.Class); err != nil {
		return err
	}
	if r.CeilingCertified {
		return fmt.Errorf("M14A must not emit the unsupported remark-level ceiling")
	}
	if r.BaselineBound.SimpleUpper != (semantic.ExactRational{Numerator: 1, Denominator: 1}) {
		return fmt.Errorf("baseline dual conversion changed")
	}
	if len(r.PreciseObstruction) < 4 || r.RHStatus != "unresolved" {
		return fmt.Errorf("M14A obstruction or RH semantics missing")
	}
	return nil
}
