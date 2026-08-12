package compiler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

//go:embed artifacts/m15_boundary_atom_witness.json
var m15WitnessJSON []byte

type M15WitnessArtifact struct {
	Family               semantic.BoundaryAtomCompletion `json:"family"`
	PositiveDefiniteness semantic.WholeLinePDCertificate `json:"positive_definiteness"`
}

type M15Experiment struct {
	SearchPath        string   `json:"search_path"`
	ArtifactPath      string   `json:"artifact_path"`
	Commands          []string `json:"commands"`
	BestNumerical     string   `json:"best_numerical_dual_candidate"`
	Counterexamples   []string `json:"counterexamples"`
	Plot              string   `json:"plot"`
	ExecutionIdentity string   `json:"execution_identity"`
	Evidence          string   `json:"evidence"`
}

type M15OctTooling struct {
	RootCause       string   `json:"root_cause"`
	Semantics       string   `json:"artifact_semantics"`
	ExecutionModes  string   `json:"execution_modes"`
	Before          string   `json:"before"`
	After           string   `json:"after"`
	FilesChanged    []string `json:"files_changed"`
	Tests           []string `json:"tests"`
	ResearchBenefit string   `json:"research_benefit"`
}

type M15Result struct {
	M14A                     M14AResult                      `json:"m14a"`
	Outcome                  string                          `json:"outcome"`
	ObjectiveNotationAudit   string                          `json:"objective_notation_audit"`
	Family                   semantic.BoundaryAtomCompletion `json:"candidate_family"`
	PDCertificate            semantic.WholeLinePDCertificate `json:"whole_line_pd_certificate"`
	Witness                  semantic.DualCompletionWitness  `json:"certified_dual_witness"`
	Bound                    semantic.CertifiedExtremalBound `json:"certified_bound"`
	FamilyCeiling            semantic.ExactRational          `json:"family_ceiling"`
	FamilyCeilingProof       string                          `json:"family_ceiling_proof"`
	CompletionRepresentation string                          `json:"completion_representation"`
	OmittedDirectionControl  string                          `json:"omitted_direction_control"`
	CertifiedCBracket        string                          `json:"certified_c_bracket"`
	CertifiedJBracket        string                          `json:"certified_j_bracket"`
	AnthropicComparison      string                          `json:"anthropic_comparison"`
	PotentialNovelty         string                          `json:"potential_novelty"`
	IndependentVerification  string                          `json:"independent_verification"`
	WitnessArtifact          string                          `json:"witness_artifact"`
	Experiment               M15Experiment                   `json:"oct_experiment"`
	OctTooling               M15OctTooling                   `json:"oct_tooling"`
	ImportedMathematics      []string                        `json:"imported_mathematics"`
	DerivedMathematics       []string                        `json:"compiler_research_derived_mathematics"`
	ArchitecturalAwkwardness string                          `json:"architectural_awkwardness"`
	CompilerTheory           string                          `json:"compiler_theory"`
	NextMilestone            string                          `json:"one_next_milestone"`
	RHStatus                 string                          `json:"rh_status"`
}

