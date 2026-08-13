package compiler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

//go:embed artifacts/m16_two_radius_witness.json
var m16WitnessJSON []byte

type M16WitnessArtifact struct {
	Family               semantic.TwoRadiusCompletion           `json:"family"`
	OriginTaylor         semantic.TwoRadiusTaylor               `json:"origin_taylor"`
	PositiveDefiniteness semantic.TwoRadiusWholeLineCertificate `json:"positive_definiteness"`
	Provenance           string                                 `json:"provenance"`
}

type M16Experiment struct {
	SearchPath        string   `json:"search_path"`
	ArtifactPath      string   `json:"artifact_path"`
	Commands          []string `json:"commands"`
	Strategy          string   `json:"strategy"`
	BestNumerical     string   `json:"best_numerical_candidate"`
	Counterexamples   []string `json:"counterexamples"`
	IndependentCheck  string   `json:"independent_check"`
	Plot              string   `json:"plot"`
	PlotSHA256        string   `json:"plot_sha256"`
	PlotStatus        string   `json:"plot_status"`
	ExecutionIdentity string   `json:"execution_identity"`
}

type M16Result struct {
	M15                      M15Result                              `json:"m15"`
	Outcome                  string                                 `json:"outcome"`
	FamilyDerivation         string                                 `json:"family_derivation"`
	ParameterConstraints     string                                 `json:"parameter_constraints"`
	Family                   semantic.TwoRadiusCompletion           `json:"candidate_family"`
	OriginTaylor             semantic.TwoRadiusTaylor               `json:"origin_taylor"`
	LocalNecessaryConditions string                                 `json:"local_necessary_conditions"`
	QuadraticCancellation    string                                 `json:"quadratic_cancellation"`
	ReducedParameterization  string                                 `json:"reduced_parameterization"`
	PDCertificate            semantic.TwoRadiusWholeLineCertificate `json:"whole_line_pd_certificate"`
	Witness                  semantic.DualCompletionWitness         `json:"certified_dual_witness"`
	Bound                    semantic.CertifiedExtremalBound        `json:"certified_bound"`
	CertifiedCBracket        string                                 `json:"certified_c_bracket"`
	CertifiedJBracket        string                                 `json:"certified_j_bracket"`
	AnthropicComparison      string                                 `json:"anthropic_comparison"`
	WitnessArtifact          string                                 `json:"witness_artifact"`
	Experiment               M16Experiment                          `json:"oct_experiment"`
	LiteratureAssessment     string                                 `json:"literature_assessment"`
	ImportedMathematics      []string                               `json:"imported_mathematics"`
	DerivedMathematics       []string                               `json:"compiler_research_derived_mathematics"`
	ArchitecturalAwkwardness string                                 `json:"architectural_awkwardness"`
	CompilerTheory           string                                 `json:"compiler_theory"`
	NextMilestone            string                                 `json:"one_next_milestone"`
	RHStatus                 string                                 `json:"rh_status"`
}

func CompileM16() (M16Result, error) {
	return compileM16(true)
}

func compileM16(verifyWholeLine bool) (M16Result, error) {
	m15, err := CompileM15()
	if err != nil {
		return M16Result{}, err
	}
	return compileM16FromM15(m15, verifyWholeLine)
}

