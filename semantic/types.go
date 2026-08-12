// Package semantic defines the deliberately small mathematical vocabulary used
// by the M0 compiler. It is not intended to be a general expression language.
package semantic

import (
	"fmt"
	"sort"
	"strings"
)

type ClaimID string
type TransformationID string
type AssumptionID string

// Function is a typed identity for a mathematical function understood by M0.
type Function uint8

const RiemannZeta Function = 1

func (f Function) String() string {
	if f == RiemannZeta {
		return "Riemann zeta function"
	}
	return fmt.Sprintf("unknown function (%d)", f)
}

// Proposition is sealed to this package so compiler passes cannot smuggle
// semantics into arbitrary strings.
type Proposition interface {
	Kind() PropositionKind
	Describe() string
	isProposition()
}

type PropositionKind string

const (
	RiemannHypothesisKind                PropositionKind = "riemann_hypothesis"
	AllNontrivialZerosOnCriticalLineKind PropositionKind = "all_nontrivial_zeros_on_critical_line"
	CriticalLineDensityOneKind           PropositionKind = "critical_line_density_one"
	NamedObligationKind                  PropositionKind = "named_obligation"
)

type RiemannHypothesis struct{ Function Function }

func (RiemannHypothesis) isProposition()        {}
func (RiemannHypothesis) Kind() PropositionKind { return RiemannHypothesisKind }
func (p RiemannHypothesis) Describe() string {
	return "Riemann Hypothesis for the " + p.Function.String()
}

type AllNontrivialZerosOnCriticalLine struct{ Function Function }

func (AllNontrivialZerosOnCriticalLine) isProposition() {}
func (AllNontrivialZerosOnCriticalLine) Kind() PropositionKind {
	return AllNontrivialZerosOnCriticalLineKind
}
func (p AllNontrivialZerosOnCriticalLine) Describe() string {
	return "every nontrivial zero ρ of the " + p.Function.String() + " has Re(ρ) = 1/2"
}

type CriticalLineDensityOne struct{ Function Function }

func (CriticalLineDensityOne) isProposition()        {}
func (CriticalLineDensityOne) Kind() PropositionKind { return CriticalLineDensityOneKind }
func (p CriticalLineDensityOne) Describe() string {
	return "the asymptotic density of nontrivial zeros of the " + p.Function.String() + " on Re(s) = 1/2 is 1"
}

// NamedObligation is intentionally only an escape hatch for proof obligations,
// not a general string-based mathematical proposition.
type NamedObligation struct{ Name string }

func (NamedObligation) isProposition()        {}
func (NamedObligation) Kind() PropositionKind { return NamedObligationKind }
func (p NamedObligation) Describe() string    { return p.Name }

type Property uint8

const ExceptionalSetSensitivity Property = 1

func (p Property) String() string {
	if p == ExceptionalSetSensitivity {
		return "exceptional-set sensitivity"
	}
	return fmt.Sprintf("unknown property (%d)", p)
}

// PropertySet is a compact immutable-ish capability set.
type PropertySet uint64

func Properties(properties ...Property) PropertySet {
	var set PropertySet
	for _, property := range properties {
		set |= 1 << property
	}
	return set
}

func (s PropertySet) Has(property Property) bool            { return s&(1<<property) != 0 }
func (s PropertySet) Contains(other PropertySet) bool       { return s&other == other }
func (s PropertySet) Without(other PropertySet) PropertySet { return s &^ other }
func (s PropertySet) Names() []string {
	var names []string
	for property := Property(1); property < 64; property++ {
		if s.Has(property) {
			names = append(names, property.String())
		}
	}
	sort.Strings(names)
	return names
}

type Exactness string

const (
	Exact       Exactness = "exact"
	Approximate Exactness = "approximate"
)

type Assumption struct {
	ID          AssumptionID `json:"id"`
	Description string       `json:"description"`
}

type EvidenceKind string

const (
	DefinitionEvidence           EvidenceKind = "definition"
	KnownTheoremEvidence         EvidenceKind = "known_theorem"
	DerivedEvidence              EvidenceKind = "derived"
	NumericalExperimentEvidence  EvidenceKind = "numerical_experiment"
	UnverifiedConjectureEvidence EvidenceKind = "unverified_conjecture"
)

type ReferenceKind string

const (
	StandardReference ReferenceKind = "standard_reference"
	CompilerRecord    ReferenceKind = "compiler_record"
	ExperimentRecord  ReferenceKind = "experiment_record"
)

type Reference struct {
	Kind     ReferenceKind `json:"kind"`
	Citation string        `json:"citation"`
	URI      string        `json:"uri,omitempty"`
}

type Evidence struct {
	Kind   EvidenceKind `json:"kind"`
	Source Reference    `json:"source"`
	Note   string       `json:"note,omitempty"`
}

type ProvenanceKind string

const (
	AuthoredProvenance ProvenanceKind = "authored"
	DerivedProvenance  ProvenanceKind = "derived"
)

type Provenance struct {
	Kind           ProvenanceKind   `json:"kind"`
	Parents        []ClaimID        `json:"parents"`
	Transformation TransformationID `json:"transformation,omitempty"`
	Source         Reference        `json:"source"`
}

type Claim struct {
	ID           ClaimID
	Proposition  Proposition
	Assumptions  []Assumption
	Evidence     []Evidence
	Requirements PropertySet
	Capabilities PropertySet
	Exactness    Exactness
	Provenance   Provenance
}

func (c Claim) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("claim ID is empty")
	}
	if c.Proposition == nil {
		return fmt.Errorf("claim %q has no proposition", c.ID)
	}
	if c.Exactness != Exact && c.Exactness != Approximate {
		return fmt.Errorf("claim %q has invalid exactness %q", c.ID, c.Exactness)
	}
	seen := make(map[AssumptionID]bool, len(c.Assumptions))
	for _, assumption := range c.Assumptions {
		if assumption.ID == "" || strings.TrimSpace(assumption.Description) == "" {
			return fmt.Errorf("claim %q has invalid assumption", c.ID)
		}
		if seen[assumption.ID] {
			return fmt.Errorf("claim %q repeats assumption %q", c.ID, assumption.ID)
		}
		seen[assumption.ID] = true
	}
	return nil
}

func CloneAssumptions(in []Assumption) []Assumption {
	return append([]Assumption(nil), in...)
}
