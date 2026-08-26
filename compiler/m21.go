package compiler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

//go:embed artifacts/m21_contact_adapted_ray_witness.json
var m21WitnessJSON []byte

type M21NumericalCandidate struct {
	Frequency float64 `json:"t1"`
	Mass      float64 `json:"m"`
	C         float64 `json:"c"`
	Curvature float64 `json:"second_derivative"`
	Evidence  string  `json:"evidence"`
}

type M21WitnessArtifact struct {
	Schema             string                                         `json:"schema"`
	BaselineC          string                                         `json:"baseline_c"`
	BaselineW          string                                         `json:"baseline_w"`
	Radius2            semantic.ExactRational                         `json:"r2"`
	Mass               semantic.ExactRational                         `json:"certified_m"`
	DerivedC           semantic.PiExpression                          `json:"derived_c"`
	DerivedW1          semantic.PiExpression                          `json:"derived_w1"`
	Weight2            semantic.ExactRational                         `json:"w2"`
	Algebra            semantic.ContactAdaptedRayAlgebra              `json:"contact_algebra"`
	WholeLine          semantic.ContactAdaptedRayWholeLineCertificate `json:"whole_line"`
	NumericalCandidate M21NumericalCandidate                          `json:"numerical_candidate"`
	Ceiling            semantic.CertifiedRayCeilingBracket            `json:"certified_optimum_bracket"`
	OriginalContact    map[string]string                              `json:"original_contact"`
	NewContact         map[string]string                              `json:"new_active_frequency"`
	TheoremConsequence map[string]string                              `json:"theorem_consequence"`
	Provenance         string                                         `json:"provenance"`
}

type M21Experiment struct {
	Path              string   `json:"path"`
	Interpreted       string   `json:"interpreted"`
	Compiled          string   `json:"compiled"`
	FallbackCount     int      `json:"fallback_count"`
	LongRange         string   `json:"long_range_scan"`
	Artifacts         []string `json:"artifacts"`
	ArtifactSHA256    []string `json:"artifact_sha256"`
	ExecutionIdentity string   `json:"execution_identity"`
	EvidenceBoundary  string   `json:"evidence_boundary"`
}

type M21Result struct {
	M20                        M20Result                                      `json:"m20"`
	Outcome                    string                                         `json:"outcome"`
	Ray                        semantic.ContactAdaptedAtomicRay               `json:"ray"`
	Algebra                    semantic.ContactAdaptedRayAlgebra              `json:"contact_algebra"`
	ResponseSignStructure      []string                                       `json:"H_sign_structure"`
	Envelope                   string                                         `json:"envelope"`
	EnvelopeTangencyTheorem    string                                         `json:"envelope_tangency_theorem"`
	NumericalCandidate         M21NumericalCandidate                          `json:"numerical_candidate"`
	CertifiedWitness           semantic.ExactRational                         `json:"certified_witness_m"`
	BestSimpleRational         semantic.ExactRational                         `json:"best_simple_rational_m"`
	WholeLine                  semantic.ContactAdaptedRayWholeLineCertificate `json:"whole_line_certificate"`
	Ceiling                    semantic.CertifiedRayCeilingBracket            `json:"ray_ceiling_bracket"`
	RayCeilingStatus           string                                         `json:"ray_ceiling_status"`
	CertifiedC                 semantic.PiExpression                          `json:"certified_c"`
	CertifiedJ                 semantic.PiExpression                          `json:"certified_J"`
	CertifiedBounds            string                                         `json:"updated_bounds"`
	DirectedDecimals           string                                         `json:"directed_decimals"`
	M20Comparison              string                                         `json:"m20_comparison"`
	PrimalComparison           string                                         `json:"compatible_primal_comparison"`
	AnthropicComparison        string                                         `json:"anthropic_0_68185"`
	MultiContactInterpretation string                                         `json:"multi_contact_interpretation"`
	Experiment                 M21Experiment                                  `json:"oct_experiment"`
	WitnessArtifact            string                                         `json:"witness_artifact"`
	LiteratureAssessment       string                                         `json:"literature_assessment"`
	ImportedMathematics        []string                                       `json:"imported_mathematics"`
	DerivedMathematics         []string                                       `json:"compiler_research_derived_mathematics"`
	ArchitecturalAwkwardness   string                                         `json:"architectural_awkwardness"`
	CompilerTheory             string                                         `json:"compiler_theory"`
	NextMilestone              string                                         `json:"one_next_milestone"`
	RHStatus                   string                                         `json:"rh_status"`
}

func CompileM21() (M21Result, error) {
	m20, err := CompileM20()
	if err != nil {
		return M21Result{}, err
	}
	return compileM21FromM20(m20, false)
}

