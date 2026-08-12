package semantic

import (
	"fmt"
	"math"
	"strings"
)

// ValueEvidenceKind is deliberately independent of Claim.Exactness. An entry
// definition may be exact while its currently attached scalar is unevaluated
// or merely approximate.
type ValueEvidenceKind string

const (
	UnevaluatedExactDefinition ValueEvidenceKind = "unevaluated_exact_definition"
	ApproximateValue           ValueEvidenceKind = "approximate_value"
	CertifiedInterval          ValueEvidenceKind = "certified_interval"
	ExactValue                 ValueEvidenceKind = "exact_value"
)

type ComplexValue struct {
	Real float64 `json:"real"`
	Imag float64 `json:"imag"`
}

type ExactScalar struct {
	Expression  string `json:"expression"`
	Numerator   int64  `json:"numerator,omitempty"`
	Denominator int64  `json:"denominator,omitempty"`
}

func (x ExactScalar) Validate() error {
	if strings.TrimSpace(x.Expression) == "" || x.Denominator <= 0 {
		return fmt.Errorf("invalid exact scalar")
	}
	return nil
}

type ExactComplex struct {
	Real ExactScalar `json:"real"`
	Imag ExactScalar `json:"imag"`
}

type ComplexInterval struct {
	RealLower float64 `json:"real_lower"`
	RealUpper float64 `json:"real_upper"`
	ImagLower float64 `json:"imag_lower"`
	ImagUpper float64 `json:"imag_upper"`
	// Exact decimal renderings preserve the directed big.Float endpoints used
	// by the certifier; the float64 fields are outward-rounded display values.
	RealLowerExact string `json:"real_lower_exact,omitempty"`
	RealUpperExact string `json:"real_upper_exact,omitempty"`
	Representation string `json:"representation,omitempty"`
}

type ErrorSemantics struct {
	Kind            string          `json:"kind"`
	Bound           string          `json:"bound,omitempty"`
	ProofObjectKind ProofObjectKind `json:"proof_object_kind,omitempty"`
	ProofObject     string          `json:"proof_object,omitempty"`
	Notes           string          `json:"notes,omitempty"`
}

type ProofObjectKind string

const (
	TheoremBackedBound       ProofObjectKind = "theorem_backed_bound"
	IndependentExactArgument ProofObjectKind = "independent_exact_argument"
	AnalyticTailBound        ProofObjectKind = "analytic_tail_bound"
	CertifiedQuadratureBound ProofObjectKind = "certified_quadrature_bound"
	OutwardRoundedFiniteSum  ProofObjectKind = "outward_rounded_finite_sum"
	OutwardRoundedArithmetic ProofObjectKind = "outward_rounded_arithmetic"
)

func (k ProofObjectKind) CertifiesInterval() bool {
	switch k {
	case TheoremBackedBound, AnalyticTailBound, CertifiedQuadratureBound, OutwardRoundedFiniteSum, OutwardRoundedArithmetic:
		return true
	default:
		return false
	}
}

type TruncationInfo struct {
	Parameter         string `json:"parameter"`
	Bound             string `json:"bound"`
	SummandDefinition string `json:"summand_definition"`
	EnumerationSource string `json:"enumeration_source"`
	TermsEvaluated    int    `json:"terms_evaluated"`
	RemainderStatus   string `json:"remainder_status"`
	SupportExhaustive bool   `json:"support_exhaustive"`
}

type QuadratureInfo struct {
	Method           string   `json:"method"`
	Tolerance        float64  `json:"tolerance"`
	DomainHandling   string   `json:"domain_handling"`
	ErrorRigorous    bool     `json:"error_rigorous"`
	Partitions       int      `json:"partitions,omitempty"`
	Breakpoints      []string `json:"breakpoints,omitempty"`
	RemainderTheorem string   `json:"remainder_theorem,omitempty"`
}

type TailBound struct {
	Start      string          `json:"start"`
	LowerBound string          `json:"lower_bound"`
	UpperBound string          `json:"upper_bound"`
	Derivation string          `json:"derivation"`
	ProofKind  ProofObjectKind `json:"proof_object_kind"`
	Exactness  string          `json:"exactness"`
	Provenance string          `json:"provenance"`
}

type EvaluationMetadata struct {
	SemanticTargetID    string                `json:"semantic_target_id"`
	Backend             string                `json:"backend"`
	BackendVersion      string                `json:"backend_version"`
	PrecisionBits       int                   `json:"precision_bits"`
	TransformConvention TransformConventionID `json:"transform_convention"`
	Truncation          *TruncationInfo       `json:"truncation,omitempty"`
	Quadrature          *QuadratureInfo       `json:"quadrature,omitempty"`
	Tail                *TailBound            `json:"tail_bound,omitempty"`
	Error               ErrorSemantics        `json:"error_semantics"`
	Provenance          []TheoremID           `json:"theorem_provenance"`
}

type TransformEvaluation struct {
	InputFunction TestFunction          `json:"input_function"`
	Convention    TransformConventionID `json:"transform_convention"`
	Point         string                `json:"point"`
	Definition    string                `json:"definition"`
	Value         EntryValue            `json:"value"`
}

type EntryValue struct {
	Kind            ValueEvidenceKind   `json:"kind"`
	DefinitionExact bool                `json:"definition_exact"`
	Approximate     *ComplexValue       `json:"approximate,omitempty"`
	Interval        *ComplexInterval    `json:"interval,omitempty"`
	Exact           *ExactComplex       `json:"exact,omitempty"`
	Metadata        *EvaluationMetadata `json:"metadata,omitempty"`
}