func CompileM15() (M15Result, error) {
	m14a, err := CompileM14A()
	if err != nil {
		return M15Result{}, err
	}
	var artifact M15WitnessArtifact
	if err := json.Unmarshal(m15WitnessJSON, &artifact); err != nil {
		return M15Result{}, fmt.Errorf("decode M15 witness artifact: %w", err)
	}
	if err := semantic.VerifyBoundaryAtomCertificate(artifact.Family, artifact.PositiveDefiniteness); err != nil {
		return M15Result{}, fmt.Errorf("verify M15 witness artifact: %w", err)
	}
	witness := semantic.DualCompletionWitness{
		ID: "boundary-atom-nine-eighths", MultiplicityLower: semantic.ExactRational{Numerator: 9, Denominator: 8},
		OutsideMeasure:         "sigma=(9/8) 1_{|x|>1} dx+(1/8)(delta_{-1}+delta_1)",
		CompletionDistribution: "P=delta_0+(|x|-9/8)1_{[-1,1]}dx+(1/8)(delta_{-1}+delta_1)",
		FourierImage:           artifact.PositiveDefiniteness.FourierDensity,
		PositivityEvidence:     semantic.GlobalAnalyticPD, ExactSupportScope: true, WholeLineControl: true, TailControl: true,
		Provenance: semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Riemann M15 boundary-atom whole-line alternating-series certificate"},
	}
	bound, err := semantic.ApplyDualCompletion(m14a.Class, witness)
	if err != nil {
		return M15Result{}, err
	}
	result := M15Result{
		M14A: m14a, Outcome: "Success A — first nontrivial certified dual improvement",
		ObjectiveNotationAudit: "The authoritative EP3.1/M14A objective used here has + integral_{-1}^{1}|alpha| hat g(alpha)dalpha, equivalently nu=delta_0+|alpha|1_{[-1,1]}dalpha. The minus sign displayed once in the M15 request conflicts with its stated dual, CGdL mapping, target scale, and repository contract; no minus-objective theorem is claimed.",
		Family:                 artifact.Family, PDCertificate: artifact.PositiveDefiniteness, Witness: witness, Bound: bound,
		FamilyCeiling:            semantic.ExactRational{Numerator: 9, Denominator: 8},
		FamilyCeilingProof:       "hat(P_c)(t/(2*pi))=((9-8c)/12)t^2+O(t^4); nonnegativity near zero forces c<=9/8. The endpoint is globally nonnegative, so 9/8 is the exact optimum of this family.",
		CompletionRepresentation: "For t=2*pi*xi, 4t^2 hat(P)(xi)=G(t)=t^2(4+cos t)-t sin t-8(1-cos t). Since hat(P)>=0 on R, Bochner-Schwartz gives positive definiteness of P.",
		OmittedDirectionControl:  artifact.PositiveDefiniteness.OmittedDirections,
		CertifiedCBracket:        "9/8 <= c_* < 1651/1250 (the upper endpoint is strict)",
		CertifiedJBracket:        "849/1250 < J_* <= 7/8 (the lower endpoint is strict)",
		AnthropicComparison:      "0.68185 lies inside the certified interval (0.6792,0.875]; the remark remains neither confirmed nor contradicted for EP3.1, and its class identity remains unsourced.",
		PotentialNovelty:         "Potentially new as a purpose-specific explicit EP3.1 dual bound; a focused post-certification search found no exact 9/8 completion in the audited sources. This is not a priority or exhaustive novelty claim.",
		IndependentVerification:  "Go consumes and exactly checks the rational proof object; Oct independently reconstructs the Fourier density, scans the endpoint, and finds negative directions immediately above 9/8. The Oct scan is corroboration, not theorem evidence.",
		WitnessArtifact:          "compiler/artifacts/m15_boundary_atom_witness.json",
		Experiment: M15Experiment{
			SearchPath: "experiments/m15_boundary_atom_completion.octest", ArtifactPath: "experiments/m15_boundary_atom_plot.octest",
			Commands:      []string{"go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m15_boundary_atom_completion.octest --execution interpreted --json", "go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m15_boundary_atom_completion.octest --execution compiled --json", "go run ./cmd/oct artifact C:/Users/yuech/source/repos/Riemann/experiments/m15_boundary_atom_plot.octest --output-root C:/Users/yuech/source/repos/Riemann/artifacts --execution compiled --json"},
			BestNumerical: "c=9/8=1.125 in the boundary-atom family", Counterexamples: []string{"every c>9/8 fails the exact origin quadratic condition", "c=1.126 has a negative near-origin numerical direction"},
			Plot: "artifacts/m15/m15_boundary_atom_spectrum.png", ExecutionIdentity: "Oct gooct-cli; build-time-interpreted artifact phase", Evidence: "numerical research visualization and independent formula reconstruction; not certificate provenance",
		},
		OctTooling: M15OctTooling{
			RootCause:       "convenience PlotLine called plotrender directly, bypassing ArtifactWriteCapability.StageArtifactOutput; the PNG was an ambient side effect unknown to the publisher",
			Semantics:       "during [Artifact], PlotLine accepts a relative logical path, stages it below --output-root, rejects absolute/traversing paths, publishes atomically, and reports package/function/kind/execution/path/identity/size/hash metadata",
			ExecutionModes:  "interpreted and --execution compiled both use the explicitly reported build-time interpreter artifact phase; compiled is a compatibility delegation, not native plot lowering",
			Before:          "1 artifact test passed; Outputs: 0 produced, 0 unchanged; requested output root empty",
			After:           "1 artifact test passed; Outputs: 1 produced, 0 unchanged; PNG and deterministic metadata under requested output root",
			FilesChanged:    []string{"internal/interpret/artifact_phase.go", "internal/interpret/interpret.go", "internal/interpret/plot.go", "internal/tester/artifact.go", "cmd/oct/artifact_plot_command_test.go"},
			Tests:           []string{"plot count and output-root placement", "deterministic attribution and two-test noncollision", "ordinary Fact excluded", "failure remains failure and does not publish", "absolute path rejection", "interpreted/compiled-request consistency"},
			ResearchBenefit: "the repaired lane retained the endpoint spectrum used to inspect cancellation near the origin and large-frequency positivity without ad hoc file hunting",
		},
		ImportedMathematics:      []string{"EP3.1 admissible class and plus-sign data measure nu", "CGdL strict primal certificate c_*<1.3208 (J_*>0.6792)", "Bochner-Schwartz positive-definite distribution criterion", "classical c=1 triangle/sinc-squared completion"},
		DerivedMathematics:       []string{"typed boundary-atom completion family", "exact local family ceiling c<=9/8", "endpoint Fourier identity G(t)", "whole-line inner alternating-series and outer coercive proof", "certified c_*>=9/8 and J_*<=7/8"},
		ArchitecturalAwkwardness: "The useful certificate is an elementary two-region inequality, not a generic SDP. Encoding this family-specific proof object is smaller and more auditable than introducing a general distribution solver.",
		CompilerTheory:           "Omitted-direction control can be a finite semantic object when the candidate exposes the right analytic decomposition: a decreasing alternating series covers the compact core and a coercive inequality covers the tail. Numerical search selects the representation; it does not discharge the theorem.",
		NextMilestone:            "M16: add a two-radius symmetric exterior-atom family and seek an exact whole-line decomposition that exceeds c=9/8, retaining origin and coercive-tail filters.",
		RHStatus:                 "unresolved",
	}
	if err := validateM15Result(result); err != nil {
		return M15Result{}, err
	}
	return result, nil
}

func validateM15Result(r M15Result) error {
	if r.Outcome == "" || r.RHStatus != "unresolved" || !r.Witness.CertifiesInfiniteClass() {
		return fmt.Errorf("incomplete M15 certified result")
	}
	if r.Bound.MultiplicityLower != (semantic.ExactRational{Numerator: 9, Denominator: 8}) || r.Bound.SimpleUpper != (semantic.ExactRational{Numerator: 7, Denominator: 8}) {
		return fmt.Errorf("M15 bound propagation failed")
	}
	if !strings.Contains(r.CertifiedCBracket, "9/8") || !strings.Contains(r.CertifiedJBracket, "7/8") || !strings.Contains(r.AnthropicComparison, "neither confirmed nor contradicted") {
		return fmt.Errorf("M15 bracket or comparison is incomplete")
	}
	return semantic.VerifyBoundaryAtomCertificate(r.Family, r.PDCertificate)
}
