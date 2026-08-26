package compiler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

type M20Candidate struct {
	ID                  string                         `json:"id"`
	Family              string                         `json:"family"`
	NewScalarParameters int                            `json:"new_scalar_parameters"`
	StructuralResult    string                         `json:"structural_result"`
	BestNumerical       semantic.M20Probe              `json:"best_numerical_probe"`
	CertificationRoute  string                         `json:"whole_line_certification_route"`
	FailureMode         string                         `json:"failure_mode"`
	Utility             semantic.RepresentationUtility `json:"utility"`
	UtilityScore        int                            `json:"utility_score"`
	Selected            bool                           `json:"selected"`
}

type M20Experiment struct {
	Path              string   `json:"path"`
	Interpreted       string   `json:"interpreted"`
	Compiled          string   `json:"compiled"`
	FallbackCount     int      `json:"fallback_count"`
	Domain            string   `json:"bounded_domain"`
	PressureFindings  []string `json:"pressure_findings"`
	Artifacts         []string `json:"artifacts"`
	ArtifactSHA256    []string `json:"artifact_sha256"`
	ArtifactExecution string   `json:"artifact_execution"`
	EvidenceBoundary  string   `json:"evidence_boundary"`
}

type M20Result struct {
	M18                      M18Result      `json:"m18_baseline"`
	Outcome                  string         `json:"outcome"`
	Baseline                 string         `json:"exact_m18_baseline"`
	FirstVariation           string         `json:"first_variation"`
	ContactSensitivity       string         `json:"contact_sensitivity"`
	ActiveConstraints        []string       `json:"active_contact_constraints"`
	MissingCapability        string         `json:"missing_capability"`
	Candidates               []M20Candidate `json:"candidate_families"`
	SelectedFamily           string         `json:"selected_family"`
	SelectionReason          string         `json:"selection_reason"`
	RejectedFamilies         []string       `json:"rejected_families"`
	CertifiedImprovement     string         `json:"certified_improvement"`
	TheoremStack             string         `json:"theorem_stack"`
	CompatiblePrimalGuard    string         `json:"compatible_primal_guard"`
	FullClassRelation        string         `json:"full_ep3_1_relation"`
	AnthropicComparison      string         `json:"anthropic_0_68185"`
	Experiment               M20Experiment  `json:"oct_experiment"`
	ImportedMathematics      []string       `json:"imported_mathematics"`
	DerivedMathematics       []string       `json:"compiler_research_derived_mathematics"`
	ArchitecturalAwkwardness string         `json:"architectural_awkwardness"`
	CompilerTheory           string         `json:"compiler_theory"`
	NextMilestone            string         `json:"one_next_milestone"`
	RHStatus                 string         `json:"rh_status"`
}

func CompileM20() (M20Result, error) {
	m18, err := CompileM18()
	if err != nil {
		return M20Result{}, err
	}
	return compileM20FromM18(m18)
}

