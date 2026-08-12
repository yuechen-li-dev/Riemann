package semantic

import (
	"fmt"
	"strings"
)

// PointTransform is the closed Klein-four action needed by zeta symmetry.
// Composition is exact bitwise composition: conjugation and one-minus commute.
type PointTransform uint8

const (
	IdentityTransform PointTransform = iota
	ConjugateTransform
	OneMinusTransform
	CriticalReflectionTransform // 1-conjugate(s)
)

func (t PointTransform) Valid() bool { return t <= CriticalReflectionTransform }
func (t PointTransform) String() string {
	switch t {
	case IdentityTransform:
		return "identity"
	case ConjugateTransform:
		return "conjugate"
	case OneMinusTransform:
		return "one_minus"
	case CriticalReflectionTransform:
		return "one_minus_conjugate"
	default:
		return fmt.Sprintf("unknown_transform_%d", t)
	}
}
func (t PointTransform) MarshalJSON() ([]byte, error) {
	if !t.Valid() {
		return nil, fmt.Errorf("invalid point transform %d", t)
	}
	return []byte(fmt.Sprintf("%q", t.String())), nil
}
func Compose(outer, inner PointTransform) PointTransform { return outer ^ inner }

type PointLocation string

const (
	UnconstrainedPoint PointLocation = "unconstrained"
	CriticalLinePoint  PointLocation = "critical_line"
	RealPoint          PointLocation = "real_axis"
)

// PointExpr stores meaning, not display syntax. Location records only the two
// fixed-set facts M3 needs and drives honest orbit deduplication.
type PointExpr struct {
	Symbol    string         `json:"symbol"`
	Transform PointTransform `json:"transform"`
	Location  PointLocation  `json:"location,omitempty"`
}

func Point(symbol string) PointExpr { return PointExpr{Symbol: symbol, Location: UnconstrainedPoint} }
func PointOnCriticalLine(symbol string) PointExpr {
	return PointExpr{Symbol: symbol, Location: CriticalLinePoint}
}
func PointOnRealAxis(symbol string) PointExpr { return PointExpr{Symbol: symbol, Location: RealPoint} }

func (p PointExpr) Apply(t PointTransform) PointExpr {
	p.Transform = Compose(t, p.Transform)
	return p.Canonical()
}
func (p PointExpr) Canonical() PointExpr {
	if p.Location == "" {
		p.Location = UnconstrainedPoint
	}
	if p.Location == CriticalLinePoint {
		other := Compose(CriticalReflectionTransform, p.Transform)
		if other < p.Transform {
			p.Transform = other
		}
	}
	if p.Location == RealPoint {
		other := Compose(ConjugateTransform, p.Transform)
		if other < p.Transform {
			p.Transform = other
		}
	}
	return p
}
func (p PointExpr) Validate() error {
	if strings.TrimSpace(p.Symbol) == "" || !p.Transform.Valid() {
		return fmt.Errorf("invalid complex point expression")
	}
	if p.Location != UnconstrainedPoint && p.Location != CriticalLinePoint && p.Location != RealPoint {
		return fmt.Errorf("invalid point location %q", p.Location)
	}
	return nil
}
func (p PointExpr) Key() string {
	p = p.Canonical()
	return fmt.Sprintf("%s|%d|%s", p.Symbol, p.Transform, p.Location)
}
func (p PointExpr) Describe() string {
	p = p.Canonical()
	switch p.Transform {
	case IdentityTransform:
		return p.Symbol
	case ConjugateTransform:
		return "conjugate(" + p.Symbol + ")"
	case OneMinusTransform:
		return "1-" + p.Symbol
	case CriticalReflectionTransform:
		return "1-conjugate(" + p.Symbol + ")"
	default:
		return "invalid-point"
	}
}

type ZeroClass string

const (
	UnclassifiedZero ZeroClass = "unclassified"
	TrivialZero      ZeroClass = "trivial"
	NontrivialZero   ZeroClass = "nontrivial"
)

type ZeroAtPoint struct {
	Object         Function  `json:"object"`
	Point          PointExpr `json:"point"`
	Classification ZeroClass `json:"classification"`
}

func (ZeroAtPoint) isProposition()        {}
func (ZeroAtPoint) Kind() PropositionKind { return ZeroAtPointKind }
func (p ZeroAtPoint) Validate() error {
	if p.Object == 0 || (p.Classification != UnclassifiedZero && p.Classification != TrivialZero && p.Classification != NontrivialZero) {
		return fmt.Errorf("invalid zero claim")
	}
	return p.Point.Validate()
}
func (p ZeroAtPoint) Describe() string {
	return fmt.Sprintf("%s(%s) = 0 [%s zero]", shortFunction(p.Object), p.Point.Describe(), p.Classification)
}

type SideConditionName string

const (
	PointInValidityDomain            SideConditionName = "point_in_validity_domain"
	CompletionFactorRegularNonzero   SideConditionName = "completion_factor_regular_nonzero"
	AnalyticContinuationAvailable    SideConditionName = "analytic_continuation_available"
	RealAnalyticConjugationAvailable SideConditionName = "real_analytic_conjugation_available"
)

