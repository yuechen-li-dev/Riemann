package compiler

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

type M23Experiment struct {
	Path              string   `json:"path"`
	Interpreted       string   `json:"interpreted"`
	Compiled          string   `json:"compiled"`
	FallbackCount     int      `json:"fallback_count"`
	Artifacts         []string `json:"artifacts"`
	ArtifactSHA256    []string `json:"artifact_sha256"`
	ExecutionIdentity string   `json:"execution_identity"`
	EvidenceBoundary  string   `json:"evidence_boundary"`
}

type M23Result struct {
	M22                      M22Result                        `json:"m22"`
	Outcome                  string                           `json:"outcome"`
	System                   semantic.TwoAtomOptimizerSystem  `json:"optimizer_system"`
	Optimizer                semantic.CertifiedOptimizerBox   `json:"optimizer_box"`
	LocalOptimality          semantic.LocalOptimalityEvidence `json:"local_optimality"`
	ExplicitWitness          semantic.ExplicitTwoAtomWitness  `json:"explicit_witness"`
	RootUniqueness           string                           `json:"root_uniqueness"`
	ContactBookkeeping       string                           `json:"contact_bookkeeping"`
	NestedEnvelope           string                           `json:"nested_envelope"`
	BranchCompetition        string                           `json:"nearby_branches"`
	SecondaryMinima          string                           `json:"secondary_minima"`
	CertifiedBounds          string                           `json:"certified_bounds"`
	OptimizerWitnessGap      string                           `json:"optimizer_witness_gap"`
	M22Gap                   string                           `json:"m22_gap"`
	PrimalComparison         string                           `json:"primal_comparison"`
	AnthropicComparison      string                           `json:"anthropic_comparison"`
	WholeLinePartition       string                           `json:"whole_line_partition"`
	TailProof                string                           `json:"tail_proof"`
	ReplayPerformance        string                           `json:"replay_performance"`
	Equioscillation          string                           `json:"equioscillation_interpretation"`
	Experiment               M23Experiment                    `json:"oct_experiment"`
	LiteratureAssessment     string                           `json:"literature_assessment"`
	ImportedMathematics      []string                         `json:"imported_mathematics"`
	DerivedMathematics       []string                         `json:"newly_derived_mathematics"`
	ArchitecturalAwkwardness string                           `json:"architectural_awkwardness"`
	CompilerTheory           string                           `json:"compiler_theory"`
	RHStatus                 string                           `json:"rh_status"`
	NextMilestone            string                           `json:"next_milestone"`
}

func CompileM23() (M23Result, error) {
	m22, err := CompileM22()
	if err != nil {
		return M23Result{}, err
	}
	return compileM23FromM22(m22, false)
}