func compileM20FromM18(m18 M18Result) (M20Result, error) {
	utilities := []semantic.RepresentationUtility{
		{Family: "A: generic two active radii", Improvement: 5, Tractability: 4, WholeLineCost: 4, ExactPotential: 4, Reuse: 4, Implementation: 2, Clarity: 3},
		{Family: "B: boundary atom plus uniform interval density", Improvement: 5, Tractability: 4, WholeLineCost: 4, ExactPotential: 3, Reuse: 5, Implementation: 3, Clarity: 3},
		{Family: "C: three-site lattice packet proxy", Improvement: 2, Tractability: 3, WholeLineCost: 3, ExactPotential: 4, Reuse: 2, Implementation: 3, Clarity: 4},
		{Family: "D: contact-adapted r=4 atomic ray", Improvement: 5, Tractability: 5, WholeLineCost: 4, ExactPotential: 5, Reuse: 4, Implementation: 5, Clarity: 5},
	}
	probes := []semantic.M20Probe{
		{Family: utilities[0].Family, CandidateC: 1.15474741707169, SampledMinimum: 3.96284314084158e-7, MinimumAt: 1.58, Evidence: semantic.NumericalEvidence},
		{Family: utilities[1].Family, CandidateC: 1.15474825257476, SampledMinimum: 3.90545516454344e-7, MinimumAt: 1.58, Evidence: semantic.NumericalEvidence},
		{Family: utilities[2].Family, CandidateC: 1.15087278613226, SampledMinimum: 3.73253628625694e-6, MinimumAt: 1.58, Evidence: semantic.NumericalEvidence},
		{Family: utilities[3].Family, CandidateC: 1.15474741707169, SampledMinimum: 3.96284314084158e-7, MinimumAt: 1.58, Evidence: semantic.NumericalEvidence},
	}
	for _, p := range probes {
		if err := p.Validate(); err != nil {
			return M20Result{}, err
		}
	}
	candidates := []M20Candidate{
		{ID: "A", Family: utilities[0].Family, NewScalarParameters: 4, StructuralResult: "The M18 point is the w2=0 boundary. A positive atom at a phase with cos(pi*r/2)>0 permits gamma>0; a boundary-weight retuning restores tangency. The coarse genuine-two-radius scan selected r2=4.", BestNumerical: probes[0], CertificationRoute: "Fix rational/integer radii, isolate the contact analytically, use exact Taylor cells on a compact interval and a trigonometric-amplitude tail.", FailureMode: "Free r2 and two weights buy three unnecessary search directions; the useful direction collapses to the contact-adapted subfamily.", Utility: utilities[0], UtilityScore: utilities[0].Score()},
		{ID: "B", Family: utilities[1].Family, NewScalarParameters: 3, StructuralResult: "A uniform positive-side density on [15/4,17/4], mirrored across zero, has response 2m sinc(t/4)cos(4t). It raises c and its response decays in the tail.", BestNumerical: probes[1], CertificationRoute: "Exact sinc/trigonometric Taylor core plus a decaying O(1/t) density tail and the inherited coercive baseline.", FailureMode: "The sampled gain is only 8.4e-7 above A at this resolution, while sinc derivatives add proof algebra and obscure exact contact curvature.", Utility: utilities[1], UtilityScore: utilities[1].Score()},
		{ID: "C", Family: utilities[2].Family, NewScalarParameters: 2, StructuralResult: "The finite equal-weight packet at radii 3,4,5 is a legal bounded proxy and has a positive contact direction. A true infinite periodic measure would produce a Fourier distribution, not the pointwise density required by the current verifier.", BestNumerical: probes[2], CertificationRoute: "Finite packet: positive trigonometric-polynomial tail plus compact Taylor cover. Infinite lattice: new distributional positivity semantics would be required.", FailureMode: "Only the r=4 site helps the contact value; the r=3 and r=5 mass mostly pays for oscillation, so improvement per total mass is poor.", Utility: utilities[2], UtilityScore: utilities[2].Score()},
		{ID: "D", Family: utilities[3].Family, NewScalarParameters: 1, StructuralResult: "Set r2=4, w2=m, c=c0+(pi/2)m, and w1=w0+(2/pi)m. These identities preserve F(pi/2)=F'(pi/2)=0 exactly for every m. The contact curvature decreases with slope 2-16/pi^2-32, predicting a first obstruction at m about 0.0040768 and c about 1.1551822.", BestNumerical: probes[3], CertificationRoute: "One exact Laurent-in-pi parameter, analytic quadratic-contact core, inherited rational Taylor subdivision away from contact, and elementary atomic tail. A later endpoint can test quartic contact.", FailureMode: "The phase-aligned atom repairs the active contact but spends curvature rapidly; the next active constraint is contact flattening, not a remote secondary minimum.", Utility: utilities[3], UtilityScore: utilities[3].Score(), Selected: true},
	}
	r := M20Result{
		M18:                m18,
		Outcome:            "Success A — selected the contact-adapted r=4 atomic ray as the smallest informative extension",
		Baseline:           "c0=1+pi/4-2/pi, w0=1/pi-4/pi^2+8/pi^3, r0=1, t0=pi/2, with exact double contact and F''(t0)>1/8",
		FirstVariation:     "For an even legal exterior perturbation with Fourier response H and c->c+epsilon*gamma, w->w+epsilon*beta, delta F(t)=-2*gamma*sin(t)/t+2*beta*cos(t)+H(t). At a=pi/2, value preservation gives gamma=a H(a)/2; derivative preservation gives beta=gamma/a^2+H'(a)/2.",
		ContactSensitivity: "Two scalar active equations must be controlled to raise c without opening a negative pocket: delta F(a)=0 and delta F'(a)=0. Once they hold, the contact moves only at higher order; local safety is governed first by delta F''(a). Without derivative retuning, the predicted shift is -delta F'(a)/F0''(a).",
		ActiveConstraints:  []string{"contact value fixes the admissible first-order c increase from H(pi/2)", "contact derivative consumes one independent retuning degree", "positive contact curvature is the remaining local margin", "the origin has about 0.0445 slack and is not active", "secondary minima start near t=9.31 with margin about 0.61; the certified tail margin is >1/3"},
		MissingCapability:  "M18 lacks a Fourier response that is positive at t=pi/2: its only atom has cos(pi/2)=0, so changing its weight cannot repair the active value constraint. The missing capability is one phase-aligned exterior mode, plus use of the existing boundary weight to preserve the derivative constraint.",
		Candidates:         candidates, SelectedFamily: utilities[3].Family,
		SelectionReason:          "`when utility` selected D with score 52 versus B=45, A=44, C=31. D matches B's sampled improvement scale, removes A's unused dimensions, preserves exact Laurent-in-pi structure, and has the shortest credible whole-line proof. Utility prioritizes research only.",
		RejectedFamilies:         []string{"A generic two-radius family: promising, but its useful point is exactly D and the free-radius search is avoidable.", "B uniform interval density: genuinely promising and slightly best on the coarse scan, but the tiny sampled edge does not justify sinc proof complexity before the exact atomic ray is understood.", "C periodic/lattice completion: the finite proxy is inefficient; the infinite version leaves the current pointwise-transform semantics."},
		CertifiedImprovement:     "None. Every M20 improvement is Structural or Numerical; the certified theorem bound remains c_*>=1+pi/4-2/pi.",
		TheoremStack:             "M0-M18 theorem semantics and certificates are unchanged. M20 adds research-only perturbation, probe, evidence, and utility records.",
		CompatiblePrimalGuard:    "Every probe enforces c<1651/1250=1.3208; none approaches the guard. The selected sampled c is about 1.1547474.",
		FullClassRelation:        "The selected atom is a legal EP3.1 exterior measure and its contact calculus generalizes to arbitrary even exterior responses H. It is therefore a reusable local basis direction, not a claim that two atoms solve the full class.",
		AnthropicComparison:      "J approximately 0.84525 for the selected numerical probe remains far above 0.68185. The Anthropic number is context only and did not enter optimization or utility scoring.",
		Experiment:               M20Experiment{Path: "experiments/m20_representation_selection.octest", Interpreted: "5/5 passed in 23.094 s", Compiled: "5/5 passed", FallbackCount: 0, Domain: "coarse family probes on 0<=t<=80 with step 1/50; baseline local-minimum pressure scan extended independently to t=60 at step 1/1000", PressureFindings: []string{"exact contact t=pi/2 is the unique active sampled region", "next local minima occur near 9.31, 15.65, 21.95 with values about 0.61, 0.64, 0.65", "tail oscillations approach the asymptotic floor 1-2w0 about 0.6579", "for D, contact curvature rather than a remote pocket is the first predicted obstruction"}, Artifacts: []string{"artifacts/m20/m20_baseline_pressure_map.png", "artifacts/m20/m20_selected_perturbation_effect.png"}, ArtifactSHA256: []string{"5014c88583683c84b8e41f859dea87ba8d3306b7af59c561f0066aa9e7b1099a", "04e35e80ecea7cf0c6d745f8cdbf0f5bc0e9b0e4b92a9756c80ff376153b5a71"}, ArtifactExecution: "Oct dev gooct-cli, build-time-interpreted [Artifact], two produced outputs with complete metadata", EvidenceBoundary: "Scans, plots, curvature forecasts, and utility scores are research evidence only; no sampled value enters theorem propagation."},
		ImportedMathematics:      []string{"M18 exact one-radius extremal and equality-aware whole-line certificate", "EP3.1 support-one/sign-cutoff dual semantics", "CGdL strict compatible primal bound c_*<1651/1250", "elementary Fourier responses of atoms and uniform intervals"},
		DerivedMathematics:       []string{"general first-variation equations at the M18 double contact", "contact-shift and curvature sensitivity", "exact phase-aligned r=4 one-parameter ray", "positive directions for atomic, uniform-density, and finite-lattice candidates", "secondary-minimum pressure map", "explicit representation utility comparison"},
		ArchitecturalAwkwardness: "M18's exact Laurent-in-pi scalar remains sufficient for D, but a future density witness would need exact sinc expressions. M20 deliberately records Float probes behind a Numerical evidence type rather than widening the theorem scalar language.",
		CompilerTheory:           "Representation selection is constraint compilation: inspect the active equation, synthesize the missing response at that equation, and spend existing slack to satisfy its derivative. A parameter is valuable only when its Fourier image is linearly independent where the certificate is tight.",
		NextMilestone:            "M21: optimize and certify or refute the one-parameter contact-adapted r=4 atomic ray, including its predicted curvature-flattening endpoint and a whole-line equality-aware certificate; do not widen the family unless that ray fails.", RHStatus: "unresolved",
	}
	if err := validateM20Result(r); err != nil {
		return M20Result{}, err
	}
	return r, nil
}

