package semantic

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

// ScalarExpression is the deliberately small arithmetic IR needed by M13.
// It is an expression tree, not display metadata and not a general CAS.
type ScalarExpressionKind string

const (
	ScalarRational  ScalarExpressionKind = "rational"
	ScalarParameter ScalarExpressionKind = "parameter"
	ScalarAdd       ScalarExpressionKind = "add"
	ScalarSubtract  ScalarExpressionKind = "subtract"
	ScalarMultiply  ScalarExpressionKind = "multiply"
	ScalarDivide    ScalarExpressionKind = "divide"
	ScalarPower     ScalarExpressionKind = "power"
	ScalarSqrt      ScalarExpressionKind = "sqrt"
	ScalarSin       ScalarExpressionKind = "sin"
	ScalarCos       ScalarExpressionKind = "cos"
	ScalarTan       ScalarExpressionKind = "tan"
	ScalarCot       ScalarExpressionKind = "cot"
)

type ScalarExpression struct {
	Kind      ScalarExpressionKind `json:"kind"`
	Rational  *ExactRational       `json:"rational,omitempty"`
	Parameter string               `json:"parameter,omitempty"`
	Left      *ScalarExpression    `json:"left,omitempty"`
	Right     *ScalarExpression    `json:"right,omitempty"`
	Exponent  int                  `json:"exponent,omitempty"`
}

func ScalarRat(n, d int64) ScalarExpression {
	r := ExactRational{Numerator: n, Denominator: d}
	return ScalarExpression{Kind: ScalarRational, Rational: &r}
}
func ScalarParam(name string) ScalarExpression {
	return ScalarExpression{Kind: ScalarParameter, Parameter: name}
}
func ScalarBinary(kind ScalarExpressionKind, left, right ScalarExpression) ScalarExpression {
	return ScalarExpression{Kind: kind, Left: &left, Right: &right}
}
func ScalarUnary(kind ScalarExpressionKind, value ScalarExpression) ScalarExpression {
	return ScalarExpression{Kind: kind, Left: &value}
}
func ScalarPow(value ScalarExpression, exponent int) ScalarExpression {
	return ScalarExpression{Kind: ScalarPower, Left: &value, Exponent: exponent}
}

func (e ScalarExpression) Validate() error {
	switch e.Kind {
	case ScalarRational:
		if e.Rational == nil || e.Parameter != "" || e.Left != nil || e.Right != nil || e.Exponent != 0 {
			return fmt.Errorf("malformed rational scalar expression")
		}
		return e.Rational.Validate()
	case ScalarParameter:
		if strings.TrimSpace(e.Parameter) == "" || e.Rational != nil || e.Left != nil || e.Right != nil || e.Exponent != 0 {
			return fmt.Errorf("malformed scalar parameter")
		}
		return nil
	case ScalarAdd, ScalarSubtract, ScalarMultiply, ScalarDivide:
		if e.Left == nil || e.Right == nil || e.Rational != nil || e.Parameter != "" || e.Exponent != 0 {
			return fmt.Errorf("malformed binary scalar expression")
		}
		if err := e.Left.Validate(); err != nil {
			return err
		}
		return e.Right.Validate()
	case ScalarPower:
		if e.Left == nil || e.Right != nil || e.Rational != nil || e.Parameter != "" || e.Exponent < 1 || e.Exponent > 2 {
			return fmt.Errorf("M13 scalar powers are restricted to 1 or 2")
		}
		return e.Left.Validate()
	case ScalarSqrt, ScalarSin, ScalarCos, ScalarTan, ScalarCot:
		if e.Left == nil || e.Right != nil || e.Rational != nil || e.Parameter != "" || e.Exponent != 0 {
			return fmt.Errorf("malformed unary scalar expression")
		}
		return e.Left.Validate()
	default:
		return fmt.Errorf("unknown scalar expression kind %q", e.Kind)
	}
}

