// Package semantic defines the deliberately small, typed mathematical
// vocabulary understood by the compiler. It is not a general logic language.
package semantic

import (
	"fmt"
	"strings"
)

type ClaimID string
type TransformationID string
type AssumptionID string
type TheoremID string

// Function identifies a mathematical object independently of its representation.
type Function uint8

const (
	RiemannZeta Function = 1
	RiemannXi   Function = 2
)

func (f Function) String() string {
	if f == RiemannZeta {
		return "Riemann zeta function"
	}
	if f == RiemannXi {
		return "Riemann xi function"
	}
	return fmt.Sprintf("unknown function (%d)", f)
}

type QuantifierKind string

const (
	ForAll     QuantifierKind = "for_all"
	Exists     QuantifierKind = "exists"
	DensityOne QuantifierKind = "density_one"
)

func (q QuantifierKind) Valid() bool {
	return q == ForAll || q == Exists || q == DensityOne
}

func (q QuantifierKind) Describe() string {
	switch q {
	case ForAll:
		return "for every"
	case Exists:
		return "there exists"
	case DensityOne:
		return "for a density-one subset of"
	default:
		return "unknown quantifier"
	}
}

// Domain is a closed, typed fragment of the set vocabulary needed by M1.
// Bound is meaningful only for ZerosBelowHeightDomain.
type Domain struct {
	Kind     DomainKind `json:"kind"`
	Function Function   `json:"function,omitempty"`
	Bound    uint64     `json:"bound,omitempty"`
}

type DomainKind string

const (
	ComplexPlaneDomain     DomainKind = "complex_plane"
	RightHalfPlaneDomain   DomainKind = "half_plane_re_gt_1"
	CriticalStripDomain    DomainKind = "critical_strip"
	NontrivialZerosDomain  DomainKind = "nontrivial_zeros"
	ZerosBelowHeightDomain DomainKind = "zeros_below_height"
	CriticalLineDomain     DomainKind = "critical_line"
	ComplexExceptOneDomain DomainKind = "complex_plane_except_one"
	BoundaryLineDomain     DomainKind = "line_re_eq_1"
)

func ComplexPlane() Domain              { return Domain{Kind: ComplexPlaneDomain} }
func HalfPlaneReGreaterThanOne() Domain { return Domain{Kind: RightHalfPlaneDomain} }
func CriticalStrip() Domain             { return Domain{Kind: CriticalStripDomain} }
func NontrivialZeros(f Function) Domain { return Domain{Kind: NontrivialZerosDomain, Function: f} }
func ZerosBelowHeight(f Function, t uint64) Domain {
	return Domain{Kind: ZerosBelowHeightDomain, Function: f, Bound: t}
}
func CriticalLine() Domain          { return Domain{Kind: CriticalLineDomain} }
func ComplexPlaneExceptOne() Domain { return Domain{Kind: ComplexExceptOneDomain} }
func LineReEqualsOne() Domain       { return Domain{Kind: BoundaryLineDomain} }

func (d Domain) Validate() error {
	switch d.Kind {
	case ComplexPlaneDomain, RightHalfPlaneDomain, CriticalStripDomain, CriticalLineDomain, ComplexExceptOneDomain, BoundaryLineDomain:
		if d.Function != 0 || d.Bound != 0 {
			return fmt.Errorf("domain %s has extraneous parameters", d.Kind)
		}
	case NontrivialZerosDomain:
		if d.Function == 0 || d.Bound != 0 {
			return fmt.Errorf("nontrivial-zero domain has invalid parameters")
		}
	case ZerosBelowHeightDomain:
		if d.Function == 0 || d.Bound == 0 {
			return fmt.Errorf("bounded-zero domain needs a function and positive height")
		}
	default:
		return fmt.Errorf("unknown domain %q", d.Kind)
	}
	return nil
}

