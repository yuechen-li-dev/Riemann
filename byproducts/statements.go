// Package byproducts inventories external mathematical statements. Its reports
// have no path into the compiler's trusted proof graph.
package byproducts

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	LeanToolchain   = "leanprover/lean4:v4.33.1"
	MathlibRevision = "0df444a360eaa60ab8c11dca51a86af692955474"
	Repository      = "https://github.com/Deicyde/riemann-hypothesis"
	Revision        = "a172c464fd4bcfc4f53e394dcba45d742b1409a6"
	ArchiveSHA256   = "2757f699bc40beff25a67722696f55be68009e00ce6b27ad6fdc95a79df38d75"
)

type Span struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	First  int    `json:"first_line"`
	Last   int    `json:"last_line"`
	SHA256 string `json:"sha256"`
}

func (s Span) Link() string {
	if s.First == s.Last {
		return fmt.Sprintf("%s/blob/%s/%s#L%d", Repository, Revision, s.Path, s.First)
	}
	return fmt.Sprintf("%s/blob/%s/%s#L%d-L%d", Repository, Revision, s.Path, s.First, s.Last)
}

type Fidelity string
type Novelty string
type Verification string
type Computation string
type Independence string
type HypothesisStatus string

const (
	Matches     Fidelity = "matches"
	SpecialCase Fidelity = "special_case"
	Stronger    Fidelity = "stronger"
	Unresolved  Fidelity = "unresolved"

	KnownMath             Novelty = "known_math"
	ExistingFormalization Novelty = "already_in_mathlib_or_other_project"
	NoveltyUnresolved     Novelty = "unresolved"

	StructurallyReviewed    Verification     = "structurally_reviewed"
	KernelReplayed          Verification     = "kernel_replayed"
	NoRHInReviewedStatement Independence     = "no_rh_in_reviewed_statement"
	Conditional             HypothesisStatus = "conditional"
	Unconditional           HypothesisStatus = "unconditional"
	NotRun                  Computation      = "not_run"
	FiniteReplay            Computation      = "exact_finite_replay"
)

type Comparison struct {
	Kind     string   `json:"kind"` // literature or formalization
	Title    string   `json:"title"`
	URL      string   `json:"url"`
	Fidelity Fidelity `json:"fidelity"`
	Finding  string   `json:"finding"`
}

// Review is human-authored and bound to source excerpts by their digests.
// Dependencies listed here are selected source uses, not elaborated constants.
type Review struct {
	Slug              string           `json:"slug"`
	Title             string           `json:"title"`
	Domain            string           `json:"domain"`
	Declaration       Span             `json:"declaration"`
	Definitions       []Span           `json:"definitions"`
	Context           []Span           `json:"context"`
	ProofUses         []string         `json:"proof_uses"`
	Assumptions       []string         `json:"assumptions"`
	Quantifiers       []string         `json:"quantifiers"`
	Interpretation    string           `json:"interpretation"`
	HumanReview       string           `json:"human_review"`
	HypothesisStatus  HypothesisStatus `json:"hypothesis_status"`
	Independence      Independence     `json:"rh_independence"`
	ExtraProofPremise bool             `json:"extra_proof_premise"`
	MathNovelty       Novelty          `json:"mathematical_novelty"`
	FormalNovelty     Novelty          `json:"formalization_novelty"`
	Verification      Verification     `json:"verification"`
	Comparisons       []Comparison     `json:"comparisons"`
	Open              []string         `json:"open_questions"`
}

type ExternalStatement struct {
	ID             string   `json:"id"`
	Repository     string   `json:"repository"`
	Revision       string   `json:"revision"`
	Statement      string   `json:"lean_statement"`
	SourceFidelity Fidelity `json:"source_span_fidelity"`
	Review
	Imports           []string            `json:"proof_module_imports"`
	DefinitionImports map[string][]string `json:"definition_module_imports"`
	Computation       Computation         `json:"independent_computation"`
	ComputationDetail string              `json:"computation_detail,omitempty"`
	RankingFeatures   []string            `json:"ranking_features"`
	Score             int                 `json:"oct_score"`
}

type Catalog struct {
	Revision   string   `json:"revision"`
	Candidates []Review `json:"candidates"`
}

func ReadCatalog(r io.Reader) (Catalog, error) {
	var catalog Catalog
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&catalog); err != nil {
		return catalog, err
	}
	if err := dec.Decode(new(json.RawMessage)); err != io.EOF {
		return catalog, fmt.Errorf("expected exactly one catalog")
	}
	if catalog.Revision != Revision {
		return catalog, fmt.Errorf("catalog revision %q, require %s", catalog.Revision, Revision)
	}
	if len(catalog.Candidates) == 0 || len(catalog.Candidates) > 32 {
		return catalog, fmt.Errorf("catalog must have 1–32 candidates")
	}
	return catalog, nil
}

var slugPattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
var rhMention = regexp.MustCompile(`(?i)\b(?:[a-z]*RiemannHypothesis[a-z]*|DirichletGRH|GRH|RH)\b`)

func validateReview(r Review, reviewedSource string) error {
	if !slugPattern.MatchString(r.Slug) || r.Title == "" || r.Domain == "" || len(r.Open) == 0 ||
		r.HumanReview == "" || len(r.Quantifiers) == 0 || r.Interpretation == "" {
		return fmt.Errorf("incomplete review")
	}
	if r.Verification != StructurallyReviewed {
		return fmt.Errorf("B1 cannot issue %s without a separate formal replay workflow", r.Verification)
	}
	switch r.HypothesisStatus {
	case Conditional:
		if len(r.Assumptions) == 0 {
			return fmt.Errorf("conditional review must list hypotheses")
		}
	case Unconditional:
		// Named propositions and section parameters can hide premises.
		// This pilot has only reviewed statements with ordinary hypotheses.
		return fmt.Errorf("B1 cannot establish an unconditional classification from a source review")
	default:
		return fmt.Errorf("missing hypothesis classification")
	}
	if r.Independence != NoRHInReviewedStatement || rhMention.MatchString(reviewedSource) {
		return fmt.Errorf("RH independence is unsupported by the reviewed statement")
	}
	literature, formal, formalMatch := false, false, false
	for _, c := range r.Comparisons {
		if c.Title == "" || !strings.HasPrefix(c.URL, "https://") || c.Finding == "" {
			return fmt.Errorf("comparison requires a link and a finding")
		}
		switch c.Fidelity {
		case Matches, SpecialCase, Stronger, Unresolved, "weaker", "incomparable":
		default:
			return fmt.Errorf("unknown fidelity %q", c.Fidelity)
		}
		switch c.Kind {
		case "literature":
			literature = true
		case "formalization":
			formal = true
			formalMatch = formalMatch || c.Fidelity == Matches
		default:
			return fmt.Errorf("unknown comparison kind %q", c.Kind)
		}
	}
	if r.MathNovelty != KnownMath || !literature {
		return fmt.Errorf("B1 requires a literature comparison for known mathematics")
	}
	switch r.FormalNovelty {
	case ExistingFormalization:
		if !formal {
			return fmt.Errorf("known formalization needs a comparison")
		}
	case NoveltyUnresolved:
		if formalMatch {
			return fmt.Errorf("matched formalization must be recorded")
		}
	default:
		return fmt.Errorf("B1 cannot issue positive formalization novelty claims")
	}
	return nil
}