func (e ScalarExpression) EvalFloat(parameters map[string]float64) (float64, error) {
	if err := e.Validate(); err != nil {
		return 0, err
	}
	var eval func(ScalarExpression) (float64, error)
	eval = func(x ScalarExpression) (float64, error) {
		switch x.Kind {
		case ScalarRational:
			return float64(x.Rational.Numerator) / float64(x.Rational.Denominator), nil
		case ScalarParameter:
			v, ok := parameters[x.Parameter]
			if !ok {
				return 0, fmt.Errorf("missing scalar parameter %q", x.Parameter)
			}
			return v, nil
		case ScalarAdd, ScalarSubtract, ScalarMultiply, ScalarDivide:
			a, err := eval(*x.Left)
			if err != nil {
				return 0, err
			}
			b, err := eval(*x.Right)
			if err != nil {
				return 0, err
			}
			switch x.Kind {
			case ScalarAdd:
				return a + b, nil
			case ScalarSubtract:
				return a - b, nil
			case ScalarMultiply:
				return a * b, nil
			default:
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				return a / b, nil
			}
		case ScalarPower:
			a, err := eval(*x.Left)
			if err != nil {
				return 0, err
			}
			return math.Pow(a, float64(x.Exponent)), nil
		case ScalarSqrt, ScalarSin, ScalarCos, ScalarTan, ScalarCot:
			a, err := eval(*x.Left)
			if err != nil {
				return 0, err
			}
			switch x.Kind {
			case ScalarSqrt:
				if a < 0 {
					return 0, fmt.Errorf("sqrt of negative scalar")
				}
				return math.Sqrt(a), nil
			case ScalarSin:
				return math.Sin(a), nil
			case ScalarCos:
				return math.Cos(a), nil
			case ScalarTan:
				return math.Tan(a), nil
			default:
				t := math.Tan(a)
				if t == 0 {
					return 0, fmt.Errorf("cotangent pole")
				}
				return 1 / t, nil
			}
		default:
			return 0, fmt.Errorf("unknown scalar expression")
		}
	}
	return eval(e)
}

// EvalExact evaluates the algebraic subset used by normalization checks.
func (e ScalarExpression) EvalExact(parameters map[string]ExactRational) (ExactRational, error) {
	if err := e.Validate(); err != nil {
		return ExactRational{}, err
	}
	var eval func(ScalarExpression) (*big.Rat, error)
	eval = func(x ScalarExpression) (*big.Rat, error) {
		switch x.Kind {
		case ScalarRational:
			return rat(*x.Rational), nil
		case ScalarParameter:
			v, ok := parameters[x.Parameter]
			if !ok {
				return nil, fmt.Errorf("missing scalar parameter %q", x.Parameter)
			}
			if err := v.Validate(); err != nil {
				return nil, err
			}
			return rat(v), nil
		case ScalarAdd, ScalarSubtract, ScalarMultiply, ScalarDivide:
			a, err := eval(*x.Left)
			if err != nil {
				return nil, err
			}
			b, err := eval(*x.Right)
			if err != nil {
				return nil, err
			}
			switch x.Kind {
			case ScalarAdd:
				return new(big.Rat).Add(a, b), nil
			case ScalarSubtract:
				return new(big.Rat).Sub(a, b), nil
			case ScalarMultiply:
				return new(big.Rat).Mul(a, b), nil
			default:
				if b.Sign() == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return new(big.Rat).Quo(a, b), nil
			}
		case ScalarPower:
			a, err := eval(*x.Left)
			if err != nil {
				return nil, err
			}
			if x.Exponent == 1 {
				return a, nil
			}
			return new(big.Rat).Mul(a, a), nil
		default:
			return nil, fmt.Errorf("transcendental scalar expression requires certified analytic evaluation")
		}
	}
	v, err := eval(e)
	if err != nil {
		return ExactRational{}, err
	}
	return exactRat(v)
}