func (d Domain) Describe() string {
	switch d.Kind {
	case ComplexPlaneDomain:
		return "the complex plane"
	case RightHalfPlaneDomain:
		return "the half-plane Re(s) > 1"
	case CriticalStripDomain:
		return "the critical strip 0 < Re(s) < 1"
	case NontrivialZerosDomain:
		return "the nontrivial zeros of the " + d.Function.String()
	case ZerosBelowHeightDomain:
		return fmt.Sprintf("the nontrivial zeros of the %s with |Im(ρ)| ≤ %d", d.Function.String(), d.Bound)
	case CriticalLineDomain:
		return "the critical line Re(s) = 1/2"
	case ComplexExceptOneDomain:
		return "the complex plane except s = 1"
	case BoundaryLineDomain:
		return "the line Re(s) = 1"
	default:
		return "an unknown domain"
	}
}

// IsSubset reports only inclusions explicitly known by the M1 domain algebra.
// false means "not established", not a general symbolic proof of non-inclusion.
func IsSubset(sub, sup Domain) bool {
	if sub == sup {
		return true
	}
	if sup.Kind == ComplexPlaneDomain {
		return true
	}
	if sub.Kind == CriticalLineDomain && sup.Kind == CriticalStripDomain {
		return true
	}
	if sub.Kind == ZerosBelowHeightDomain && sup.Kind == NontrivialZerosDomain {
		return sub.Function == sup.Function
	}
	if sub.Kind == ZerosBelowHeightDomain && sup.Kind == ZerosBelowHeightDomain {
		return sub.Function == sup.Function && sub.Bound <= sup.Bound
	}
	return false
}

type PredicateKind string

const (
	RealPartEqualsHalfPredicate      PredicateKind = "real_part_equals_one_half"
	FunctionNonzeroPredicate         PredicateKind = "function_nonzero"
	CriticalReflectionFixedPredicate PredicateKind = "critical_reflection_fixed"
)

type Predicate struct {
	Kind     PredicateKind `json:"kind"`
	Function Function      `json:"function"`
}

func (p Predicate) Validate() error {
	if p.Function == 0 {
		return fmt.Errorf("predicate has no function")
	}
	if p.Kind != RealPartEqualsHalfPredicate && p.Kind != FunctionNonzeroPredicate && p.Kind != CriticalReflectionFixedPredicate {
		return fmt.Errorf("unknown predicate %q", p.Kind)
	}
	return nil
}

func (p Predicate) Describe(variable string) string {
	switch p.Kind {
	case RealPartEqualsHalfPredicate:
		return fmt.Sprintf("Re(%s) = 1/2", variable)
	case FunctionNonzeroPredicate:
		return fmt.Sprintf("ζ(%s) ≠ 0", variable)
	case CriticalReflectionFixedPredicate:
		return fmt.Sprintf("1-conjugate(%s) = %s", variable, variable)
	default:
		return "unknown predicate"
	}
}

// Proposition is sealed so compiler passes cannot smuggle semantics into strings.
type Proposition interface {
	Kind() PropositionKind
	Describe() string
	isProposition()
}

type PropositionKind string

const (
	QuantifiedStatementKind           PropositionKind = "quantified_statement"
	RepresentationKind                PropositionKind = "representation"
	RepresentationIdentityKind        PropositionKind = "representation_identity"
	AnalyticFactKind                  PropositionKind = "analytic_fact"
	NamedObligationKind               PropositionKind = "named_obligation"
	ZeroAtPointKind                   PropositionKind = "zero_at_point"
	SideConditionKind                 PropositionKind = "side_condition"
	FunctionalIdentityKind            PropositionKind = "functional_identity"
	ZeroSetPropertyKind               PropositionKind = "zero_set_property"
	ZeroClassificationKind            PropositionKind = "zero_classification"
	FunctionalDefinitionKind          PropositionKind = "functional_definition"
	UniversalFunctionalStatementKind  PropositionKind = "universal_functional_statement"
	TestFunctionAdmissibilityKind     PropositionKind = "test_function_admissibility"
	ExplicitFormulaIdentityKind       PropositionKind = "explicit_formula_identity"
	FiniteSpanDefinitionKind          PropositionKind = "finite_span_definition"
	QuadraticFormStructureKind        PropositionKind = "quadratic_form_structure"
	HermitianFormDefinitionKind       PropositionKind = "hermitian_form_definition"
	HermitianMatrixDefinitionKind     PropositionKind = "hermitian_matrix_definition"
	FiniteSpanFunctionalStatementKind PropositionKind = "finite_span_functional_statement"
	CoordinatePositivityKind          PropositionKind = "coordinate_quadratic_positivity"
	QuadraticMatrixIdentityKind       PropositionKind = "quadratic_matrix_identity"
	MatrixPropertyStatementKind       PropositionKind = "matrix_property"
	TwoByTwoMinorCertificateKind      PropositionKind = "two_by_two_principal_minor_certificate"
)