type SideCondition struct {
	Condition SideConditionName `json:"condition"`
	Object    Function          `json:"object"`
	Point     PointExpr         `json:"point,omitempty"`
	Domain    Domain            `json:"domain,omitempty"`
}

func (SideCondition) isProposition()        {}
func (SideCondition) Kind() PropositionKind { return SideConditionKind }
func (p SideCondition) Validate() error {
	if p.Object == 0 {
		return fmt.Errorf("side condition has no object")
	}
	switch p.Condition {
	case PointInValidityDomain:
		if err := p.Point.Validate(); err != nil {
			return err
		}
		if err := p.Domain.Validate(); err != nil {
			return err
		}
	case CompletionFactorRegularNonzero:
		if err := p.Point.Validate(); err != nil {
			return err
		}
	case AnalyticContinuationAvailable, RealAnalyticConjugationAvailable:
		if err := p.Domain.Validate(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown side condition %q", p.Condition)
	}
	return nil
}
func (p SideCondition) Describe() string {
	switch p.Condition {
	case PointInValidityDomain:
		return fmt.Sprintf("%s lies in %s", p.Point.Describe(), p.Domain.Describe())
	case CompletionFactorRegularNonzero:
		return fmt.Sprintf("the xi completion factor is finite and nonzero at %s", p.Point.Describe())
	case AnalyticContinuationAvailable:
		return fmt.Sprintf("analytic continuation of %s is available on %s", p.Object.String(), p.Domain.Describe())
	case RealAnalyticConjugationAvailable:
		return fmt.Sprintf("%s respects complex conjugation on its finite points in %s", p.Object.String(), p.Domain.Describe())
	default:
		return "unknown side condition"
	}
}

type FunctionalIdentity struct {
	Object  Function       `json:"object"`
	Left    PointTransform `json:"left"`
	Right   PointTransform `json:"right"`
	Domain  Domain         `json:"domain"`
	Formula string         `json:"formula"`
}

func (FunctionalIdentity) isProposition()        {}
func (FunctionalIdentity) Kind() PropositionKind { return FunctionalIdentityKind }
func (p FunctionalIdentity) Validate() error {
	if p.Object == 0 || !p.Left.Valid() || !p.Right.Valid() || strings.TrimSpace(p.Formula) == "" {
		return fmt.Errorf("invalid functional identity")
	}
	return p.Domain.Validate()
}
func (p FunctionalIdentity) Describe() string { return p.Formula + " on " + p.Domain.Describe() }

type ZeroSetPropertyName string

const (
	InvariantUnderTransform ZeroSetPropertyName = "invariant_under_transform"
	ConfinedToRegion        ZeroSetPropertyName = "confined_to_region"
)

type ZeroSetProperty struct {
	Set       Domain              `json:"set"`
	Property  ZeroSetPropertyName `json:"property"`
	Transform PointTransform      `json:"transform,omitempty"`
	Region    Domain              `json:"region,omitempty"`
}

func (ZeroSetProperty) isProposition()        {}
func (ZeroSetProperty) Kind() PropositionKind { return ZeroSetPropertyKind }
func (p ZeroSetProperty) Validate() error {
	if err := p.Set.Validate(); err != nil {
		return err
	}
	if p.Property == InvariantUnderTransform {
		if !p.Transform.Valid() || p.Transform == IdentityTransform {
			return fmt.Errorf("invalid invariance transform")
		}
	} else if p.Property == ConfinedToRegion {
		return p.Region.Validate()
	} else {
		return fmt.Errorf("unknown zero-set property")
	}
	return nil
}
func (p ZeroSetProperty) Describe() string {
	if p.Property == InvariantUnderTransform {
		return fmt.Sprintf("%s is invariant under %s", p.Set.Describe(), p.Transform)
	}
	return fmt.Sprintf("%s is confined to %s", p.Set.Describe(), p.Region.Describe())
}

type ZeroClassificationName string

const (
	SimplePoleAtOne                 ZeroClassificationName = "simple_pole_at_one"
	TrivialZerosExactlyNegativeEven ZeroClassificationName = "trivial_zeros_negative_even"
)

type ZeroClassification struct {
	Object         Function               `json:"object"`
	Classification ZeroClassificationName `json:"classification"`
}

func (ZeroClassification) isProposition()        {}
func (ZeroClassification) Kind() PropositionKind { return ZeroClassificationKind }
func (p ZeroClassification) Validate() error {
	if p.Object == 0 {
		return fmt.Errorf("classification has no object")
	}
	if p.Classification != SimplePoleAtOne && p.Classification != TrivialZerosExactlyNegativeEven {
		return fmt.Errorf("unknown classification")
	}
	return nil
}
func (p ZeroClassification) Describe() string {
	switch p.Classification {
	case SimplePoleAtOne:
		return "ζ has its only simple pole at s=1"
	case TrivialZerosExactlyNegativeEven:
		return "the trivial zeros of ζ are exactly the negative even integers"
	}
	return "unknown classification"
}

func shortFunction(f Function) string {
	if f == RiemannZeta {
		return "ζ"
	}
	if f == RiemannXi {
		return "ξ"
	}
	return f.String()
}
