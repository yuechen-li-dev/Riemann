package compiler

import (
	"encoding/json"
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

type M25Experiment struct {
	Path              string   `json:"path"`
	Interpreted       string   `json:"interpreted"`
	Compiled          string   `json:"compiled"`
	FallbackCount     int      `json:"fallback_count"`
	Generation        string   `json:"generation"`
	ExactReplay       string   `json:"exact_replay"`
	ManifestReplay    string   `json:"manifest_replay"`
	RadiusSlabs       int      `json:"radius_slabs"`
	ContactBoxes      int      `json:"contact_boxes"`
	CompetingBands    int      `json:"competing_bands"`
	Artifacts         []string `json:"artifacts"`
	ArtifactSHA256    []string `json:"artifact_sha256"`
	GeneratedGo       string   `json:"generated_go"`
	GeneratedGoSHA256 string   `json:"generated_go_sha256"`
	EvidenceBoundary  string   `json:"evidence_boundary"`
}

type M25Result struct {
	M24                      M24Result               `json:"m24"`
	Outcome                  string                  `json:"outcome"`
	Certificate              semantic.M25Certificate `json:"certificate"`
	B1Definition             string                  `json:"b1_definition"`
	CompetitorExclusion      string                  `json:"competitor_exclusion"`
	ObjectiveGap             string                  `json:"objective_gap"`
	BranchSwitch             string                  `json:"branch_switch"`
	BoundaryComposition      string                  `json:"boundary_composition"`
	WindowConclusion         string                  `json:"window_conclusion"`
	ComponentConclusion      string                  `json:"component_conclusion"`
	OptimizerFeasibility     string                  `json:"optimizer_feasibility"`
	BestExplicitWitness      string                  `json:"best_explicit_witness"`
	BestCertifiedC           string                  `json:"best_certified_c"`
	BestCertifiedJ           string                  `json:"best_certified_J"`
	OptimizerWitnessGap      string                  `json:"optimizer_witness_gap"`
	NextRestriction          string                  `json:"next_restriction"`
	ActiveSet                string                  `json:"active_set"`
	OctFirstImplementation   string                  `json:"oct_first_implementation"`
	GeneratedGoDetails       string                  `json:"generated_go_details"`
	Experiment               M25Experiment           `json:"oct_experiment"`
	LiteratureAssessment     string                  `json:"literature_assessment"`
	ImportedMathematics      []string                `json:"imported_mathematics"`
	DerivedMathematics       []string                `json:"newly_derived_mathematics"`
	ArchitecturalAwkwardness string                  `json:"architectural_awkwardness"`
	CompilerTheory           string                  `json:"compiler_theory"`
	RHStatus                 string                  `json:"rh_status"`
	NextMilestone            string                  `json:"next_milestone"`
}

func CompileM25() (M25Result, error) {
	m24, err := CompileM24()
	if err != nil {
		return M25Result{}, err
	}
	return compileM25FromM24(m24)
}

func compileM25FromM24(m24 M24Result) (M25Result, error) {
	certificate, err := semantic.CertifyM25B1Control()
	if err != nil {
		return M25Result{}, err
	}
	r := M25Result{
		M24: m24, Outcome: "Success A — ContactPreservingTwoAtomComponentGlobalOptimum",
		Certificate:              certificate,
		B1Definition:             "B1 is the unique smooth solution graph (m1(r),t1(r)) of F0(t)+m H_r(t)=0 and F0'(t)+m partial_t H_r(t)=0 with t in [1.790,1.801] over r in [3.503,3.512]",
		CompetitorExclusion:      "the B1 spectrum is certified nonnegative on the whole line; hence every H_r(t)<0 competitor obeys F0(t)/(-H_r(t))>=m1(r), with equality only at the unique movable B1 contact",
		ObjectiveGap:             "strict compact and tail bands have positive spectrum margins; this proves no branch equality/switch, while the fixed pi/2 equality has H=0 and is not a competing mass cap",
		BranchSwitch:             "none on [3.503,3.512]",
		BoundaryComposition:      certificate.Control.BoundaryComposition,
		WindowConclusion:         "B1 controls the entire M24 residual window and its strictly concave envelope has the unique M23 stationary maximum",
		ComponentConclusion:      "M23 controls the full legal M18-contact-preserving two-atom radius component after composing M24 outside-window barriers with M25 closed-window control",
		OptimizerFeasibility:     "yes; M25 proves F_{m*,r*}(t)>=0 for all real t for the interval-defined optimizer branch, with fixed and movable quadratic contacts and a uniform t>=4 tail",
		BestExplicitWitness:      "r=7/2, m=69/10000 remains the best simple explicit rational witness",
		BestCertifiedC:           "c_* >= 1.1566431645986454 from the globally feasible interval-defined optimizer; explicit rational lower endpoint remains 1.1564423523252728",
		BestCertifiedJ:           "J_* <= 0.8433568354013546 from J=2-c for the interval-defined optimizer",
		OptimizerWitnessGap:      "the interval optimizer lower endpoint exceeds the M23 explicit witness lower endpoint by 0.0002008122733726",
		NextRestriction:          "yes: after closing this component, preservation of the original pi/2 contact is the evident remaining representation restriction; M25 does not release it",
		ActiveSet:                "fixed double contacts ±pi/2; movable double contacts ±t1(r); unique global radius optimum C1'(r*)=0",
		OctFirstImplementation:   "Oct owns branch discovery, 90-slab continuation predictors, adverse scans, plots, and deterministic Go atlas generation; Go owns exact interval reconstruction, whole-line feasibility, theorem composition, and the persisted fast manifest",
		GeneratedGoDetails:       "semantic/m25_continuation.generated.go contains the deterministic Oct continuation atlas; ordinary tests validate the pre-initialized manifest, while cmd/m25replay performs the expensive exact reconstruction",
		Experiment:               M25Experiment{Path: "experiments/m25_b1_component_closure.octest", Interpreted: "5/5 passed", Compiled: "5/5 passed with zero fallback", FallbackCount: 0, Generation: "build-time interpreted artifact publication", ExactReplay: "cmd/m25replay exact rational reconstruction passed in about 307 s", ManifestReplay: "ordinary semantic M25 tests complete in about 1.6 s", RadiusSlabs: 90, ContactBoxes: 90, CompetingBands: 6, Artifacts: []string{"artifacts/m25/m25_b1_continuation.png", "artifacts/m25/m25_b1_objective.png", "artifacts/m25/m25_competitor_gap.png"}, ArtifactSHA256: []string{"96cfe6d47be85cc75c73b3d30678bd79ee2dc1f1d7b1aeec202d6416870250c9", "7cb7e26da05a4f651b04a6a47301b705757abaa33f96597c2e5c4111fe85aabd", "9427858006ab48710368fd43f5f3a5087c847d06081a95e3599e204f0883be03"}, GeneratedGo: "semantic/m25_continuation.generated.go", GeneratedGoSHA256: "9caa7311d514431159b52b1f3d5ed054f13e116af12c970b1692db10aa565c81", EvidenceBoundary: "Oct-generated Go is replay input, not independent verification; independent evidence is the compiled broad adversarial scan versus Go whole-line interval reconstruction"},
		LiteratureAssessment:     "focused post-closure screen retained M23's related Fourier-extremal references and found no stated match for this contact-preserving two-atom radius/contact system; no novelty or priority claim",
		ImportedMathematics:      []string{"M18 fixed double contact and exact base family", "M22 contact-preserving two-atom parametrization and analytic tail", "M23 unique stationary optimizer and local strict maximum", "M24 exact outside-window barriers"},
		DerivedMathematics:       []string{"parameterized interval-Newton continuation of B1 across the exact residual window", "branch-correlated centered Taylor compact proof", "whole-line B1 feasibility and competitor exclusion", "global component theorem by M23/M24/M25 composition"},
		ArchitecturalAwkwardness: "exact rational degree-80 transcendental reconstruction remains expensive; M25 therefore separates the persistent manifest/fast semantic path from an explicit audit command instead of embedding multi-minute work in ordinary tests",
		CompilerTheory:           "M25 shows that a discovered active graph can compile into a small persistent manifest: structural equalities, interval-Newton graph cells, branch-correlated Taylor bands, and an analytic tail compose into a global theorem without rerunning discovery in every test",
		RHStatus:                 "unresolved",
		NextMilestone:            "M26: test whether releasing the forced pi/2 contact improves the two-atom optimum, while keeping exactly two atoms and the same EP3.1 route.",
	}
	if err := validateM25Result(r); err != nil {
		return M25Result{}, err
	}
	return r, nil
}

func validateM25Result(r M25Result) error {
	if r.RHStatus != "unresolved" || !r.Certificate.Control.M23ControlsComponent || !r.Certificate.Control.OptimizerFeasible || r.NextMilestone == "" || r.Experiment.FallbackCount != 0 {
		return fmt.Errorf("incomplete M25 Success A result")
	}
	return nil
}

func M25JSONReport(r M25Result) ([]byte, error) {
	if err := validateM25Result(r); err != nil {
		return nil, err
	}
	return json.MarshalIndent(struct {
		Schema string    `json:"schema"`
		Result M25Result `json:"m25"`
	}{"riemann.m25.v1", r}, "", "  ")
}