type QuantifiedStatement struct {
	Quantifier QuantifierKind `json:"quantifier"`
	Domain     Domain         `json:"domain"`
	Predicate  Predicate      `json:"predicate"`
}

func (QuantifiedStatement) isProposition()        {}
func (QuantifiedStatement) Kind() PropositionKind { return QuantifiedStatementKind }
func (p QuantifiedStatement) Validate() error {
	if !p.Quantifier.Valid() {
		return fmt.Errorf("invalid quantifier %q", p.Quantifier)
	}
	if err := p.Domain.Validate(); err != nil {
		return err
	}
	return p.Predicate.Validate()
}
func (p QuantifiedStatement) Describe() string {
	variable := "s"
	if p.Domain.Kind == NontrivialZerosDomain || p.Domain.Kind == ZerosBelowHeightDomain {
		variable = "ρ"
	}
	return fmt.Sprintf("%s %s in %s: %s", p.Quantifier.Describe(), variable, p.Domain.Describe(), p.Predicate.Describe(variable))
}

type RepresentationName string

const (
	DirichletSeriesRepresentation      RepresentationName = "dirichlet_series"
	EulerProductRepresentation         RepresentationName = "euler_product"
	AnalyticContinuationRepresentation RepresentationName = "analytic_continuation"
	CompletedXiRepresentation          RepresentationName = "completed_xi"
)

type Representation struct {
	Object      Function           `json:"object"`
	BaseObject  Function           `json:"base_object,omitempty"`
	Name        RepresentationName `json:"name"`
	ValidOn     Domain             `json:"valid_on"`
	Formula     string             `json:"formula"`
	Affordances []string           `json:"affordances"`
}

func (r Representation) Validate() error {
	if r.Object == 0 || (r.Name != DirichletSeriesRepresentation && r.Name != EulerProductRepresentation && r.Name != AnalyticContinuationRepresentation && r.Name != CompletedXiRepresentation) {
		return fmt.Errorf("invalid representation identity")
	}
	if r.Name == CompletedXiRepresentation {
		if r.Object != RiemannXi || r.BaseObject != RiemannZeta {
			return fmt.Errorf("completed xi must structurally reference zeta")
		}
	} else if r.BaseObject != 0 {
		return fmt.Errorf("base object is only valid for completed xi")
	}
	if err := r.ValidOn.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(r.Formula) == "" {
		return fmt.Errorf("representation formula is empty")
	}
	return nil
}

type RepresentationProposition struct {
	Representation Representation `json:"representation"`
}

func (RepresentationProposition) isProposition()        {}
func (RepresentationProposition) Kind() PropositionKind { return RepresentationKind }
func (p RepresentationProposition) Describe() string {
	return fmt.Sprintf("%s represented by %s on %s", p.Representation.Object.String(), p.Representation.Name, p.Representation.ValidOn.Describe())
}

type RepresentationIdentity struct {
	Object Function           `json:"object"`
	Left   RepresentationName `json:"left"`
	Right  RepresentationName `json:"right"`
	Domain Domain             `json:"domain"`
}

func (p RepresentationIdentity) Validate() error {
	if p.Object == 0 || p.Left == p.Right {
		return fmt.Errorf("invalid representation identity")
	}
	if (p.Left != DirichletSeriesRepresentation && p.Left != EulerProductRepresentation) ||
		(p.Right != DirichletSeriesRepresentation && p.Right != EulerProductRepresentation) {
		return fmt.Errorf("representation identity contains an unknown representation")
	}
	return p.Domain.Validate()
}