func NewUnevaluatedEntryValue() EntryValue {
	return EntryValue{Kind: UnevaluatedExactDefinition, DefinitionExact: true}
}

func (v EntryValue) Validate() error {
	if !v.DefinitionExact {
		return fmt.Errorf("entry value lost exact semantic definition")
	}
	count := 0
	if v.Approximate != nil {
		count++
	}
	if v.Interval != nil {
		count++
	}
	if v.Exact != nil {
		count++
	}
	switch v.Kind {
	case UnevaluatedExactDefinition:
		if count != 0 || v.Metadata != nil {
			return fmt.Errorf("unevaluated definition carries a value")
		}
	case ApproximateValue:
		if count != 1 || v.Approximate == nil || v.Metadata == nil || v.Metadata.PrecisionBits <= 0 || !finiteComplex(*v.Approximate) {
			return fmt.Errorf("invalid approximate value")
		}
	case CertifiedInterval:
		if count != 1 || v.Interval == nil || v.Metadata == nil || !v.Metadata.Error.ProofObjectKind.CertifiesInterval() || v.Metadata.Error.ProofObject == "" || v.Metadata.PrecisionBits <= 0 || !finiteInterval(*v.Interval) || v.Interval.RealLower > v.Interval.RealUpper || v.Interval.ImagLower > v.Interval.ImagUpper {
			return fmt.Errorf("invalid certified interval")
		}
		if v.Metadata.Error.ProofObjectKind != TheoremBackedBound && (v.Interval.Representation == "" || v.Interval.RealLowerExact == "" || v.Interval.RealUpperExact == "") {
			return fmt.Errorf("certified numerical interval lost representation or directed endpoints")
		}
		if v.Metadata.Error.ProofObjectKind == CertifiedQuadratureBound && (v.Metadata.Quadrature == nil || !v.Metadata.Quadrature.ErrorRigorous || v.Metadata.Tail == nil || v.Metadata.Tail.ProofKind != AnalyticTailBound) {
			return fmt.Errorf("certified quadrature requires rigorous partition and analytic tail proof objects")
		}
	case ExactValue:
		if count != 1 || v.Exact == nil || v.Metadata == nil || v.Metadata.Error.ProofObjectKind != IndependentExactArgument || v.Metadata.Error.ProofObject == "" {
			return fmt.Errorf("invalid exact value")
		}
		if err := v.Exact.Real.Validate(); err != nil {
			return err
		}
		if err := v.Exact.Imag.Validate(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown entry value evidence kind %q", v.Kind)
	}
	return nil
}

func finiteComplex(v ComplexValue) bool {
	return !math.IsNaN(v.Real) && !math.IsInf(v.Real, 0) && !math.IsNaN(v.Imag) && !math.IsInf(v.Imag, 0)
}

func finiteInterval(v ComplexInterval) bool {
	return finiteComplex(ComplexValue{Real: v.RealLower, Imag: v.RealUpper}) && finiteComplex(ComplexValue{Real: v.ImagLower, Imag: v.ImagUpper})
}

// UpgradeEntryValue encodes the evidence lattice. Precision is never an exact
// argument, and heuristic error estimates are never interval certificates.
func UpgradeEntryValue(from, to EntryValue) (EntryValue, error) {
	if err := from.Validate(); err != nil {
		return EntryValue{}, err
	}
	if err := to.Validate(); err != nil {
		return EntryValue{}, err
	}
	legal := false
	switch from.Kind {
	case UnevaluatedExactDefinition:
		legal = to.Kind == ApproximateValue || to.Kind == CertifiedInterval || to.Kind == ExactValue
	case CertifiedInterval:
		legal = to.Kind == ExactValue
	case ApproximateValue, ExactValue:
		legal = false
	}
	if !legal {
		return EntryValue{}, fmt.Errorf("illegal value-evidence upgrade %s -> %s", from.Kind, to.Kind)
	}
	return CloneEntryValue(to), nil
}

// WeakestValueKind is the conservative arithmetic join used during assembly.
func WeakestValueKind(values ...EntryValue) ValueEvidenceKind {
	result := ExactValue
	for _, value := range values {
		switch value.Kind {
		case UnevaluatedExactDefinition:
			return UnevaluatedExactDefinition
		case ApproximateValue:
			return ApproximateValue
		case CertifiedInterval:
			if result == ExactValue {
				result = CertifiedInterval
			}
		}
	}
	return result
}

func CloneEntryValue(v EntryValue) EntryValue {
	if v.Approximate != nil {
		x := *v.Approximate
		v.Approximate = &x
	}
	if v.Interval != nil {
		x := *v.Interval
		v.Interval = &x
	}
	if v.Exact != nil {
		x := *v.Exact
		v.Exact = &x
	}
	if v.Metadata != nil {
		x := *v.Metadata
		x.Provenance = append([]TheoremID(nil), v.Metadata.Provenance...)
		if v.Metadata.Truncation != nil {
			y := *v.Metadata.Truncation
			x.Truncation = &y
		}
		if v.Metadata.Quadrature != nil {
			y := *v.Metadata.Quadrature
			y.Breakpoints = append([]string(nil), v.Metadata.Quadrature.Breakpoints...)
			x.Quadrature = &y
		}
		if v.Metadata.Tail != nil {
			y := *v.Metadata.Tail
			x.Tail = &y
		}
		v.Metadata = &x
	}
	return v
}