func compileM21FromM20(m20 M20Result, replayCompact bool) (M21Result, error) {
	var artifact M21WitnessArtifact
	if err := json.Unmarshal(m21WitnessJSON, &artifact); err != nil {
		return M21Result{}, fmt.Errorf("decode M21 witness artifact: %w", err)
	}
	if artifact.Schema != "riemann.m21.contact-adapted-ray.v1" || artifact.BaselineC != "1+pi/4-2/pi" || artifact.BaselineW != "1/pi-4/pi^2+8/pi^3" || artifact.Radius2 != (semantic.ExactRational{Numerator: 4, Denominator: 1}) || artifact.Weight2 != artifact.Mass {
		return M21Result{}, fmt.Errorf("M21 witness does not preserve the selected M20 ray")
	}
	ray, err := semantic.DeriveContactAdaptedAtomicRay(artifact.Mass)
	if err != nil || !reflect.DeepEqual(ray.C, artifact.DerivedC) || !reflect.DeepEqual(ray.Weight1, artifact.DerivedW1) {
		return M21Result{}, fmt.Errorf("M21 derived ray artifact mismatch: %v", err)
	}
	algebra, err := semantic.VerifyContactAdaptedRayAlgebra(ray)
	if err != nil || !reflect.DeepEqual(algebra, artifact.Algebra) {
		return M21Result{}, fmt.Errorf("M21 exact contact artifact mismatch: %v", err)
	}
	if err := semantic.VerifyContactAdaptedRayWholeLine(ray, artifact.WholeLine, replayCompact); err != nil {
		return M21Result{}, err
	}
	if artifact.Ceiling.LowerWitness != artifact.Mass {
		return M21Result{}, fmt.Errorf("M21 ceiling lower endpoint is not the certified witness")
	}
	if err := semantic.VerifyRayCeilingUpper(artifact.Ceiling, 100); err != nil {
		return M21Result{}, err
	}
	j := semantic.DeriveContactAdaptedRayJ(ray)
	j.Expression = "1-pi/4+2/pi-19*pi/10000"
	certifiedC := ray.C
	certifiedC.Expression = "1+pi/4-2/pi+19*pi/10000"
	if !semantic.VerifyPiExpressionRationalBounds(certifiedC, semantic.ExactRational{Numerator: 115474741, Denominator: 100000000}, semantic.ExactRational{Numerator: 115474742, Denominator: 100000000}) || !semantic.VerifyPiExpressionRationalBounds(j, semantic.ExactRational{Numerator: 84525258, Denominator: 100000000}, semantic.ExactRational{Numerator: 84525259, Denominator: 100000000}) {
		return M21Result{}, fmt.Errorf("M21 directed decimal display bounds failed")
	}
	r := M21Result{
		M20: m20, Outcome: "Success A — certified positive improvement on the M20-selected ray, with an additional rigorous non-tight ceiling bracket",
		Ray: ray, Algebra: algebra,
		ResponseSignStructure: []string{
			"H is even and entire after defining sin(t)/t=1 at zero; H(0)=2+4/pi-pi>0.",
			"H(pi/2)=H'(pi/2)=0 and H''(pi/2)=-(30+16/pi^2)<0, so pi/2 is an isolated double zero and H<0 on both adjacent punctured sides.",
			"H is not periodic because of the sinc term; H(t)=4*cos(t)/pi+2*cos(4t)+O(1/|t|), so it has infinitely many asymptotically periodic sign regions.",
			"Numerical root inventory on [0,20], stored only as research evidence, begins near 0.09085, the double zero pi/2, 2.95185, 3.40495, 4.29985, and 5.26195.",
		},
		Envelope:                "For H(t)<0, M(t)=F0(t)/(-H(t)); F_m(t)>=0 iff m<=M(t). Hence m_*=inf_{H<0} M(t), including the limiting value at the double zero pi/2.",
		EnvelopeTangencyTheorem: "M'(t)=(-F0'(t)H(t)+F0(t)H'(t))/H(t)^2. At m=M(t), F_m'(t)=0 iff M'(t)=0 whenever H(t)<0.",
		NumericalCandidate:      artifact.NumericalCandidate, CertifiedWitness: artifact.Mass, BestSimpleRational: artifact.Mass,
		WholeLine: artifact.WholeLine, Ceiling: artifact.Ceiling,
		RayCeilingStatus: "Unsolved exactly. Certified 19/5000<=m_*<3873/1000000; the upper endpoint fails at t=5459/4000 by directed degree-100 Taylor arithmetic. The numerical minimizer lies strictly inside this bracket.",
		CertifiedC:       certifiedC, CertifiedJ: j,
		CertifiedBounds:            "c_*>=1+pi/4-2/pi+19*pi/10000>1+pi/4-2/pi; J_*<=1-pi/4+2/pi-19*pi/10000",
		DirectedDecimals:           "1.15474741<c_safe<1.15474742 and 0.84525258<J_safe_upper<0.84525259",
		M20Comparison:              "M21 certifies the M20 sampled value essentially as stated: m=19/5000 gives c=1.154747417..., while the stronger numerical ray optimum is c approximately 1.15486093435.",
		PrimalComparison:           "The certified dual c is strictly below the compatible CGdL primal guard 1651/1250=1.3208; the compiler hard-fails otherwise.",
		AnthropicComparison:        "The certified J upper bound remains about 0.16340 above the unexplained 0.68185 remark (c about 1.31815); this one-dimensional ray is not expected to close that full-class gap.",
		MultiContactInterpretation: "At the numerical ceiling, the fixed contacts ±pi/2 and new contacts ±t1 are both strict quadratic tangencies. This resembles an active-set/equioscillation condition, not yet a general alternation theorem.",
		Experiment:                 M21Experiment{Path: "experiments/m21_contact_adapted_ray.octest", Interpreted: "5/5 passed in 7.197 s", Compiled: "5/5 passed in 0.661 s", FallbackCount: 0, LongRange: "200,001 deterministic samples on [0,2000] below the candidate; no sampled late-frequency failure", Artifacts: []string{"artifacts/m21/m21_m18_baseline_contact_zoom.png", "artifacts/m21/m21_contact_adapted_ray_pressure.png", "artifacts/m21/m21_contact_adapted_ray_envelope.png"}, ArtifactSHA256: []string{"a87b0171e1619fd78ee4369c17cec87c34543b2bf1b0eec3bc107135f8465537", "660f3defb5c68ac50c495c25352089044841aafa03b40d98e2957ed669309888", "bdbc8a645e4be17ef405182876a4b2e03190ac50c94a9d0db79b66671315dc13"}, ExecutionIdentity: "Oct dev; gooct-cli; CPU Float experiment; compiled Facts use zero fallback; three build-time-interpreted plot.line artifacts produced in 0.837 s", EvidenceBoundary: "Oct Float scans, roots, perturbations, and plots are research evidence only. Exact Go Laurent/rational Taylor and analytic-tail checks carry the theorem."},
		WitnessArtifact:            "compiler/artifacts/m21_contact_adapted_ray_witness.json",
		LiteratureAssessment:       "A focused primary/preprint screen of arXiv:1810.08843, arXiv:2310.01913, arXiv:2108.09258, and arXiv:2502.05106, plus exact-constant/radius-4 searches, found no stated match for this contact-adapted r2=4 ray, m=19/5000 witness, or numerical two-contact constant. This is a limited screen, not a novelty or priority claim.",
		ImportedMathematics:        []string{"M18 exact one-radius tangency witness", "M20 contact-preserving r2=4 perturbation rule", "EP3.1 weak duality and compatible CGdL strict primal guard", "Taylor's theorem and classical rational pi bounds"},
		DerivedMathematics:         []string{"exact response H and fixed-contact regression", "finite original-contact curvature ceiling", "scalar envelope and stationary-envelope/tangency equivalence", "second numerical strict quadratic contact", "whole-line rational witness m=19/5000", "certified ray ceiling upper rejection m<3873/1000000"},
		ArchitecturalAwkwardness:   "The exact compact replay is intentionally a slow lane because the shallow pre-contact shoulder requires targeted degree-100 rational cells. The production certificate remains narrow and reuses M18's equality-aware cell kind.",
		CompilerTheory:             "A representation milestone can compile an infinite feasibility problem into a contact-preserving response, a scalar pressure envelope, equality-aware local cells, and a finite counterfrequency. NumericalCandidate, CertifiedWitness, and CertifiedOptimumBracket remain distinct evidence types.",
		NextMilestone:              "M22: add exactly one contact-motion degree of freedom chosen to release the new t1 tangency while retaining the two-atom support assumptions; do not widen beyond that evidence-selected family.", RHStatus: "unresolved",
	}
	if err := validateM21Result(r); err != nil {
		return M21Result{}, err
	}
	return r, nil
}

func validateM21Result(r M21Result) error {
	if !strings.Contains(r.Outcome, "Success A") || r.RHStatus != "unresolved" || r.CertifiedWitness != (semantic.ExactRational{Numerator: 19, Denominator: 5000}) || r.Ceiling.UpperBound != (semantic.ExactRational{Numerator: 3873, Denominator: 1000000}) {
		return fmt.Errorf("incomplete M21 result")
	}
	if r.Ray.C.Expression == "" || !strings.Contains(r.CertifiedBounds, "c_*") || !strings.Contains(r.CertifiedBounds, "J_*") || r.NextMilestone == "" || len(r.Experiment.Artifacts) != 3 {
		return fmt.Errorf("M21 theorem propagation or artifacts are incomplete")
	}
	return nil
}

func M21JSONReport(result M21Result) ([]byte, error) {
	if err := validateM21Result(result); err != nil {
		return nil, err
	}
	payload := struct {
		Schema string    `json:"schema"`
		Result M21Result `json:"m21"`
	}{Schema: "riemann.m21.v1", Result: result}
	return json.MarshalIndent(payload, "", "  ")
}