func (e ScalarExpression) String() string {
	switch e.Kind {
	case ScalarRational:
		if e.Rational == nil {
			return "<invalid>"
		}
		if e.Rational.Denominator == 1 {
			return fmt.Sprintf("%d", e.Rational.Numerator)
		}
		return fmt.Sprintf("(%d/%d)", e.Rational.Numerator, e.Rational.Denominator)
	case ScalarParameter:
		return e.Parameter
	case ScalarAdd:
		return "(" + e.Left.String() + "+" + e.Right.String() + ")"
	case ScalarSubtract:
		return "(" + e.Left.String() + "-" + e.Right.String() + ")"
	case ScalarMultiply:
		return "(" + e.Left.String() + "*" + e.Right.String() + ")"
	case ScalarDivide:
		return "(" + e.Left.String() + "/" + e.Right.String() + ")"
	case ScalarPower:
		return fmt.Sprintf("(%s^%d)", e.Left.String(), e.Exponent)
	case ScalarSqrt:
		return "sqrt(" + e.Left.String() + ")"
	case ScalarSin:
		return "sin(" + e.Left.String() + ")"
	case ScalarCos:
		return "cos(" + e.Left.String() + ")"
	case ScalarTan:
		return "tan(" + e.Left.String() + ")"
	case ScalarCot:
		return "cot(" + e.Left.String() + ")"
	default:
		return "<invalid>"
	}
}

type ScalarDomain struct {
	Lower         ExactRational `json:"lower"`
	Upper         ExactRational `json:"upper"`
	LowerIncluded bool          `json:"lower_included"`
	UpperIncluded bool          `json:"upper_included"`
}

func (d ScalarDomain) Validate() error {
	if err := d.Lower.Validate(); err != nil {
		return err
	}
	if err := d.Upper.Validate(); err != nil {
		return err
	}
	if rat(d.Lower).Cmp(rat(d.Upper)) >= 0 {
		return fmt.Errorf("scalar domain is empty")
	}
	return nil
}
func (d ScalarDomain) ContainsFloat(x float64) bool {
	lo := float64(d.Lower.Numerator) / float64(d.Lower.Denominator)
	hi := float64(d.Upper.Numerator) / float64(d.Upper.Denominator)
	return (x > lo || d.LowerIncluded && x == lo) && (x < hi || d.UpperIncluded && x == hi)
}

type WindowParameter struct {
	Symbol  string       `json:"symbol"`
	Domain  ScalarDomain `json:"domain"`
	Meaning string       `json:"meaning"`
}

func (p WindowParameter) Validate() error {
	if strings.TrimSpace(p.Symbol) == "" || strings.TrimSpace(p.Meaning) == "" {
		return fmt.Errorf("invalid window parameter")
	}
	return p.Domain.Validate()
}

type WindowAdmissibility struct {
	Even                 bool          `json:"even"`
	NonnegativeProfile   bool          `json:"nonnegative_profile"`
	CompactSupport       bool          `json:"compact_support"`
	MonotoneInAbs        bool          `json:"nonincreasing_in_absolute_value"`
	TwiceDifferentiable  bool          `json:"c2_after_endpoint_mollification"`
	FixedWidthRamp       bool          `json:"fixed_width_endpoint_ramp"`
	FourierSupportAtMost ExactRational `json:"normalized_fourier_support_at_most"`
}

type TestWindowFamily struct {
	ID                  string              `json:"id"`
	Parameter           WindowParameter     `json:"parameter"`
	WindowObjectID      string              `json:"window_object_id"`
	SquaredProfileID    string              `json:"squared_profile_id"`
	ProfileArgument     ScalarExpression    `json:"profile_argument"`
	SupportScale        ScalarExpression    `json:"support_scale"`
	TransformConvention string              `json:"transform_convention"`
	Normalization       string              `json:"normalization"`
	Admissibility       WindowAdmissibility `json:"admissibility"`
	Theorem             TheoremID           `json:"theorem"`
	Provenance          Reference           `json:"provenance"`
}

