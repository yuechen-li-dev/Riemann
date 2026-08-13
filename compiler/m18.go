package compiler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

//go:embed artifacts/m18_exact_one_radius_extremal.json
var m18WitnessJSON []byte

type M18WitnessArtifact struct {
	Candidate          semantic.OneRadiusTangencyCandidate            `json:"candidate"`
	Contact            semantic.ContactEquationCertificate            `json:"contact"`
	OriginTaylor       semantic.OneRadiusExactTaylor                  `json:"origin_taylor"`
	WholeLine          semantic.OneRadiusTangencyWholeLineCertificate `json:"whole_line"`
	M17UpperTheorem    string                                         `json:"m17_upper_theorem"`
	ExactFamilyOptimum string                                         `json:"exact_family_optimum"`
	ContactMetadata    struct {
		T            string `json:"t"`
		G            string `json:"G"`
		GPrime       string `json:"G_prime"`
		GSecond      string `json:"G_second"`
		Multiplicity int    `json:"multiplicity"`
	} `json:"contact_metadata"`
	Provenance string `json:"provenance"`
}

type M18Experiment struct {
	Path               string `json:"path"`
	PlotSource         string `json:"plot_source"`
	InterpretedFacts   string `json:"interpreted_facts"`
	CompiledFacts      string `json:"compiled_facts"`
	FallbackCount      int    `json:"fallback_count"`
	BestSampledMinimum string `json:"best_sampled_minimum"`
	ContactLocation    string `json:"contact_location"`
	PerturbedFailure   string `json:"perturbed_candidate_failure"`
	LongRangeScan      string `json:"long_range_scan"`
	ArtifactPath       string `json:"artifact_path"`
	ArtifactSHA256     string `json:"artifact_sha256"`
	ExecutionIdentity  string `json:"execution_identity"`
	EvidenceBoundary   string `json:"evidence_boundary"`
}

type M18Result struct {
	M17                      M17Result                                      `json:"m17"`
	Outcome                  string                                         `json:"outcome"`
	SignAudit                string                                         `json:"sign_audit"`
	Candidate                semantic.OneRadiusTangencyCandidate            `json:"candidate"`
	Contact                  semantic.ContactEquationCertificate            `json:"contact"`
	ContactMultiplicity      string                                         `json:"contact_multiplicity"`
	LocalPositivity          string                                         `json:"local_positivity"`
	OriginTaylor             semantic.OneRadiusExactTaylor                  `json:"origin_taylor"`
	OriginBehavior           string                                         `json:"origin_behavior"`
	WholeLine                semantic.OneRadiusTangencyWholeLineCertificate `json:"whole_line_certificate"`
	CompactVerifier          string                                         `json:"compact_verifier"`
	TailVerifier             string                                         `json:"tail_verifier"`
	WholeLineCoverage        string                                         `json:"whole_line_coverage"`
	Optimum                  semantic.ExactOneRadiusOptimum                 `json:"one_radius_optimum"`
	ExactAttainmentStatus    string                                         `json:"exact_attainment_status"`
	CertifiedCBracket        string                                         `json:"certified_c_bracket"`
	CertifiedJBracket        string                                         `json:"certified_j_bracket"`
	DirectedDecimals         string                                         `json:"directed_decimal_displays"`
	M17Comparison            string                                         `json:"m17_comparison"`
	AnthropicComparison      string                                         `json:"anthropic_comparison"`
	Experiment               M18Experiment                                  `json:"oct_experiment"`
	TangencyPlot             string                                         `json:"tangency_plot"`
	WitnessArtifact          string                                         `json:"witness_artifact"`
	IndependentVerification  string                                         `json:"independent_verification"`
	LiteratureAssessment     string                                         `json:"literature_assessment"`
	ImportedMathematics      []string                                       `json:"imported_mathematics"`
	DerivedMathematics       []string                                       `json:"compiler_research_derived_mathematics"`
	ArchitecturalAwkwardness string                                         `json:"architectural_awkwardness"`
	CompilerTheory           string                                         `json:"compiler_theory"`
	ContactMetadata          map[string]string                              `json:"contact_metadata_for_future_aetheris_audit"`
	NextMilestone            string                                         `json:"one_next_milestone"`
	RHStatus                 string                                         `json:"rh_status"`
}