func validateM20Result(r M20Result) error {
	if r.RHStatus != "unresolved" || r.M18.Optimum.ExactValue.Expression != "1+pi/4-2/pi" || !strings.Contains(r.CertifiedImprovement, "None") {
		return fmt.Errorf("M20 changed or obscured the M18 theorem baseline")
	}
	if len(r.Candidates) != 4 || len(r.RejectedFamilies) != 3 || r.NextMilestone == "" {
		return fmt.Errorf("incomplete representation selection")
	}
	selected := 0
	maxScore := -1
	for _, c := range r.Candidates {
		if err := c.BestNumerical.Validate(); err != nil {
			return err
		}
		if c.UtilityScore != c.Utility.Score() {
			return fmt.Errorf("stored utility score drift for %s", c.ID)
		}
		if c.UtilityScore > maxScore {
			maxScore = c.UtilityScore
		}
		if c.Selected {
			selected++
			if c.Family != r.SelectedFamily {
				return fmt.Errorf("selected family mismatch")
			}
		}
	}
	if selected != 1 || r.Candidates[3].UtilityScore != maxScore {
		return fmt.Errorf("utility selection is incomplete or nondeterministic")
	}
	return nil
}

func M20JSONReport(result M20Result) ([]byte, error) {
	if err := validateM20Result(result); err != nil {
		return nil, err
	}
	payload := struct {
		Schema string    `json:"schema"`
		Result M20Result `json:"m20"`
	}{Schema: "riemann.m20.v1", Result: result}
	return json.MarshalIndent(payload, "", "  ")
}