func (f TestWindowFamily) Validate() error {
	if f.ID == "" || f.WindowObjectID == "" || f.SquaredProfileID == "" || f.TransformConvention == "" || f.Normalization == "" || f.Theorem == "" || f.Provenance.Citation == "" {
		return fmt.Errorf("incomplete test-window family")
	}
	if err := f.Parameter.Validate(); err != nil {
		return err
	}
	if err := f.ProfileArgument.Validate(); err != nil {
		return err
	}
	if err := f.SupportScale.Validate(); err != nil {
		return err
	}
	a := f.Admissibility
	if !a.Even || !a.NonnegativeProfile || !a.CompactSupport || !a.MonotoneInAbs || !a.TwiceDifferentiable || !a.FixedWidthRamp {
		return fmt.Errorf("window family lacks a theorem admissibility condition")
	}
	if err := a.FourierSupportAtMost.Validate(); err != nil {
		return err
	}
	if rat(a.FourierSupportAtMost).Cmp(big.NewRat(1, 1)) > 0 {
		return fmt.Errorf("window exceeds imported support theorem")
	}
	return nil
}

type WindowMomentCoefficient struct {
	ID            string           `json:"id"`
	Meaning       string           `json:"meaning"`
	Expression    ScalarExpression `json:"exact_expression"`
	Parameter     string           `json:"parameter"`
	Normalization string           `json:"normalization"`
	Theorem       TheoremID        `json:"source_theorem"`
	ErrorTerm     string           `json:"error_term"`
	Provenance    Reference        `json:"provenance"`
}

func (c WindowMomentCoefficient) Validate() error {
	if c.ID == "" || c.Meaning == "" || c.Parameter == "" || c.Normalization == "" || c.Theorem == "" || c.ErrorTerm == "" || c.Provenance.Citation == "" {
		return fmt.Errorf("incomplete window moment coefficient")
	}
	return c.Expression.Validate()
}

type WindowMomentCoefficients struct {
	MassA       WindowMomentCoefficient `json:"a_mass"`
	SquareMassB WindowMomentCoefficient `json:"b_square_mass"`
	DistanceJ   WindowMomentCoefficient `json:"j_distance_moment"`
}

func (c WindowMomentCoefficients) Validate() error {
	for _, x := range []WindowMomentCoefficient{c.MassA, c.SquareMassB, c.DistanceJ} {
		if err := x.Validate(); err != nil {
			return err
		}
	}
	if c.MassA.Parameter != c.SquareMassB.Parameter || c.MassA.Parameter != c.DistanceJ.Parameter {
		return fmt.Errorf("window coefficients use different parameters")
	}
	return nil
}

type ScaleChange struct {
	SourceObject          string           `json:"source_object"`
	TargetObject          string           `json:"target_object"`
	Factor                ScalarExpression `json:"factor"`
	ParameterDependencies []string         `json:"parameter_dependencies"`
	Provenance            Reference        `json:"provenance"`
}

func (s ScaleChange) Validate() error {
	if s.SourceObject == "" || s.TargetObject == "" || s.SourceObject == s.TargetObject || len(s.ParameterDependencies) == 0 || s.Provenance.Citation == "" {
		return fmt.Errorf("invalid scale change")
	}
	return s.Factor.Validate()
}

func (s ScaleChange) ApplyExact(kind SpectralMomentKind, value ExactRational, parameters map[string]ExactRational) (ExactRational, error) {
	if err := s.Validate(); err != nil {
		return ExactRational{}, err
	}
	if err := value.Validate(); err != nil {
		return ExactRational{}, err
	}
	f, err := s.Factor.EvalExact(parameters)
	if err != nil {
		return ExactRational{}, err
	}
	power := int64(1)
	if kind == FrobeniusNormSquared {
		power = 2
	} else if kind != Trace {
		return ExactRational{}, fmt.Errorf("M13 scale change only supports trace and Frobenius square")
	}
	scaled := rat(value)
	fr := rat(f)
	if power == 2 {
		fr.Mul(fr, fr)
	}
	scaled.Mul(scaled, fr)
	return exactRat(scaled)
}

