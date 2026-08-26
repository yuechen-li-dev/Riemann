package compiler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

//go:embed artifacts/m22_contact_preserving_two_atom_witness.json
var m22WitnessJSON []byte

type M22WitnessArtifact struct {
	Schema             string                                  `json:"schema"`
	Family             semantic.ContactPreservingTwoAtomFamily `json:"family"`
	WholeLine          semantic.TwoAtomWholeLineCertificate    `json:"whole_line"`
	CExact             string                                  `json:"c_exact"`
	JExact             string                                  `json:"J_exact"`
	Verified           bool                                    `json:"verified"`
	NumericalCandidate semantic.RadiusCandidate                `json:"numerical_radius_candidate"`
	ContactMotion      semantic.ContactMotionEvidence          `json:"contact_motion"`
	Provenance         string                                  `json:"provenance"`
}

type M22Experiment struct {
	Path              string   `json:"path"`
	SearchTool        string   `json:"search_tool"`
	Interpreted       string   `json:"interpreted"`
	Compiled          string   `json:"compiled"`
	FallbackCount     int      `json:"fallback_count"`
	Artifacts         []string `json:"artifacts"`
	ArtifactSHA256    []string `json:"artifact_sha256"`
	ExecutionIdentity string   `json:"execution_identity"`
	EvidenceBoundary  string   `json:"evidence_boundary"`
}

type M22Result struct {
	M21                        M21Result                               `json:"m21"`
	Outcome                    string                                  `json:"outcome"`
	Family                     semantic.ContactPreservingTwoAtomFamily `json:"certified_family"`
	ExactFamily                string                                  `json:"exact_two_atom_family"`
	ContactFormulas            []string                                `json:"contact_preserving_formulas"`
	M21Regression              string                                  `json:"m21_structural_regression"`
	SpecialRadii               []string                                `json:"special_radii"`
	LegalRegion                string                                  `json:"legal_region"`
	SingularBranches           string                                  `json:"singular_branches"`
	ConstraintCount            []string                                `json:"constraint_count"`
	Sensitivity                semantic.ContactMotionEvidence          `json:"first_order_sensitivity"`
	SecondContactMotion        string                                  `json:"second_contact_motion"`
	ActiveContactSystem        []string                                `json:"active_contact_system"`
	Envelope                   string                                  `json:"nested_envelope"`
	RadiusEnvelope             string                                  `json:"radius_envelope"`
	NumericalCandidate         semantic.RadiusCandidate                `json:"numerical_radius_candidate"`
	NumericalContactFrequency  float64                                 `json:"numerical_contact_frequency"`
	CertifiedFamilyBracket     string                                  `json:"certified_family_bracket"`
	ContactTrack               []string                                `json:"contact_track"`
	ThirdConstraint            string                                  `json:"third_near_active_constraint"`
	CertifiedWitness           semantic.CertifiedTwoAtomWitness        `json:"certified_witness"`
	CertifiedC                 string                                  `json:"certified_c"`
	CertifiedJ                 string                                  `json:"certified_J"`
	DirectedDecimals           string                                  `json:"directed_decimals"`
	OriginalContact            string                                  `json:"original_contact_certificate"`
	NewContact                 string                                  `json:"new_contact_certificate"`
	CompactProof               string                                  `json:"compact_proof"`
	TailProof                  string                                  `json:"tail_proof"`
	CertifiedBounds            string                                  `json:"updated_bounds"`
	M21Comparison              string                                  `json:"m21_comparison"`
	PrimalComparison           string                                  `json:"compatible_primal_comparison"`
	AnthropicComparison        string                                  `json:"anthropic_0_68185"`
	MultiContactInterpretation string                                  `json:"multi_contact_interpretation"`
	Experiment                 M22Experiment                           `json:"oct_experiment"`
	WitnessArtifact            string                                  `json:"witness_artifact"`
	LiteratureAssessment       string                                  `json:"literature_assessment"`
	ImportedMathematics        []string                                `json:"imported_mathematics"`
	DerivedMathematics         []string                                `json:"compiler_research_derived_mathematics"`
	ArchitecturalAwkwardness   string                                  `json:"architectural_awkwardness"`
	CompilerTheory             string                                  `json:"compiler_theory"`
	NextMilestone              string                                  `json:"one_next_milestone"`
	RHStatus                   string                                  `json:"rh_status"`
}