func (RepresentationIdentity) isProposition()        {}
func (RepresentationIdentity) Kind() PropositionKind { return RepresentationIdentityKind }
func (p RepresentationIdentity) Describe() string {
	return fmt.Sprintf("%s and %s denote the same %s on %s", p.Left, p.Right, p.Object.String(), p.Domain.Describe())
}

type AnalyticFactName string

const (
	EulerProductConvergesAbsolutely AnalyticFactName = "euler_product_converges_absolutely"
	EulerProductFactorsNonzero      AnalyticFactName = "euler_product_factors_nonzero"
	ConvergentProductNonzeroLimit   AnalyticFactName = "convergent_product_has_nonzero_limit"
)

type AnalyticFact struct {
	Fact   AnalyticFactName `json:"fact"`
	Object Function         `json:"object"`
	Domain Domain           `json:"domain"`
}

func (p AnalyticFact) Validate() error {
	if p.Object == 0 {
		return fmt.Errorf("analytic fact has no mathematical object")
	}
	if p.Fact != EulerProductConvergesAbsolutely && p.Fact != EulerProductFactorsNonzero && p.Fact != ConvergentProductNonzeroLimit {
		return fmt.Errorf("unknown analytic fact %q", p.Fact)
	}
	return p.Domain.Validate()
}

func (AnalyticFact) isProposition()        {}
func (AnalyticFact) Kind() PropositionKind { return AnalyticFactKind }
func (p AnalyticFact) Describe() string {
	switch p.Fact {
	case EulerProductConvergesAbsolutely:
		return "the Euler product for ζ converges absolutely on " + p.Domain.Describe()
	case EulerProductFactorsNonzero:
		return "every Euler factor (1 - p^-s)^-1 is finite and nonzero on " + p.Domain.Describe()
	case ConvergentProductNonzeroLimit:
		return "an absolutely convergent product of nonzero factors has a nonzero limit"
	default:
		return "unknown analytic fact"
	}
}

// NamedObligation remains only as a test/obligation escape hatch.
type NamedObligation struct{ Name string }

func (NamedObligation) isProposition()        {}
func (NamedObligation) Kind() PropositionKind { return NamedObligationKind }
func (p NamedObligation) Describe() string    { return p.Name }

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
	CertifiedComputationEvidence EvidenceKind = "certified_computation"
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
	Theorem        TheoremID        `json:"theorem,omitempty"`
	Source         Reference        `json:"source"`
}

type Claim struct {
	ID          ClaimID
	Proposition Proposition
	Assumptions []Assumption
	Evidence    []Evidence
	Exactness   Exactness
	Provenance  Provenance
}

func (c Claim) Validate() error {
	if c.ID == "" || c.Proposition == nil {
		return fmt.Errorf("claim must have an ID and proposition")
	}
	if c.Exactness != Exact && c.Exactness != Approximate {
		return fmt.Errorf("claim %q has invalid exactness %q", c.ID, c.Exactness)
	}
	switch p := c.Proposition.(type) {
	case QuantifiedStatement:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case RepresentationProposition:
		if err := p.Representation.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case RepresentationIdentity:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case AnalyticFact:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case NamedObligation:
		if strings.TrimSpace(p.Name) == "" {
			return fmt.Errorf("claim %q has an empty named obligation", c.ID)
		}
	case ZeroAtPoint:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case SideCondition:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case FunctionalIdentity:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case ZeroSetProperty:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case ZeroClassification:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case FunctionalDefinition:
		if err := p.Functional.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case UniversalFunctionalStatement:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case TestFunctionAdmissibility:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case ExplicitFormulaIdentity:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case FiniteSpanDefinition:
		if err := p.Span.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case QuadraticFormStructure:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case HermitianFormDefinition:
		if err := p.Form.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case HermitianMatrixDefinition:
		if err := p.Matrix.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case FiniteSpanFunctionalStatement:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case CoordinateQuadraticPositivity:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case QuadraticMatrixIdentity:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case MatrixProperty:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	case TwoByTwoPrincipalMinorCertificate:
		if err := p.Validate(); err != nil {
			return fmt.Errorf("claim %q: %w", c.ID, err)
		}
	}
	seen := make(map[AssumptionID]bool, len(c.Assumptions))
	for _, assumption := range c.Assumptions {
		if assumption.ID == "" || strings.TrimSpace(assumption.Description) == "" || seen[assumption.ID] {
			return fmt.Errorf("claim %q has invalid or repeated assumption", c.ID)
		}
		seen[assumption.ID] = true
	}
	return nil
}

