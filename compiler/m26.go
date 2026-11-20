package compiler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

type M26Experiment struct {
	Path              string `json:"path"`
	Interpreted       string `json:"interpreted"`
	Compiled          string `json:"compiled"`
	FallbackCount     int    `json:"fallback_count"`
	ExecutionIdentity string `json:"execution_identity"`
	CompilerIdentity  string `json:"compiler_identity"`
	Timing            string `json:"timing"`
	Limits            string `json:"limits"`
	Artifact          string `json:"artifact"`
	ArtifactSHA256    string `json:"artifact_sha256"`
	ExactReplay       string `json:"exact_replay"`
	EvidenceBoundary  string `json:"evidence_boundary"`
}

type M26Result struct {
	M25                  M25Result                                      `json:"m25"`
	Outcome              string                                         `json:"outcome"`
	Family               string                                         `json:"family"`
	Support              string                                         `json:"support"`
	Embedding            semantic.ContactPreservingEmbeddingCertificate `json:"embedding"`
	FirstVariation       string                                         `json:"first_variation"`
	RayBoundaryCandidate string                                         `json:"release_ray_boundary_candidate"`
	Witness              semantic.M26ReleasedWitness                    `json:"certified_witness"`
	ObjectiveComparison  string                                         `json:"objective_comparison"`
	BestCertifiedC       string                                         `json:"best_certified_c"`
	BestCertifiedJ       string                                         `json:"best_certified_J"`
	ActiveSet            string                                         `json:"active_set"`
	CentralAnswer        string                                         `json:"central_answer"`
	TheoremBoundary      string                                         `json:"theorem_boundary"`
	OctFirstWorkflow     string                                         `json:"oct_first_workflow"`
	Experiment           M26Experiment                                  `json:"oct_experiment"`
	ImportedMathematics  []string                                       `json:"imported_mathematics"`
	DerivedMathematics   []string                                       `json:"newly_derived_mathematics"`
	RHStatus             string                                         `json:"rh_status"`
	NextMilestone        string                                         `json:"next_milestone"`
}

func CompileM26() (M26Result, error) {
	m25, err := CompileM25()
	if err != nil {
		return M26Result{}, err
	}
	return compileM26FromM25(m25)
}

func compileM26FromM25(m25 M25Result) (M26Result, error) {
	witness, err := semantic.CertifyM26ReleasedWitness()
	if err != nil {
		return M26Result{}, err
	}
	r := M26Result{
		M25: m25, Outcome: "Success A — ForcedPiOverTwoContactReleased",
		Family:               semantic.UnrestrictedTwoAtomTransform,
		Support:              "exactly two symmetric atomic pairs at +-1 and +-r with r>1; the first radius remains fixed at the support boundary",
		Embedding:            witness.Embedding,
		FirstVariation:       "at the M25 active set, (dc,dw1,dw2,dr)=(1,0,1,0) has D(t)=-2*sin(t)/t+2*cos(r*t); Oct finds D(pi/2)>0.1 and D(t1)>0.8, so both forced contacts lift while c increases",
		RayBoundaryCandidate: "on r=3.50785, m=0.00699 and c=c_CP+epsilon, w1=w1_CP, w2=m+epsilon, Oct finds epsilon_cap≈0.0000946358742382 and the replacement double contact t≈1.523146843527; this is active-set discovery, not theorem state",
		Witness:              witness,
		ObjectiveComparison:  "the certified released witness has c>1.1567274714890000, strictly above the full M25 contact-preserving optimum upper endpoint 1.1566433083347464",
		BestCertifiedC:       "c_* >= 1.1567274714890000",
		BestCertifiedJ:       "J_* <= 0.8432725285110000 from J=2-c",
		ActiveSet:            "the certified epsilon=9/100000 witness is strict on the whole line, so the old pi/2 and t1 contacts are inactive; continuation to the release-ray boundary creates a new double contact near t=1.52315",
		CentralAnswer:        "yes: preserving the original pi/2 contact was an artificial active constraint; releasing only its two algebraic equations strictly improves the legal two-atom dual objective",
		TheoremBoundary:      "M26 proves strict improvement in the unrestricted {+-1,+-r} two-atom family and identifies the next release-ray contact. It does not certify the global optimum of that unrestricted family and does not free the first radius.",
		OctFirstWorkflow:     "Oct checked the embedding numerically, discovered the positive first variation, searched the release ray, located the replacement contact, scanned [0,40], and emitted the seed; Go then replayed a rational interior witness through the authoritative certified semantics",
		Experiment: M26Experiment{
			Path: "experiments/m26_release_forced_contact.octest", Interpreted: "4/4 passed", Compiled: "4/4 passed with zero fallback", FallbackCount: 0,
			ExecutionIdentity: "gooct-cli; CPU Float; artifact build-time-interpreted", CompilerIdentity: "oct dev", Timing: "interpreted 5443 ms; compiled 751 ms; exact Go replay 288.06 s", Limits: "one bounded four-fact file; deterministic [0,40] scan; no external data, GPU, or hosted execution ID",
			Artifact: "artifacts/m26/m26_release_seed.txt", ArtifactSHA256: "3f172eef56eb4ba26937ccd6336530f6f33be3229d8c9b4b7c7b964e8d08f48a", ExactReplay: "cmd/m26replay or go test -tags=slow ./semantic -run M26Slow",
			EvidenceBoundary: "Oct binary64 values are discovery/seed evidence only; Go certifies the rationalized seed, Taylor cells, contact uplift, and analytic tail. This is staged discovery/replay, not two independent implementations",
		},
		ImportedMathematics: []string{"M25 global optimum of the legal pi/2-contact-preserving component", "M22 exact contact-preserving embedding formulas", "EP3.1 weak duality and J=2-c", "Taylor's theorem and elementary atomic tail bounds"},
		DerivedMathematics:  []string{"unrestricted two-atom semantic family with fixed first radius", "mechanical M22-M25 subfamily embedding", "positive release first variation", "replacement-contact discovery", "strict rational whole-line witness exceeding the M25 optimum"},
		RHStatus:            "unresolved", NextMilestone: "determine and certify the controlling active set of the unrestricted fixed-first-radius two-atom family; do not free the first radius until that representation gives an evidence-backed reason",
	}
	if err := validateM26Result(r); err != nil {
		return M26Result{}, err
	}
	return r, nil
}

func validateM26Result(r M26Result) error {
	if !strings.Contains(r.Outcome, "Success A") || !r.Witness.WholeLine || !r.Embedding.Verified || r.Experiment.FallbackCount != 0 || r.RHStatus != "unresolved" || !strings.Contains(r.BestCertifiedJ, "J=2-c") {
		return fmt.Errorf("incomplete M26 release result")
	}
	return nil
}

func M26JSONReport(r M26Result) ([]byte, error) {
	if err := validateM26Result(r); err != nil {
		return nil, err
	}
	return json.MarshalIndent(struct {
		Schema string    `json:"schema"`
		Result M26Result `json:"m26"`
	}{"riemann.m26.v1", r}, "", "  ")
}