func CompileM22() (M22Result, error) {
	m21, err := CompileM21()
	if err != nil {
		return M22Result{}, err
	}
	return compileM22FromM21(m21, false)
}

func compileM22FromM21(m21 M21Result, replayCompact bool) (M22Result, error) {
	var artifact M22WitnessArtifact
	if err := json.Unmarshal(m22WitnessJSON, &artifact); err != nil {
		return M22Result{}, fmt.Errorf("decode M22 witness artifact: %w", err)
	}
	if artifact.Schema != "riemann.m22.contact-preserving-two-atom.v1" {
		return M22Result{}, fmt.Errorf("wrong M22 witness schema")
	}
	derived, err := semantic.DeriveContactPreservingTwoAtomFamily(artifact.Family.Mass, artifact.Family.Radius)
	if err != nil || derived != artifact.Family {
		return M22Result{}, fmt.Errorf("M22 family artifact is not derived: %v", err)
	}
	witness := semantic.CertifiedTwoAtomWitness{Family: derived, WholeLine: artifact.WholeLine, CExact: artifact.CExact, JExact: artifact.JExact, Verified: artifact.Verified}
	if err := semantic.VerifyCertifiedTwoAtomWitness(witness, replayCompact); err != nil {
		return M22Result{}, err
	}
	motion := semantic.DeriveM21ContactMotionEvidence(m21.NumericalCandidate.Mass, m21.NumericalCandidate.Frequency)
	if math.Abs(motion.DeltaFAtOldContact-artifact.ContactMotion.DeltaFAtOldContact) > 2e-12 || motion.DeltaFAtOldContact >= 0 || artifact.NumericalCandidate.C <= m21.NumericalCandidate.C {
		return M22Result{}, fmt.Errorf("M22 contact-motion or numerical-improvement evidence mismatch")
	}
	r := M22Result{
		M21: m21, Outcome: "Success A — releasing the second radius yields a stronger certified two-atom witness",
		Family:                    derived,
		ExactFamily:               "F_{m,r}=F0+m*H_r, H_r(t)=-pi*cos(pi*r/2)*sin(t)/t+2*((2/pi)*cos(pi*r/2)-r*sin(pi*r/2))*cos(t)+2*cos(r*t)",
		ContactFormulas:           []string{"c(m,r)=c0+(pi/2)*m*cos(pi*r/2)", "w1(m,r)=w0+m*((2/pi)*cos(pi*r/2)-r*sin(pi*r/2))", "w2=m"},
		M21Regression:             "At r=4, cos(2pi)=1 and sin(2pi)=0, so c=c0+(pi/2)m, w1=w0+(2/pi)m, and H_r is exactly the M21 response.",
		SpecialRadii:              []string{"r=2: cos(pi)=-1, sin(pi)=0; positive mass decreases c", "r=3: cos(3pi/2)=0, sin(3pi/2)=-1; c=c0 and w1=w0+3m", "r=4: the exact M21 subfamily", "r=5: cos(5pi/2)=0, sin(5pi/2)=1; c=c0 and w1=w0-5m", "nearby nonintegers with cos(pi*r/2)>0 can improve c; the local optimum is nonintegral"},
		LegalRegion:               "r>1, m>=0, c0+(pi/2)m*cos(pi*r/2)>0, w0+m*((2/pi)cos(pi*r/2)-r sin(pi*r/2))>=0, plus c<1651/1250 before theorem propagation; tail feasibility is 3/4-(c-1)/2-2w1-2m>0 for the reused |t|>=4 proof.",
		SingularBranches:          semantic.ContactPreservationSingularRadii(),
		ConstraintCount:           []string{"4 original scalars (c,w1,m,r)", "2 built-in contact equations determine (c,w1)", "2-dimensional contact-preserving surface (m,r)", "a movable double contact adds F(t1)=0 and F_t(t1)=0 with t1 free, leaving one local degree", "radius-envelope stationarity supplies the third finite equation at an interior optimum"},
		Sensitivity:               motion,
		SecondContactMotion:       "At fixed M21 mass, partial_r F(t1)<0, so rho<0 raises the old t1 contact. A frequency shift alone cannot preserve F=0 to first order because F_t(t1)=0; coupled dm/dr=-m H_r/H restores the value, then dt1/dr=-(H_t dm/dr+m H_tr)/F_tt moves the double contact.",
		ActiveContactSystem:       []string{"Q1=F0(t1)+m*H_r(t1)=0", "Q2=F0'(t1)+m*partial_t H_r(t1)=0", "Q3=C'(r)*H_r(t1)-C(r)*partial_r H_r(t1)=0, C(r)=cos(pi*r/2)"},
		Envelope:                  "For fixed r and H_r(t)<0, M_r(t)=F0(t)/(-H_r(t)); subject to parameter legality, m_*(r)=inf_t M_r(t). Interior minimizers satisfy Q1=Q2=0.",
		RadiusEnvelope:            "C(r)=c0+(pi/2)*cos(pi*r/2)*m_*(r). The envelope theorem gives m_*'(r)=-m_* partial_r H_r/H_r and C'(r)=0 exactly when Q3=0.",
		NumericalCandidate:        artifact.NumericalCandidate,
		NumericalContactFrequency: artifact.NumericalCandidate.ContactFrequency,
		CertifiedFamilyBracket:    "not claimed in M22: Success A certifies the explicit r=7/2 witness; the optimizer remains numerical inside the searched local radius window [3,5]",
		ContactTrack:              []string{"r=4: t1 approximately 1.36474957", "r=3.875: t1 approximately 1.46138524", "r=3.75: t1 approximately 1.56609793", "r=3.5: t1 approximately 1.80339636", "the branch moves smoothly away from pi/2 as r decreases"},
		ThirdConstraint:           "No third near-active minimum appears in the bounded best-spectrum scan: the first remote minima near 9.44 and 15.54 are about 0.571 and 0.597.",
		CertifiedWitness:          witness, CertifiedC: artifact.CExact, CertifiedJ: artifact.JExact,
		DirectedDecimals:           "c_safe approximately 1.15633129202474 and J_safe_upper approximately 0.843668707975264",
		OriginalContact:            "F(pi/2)=F'(pi/2)=0 exactly; degree-100 curvature replay on [154/100,160/100] proves F''>1/1000, while direct numerical curvature is about 0.0129312.",
		NewContact:                 "The certified witness is deliberately below the r=7/2 envelope ceiling: its released minimum near t=1.8034 is strictly positive (about 1.53e-4) and lies in strict rational Taylor cells.",
		CompactProof:               "Degree-100 directed rational Taylor evaluation plus |F'|<19/20 covers [0,154/100] and [160/100,4]; the only equality cell is the built-in contact neighborhood.",
		TailProof:                  "For |t|>=4, |2(cos t-1)/t^2|<=1/4, |sin t/t|<=1/4, and both atom phases have magnitude at most one, hence F>=3/4-(c-1)/2-2w1-2m>1/4. The bound is independent of r.",
		CertifiedBounds:            "c_*>=1+pi/4-2/pi+17*pi*sqrt(2)/10000; J_*<=1-pi/4+2/pi-17*pi*sqrt(2)/10000",
		M21Comparison:              "The new certified c is about 1.15633129, exceeding M21's certified 1.15474742 and even its numerical ray ceiling 1.15486093.",
		PrimalComparison:           "The certified dual remains strictly below the compatible CGdL primal guard c_*<1651/1250=1.3208.",
		AnthropicComparison:        "The resulting J upper bound is about 0.843669, still far above the unexplained 0.68185 remark; no identification is claimed.",
		MultiContactInterpretation: "The numerical variable-radius optimum retains the fixed pi/2 double contact and one movable double contact. No emerging third active constraint was found, so the active-set picture is strengthened only to a smooth two-contact envelope.",
		Experiment:                 M22Experiment{Path: "experiments/m22_contact_preserving_two_atom.octest", SearchTool: "cmd/m22search", Interpreted: "5/5 passed in 5.029 s", Compiled: "5/5 passed in 0.624 s", FallbackCount: 0, Artifacts: []string{"artifacts/m22/m22_radius_response.png", "artifacts/m22/m22_contact_motion.png", "artifacts/m22/m22_best_spectrum.png"}, ArtifactSHA256: []string{"8cdcd23c4606ce248398f0eb9b00f6b5b4619325677aef093379430c91d49f35", "044871b4287c8e62dca49a6ca61cddcc2984a04222d6cd7f42da8bbac330c6d0", "364a850d58f57f4b8ccc680a53435a5db657a2946ec6eea7651385f5a52ebfcb"}, ExecutionIdentity: "Oct dev; gooct-cli; CPU Float; compiled Facts use zero fallback; compiled artifact spelling delegates to build-time-interpreted publication; deterministic replay returned 3 unchanged", EvidenceBoundary: "Oct and Go Float searches/plots are research evidence only; the exact sqrt(2)/pi intervals, rational Taylor cells, contact convexity, and analytic tail carry the theorem."},
		WitnessArtifact:            "compiler/artifacts/m22_contact_preserving_two_atom_witness.json",
		LiteratureAssessment:       "A post-certification primary/preprint screen of arXiv:1810.08843 and arXiv:2502.05106, plus targeted atomic/Dirac/variable-radius searches, found SDP Gaussian-polynomial and reproducing-kernel/Fourier-optimization frameworks but no stated two-atom completion, variable-radius contact-preserving family, 7/2 witness, or matching three-equation active system. This is a limited screen, not a novelty or priority claim.",
		ImportedMathematics:        []string{"M18 exact one-radius witness and equality-aware contact method", "M21 contact-adapted ray and degree-100 rational Taylor infrastructure", "EP3.1 weak duality and compatible CGdL strict primal guard", "Taylor's theorem and classical rational pi bounds"},
		DerivedMathematics:         []string{"exact variable-radius contact-preserving surface", "absence of singular radius branches", "nested scalar envelope and three-equation interior optimizer", "first-order release and contact-motion equations", "nonintegral numerical optimum near r=3.50785", "whole-line algebraic witness r=7/2, m=17/2500", "radius-uniform tail"},
		ArchitecturalAwkwardness:   "The existing PiExpression IR cannot represent sqrt(2), so M22 keeps the general family symbolic and adds a narrowly scoped Pell-interval verifier for the sole certified half-integer radius instead of widening the scalar language globally.",
		CompilerTheory:             "One released geometric parameter compiles an apparent obstruction into a nested envelope: preserve known equalities structurally, let the next contact determine mass and frequency, and use stationarity to select geometry. NumericalRadiusCandidate and CertifiedTwoAtomWitness remain separate evidence types.",
		NextMilestone:              "M23: certify a tight box for the three-equation variable-radius optimizer and determine whether its two-contact configuration is the local two-atom optimum; do not add another atom unless that certification exposes a new active constraint.", RHStatus: "unresolved",
	}
	if err := validateM22Result(r); err != nil {
		return M22Result{}, err
	}
	return r, nil
}

func validateM22Result(r M22Result) error {
	if !strings.Contains(r.Outcome, "Success A") || r.RHStatus != "unresolved" || r.Family.Radius != (semantic.ExactRational{Numerator: 7, Denominator: 2}) || r.Family.Mass != (semantic.ExactRational{Numerator: 17, Denominator: 2500}) || !r.CertifiedWitness.Verified {
		return fmt.Errorf("incomplete M22 result")
	}
	if len(r.ActiveContactSystem) != 3 || len(r.Experiment.Artifacts) != 3 || r.NextMilestone == "" || !strings.Contains(r.CertifiedBounds, "sqrt(2)") || r.NumericalContactFrequency <= 0 || r.CertifiedFamilyBracket == "" {
		return fmt.Errorf("M22 report propagation is incomplete")
	}
	return nil
}

func M22JSONReport(result M22Result) ([]byte, error) {
	if err := validateM22Result(result); err != nil {
		return nil, err
	}
	payload := struct {
		Schema string    `json:"schema"`
		Result M22Result `json:"m22"`
	}{"riemann.m22.v1", result}
	return json.MarshalIndent(payload, "", "  ")
}