func CloneAssumptions(in []Assumption) []Assumption { return append([]Assumption(nil), in...) }
func CloneRepresentation(r Representation) Representation {
	r.Affordances = append([]string(nil), r.Affordances...)
	return r
}

// Quantified extracts structural quantifier/domain semantics when present.
func Quantified(p Proposition) (QuantifiedStatement, bool) {
	q, ok := p.(QuantifiedStatement)
	return q, ok
}

// SemanticKey is canonical identity for the typed meaning of a proposition.
// Display-only representation formulae and affordances are deliberately absent.
func SemanticKey(p Proposition) string {
	switch v := p.(type) {
	case QuantifiedStatement:
		return fmt.Sprintf("q|%s|%s|%d|%d|%s|%d", v.Quantifier, v.Domain.Kind, v.Domain.Function, v.Domain.Bound, v.Predicate.Kind, v.Predicate.Function)
	case RepresentationProposition:
		r := v.Representation
		return fmt.Sprintf("r|%d|%d|%s|%s|%d|%d", r.Object, r.BaseObject, r.Name, r.ValidOn.Kind, r.ValidOn.Function, r.ValidOn.Bound)
	case RepresentationIdentity:
		return fmt.Sprintf("i|%d|%s|%s|%s|%d|%d", v.Object, v.Left, v.Right, v.Domain.Kind, v.Domain.Function, v.Domain.Bound)
	case AnalyticFact:
		return fmt.Sprintf("a|%s|%d|%s|%d|%d", v.Fact, v.Object, v.Domain.Kind, v.Domain.Function, v.Domain.Bound)
	case NamedObligation:
		return "n|" + v.Name
	case ZeroAtPoint:
		return fmt.Sprintf("z|%d|%s|%s", v.Object, v.Point.Key(), v.Classification)
	case SideCondition:
		return fmt.Sprintf("s|%s|%d|%s|%s", v.Condition, v.Object, v.Point.Key(), domainKey(v.Domain))
	case FunctionalIdentity:
		return fmt.Sprintf("f|%d|%s|%s|%s", v.Object, v.Left, v.Right, domainKey(v.Domain))
	case ZeroSetProperty:
		return fmt.Sprintf("g|%s|%s|%s|%s", domainKey(v.Set), v.Property, v.Transform, domainKey(v.Region))
	case ZeroClassification:
		return fmt.Sprintf("c|%d|%s", v.Object, v.Classification)
	case FunctionalDefinition:
		return "fd|" + string(v.Functional.ID) + "|" + string(v.Functional.TransformConvention)
	case UniversalFunctionalStatement:
		return fmt.Sprintf("uf|%s|%s|%s|%s|%s", v.Quantifier, v.FunctionClass.Key(), v.Functional, v.Predicate, v.TransformConvention)
	case FiniteSpanDefinition:
		return "span-definition|" + v.Span.Key()
	case QuadraticFormStructure:
		return v.Key()
	case HermitianFormDefinition:
		return "hermitian-form|" + v.Form.Key()
	case HermitianMatrixDefinition:
		return "hermitian-matrix|" + v.Matrix.Key()
	case FiniteSpanFunctionalStatement:
		return v.Key()
	case CoordinateQuadraticPositivity:
		return v.Key()
	case QuadraticMatrixIdentity:
		return v.Key()
	case MatrixProperty:
		return v.Key()
	case TestFunctionAdmissibility:
		return "ta|" + v.Function.Key() + "|" + v.Class.Key()
	case ExplicitFormulaIdentity:
		return "ef|" + string(v.Functional) + "|" + v.FunctionClass.Key() + "|" + string(v.TransformConvention.ID) + "|" + string(v.Theorem)
	default:
		return fmt.Sprintf("unknown|%T", p)
	}
}

func domainKey(d Domain) string { return fmt.Sprintf("%s|%d|%d", d.Kind, d.Function, d.Bound) }
