package compiler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

//go:embed artifacts/m17_one_radius_witness.json
var m17WitnessJSON []byte

type M17WitnessArtifact struct {
	Family               semantic.UnsaturatedOneRadiusCompletion `json:"family"`
	PositiveDefiniteness semantic.OneRadiusWholeLineCertificate  `json:"positive_definiteness"`
	FamilyCeiling        semantic.OneRadiusFamilyCeiling         `json:"family_ceiling"`
	Provenance           string                                  `json:"provenance"`
}

type M17Experiment struct {
	SearchPath        string   `json:"search_path"`
	ArtifactPath      string   `json:"artifact_path"`
	Commands          []string `json:"commands"`
	AnalyticReduction string   `json:"analytic_reduction"`
	Strategy          string   `json:"strategy"`
	WhenUtility       string   `json:"when_utility"`
	BestNumerical     string   `json:"best_numerical_candidate"`
	DangerousMinimum  string   `json:"dangerous_frequency_and_minimum"`
	Counterexamples   []string `json:"counterexamples"`
	IndependentCheck  string   `json:"independent_check"`
	Plot              string   `json:"plot"`
	PlotSHA256        string   `json:"plot_sha256"`
	PlotStatus        string   `json:"plot_status"`
	ExecutionIdentity string   `json:"execution_identity"`
}

type M17Result struct {
	M16                      M16Result                               `json:"m16"`
	Outcome                  string                                  `json:"outcome"`
	SignAudit                string                                  `json:"sign_audit"`
	ExactFamily              string                                  `json:"exact_one_radius_family"`
	ExactTransform           string                                  `json:"exact_transform"`
	Family                   semantic.UnsaturatedOneRadiusCompletion `json:"candidate_family"`
	SaturatedRegression      semantic.SaturatedOneRadiusCompletion   `json:"saturated_subfamily_regression"`
	SubfamilyRelation        string                                  `json:"subfamily_relation"`
	M15Regression            string                                  `json:"m15_regression"`
	M16Regression            string                                  `json:"m16_one_radius_regression"`
	LocalAnalysis            string                                  `json:"local_coefficient_analysis"`
	ParameterReduction       string                                  `json:"parameter_reduction"`
	RadiusConclusion         string                                  `json:"radius_conclusion"`
	PDCertificate            semantic.OneRadiusWholeLineCertificate  `json:"whole_line_pd_certificate"`
	Witness                  semantic.DualCompletionWitness          `json:"certified_dual_witness"`
	Bound                    semantic.CertifiedExtremalBound         `json:"certified_bound"`
	FamilyCeiling            semantic.OneRadiusFamilyCeiling         `json:"one_radius_family_ceiling"`
	FamilyCeilingBracket     string                                  `json:"one_radius_family_ceiling_bracket"`
	CertifiedCBracket        string                                  `json:"certified_c_bracket"`
	CertifiedJBracket        string                                  `json:"certified_j_bracket"`
	NineEighthsComparison    string                                  `json:"nine_eighths_comparison"`
	AnthropicComparison      string                                  `json:"anthropic_comparison"`
	WitnessArtifact          string                                  `json:"witness_artifact"`
	CeilingArtifact          string                                  `json:"ceiling_artifact"`
	Experiment               M17Experiment                           `json:"oct_experiment"`
	LiteratureAssessment     string                                  `json:"literature_assessment"`
	ImportedMathematics      []string                                `json:"imported_mathematics"`
	DerivedMathematics       []string                                `json:"compiler_research_derived_mathematics"`
	ArchitecturalAwkwardness string                                  `json:"architectural_awkwardness"`
	CompilerTheory           string                                  `json:"compiler_theory"`
	NextMilestone            string                                  `json:"one_next_milestone"`
	RHStatus                 string                                  `json:"rh_status"`
}

func CompileM17() (M17Result, error) {
	return compileM17(true)
}

// compileM17 separates cheap result construction from whole-line proof replay.
// Production compilation verifies. Normal unit tests use the structural path;
// the slow-tag suite exercises this function with verifyWholeLine=true.
func compileM17(verifyWholeLine bool) (M17Result, error) {
	m16, err := CompileM16()
	if err != nil {
		return M17Result{}, err
	}
	return compileM17FromM16(m16, verifyWholeLine)
}