func CompileM18() (M18Result, error) { return compileM18(true, true) }

func compileM18(verifyM17, verifyM18 bool) (M18Result, error) {
	m17, err := compileM17(verifyM17)
	if err != nil {
		return M18Result{}, err
	}
	return compileM18FromM17(m17, verifyM18)
}

func compileM18FromM17(m17 M17Result, verifyWholeLine bool) (M18Result, error) {
	var artifact M18WitnessArtifact
	if err := json.Unmarshal(m18WitnessJSON, &artifact); err != nil {
		return M18Result{}, fmt.Errorf("decode M18 witness artifact: %w", err)
	}
	derived := semantic.DeriveOneRadiusTangencyCandidate()
	if err := artifact.Candidate.Validate(); err != nil || !reflect.DeepEqual(artifact.Candidate, derived) {
		return M18Result{}, fmt.Errorf("M18 candidate is not the derived tangency solution: %v", err)
	}
	contact, err := semantic.VerifyTangencyAlgebra(artifact.Candidate)
	if err != nil || !reflect.DeepEqual(contact, artifact.Contact) {
		return M18Result{}, fmt.Errorf("M18 exact contact artifact mismatch: %v", err)
	}
	origin := semantic.TangencyOriginTaylor(artifact.Candidate)
	if !reflect.DeepEqual(origin, artifact.OriginTaylor) {
		return M18Result{}, fmt.Errorf("M18 origin Taylor artifact mismatch")
	}
	if artifact.M17UpperTheorem != "C_1R<=1+pi/4-2/pi" || artifact.ExactFamilyOptimum != "C_1R=1+pi/4-2/pi" || artifact.ContactMetadata.T != "pi/2" || artifact.ContactMetadata.G != "0" || artifact.ContactMetadata.GPrime != "0" || artifact.ContactMetadata.GSecond != ">1/8" || artifact.ContactMetadata.Multiplicity != 2 {
		return M18Result{}, fmt.Errorf("M18 theorem/contact metadata is incomplete")
	}
	if m17.FamilyCeiling != semantic.DeriveOneRadiusFamilyCeiling() || m17.FamilyCeiling.ExactUpperExpression != artifact.Candidate.ExactCExpression.Expression {
		return M18Result{}, fmt.Errorf("M17 ceiling cannot compose with M18 witness")
	}
	optimum := semantic.ExactOneRadiusOptimum{Family: "UnsaturatedOneRadius", ExactValue: artifact.Candidate.ExactCExpression, LowerRoute: "M18 globally nonnegative exact r=1 tangency witness", UpperRoute: "M17 one-radius frequency envelope", EqualityDerived: true}
	if verifyWholeLine {
		optimum, err = semantic.ComposeExactOneRadiusOptimum(artifact.Candidate, artifact.WholeLine, m17.FamilyCeiling)
		if err != nil {
			return M18Result{}, err
		}
	}
	j := semantic.DeriveTangencySimpleUpper(artifact.Candidate)
	if !semantic.VerifyPiExpressionRationalBounds(artifact.Candidate.ExactCExpression, semantic.ExactRational{Numerator: 114877839, Denominator: 100000000}, semantic.ExactRational{Numerator: 114877840, Denominator: 100000000}) || !semantic.VerifyPiExpressionRationalBounds(artifact.Candidate.ExactWExpression, semantic.ExactRational{Numerator: 17103742, Denominator: 100000000}, semantic.ExactRational{Numerator: 17103744, Denominator: 100000000}) || !semantic.VerifyPiExpressionRationalBounds(j, semantic.ExactRational{Numerator: 85122160, Denominator: 100000000}, semantic.ExactRational{Numerator: 85122161, Denominator: 100000000}) {
		return M18Result{}, fmt.Errorf("directed decimal display bounds failed")
	}

	r := M18Result{
		M17: m17, Outcome: "Success A — exact attainment of the unsaturated one-radius ceiling",
		SignAudit: "M18 retains M17's nonnegative-measure convention sigma=c 1_{|x|>1}dx+w(delta_-1+delta_1) and its plus-sign Fourier density. The sign-reversed displayed G in the request would not yield the stated ceiling at pi/2.",
		Candidate: artifact.Candidate, Contact: artifact.Contact, ContactMultiplicity: "F(pi/2)=F'(pi/2)=0 and F''(pi/2)>1/8, so the zero has multiplicity exactly two and is a strict local minimum.",
		LocalPositivity: "On 3/2<=t<=17/10, |t-pi/2|<13/100. The integral representation gives |F'''|<11/20; hence F''>1/8-(11/20)(13/100)=107/2000>0. Exact contact and strict convexity imply F>=0, with equality only at pi/2 in this region.",
		OriginTaylor:    origin, OriginBehavior: "a0=2*(1+w-c)>1/25 (numerically about 0.04451807), so the witness is safely unsaturated. The subsequent coefficients are retained exactly as Laurent expressions in pi.",
		WholeLine:         artifact.WholeLine,
		CompactVerifier:   "Equality-aware partition: the contact cell [3/2,17/10] is discharged analytically; [0,3/2] and [17/10,4] use exact degree-40 sine/cosine Taylor enclosures on a rational 1/2500 grid and the certified |F'|<83/100 loss. An interval merely containing zero is never accepted.",
		TailVerifier:      "For |t|>=4, F>=1-2w-4/t^2-2(c-1)/|t|>=3/4-2w-(c-1)/2>1/3, using certified pi bounds.",
		WholeLineCoverage: artifact.WholeLine.Coverage + "; evenness supplies t<0 and endpoints overlap without gaps.",
		Optimum:           optimum, ExactAttainmentStatus: "attained: certified witness lower bound plus the M17 family upper theorem derives C_1R=1+pi/4-2/pi; this is not the full EP3.1 optimum.",
		CertifiedCBracket:   "1+pi/4-2/pi <= c_* < 1651/1250",
		CertifiedJBracket:   "849/1250 < J_* <= 1-pi/4+2/pi",
		DirectedDecimals:    "1.14877839 < 1+pi/4-2/pi < 1.14877840; 0.17103742 < w < 0.17103744; 0.85122160 < 1-pi/4+2/pi < 0.85122161 (the displayed J upper endpoint is rounded upward)",
		M17Comparison:       "M17's best numerical c approximately 1.1487788 was within about 4.1e-7 of the exact value; its retained rational 2297/2000 was deliberately below the zero-margin endpoint.",
		AnthropicComparison: "The exact one-radius J upper bound is approximately 0.851221609, about 0.16937 above 0.68185. The latter remains unsupported for the full EP3.1 problem and is neither proved nor refuted here.",
		Experiment: M18Experiment{
			Path: "experiments/m18_exact_tangency.octest", PlotSource: "experiments/m18_tangency_plot.octest",
			InterpretedFacts: "6/6 passed", CompiledFacts: "6/6 passed", FallbackCount: 0,
			BestSampledMinimum: "about 7.16e-15 at t=1.570796 on the 1e-6 local grid; about 2.67e-9 at t=1.571 on the 1e-3 [0,2000] grid",
			ContactLocation:    "t=pi/2 approximately 1.5707963267948966; sampled F and F' within 2e-15 of zero; F'' approximately 0.12891354",
			PerturbedFailure:   "c -> c+1e-6 gives F(pi/2) approximately -1.2732395e-6=-4e-6/pi; changing w cannot alter this frequency",
			LongRangeScan:      "2,000,001 samples on [0,2000], no sampled negative pocket",
			ArtifactPath:       "artifacts/m18/m18_exact_tangency.png", ArtifactSHA256: "f37dd4e48fd82b7c4b4a2fe3097decc7338b11b8616a7ff03b810b51c3e19b3d",
			ExecutionIdentity: "Oct dev gooct-cli; compiled test backend for 6 Facts; artifact build-time-interpreted compatibility path; plot.line, 13,794 bytes",
			EvidenceBoundary:  "independent Float reconstruction, perturbation, plotting, and adversarial scans only; no Oct sample enters the theorem route",
		},
		TangencyPlot: "artifacts/m18/m18_exact_tangency.png", WitnessArtifact: "compiler/artifacts/m18_exact_one_radius_extremal.json",
		IndependentVerification:  "Go exact Laurent algebra plus rational Taylor/interval and analytic-tail verifier is the theorem route; Oct independently reconstructs formulas, contact, curvature, perturbation failure, and long-range samples.",
		LiteratureAssessment:     "A focused primary/preprint screen of arXiv:1810.08843 (CGdL SDP pair-correlation bounds), arXiv:2310.01913 (Cohn-Elkies Fourier optimization for Montgomery averages), and arXiv:2502.05106 (general pair-correlation Fourier problems), plus exact-constant searches, found no stated one-boundary-atom completion, weight 1/pi-4/pi^2+8/pi^3, or constant 1+pi/4-2/pi. This is a limited novelty screen, not a priority claim.",
		ImportedMathematics:      []string{"M17 one-radius envelope C_1R<=1+pi/4-2/pi", "EP3.1 weak duality and CGdL strict primal bound c_*<1651/1250", "Taylor's theorem and classical rational bounds for pi"},
		DerivedMathematics:       []string{"exact tangency weight w=1/pi-4/pi^2+8/pi^3", "symbolic double-contact identities and F''(pi/2)>1/8", "certified local convexity around the zero", "equality-aware compact certificate and strict t>=4 tail", "exact composition C_1R=1+pi/4-2/pi"},
		ArchitecturalAwkwardness: "The older witness and extremal-bound IR stores only rationals. M18 adds a narrow Laurent-in-pi scalar rather than decimalizing the endpoint or introducing a generic CAS; older milestones remain unchanged.",
		CompilerTheory:           "M15 exposed an accidental saturation boundary, M16 escaped it, M17 isolated the one-radius envelope, and M18 shows that this infinite nonnegativity problem really does compile to a finite double-contact condition plus a global equality-aware verifier.",
		ContactMetadata:          map[string]string{"t": "pi/2", "F": "0", "F_prime": "0", "F_second": ">1/8", "multiplicity": "2", "semantics": "strict quadratic analytic contact; no geometry or STEP lowering"},
		NextMilestone:            "M19: audit Aetheris for mathematical-research semantics, focusing on exact analytic contact/tangency metadata and faithful analytic-to-geometric lowering; do not integrate it until the audit establishes suitability.",
		RHStatus:                 "unresolved",
	}
	if err := validateM18Result(r); err != nil {
		return M18Result{}, err
	}
	return r, nil
}

func validateM18Result(r M18Result) error {
	if r.Outcome == "" || r.RHStatus != "unresolved" || !r.Optimum.EqualityDerived || r.Optimum.ExactValue.Expression != "1+pi/4-2/pi" {
		return fmt.Errorf("incomplete M18 exact optimum")
	}
	if !strings.Contains(r.CertifiedCBracket, "c_*") || !strings.Contains(r.CertifiedJBracket, "J_*") || !strings.Contains(r.ExactAttainmentStatus, "not the full EP3.1") {
		return fmt.Errorf("M18 theorem scope or bracket is incomplete")
	}
	if len(r.ContactMetadata) != 6 || r.ContactMetadata["multiplicity"] != "2" || r.NextMilestone == "" {
		return fmt.Errorf("M18 contact metadata or next milestone is incomplete")
	}
	return r.Candidate.Validate()
}
