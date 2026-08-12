package semantic

import (
	"fmt"
	"math/big"
	"strings"
)

// SignCutoff is the exact support-one semantic used by the broader LP.  It is
// deliberately distinct from compact Fourier support: values outside the data
// interval may exist, but must have the sign that lets the unknown pair-
// correlation contribution be discarded.
type FourierTailSemantics string

const (
	ExactFourierSupport FourierTailSemantics = "exact_fourier_support"
	NonpositiveTail     FourierTailSemantics = "nonpositive_fourier_tail"
)

type SupportOneExtremalClass struct {
	ID                  string               `json:"id"`
	DataRadius          ExactRational        `json:"data_radius"`
	TailSemantics       FourierTailSemantics `json:"tail_semantics"`
	TransformConvention string               `json:"transform_convention"`
	Even                bool                 `json:"even"`
	Continuous          bool                 `json:"continuous"`
	FunctionL1          bool                 `json:"function_l1"`
	TransformL1         bool                 `json:"transform_l1"`
	FunctionNonnegative bool                 `json:"function_nonnegative"`
	Normalization       string               `json:"normalization"`
	ObjectiveID         string               `json:"objective_id"`
	Provenance          []Reference          `json:"provenance"`
}

func (c SupportOneExtremalClass) Validate() error {
	if c.ID == "" || c.TransformConvention == "" || c.Normalization == "" || c.ObjectiveID == "" || len(c.Provenance) == 0 {
		return fmt.Errorf("incomplete support-one extremal class")
	}
	if err := c.DataRadius.Validate(); err != nil {
		return err
	}
	if rat(c.DataRadius).Cmp(big.NewRat(1, 1)) != 0 {
		return fmt.Errorf("M14A class must use the authoritative radius one")
	}
	if c.TailSemantics != NonpositiveTail || !c.Even || !c.Continuous || !c.FunctionL1 || !c.TransformL1 || !c.FunctionNonnegative {
		return fmt.Errorf("M14A class is missing an EP3.1 admissibility condition")
	}
	return nil
}

type ExtremalCandidateMembership struct {
	CandidateID         string               `json:"candidate_id"`
	Even                bool                 `json:"even"`
	Continuous          bool                 `json:"continuous"`
	FunctionL1          bool                 `json:"function_l1"`
	TransformL1         bool                 `json:"transform_l1"`
	FunctionNonnegative bool                 `json:"function_nonnegative"`
	TailSemantics       FourierTailSemantics `json:"tail_semantics"`
	TailStartsAt        ExactRational        `json:"tail_starts_at"`
	ValueAtZero         ExactRational        `json:"value_at_zero"`
	Normalization       string               `json:"normalization"`
	Proof               string               `json:"proof"`
}

func (m ExtremalCandidateMembership) ValidateFor(c SupportOneExtremalClass) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if m.CandidateID == "" || m.Normalization != c.Normalization || strings.TrimSpace(m.Proof) == "" {
		return fmt.Errorf("candidate normalization or proof is missing")
	}
	if !m.Even || !m.Continuous || !m.FunctionL1 || !m.TransformL1 || !m.FunctionNonnegative {
		return fmt.Errorf("candidate lacks a class admissibility condition")
	}
	if err := m.TailStartsAt.Validate(); err != nil {
		return err
	}
	if rat(m.TailStartsAt).Cmp(rat(c.DataRadius)) > 0 {
		return fmt.Errorf("candidate sign control starts after the support-one boundary")
	}
	if m.TailSemantics != NonpositiveTail && m.TailSemantics != ExactFourierSupport {
		return fmt.Errorf("candidate has no legal Fourier-tail theorem")
	}
	if err := m.ValueAtZero.Validate(); err != nil {
		return err
	}
	if rat(m.ValueAtZero).Cmp(big.NewRat(1, 1)) != 0 {
		return fmt.Errorf("candidate must be normalized by g(0)=1")
	}
	return nil
}

type ExtremalObjective struct {
	ID                string      `json:"id"`
	MultiplicityRatio string      `json:"multiplicity_ratio"`
	SimpleProportion  string      `json:"simple_proportion"`
	Homogeneous       bool        `json:"homogeneous"`
	DataMeasure       string      `json:"data_measure"`
	PipelineContract  string      `json:"pipeline_contract"`
	Provenance        []Reference `json:"provenance"`
}

func (o ExtremalObjective) Validate() error {
	if o.ID == "" || o.MultiplicityRatio == "" || o.SimpleProportion == "" || !o.Homogeneous || o.DataMeasure == "" || o.PipelineContract == "" || len(o.Provenance) == 0 {
		return fmt.Errorf("incomplete M14A extremal objective")
	}
	return nil
}

type PositivityEvidence string

const (
	GridPositivityOnly PositivityEvidence = "grid_only"
	GlobalAnalyticPD   PositivityEvidence = "global_analytic_positive_definiteness"
)

type DualCompletionWitness struct {
	ID                     string             `json:"id"`
	MultiplicityLower      ExactRational      `json:"multiplicity_lower"`
	OutsideMeasure         string             `json:"outside_nonnegative_measure"`
	CompletionDistribution string             `json:"completion_distribution"`
	FourierImage           string             `json:"fourier_image"`
	PositivityEvidence     PositivityEvidence `json:"positivity_evidence"`
	ExactSupportScope      bool               `json:"exact_support_scope"`
	WholeLineControl       bool               `json:"whole_line_control"`
	TailControl            bool               `json:"tail_control"`
	Provenance             Reference          `json:"provenance"`
}

func (w DualCompletionWitness) Validate() error {
	if w.ID == "" || w.OutsideMeasure == "" || w.CompletionDistribution == "" || w.FourierImage == "" || w.Provenance.Citation == "" {
		return fmt.Errorf("incomplete dual completion witness")
	}
	if err := w.MultiplicityLower.Validate(); err != nil {
		return err
	}
	if rat(w.MultiplicityLower).Sign() <= 0 {
		return fmt.Errorf("dual lower constant must be positive")
	}
	return nil
}

func (w DualCompletionWitness) CertifiesInfiniteClass() bool {
	return w.Validate() == nil && w.PositivityEvidence == GlobalAnalyticPD && w.ExactSupportScope && w.WholeLineControl && w.TailControl
}

type CertifiedExtremalBound struct {
	MultiplicityLower ExactRational `json:"multiplicity_lower"`
	SimpleUpper       ExactRational `json:"simple_upper"`
	WitnessID         string        `json:"witness_id"`
}

// ApplyDualCompletion is the narrow weak-duality verifier.  Numerical samples
// never enter this route.
func ApplyDualCompletion(c SupportOneExtremalClass, w DualCompletionWitness) (CertifiedExtremalBound, error) {
	if err := c.Validate(); err != nil {
		return CertifiedExtremalBound{}, err
	}
	if !w.CertifiesInfiniteClass() {
		return CertifiedExtremalBound{}, fmt.Errorf("dual witness lacks global positivity, exact scope, or tail control")
	}
	upper := new(big.Rat).Sub(big.NewRat(2, 1), rat(w.MultiplicityLower))
	u, err := exactRat(upper)
	if err != nil {
		return CertifiedExtremalBound{}, err
	}
	return CertifiedExtremalBound{MultiplicityLower: w.MultiplicityLower, SimpleUpper: u, WitnessID: w.ID}, nil
}