func compileM23FromM22(m22 M22Result, replay bool) (M23Result, error) {
	s, b, l, err := semantic.CertifyVariableRadiusOptimizer()
	if err != nil {
		return M23Result{}, err
	}
	w := semantic.M23ExplicitWitness()
	if err := semantic.VerifyM23ExplicitWitness(w, replay); err != nil {
		return M23Result{}, err
	}
	cnum := 1 + math.Pi/4 - 2/math.Pi + 69*math.Pi*math.Sqrt2/40000
	m22c := 1 + math.Pi/4 - 2/math.Pi + 17*math.Pi*math.Sqrt2/10000
	r := M23Result{
		M22: m22, Outcome: "Success A — CertifiedLocalTwoAtomOptimum plus an improved CertifiedGlobalDualWitness",
		System: s, Optimizer: b, LocalOptimality: l, ExplicitWitness: w,
		RootUniqueness:           "the rational Krawczyk image is strictly inside the final 3-box; the interval Jacobian determinant excludes zero, proving exactly one joint Q-root",
		ContactBookkeeping:       "fixed contacts ±pi/2: active double zeros with positive curvature; movable contacts ±t1: active double zeros with positive curvature; no other active contact in the certified local branch",
		NestedEnvelope:           "for H_r(t)<0, m_*(r)=inf_t F0(t)/(-H_r(t)); positive F_tt makes the isolated t1 the strict interior envelope minimizer, and H*F_tt!=0 continues it smoothly",
		BranchCompetition:        "bounded numerical continuation on r in [3.50,3.52] finds only the isolated movable branch near t=1.79531; this scan is research evidence, while local branch uniqueness is certified",
		SecondaryMinima:          "deterministic scan locates the next minima near t=9.44 and 15.54 with values about 0.571 and 0.597; neither is near-active (research evidence)",
		CertifiedBounds:          "c_* >= 1+pi/4-2/pi+69*pi*sqrt(2)/40000; J_* <= 1-pi/4+2/pi-69*pi*sqrt(2)/40000",
		OptimizerWitnessGap:      fmt.Sprintf("local optimizer c≈1.156643235727947 minus explicit witness c≈%.15f is about %.9g", cnum, 1.156643235727947-cnum),
		M22Gap:                   fmt.Sprintf("local optimizer minus M22 r=7/2,m=17/2500 witness is about %.9g; explicit M23 witness improves M22 by about %.9g", 1.156643235727947-m22c, cnum-m22c),
		PrimalComparison:         "the CGdL-compatible guard c<1651/1250 is enforced; the local family optimum does not alter the primal side or prove full-class optimality",
		AnthropicComparison:      "the external 0.68185 quantity remains contextual and is not identified with this EP3.1 dual c or J=2-c",
		WholeLinePartition:       "Region A: equality-aware fixed contact [1.54,1.61]; Region B: adaptive strict cells through the movable danger [1.7,1.9]; Region C: remaining exact compact cells to 4; Region D: analytic tail",
		TailProof:                "box-wide and witness-specific F(t)>=3/4-(c-1)/2-2w1-2m>0 for |t|>=4; the optimizer is not tail-limited",
		ReplayPerformance:        "M22 reported about 146 s; the stronger M23 witness replay took 310.9 s because [1.7,1.9] was refined to step 1/200000; exact coverage and degree-100 semantics are unchanged",
		Equioscillation:          "two free family parameters plus the movable frequency satisfy F=F_t=0 and one structured envelope-stationarity relation; this exhausts the local degrees of freedom, an analogy rather than a general equioscillation theorem",
		Experiment:               M23Experiment{Path: "experiments/m23_variable_radius_optimizer.octest", Interpreted: "5/5 passed in 666 ms", Compiled: "5/5 passed in 612 ms", FallbackCount: 0, Artifacts: []string{"artifacts/m23/m23_radius_response.png", "artifacts/m23/m23_contact_branch.png", "artifacts/m23/m23_optimizer_spectrum.png"}, ArtifactSHA256: []string{"0d4abaa3b1488ce9511066e5e5aa2ccd58964819e51b352a9a0189f1db53c5c7", "e4521aef10b9ca96323cdc8373e6c9a0a2028dfe53692f853ca6ff4208bf2935", "34c9d21307c32f05ba431139bf30c91d6ea782d0b9a7b50ea931fe5f897926bd"}, ExecutionIdentity: "Oct dev gooct-cli; CPU Float; compiled Facts with zero fallback; artifacts build-time-interpreted and unchanged on deterministic replay", EvidenceBoundary: "independent binary64 reconstruction, Jacobian check, scans and plots only; no Oct value enters theorem state"},
		LiteratureAssessment:     "post-certification screen of arXiv:1810.08843, 2012.07781, and 2502.05106 found SDP, Fourier-extremal, and pair-correlation frameworks but no stated match for the two-atom radius 3.50785, contact 1.79531, or this active system; no novelty or priority claim without expert review",
		ImportedMathematics:      []string{"M18 exact fixed double contact and whole-line base witness", "M22 contact-preserving two-atom surface, Pell sqrt(2) enclosure, rational Taylor cells and analytic tail", "interval Newton/Krawczyk and implicit-function theorems"},
		DerivedMathematics:       []string{"two independent derivations of Q3=A'H-AH_r", "purpose-specific three-variable rational interval jet and Krawczyk inclusion", "certified local branch and negative radius-envelope second derivative", "improved r=7/2,m=69/10000 whole-line witness with adaptive compact replay"},
		ArchitecturalAwkwardness: "the narrow rational interval jet repeats scalar formulas because the existing PiExpression handles Laurent polynomials in pi but not variable trigonometric boxes; widening that IR would be disproportionate for M23",
		CompilerTheory:           "preserved contacts act like compiled invariants: after structural elimination, the next active double contact and envelope stationarity become a finite proof object, while local optimizer evidence and globally feasible theorem witnesses remain different semantic types",
		RHStatus:                 "unresolved", NextMilestone: "M24: certify or refute whether this local two-atom branch remains the controlling envelope across the full legal radius component, without adding support or changing EP3.1.",
	}
	if err := validateM23Result(r); err != nil {
		return M23Result{}, err
	}
	return r, nil
}

func validateM23Result(r M23Result) error {
	if r.RHStatus != "unresolved" || !r.Optimizer.ExistsUnique || !r.LocalOptimality.CertifiedLocalMaximum || !r.ExplicitWitness.ImprovesM22 || r.NextMilestone == "" {
		return fmt.Errorf("incomplete M23 result")
	}
	return nil
}
func M23JSONReport(r M23Result) ([]byte, error) {
	if err := validateM23Result(r); err != nil {
		return nil, err
	}
	return json.MarshalIndent(struct {
		Schema string    `json:"schema"`
		Result M23Result `json:"m23"`
	}{"riemann.m23.v1", r}, "", "  ")
}