func compileM16FromM15(m15 M15Result, verifyWholeLine bool) (M16Result, error) {
	var artifact M16WitnessArtifact
	if err := json.Unmarshal(m16WitnessJSON, &artifact); err != nil {
		return M16Result{}, fmt.Errorf("decode M16 witness artifact: %w", err)
	}
	derivedTaylor, err := semantic.TwoRadiusOriginTaylor(artifact.Family)
	if err != nil {
		return M16Result{}, err
	}
	if derivedTaylor != artifact.OriginTaylor {
		return M16Result{}, fmt.Errorf("M16 artifact Taylor coefficients do not match exact derivation: got %+v", derivedTaylor)
	}
	if verifyWholeLine {
		if err := semantic.VerifyTwoRadiusCertificate(artifact.Family, artifact.PositiveDefiniteness); err != nil {
			return M16Result{}, fmt.Errorf("verify M16 witness artifact: %w", err)
		}
	}
	witness := semantic.DualCompletionWitness{
		ID: "two-radius-573-over-500", MultiplicityLower: artifact.Family.Constant,
		OutsideMeasure:         "sigma=(573/500)1_{|x|>1}dx+(21/125)(delta_-1+delta_1)+(1/1000)(delta_-2+delta_2)",
		CompletionDistribution: "P=delta_0+(|x|-573/500)1_{[-1,1]}dx+(21/125)(delta_-1+delta_1)+(1/1000)(delta_-2+delta_2)",
		FourierImage:           artifact.PositiveDefiniteness.FourierDensity, PositivityEvidence: semantic.GlobalAnalyticPD,
		ExactSupportScope: true, WholeLineControl: true, TailControl: true,
		Provenance: semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Riemann M16 exact Taylor-enclosure and coercive-tail certificate"},
	}
	bound, err := semantic.ApplyDualCompletion(m15.M14A.Class, witness)
	if err != nil {
		return M16Result{}, err
	}
	r := M16Result{
		M15: m15, Outcome: "Success A — certified two-radius improvement beyond 9/8",
		FamilyDerivation:     "From P=nu-c dx+sigma and nu=delta_0+|x|1_{[-1,1]}dx, cancellation of -c dx outside [-1,1] forces exterior density +c. Legal symmetric atom pairs enter sigma with nonnegative masses, giving the displayed two-radius family; no atom normalization is assumed.",
		ParameterConstraints: "c>0; 1<=r1<=r2; w1,w2>=0; sigma is supported in R\\(-1,1). Radius permutation is canonicalized and coincident radii combine by adding weights.",
		Family:               artifact.Family, OriginTaylor: artifact.OriginTaylor,
		LocalNecessaryConditions: "a0=2(1+w1+w2-c)>=0. If a0=0, the first subsequent nonzero coefficient must be positive. On the saturated branch w1+w2=c-1, a2=c/3-1/4-w1*r1^2-w2*r2^2<=(9-8c)/12, hence c<=9/8. More generally a2=0 and a0>=0 imply c-1<=w1+w2<=sum(wj*rj^2)=c/3-1/4, again forcing c<=9/8.",
		QuadraticCancellation:    "No locally legal candidate with c>9/8 can cancel the quadratic coefficient. The successful witness instead has positive origin slack a0=23/500, so its negative a2=-1/25 is locally legal. Exact verification still passes after setting w2=0, proving that slack—not the second radius—is what removes M15's obstruction.",
		ReducedParameterization:  "Exact algebra splits the search into saturated and slack branches. The saturated branch is eliminated at 9/8. In the slack branch integer radii (1,2) make the asymptotic trigonometric polynomial and tail elementary, leaving only (c,w1,w2) for a bounded Oct scan.",
		PDCertificate:            artifact.PositiveDefiniteness, Witness: witness, Bound: bound,
		CertifiedCBracket:   "573/500 <= c_* < 1651/1250 (the upper endpoint is strict)",
		CertifiedJBracket:   "849/1250 < J_* <= 427/500 (the lower endpoint is strict)",
		AnthropicComparison: "0.68185 remains inside the certified interval (0.6792,0.854]; it is neither confirmed nor contradicted, and the unnamed Anthropic class is still not sourced as identical to EP3.1.",
		WitnessArtifact:     "compiler/artifacts/m16_two_radius_witness.json",
		Experiment: M16Experiment{
			SearchPath: "experiments/m16_two_radius_completion.octest", ArtifactPath: "experiments/m16_two_radius_plot.octest",
			Commands:         []string{"go run ./cmd/oct test .../experiments/m16_two_radius_completion.octest --execution interpreted --json", "go run ./cmd/oct test .../experiments/m16_two_radius_completion.octest --execution compiled --json", "go run ./cmd/oct artifact .../experiments/m16_two_radius_plot.octest --output-root .../artifacts --execution compiled --json"},
			Strategy:         "local filters first; eliminate saturation; choose integer radii 1 and 2; scan the reduced rational neighborhood; adversarially probe t in [0,200] and a dense local mesh",
			BestNumerical:    "c=1.147, r1=1, r2=2, w1=0.168, w2=0.001, sampled minimum about 0.00018545; rationalized conservatively to certified c=573/500=1.146, sampled minimum about 0.00148322 at t about 1.5422",
			Counterexamples:  []string{"every saturated c>9/8 candidate fails at the origin", "c=1.148 with the retained rational weights develops a negative pocket near t=1.54", "coincident radii add no representational degree"},
			IndependentCheck: "Oct independently reconstructs the transform, exact-coefficient formulas in binary64, legality filters, the local minimum, and the long-range scan; it is corroboration only.",
			Plot:             "artifacts/m16/m16_two_radius_spectrum.png", PlotSHA256: "942e1e969e550d391b6fe44e62e319eb7233216efc8299fb7edd508fef7ffc46", PlotStatus: "produced on first run and unchanged on deterministic rerun; 11,782 bytes; identity M16TwoRadiusPlot.PlotCandidateAndRegionHandoff:m16/m16_two_radius_spectrum.png", ExecutionIdentity: "Oct gooct-cli; 6 compiled Facts with zero fallback and build-time-interpreted artifact phase",
		},
		LiteratureAssessment:     "A focused post-certification check of arXiv:1810.08843, arXiv:2310.01913, and arXiv:2502.05106 found the relevant SDP/EP extremal formulations and the 1.3208 primal scale, but no matching two-atom completion or 573/500 dual value. This is a limited novelty screen, not a priority claim.",
		ImportedMathematics:      []string{"EP3.1 weak dual and plus-sign data measure", "CGdL strict primal certificate c_*<1651/1250", "Bochner-Schwartz criterion", "Taylor's theorem with bounded trigonometric derivatives"},
		DerivedMathematics:       []string{"correct unsaturated two-radius family", "exact origin expansion through t^6", "9/8 ceiling for the saturated two-radius branch", "rational (1,2)-radius witness at c=573/500", "exact compact enclosure and coercive tail", "certified w2=0 degeneration showing positive origin slack, not radius count, is decisive"},
		ArchitecturalAwkwardness: "The compact proof is a deterministic family-specific rational enclosure over 4,001 centers. It is auditable but intentionally not generalized into interval-arithmetic or nonlinear-optimization infrastructure.",
		CompilerTheory:           "The broader family certifies new information beyond 9/8, but the answer to whether the additional radius itself is essential is no: the exact verifier also accepts the w2=0 degeneration. M15's obstruction generalizes only to origin-saturated nonnegative exterior-atom families. The useful new representational state is explicit positive origin slack plus a strictly positive tail trigonometric polynomial.",
		NextMilestone:            "M17: characterize and optimize the unsaturated one-radius boundary-atom family analytically, correcting the representation boundary exposed by M16 before adding further radii.",
		RHStatus:                 "unresolved",
	}
	if err := validateM16Result(r); err != nil {
		return M16Result{}, err
	}
	return r, nil
}

func validateM16Result(r M16Result) error {
	if r.Outcome == "" || r.RHStatus != "unresolved" || !r.Witness.CertifiesInfiniteClass() {
		return fmt.Errorf("incomplete M16 certified result")
	}
	if r.Bound.MultiplicityLower != (semantic.ExactRational{Numerator: 573, Denominator: 500}) || r.Bound.SimpleUpper != (semantic.ExactRational{Numerator: 427, Denominator: 500}) {
		return fmt.Errorf("M16 bound propagation failed")
	}
	if !strings.Contains(r.CertifiedCBracket, "573/500") || !strings.Contains(r.CertifiedJBracket, "427/500") {
		return fmt.Errorf("M16 bracket is incomplete")
	}
	return r.Family.Validate()
}