func compileM17FromM16(m16 M16Result, verifyWholeLine bool) (M17Result, error) {
	var artifact M17WitnessArtifact
	if err := json.Unmarshal(m17WitnessJSON, &artifact); err != nil {
		return M17Result{}, fmt.Errorf("decode M17 witness artifact: %w", err)
	}
	if err := artifact.Family.Validate(); err != nil {
		return M17Result{}, fmt.Errorf("verify M17 typed family: %w", err)
	}
	if verifyWholeLine {
		if err := semantic.VerifyOneRadiusCertificate(artifact.Family, artifact.PositiveDefiniteness); err != nil {
			return M17Result{}, fmt.Errorf("verify M17 witness artifact: %w", err)
		}
	}
	if artifact.FamilyCeiling != semantic.DeriveOneRadiusFamilyCeiling() {
		return M17Result{}, fmt.Errorf("M17 ceiling artifact does not match the derived theorem")
	}

	m15Family := semantic.NewUnsaturatedOneRadiusCompletion(
		semantic.ExactRational{Numerator: 9, Denominator: 8}, semantic.ExactRational{Numerator: 1, Denominator: 1}, semantic.ExactRational{Numerator: 1, Denominator: 8},
		semantic.Reference{Kind: semantic.CompilerRecord, Citation: "M17 structural M15 regression"})
	saturated, err := semantic.AsSaturatedOneRadius(m15Family)
	if err != nil || semantic.SaturatedOneRadiusCeiling(saturated) != nil {
		return M17Result{}, fmt.Errorf("M15 saturated regression failed: %v", err)
	}

	// Re-run M16's w2=0 degeneration through the new one-radius verifier.
	m16One := semantic.NewUnsaturatedOneRadiusCompletion(
		semantic.ExactRational{Numerator: 573, Denominator: 500}, semantic.ExactRational{Numerator: 1, Denominator: 1}, semantic.ExactRational{Numerator: 21, Denominator: 125},
		semantic.Reference{Kind: semantic.CompilerRecord, Citation: "M17 exact M16 one-radius regression"})
	m16PD := semantic.OneRadiusWholeLineCertificate{
		ID: "m16-one-radius-regression", FourierVariable: "t=2*pi*xi", FourierDensity: semantic.OneRadiusFourierDensity,
		CompactInterval: "0<=|t|<=4", GridStep: semantic.ExactRational{Numerator: 1, Denominator: 1000}, TaylorDegree: 40,
		LipschitzBound: semantic.ExactRational{Numerator: 1223, Denominator: 1500}, TailInterval: "|t|>=4",
		TailLowerBound: "hat(P)>=1-2*w-4/t^2-2*(c-1)/|t|>0 for |t|>=4", TailAnchor: semantic.ExactRational{Numerator: 341, Denominator: 1000},
		OmittedDirections: "exact Taylor enclosures plus a Lipschitz cover on |t|<=4 and a radius-uniform analytic bound on |t|>=4 cover every real t", WholeLine: true,
	}
	if verifyWholeLine {
		if err := semantic.VerifyOneRadiusCertificate(m16One, m16PD); err != nil {
			return M17Result{}, fmt.Errorf("M16 one-radius regression failed: %w", err)
		}
	}

	witness := semantic.DualCompletionWitness{
		ID: "unsaturated-one-radius-2297-over-2000", MultiplicityLower: artifact.Family.Constant,
		OutsideMeasure:         "sigma=(2297/2000)1_{|x|>1}dx+(171/1000)(delta_-1+delta_1)",
		CompletionDistribution: "P=delta_0+(|x|-2297/2000)1_{[-1,1]}dx+(171/1000)(delta_-1+delta_1)",
		FourierImage:           artifact.PositiveDefiniteness.FourierDensity, PositivityEvidence: semantic.GlobalAnalyticPD,
		ExactSupportScope: true, WholeLineControl: true, TailControl: true,
		Provenance: semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Riemann M17 exact one-radius Taylor-enclosure and coercive-tail certificate"},
	}
	bound, err := semantic.ApplyDualCompletion(m16.M15.M14A.Class, witness)
	if err != nil {
		return M17Result{}, err
	}

	r := M17Result{
		M16: m16, Outcome: "Success B — improved witness and narrow one-radius family ceiling",
		SignAudit:      "The nonnegative measure has +w(delta_-r+delta_r). The isolated minus sign in the request conflicts with w>=0, the supplied Taylor coefficients, M16, and the final conceptual block; no signed-atom family is introduced.",
		ExactFamily:    "sigma=c 1_{|x|>1}dx+w(delta_-r+delta_r), c>0, r>=1, w>=0",
		ExactTransform: semantic.OneRadiusFourierDensity,
		Family:         artifact.Family, SaturatedRegression: saturated,
		SubfamilyRelation:  "SaturatedOneRadius is a checked strict subtype of UnsaturatedOneRadius obtained only when a0=2(1+w-c)=0; saturation is not family identity.",
		M15Regression:      "On a0=0, w=c-1. Since r>=1, a2=c/3-1/4-w*r^2<=(9-8c)/12, so local nonnegativity forces c<=9/8; (c,r,w)=(9/8,1,1/8) recovers M15 exactly.",
		M16Regression:      "The new one-radius verifier independently certifies (c,r,w)=(573/500,1,21/125), with a0=11/250>0 and a2=-9/250; it does not construct or depend on TwoRadiusCompletion.",
		LocalAnalysis:      "a0=2(1+w-c), a2=c/3-1/4-w*r^2, a4=-1/360+(1-c)/60+w*r^4/12, a6=1/20160+(c-1)/2520-w*r^6/360. Thus a0>=0 is necessary. If a0>0, a2<0 is legal; the quartic truncation a0+a2*u+a4*u^2 (u=t^2) has positive minimum on u>=0 when a4>0 and 4*a0*a4>=a2^2, a useful sufficient local screen but not a global certificate. The witness has (a0,a2,a4,a6)=(9/200,-229/6000,3239/360000,-1847/5040000).",
		ParameterReduction: "For fixed (r,w), F is affine decreasing in c at sin(t)>0, so c<=inf_{sin(t)>0} [1+2(cos(t)-1)/t^2+2*sin(t)/t+2*w*cos(r*t)]*t/(2*sin(t)), together with c<=1+w at t=0; sin(t)<=0 frequencies are feasibility checks. Tail subsequences with cos(r*t)->-1 force w<=1/2. The contact t=pi/(2r) removes w entirely and yields the global ceiling theorem.",
		RadiusConclusion:   "The reduced scan selects r=1. Analytically, h(pi/(2r)) is strictly smaller than the universal bound for every r>1, so only the boundary radius can approach that ceiling. Exact attainment at the transcendental endpoint was not needed or claimed.",
		PDCertificate:      artifact.PositiveDefiniteness, Witness: witness, Bound: bound, FamilyCeiling: artifact.FamilyCeiling,
		FamilyCeilingBracket:  "2297/2000 <= C_1R <= 1+pi/4-2/pi < 1149/1000",
		CertifiedCBracket:     "2297/2000 <= c_* < 1651/1250 (the upper endpoint is the imported strict primal bound)",
		CertifiedJBracket:     "849/1250 < J_* <= 1703/2000",
		NineEighthsComparison: "2297/2000-9/8=47/2000: unsaturation gains a certified 0.0235 over the exact saturated ceiling.",
		AnthropicComparison:   "The one-radius ceiling corresponds to J>=1-(pi/4)+2/pi approximately 0.85122 within this family, far above 0.68185. Anthropic's value remains unresolved for EP3.1 and is neither confirmed nor contradicted.",
		WitnessArtifact:       "compiler/artifacts/m17_one_radius_witness.json", CeilingArtifact: "compiler/artifacts/m17_one_radius_witness.json#family_ceiling",
		Experiment: M17Experiment{
			SearchPath: "experiments/m17_one_radius_completion.octest", ArtifactPath: "experiments/m17_one_radius_plot.octest",
			Commands:          []string{"go run ./cmd/oct test .../experiments/m17_one_radius_completion.octest --execution interpreted --json", "go run ./cmd/oct test .../experiments/m17_one_radius_completion.octest --execution compiled --json", "go run ./cmd/oct artifact .../experiments/m17_one_radius_plot.octest --output-root .../artifacts --execution compiled --json"},
			AnalyticReduction: "eliminate c by the sin-positive frequency envelope; bound 0<=w<=1/2; use t=pi/(2r) for a radius-dependent ceiling before scanning",
			Strategy:          "incremental radius scan [1,12], fixed-r concave envelope optimization in w, dense frequency feasibility scan, then rational candidates ranked by margin and exact proof cost",
			WhenUtility:       "used after the envelope scan: 2297/2000 was preferred to a knife-edge 1.14875 candidate because its 3.5e-4 margin certifies on a 1/5000 grid with simple weight 171/1000",
			BestNumerical:     "c approximately 1.1487788, r=1, w approximately 0.1713, contact t approximately pi/2; this estimates but does not certify the endpoint",
			DangerousMinimum:  "retained rational witness sampled minimum approximately 0.00035437 at t approximately 1.57195 on [0,80]",
			Counterexamples:   []string{"saturated c>9/8 fails at the origin", "excessive w>1/2 fails along atom-negative tail subsequences", "r>1 pays the strict contact bound h(pi/(2r))<h(pi/2)", "c=1.149 exceeds the rational family ceiling", "locally legal candidates can still develop a first-oscillation negative pocket"},
			IndependentCheck:  "Oct independently reconstructs exact rational inputs in Float, coefficients, M15/M16 regressions, contact geometry, dangerous minimum, and tail samples; all numerical statements remain non-certifying.",
			Plot:              "artifacts/m17/m17_one_radius_spectrum.png", PlotSHA256: "1f0dd33ec965e8be79acc0faa0dc4fb853cd1273999908b23bdec46975ab4e82", PlotStatus: "produced on first run and unchanged on deterministic rerun; 11,866 bytes; identity M17OneRadiusPlot.PlotCandidateMinimumAndTailHandoff:m17/m17_one_radius_spectrum.png", ExecutionIdentity: "Oct gooct-cli; interpreted 6/6; compiled 6/6 with zero fallback; artifact build-time-interpreted compatibility phase",
		},
		LiteratureAssessment:     "A focused post-certification search of arXiv:1810.08843, arXiv:2310.01913, and arXiv:2502.05106 found the relevant SDP, Cohn-Elkies, EP, and pair-correlation extremal frameworks, but no matching one-boundary-atom completion, constant 1+pi/4-2/pi, or saturated/unsaturated representation theorem. This is a limited novelty screen, not a priority claim.",
		ImportedMathematics:      []string{"EP3.1 weak dual and plus-sign data measure", "CGdL strict primal certificate c_*<1651/1250", "Bochner-Schwartz positive-definite distribution criterion", "Taylor's theorem with bounded trigonometric derivatives", "classical rational bounds 333/106<pi<355/113"},
		DerivedMathematics:       []string{"typed UnsaturatedOneRadiusCompletion and checked SaturatedOneRadius subtype", "exact one-radius transform and Taylor coefficients", "M15 9/8 saturated regression", "fixed-(r,w) c-envelope and w<=1/2 compactness", "contact ceiling C_1R<=1+pi/4-2/pi<1149/1000", "r=1 rational witness c=2297/2000 with exact compact and tail certificate"},
		ArchitecturalAwkwardness: "The exact ceiling naturally contains pi while the witness language is rational. M17 keeps the analytic ceiling theorem and rational enclosure distinct instead of expanding the scalar IR or pretending the endpoint is rational.",
		CompilerTheory:           "M15's 9/8 barrier was a property of the accidentally origin-saturated subtype, not of one-radius representation. Making saturation an explicit refinement exposes additional legal states and raises the certified dual to 2297/2000; a representation compiler must distinguish an object's identity from an active boundary condition.",
		NextMilestone:            "M18: certify or refute attainment of the one-radius endpoint c=1+pi/4-2/pi at r=1 with the tangency-determined weight, using a purpose-specific analytic whole-line proof and no additional radii.",
		RHStatus:                 "unresolved",
	}
	if err := validateM17Result(r); err != nil {
		return M17Result{}, err
	}
	return r, nil
}

func validateM17Result(r M17Result) error {
	if r.Outcome == "" || r.RHStatus != "unresolved" || !r.Witness.CertifiesInfiniteClass() {
		return fmt.Errorf("incomplete M17 certified result")
	}
	if r.Bound.MultiplicityLower != (semantic.ExactRational{Numerator: 2297, Denominator: 2000}) || r.Bound.SimpleUpper != (semantic.ExactRational{Numerator: 1703, Denominator: 2000}) {
		return fmt.Errorf("M17 bound propagation failed")
	}
	if !(semantic.ExactRational{Numerator: 1651, Denominator: 1250}).GreaterThan(r.Bound.MultiplicityLower) {
		return fmt.Errorf("certified dual crosses the compatible strict primal upper bound")
	}
	if !strings.Contains(r.FamilyCeilingBracket, "C_1R") || !strings.Contains(r.M15Regression, "9/8") || !strings.Contains(r.M16Regression, "573/500") {
		return fmt.Errorf("M17 family boundary report is incomplete")
	}
	return r.Family.Validate()
}