type ScalarObjective struct {
	ID              string           `json:"id"`
	Parameter       WindowParameter  `json:"parameter"`
	ExactExpression ScalarExpression `json:"exact_expression"`
	DerivedFrom     []string         `json:"derived_from"`
	Provenance      Reference        `json:"provenance"`
}

func (o ScalarObjective) Validate() error {
	if o.ID == "" || len(o.DerivedFrom) < 2 || o.Provenance.Citation == "" {
		return fmt.Errorf("invalid scalar objective")
	}
	if err := o.Parameter.Validate(); err != nil {
		return err
	}
	return o.ExactExpression.Validate()
}

type RationalInterval struct {
	Lower ExactRational `json:"lower"`
	Upper ExactRational `json:"upper"`
}

func (i RationalInterval) Validate() error {
	if err := i.Lower.Validate(); err != nil {
		return err
	}
	if err := i.Upper.Validate(); err != nil {
		return err
	}
	if rat(i.Lower).Cmp(rat(i.Upper)) >= 0 {
		return fmt.Errorf("invalid rational interval")
	}
	return nil
}

type TaylorCertificate struct {
	VariableIdentity  string           `json:"variable_identity"`
	Terms             int              `json:"terms"`
	CosineInterval    RationalInterval `json:"cosine_interval"`
	SincInterval      RationalInterval `json:"sinc_interval"`
	XCotXInterval     RationalInterval `json:"x_cot_x_interval"`
	ObjectiveInterval RationalInterval `json:"objective_interval"`
	Method            string           `json:"method"`
}

// MontgomeryTaylorTaylorCertificate rigorously encloses
// J(1)=3/2-x*cot(x), x^2=1/2, using alternating rational series for
// cos(x) and sin(x)/x. No floating-point result enters the certificate.
func MontgomeryTaylorTaylorCertificate() (TaylorCertificate, error) {
	y := big.NewRat(1, 2)
	partial := func(sinc bool, n int) *big.Rat {
		s := new(big.Rat)
		pow := big.NewRat(1, 1)
		for k := 0; k <= n; k++ {
			factN := 2 * k
			if sinc {
				factN++
			}
			fact := big.NewInt(1)
			for j := 2; j <= factN; j++ {
				fact.Mul(fact, big.NewInt(int64(j)))
			}
			term := new(big.Rat).Quo(new(big.Rat).Set(pow), new(big.Rat).SetInt(fact))
			if k%2 == 0 {
				s.Add(s, term)
			} else {
				s.Sub(s, term)
			}
			pow.Mul(pow, y)
		}
		return s
	}
	// Odd truncations are lower bounds and even truncations are upper bounds.
	cl, cu := partial(false, 5), partial(false, 6)
	sl, su := partial(true, 5), partial(true, 6)
	ql := new(big.Rat).Quo(cl, su)
	qu := new(big.Rat).Quo(cu, sl)
	jl := new(big.Rat).Sub(big.NewRat(3, 2), qu)
	ju := new(big.Rat).Sub(big.NewRat(3, 2), ql)
	toInterval := func(a, b *big.Rat) (RationalInterval, error) {
		x, e := exactRat(a)
		if e != nil {
			return RationalInterval{}, e
		}
		z, e := exactRat(b)
		return RationalInterval{Lower: x, Upper: z}, e
	}
	ci, err := toInterval(cl, cu)
	if err != nil {
		return TaylorCertificate{}, err
	}
	si, err := toInterval(sl, su)
	if err != nil {
		return TaylorCertificate{}, err
	}
	qi, err := toInterval(ql, qu)
	if err != nil {
		return TaylorCertificate{}, err
	}
	ji, err := toInterval(jl, ju)
	if err != nil {
		return TaylorCertificate{}, err
	}
	return TaylorCertificate{VariableIdentity: "x=1/sqrt(2), so x^2=1/2", Terms: 7, CosineInterval: ci, SincInterval: si, XCotXInterval: qi, ObjectiveInterval: ji, Method: "alternating-series bounds: cos_lower/sinc_upper < x*cot(x) < cos_upper/sinc_lower"}, nil
}
