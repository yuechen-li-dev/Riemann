package compiler

import (
	"encoding/json"
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

type M24Experiment struct {
	Path              string   `json:"path"`
	Compiled          string   `json:"compiled"`
	FallbackCount     int      `json:"fallback_count"`
	ArtifactReplay    string   `json:"artifact_replay"`
	Artifacts         []string `json:"artifacts"`
	ArtifactSHA256    []string `json:"artifact_sha256"`
	GeneratedGo       string   `json:"generated_go"`
	GeneratedGoSHA256 string   `json:"generated_go_sha256"`
	ExecutionIdentity string   `json:"execution_identity"`
	EvidenceBoundary  string   `json:"evidence_boundary"`
}

type M24Result struct {
	M23                           M23Result                            `json:"m23"`
	Outcome                       string                               `json:"outcome"`
	Certificate                   semantic.ComponentOptimumCertificate `json:"component_certificate"`
	ComponentBoundaries           string                               `json:"component_boundaries"`
	BranchDiscovery               string                               `json:"branch_discovery"`
	BranchSwitches                string                               `json:"branch_switches"`
	FrequencyCompactification     string                               `json:"frequency_compactification"`
	LargeRadiusBound              string                               `json:"large_radius_bound"`
	M23BranchContinuation         string                               `json:"m23_branch_continuation"`
	DerivativeSignStructure       string                               `json:"derivative_sign_structure"`
	CompetingBranchBounds         string                               `json:"competing_branch_bounds"`
	ComponentConclusion           string                               `json:"component_conclusion"`
	OptimizerFeasibility          string                               `json:"optimizer_feasibility"`
	BestGlobalWitness             string                               `json:"best_global_witness"`
	BestCertifiedC                string                               `json:"best_certified_c"`
	BestCertifiedJ                string                               `json:"best_certified_J"`
	OptimizerWitnessGap           string                               `json:"optimizer_witness_gap"`
	ContactPreservationScope      string                               `json:"contact_preservation_scope"`
	ActiveSetDimension            string                               `json:"active_set_dimension"`
	EquioscillationInterpretation string                               `json:"equioscillation_interpretation"`
	Experiment                    M24Experiment                        `json:"oct_experiment"`
	OctFirstImplementation        string                               `json:"oct_first_implementation"`
	HandwrittenGo                 string                               `json:"handwritten_go"`
	DuplicatedLogic               string                               `json:"duplicated_logic"`
	ReplayPerformance             string                               `json:"replay_performance"`
	Artifacts                     []string                             `json:"artifacts"`
	LiteratureAssessment          string                               `json:"literature_assessment"`
	ImportedMathematics           []string                             `json:"imported_mathematics"`
	DerivedMathematics            []string                             `json:"newly_derived_mathematics"`
	ArchitecturalAwkwardness      string                               `json:"architectural_awkwardness"`
	CompilerTheory                string                               `json:"compiler_theory"`
	RHStatus                      string                               `json:"rh_status"`
	NextMilestone                 string                               `json:"next_milestone"`
}

func CompileM24() (M24Result, error) {
	m23, err := CompileM23()
	if err != nil {
		return M24Result{}, err
	}
	return compileM24FromM23(m23)
}

func compileM24FromM23(m23 M23Result) (M24Result, error) {
	certificate, err := semantic.CertifyM24ComponentControl()
	if err != nil {
		return M24Result{}, err
	}
	result := M24Result{
		M23:                           m23,
		Outcome:                       "Success C — the M23 branch is globally controlling outside one exact unresolved radius/contact window",
		Certificate:                   certificate,
		ComponentBoundaries:           "the legal component is defined exactly as the connected component of the finite/nonnegative-envelope and nonnegative-weight set containing M23; its endpoints and legality-induced disconnections are not certified, and no numerical truncation is substituted",
		BranchDiscovery:               "compiled Oct scans H<0 minima on 1.01<=r<=22.99 and 0<t<=40, refines the two lowest minima, audits integer and half-integer radii, and recovers the M23 branch; scans and plots are research evidence",
		BranchSwitches:                "no controlling switch was sampled near M23; remote oscillatory pockets exist, but a switch/root-count theorem remains missing only on r in [3.503,3.512]",
		FrequencyCompactification:     "not obtained component-wide; the theorem instead uses one adverse frequency per compact radius cell and the structural fixed contact for the entire tail, so no numerical search to infinite frequency enters certification",
		LargeRadiusBound:              "for r>=181/50, fixed-contact curvature gives c<=c0+(pi/4)F0_tt(pi/2)/(r^2-1+8/pi^2)<the M23 lower bound on positive-cosine bands; nonpositive-cosine bands satisfy c<=c0",
		M23BranchContinuation:         "research continuation tracks B1 on [3,3.62] with t about 2.424 down to 1.685; theorem continuation remains the M23 Krawczyk/IFT box only",
		DerivativeSignStructure:       "M23 certifies C'>0 immediately left and C'<0 immediately right of its unique stationary root; the centered M24 envelope derivative controls barrier-cell Taylor errors but does not certify global B1 monotonicity in the unresolved window",
		CompetingBranchBounds:         "all possible branches are eliminated outside [3.503,3.512] by feasibility caps that do not assume which branch is active; near M23 the next sampled mass exceeds B1 by more than 0.01, but root counting remains research-only",
		ComponentConclusion:           "M23 has not yet been proved the ruler of the whole legal component, and no stronger branch was found; only [3.503,3.512] can still affect the comparison",
		OptimizerFeasibility:          "not globally certified on the whole line; M23's explicit rational witness remains the theorem carrier",
		BestGlobalWitness:             "r=7/2, m=69/10000 with the existing M23 whole-line certificate",
		BestCertifiedC:                "c_* >= 1+pi/4-2/pi+69*pi*sqrt(2)/40000, interval lower 1.1564423523252728",
		BestCertifiedJ:                "J_* <= 1-pi/4+2/pi-69*pi*sqrt(2)/40000",
		OptimizerWitnessGap:           "M23 local family candidate minus explicit theorem witness is about 0.00020088; M24 does not close it",
		ContactPreservationScope:      "M24 is global-control work only inside the M18-contact-preserving family; whether releasing the pi/2 contact helps is not answered, and current evidence does not justify a direction",
		ActiveSetDimension:            "family freedoms m,r; structural equations F(pi/2)=F_t(pi/2)=0; movable equations F(t1)=F_t(t1)=0; radius stationarity C'(r)=0. Whole-line positivity is reduced to finite active equalities plus strict-region/tail inequalities",
		EquioscillationInterpretation: "the evidence is consistent with each useful parameter being consumed by another active contact or stationarity equation, but M24 does not promote that pattern to a theorem",
		Experiment: M24Experiment{
			Path:              "experiments/m24_two_atom_component.octest",
			Compiled:          "5/5 passed in about 1.19 s",
			FallbackCount:     0,
			ArtifactReplay:    "relocated publication passed in about 87.0 s; the final deterministic replay took about 91.7 s with all four outputs unchanged",
			Artifacts:         []string{"artifacts/m24/m24_global_radius_envelope.png", "artifacts/m24/m24_active_contact_atlas.png", "artifacts/m24/m24_competitor_mass_gap.png"},
			ArtifactSHA256:    []string{"2b6c36de1cd460c422828021fd991082fc26e810bd172064162216c4480f2d82", "d71d3711429e567447045f381c51499c453afac76b5909013be4ebe0ada5c879", "04733e17e3c25f35a433084a0f0883ddf35afa01dc6291d30c1e26dc8cd9aaa2"},
			GeneratedGo:       "semantic/m24_radius_barriers.generated.go",
			GeneratedGoSHA256: "68592d9c19e79e237c1ce5f7f3f74aa45ec81b3011fe778604005f05abbbf4af",
			ExecutionIdentity: "Oct dev; gooct-cli; CPU Float; compiled Facts with zero fallback; artifacts build-time-interpreted by documented design",
			EvidenceBoundary:  "Oct branch discovery and plots choose candidate witnesses only; Go rational intervals independently replay every generated decimal as an exact rational",
		},
		OctFirstImplementation:   "Oct owns branch scanning, special-radius audits, competitor discovery, plots, contact witness selection, and Go table generation; the stable Go boundary owns only rational interval replay and semantic theorem state",
		HandwrittenGo:            "purpose-specific fixed-contact curvature theorem, scalar rational interval response/derivative evaluator, centered Taylor replay, narrow result types, validation, and reports",
		DuplicatedLogic:          "the family formula appears in Oct for research and in the pre-existing Go certification kernel; M24 adds no second handwritten Go branch finder. The generated table is not claimed as independent evidence; interval replay versus the broad compiled scan is the structural independence",
		ReplayPerformance:        "M23 old compact witness replay remains about 310.9 s and measured 307.3 s. M24 compiled Oct discovery is about 1.16 s, deterministic artifact replay about 91.7 s, semantic component replay about 34.0 s, and the cached compiler M24 suite about 47.4 s",
		Artifacts:                []string{"experiments/m24_two_atom_component.octest", "semantic/m24_radius_barriers.generated.go", "artifacts/m24/m24_global_radius_envelope.png", "artifacts/m24/m24_active_contact_atlas.png", "artifacts/m24/m24_competitor_mass_gap.png", "docs/m24-two-atom-component-control.md"},
		LiteratureAssessment:     "Success C neither solves the component nor finds a stronger branch, so the conditional post-result literature screen was not triggered; M23's limited primary/preprint screen remains the current context and no priority claim is made",
		ImportedMathematics:      []string{"M18 structural fixed double contact", "M22 contact-preserving two-atom family and analytic tail semantics", "M23 Krawczyk root, IFT branch, local C''<0, and explicit whole-line witness", "Taylor's theorem and the necessary nonnegative second derivative at a smooth global minimum"},
		DerivedMathematics:       []string{"fixed-contact identity F_tt(pi/2)=F0_tt(pi/2)-2m cos(pi r/2)(r^2-1+8/pi^2)", "all-radius large-r objective barrier with cutoff 181/50", "purpose-specific envelope derivative for centered radius-cell bounds", "certified exclusion of every legal branch outside r in [3.503,3.512]"},
		ArchitecturalAwkwardness: "Oct compiled-data publication intentionally emits package main; M24 therefore uses readable Artifact.WriteText metaprogramming for a destination-specific semantic table. The rational interval kernel still lacks cached immutable trig point enclosures, making replay slower than its mathematical size suggests",
		CompilerTheory:           "a global continuum question can compile into three layers: structural invariants eliminate infinite tails, generated adverse witnesses eliminate compact cells without naming active branches, and the residual box states the exact missing theorem. Failure to close the last box is useful semantic output, not failed optimization",
		RHStatus:                 "unresolved",
		NextMilestone:            "Certify the B1 contact branch and exclude all competing tangencies on r in [3.503,3.512], then combine that result with M23 stationarity; keep the same two-atom contact-preserving family.",
	}
	if err := validateM24Result(result); err != nil {
		return M24Result{}, err
	}
	return result, nil
}

func validateM24Result(result M24Result) error {
	if result.RHStatus != "unresolved" || result.Certificate.M23ControlsComponent || result.Certificate.MissingTheorem == "" || result.NextMilestone == "" {
		return fmt.Errorf("incomplete M24 Success C result")
	}
	if result.Experiment.FallbackCount != 0 || len(result.Artifacts) == 0 || len(result.DerivedMathematics) == 0 {
		return fmt.Errorf("incomplete M24 implementation evidence")
	}
	return nil
}

func M24JSONReport(result M24Result) ([]byte, error) {
	if err := validateM24Result(result); err != nil {
		return nil, err
	}
	return json.MarshalIndent(struct {
		Schema string    `json:"schema"`
		Result M24Result `json:"m24"`
	}{Schema: "riemann.m24.v1", Result: result}, "", "  ")
}
