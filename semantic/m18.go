package semantic

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
)

// PiTerm and PiExpression are a deliberately narrow exact scalar language.
// They represent finite Laurent polynomials in pi; decimal approximations never
// participate in M18 theorem checking.
type PiTerm struct {
	Power       int           `json:"pi_power"`
	Coefficient ExactRational `json:"coefficient"`
}

type PiExpression struct {
	Expression string   `json:"expression"`
	Terms      []PiTerm `json:"laurent_terms"`
}

func piExpr(display string, coefficients map[int]*big.Rat) PiExpression {
	powers := make([]int, 0, len(coefficients))
	for power, coefficient := range coefficients {
		if coefficient.Sign() != 0 {
			powers = append(powers, power)
		}
	}
	sort.Ints(powers)
	terms := make([]PiTerm, 0, len(powers))
	for _, power := range powers {
		coefficient, err := exactRat(coefficients[power])
		if err != nil {
			panic(err)
		}
		terms = append(terms, PiTerm{Power: power, Coefficient: coefficient})
	}
	return PiExpression{Expression: display, Terms: terms}
}

func piConstant(q *big.Rat) PiExpression { return piExpr(q.RatString(), map[int]*big.Rat{0: q}) }
func piMonomial(power int, q *big.Rat) PiExpression {
	return piExpr("", map[int]*big.Rat{power: q})
}

func piCoefficients(e PiExpression) map[int]*big.Rat {
	out := make(map[int]*big.Rat, len(e.Terms))
	for _, term := range e.Terms {
		out[term.Power] = rat(term.Coefficient)
	}
	return out
}

func piAdd(a, b PiExpression) PiExpression {
	out := piCoefficients(a)
	for power, coefficient := range piCoefficients(b) {
		if out[power] == nil {
			out[power] = new(big.Rat)
		}
		out[power].Add(out[power], coefficient)
	}
	return piExpr("", out)
}

func piScale(a PiExpression, q *big.Rat) PiExpression {
	out := piCoefficients(a)
	for power := range out {
		out[power].Mul(out[power], q)
	}
	return piExpr("", out)
}

func piMul(a, b PiExpression) PiExpression {
	out := map[int]*big.Rat{}
	for ap, ac := range piCoefficients(a) {
		for bp, bc := range piCoefficients(b) {
			power := ap + bp
			if out[power] == nil {
				out[power] = new(big.Rat)
			}
			out[power].Add(out[power], new(big.Rat).Mul(ac, bc))
		}
	}
	return piExpr("", out)
}

func samePiValue(a, b PiExpression) bool {
	d := piAdd(a, piScale(b, big.NewRat(-1, 1)))
	return len(d.Terms) == 0
}

// OneRadiusTangencyCandidate is exact contact data, not a decimalized M17
// rational witness and not a general transcendental expression system.
type OneRadiusTangencyCandidate struct {
	Radius              ExactRational `json:"r"`
	ContactFrequency    PiExpression  `json:"contact_frequency"`
	ExactCExpression    PiExpression  `json:"c"`
	ExactWExpression    PiExpression  `json:"w"`
	ContactMultiplicity int           `json:"contact_multiplicity"`
	Transform           string        `json:"transform"`
	Derivation          string        `json:"derivation"`
}

// DeriveOneRadiusTangencyCandidate solves F(a)=F'(a)=0 at a=pi/2.
// With F=1+2(cos(t)-1)/t^2+2(1-c)sin(t)/t+2w cos(t), the first equation is
// c=1+a/2-1/a.  The derivative equation is w=(c-2)/a^2+2/a^3.
func DeriveOneRadiusTangencyCandidate() OneRadiusTangencyCandidate {
	one := piConstant(big.NewRat(1, 1))
	a := piMonomial(1, big.NewRat(1, 2))
	invA := piMonomial(-1, big.NewRat(2, 1))
	c := piAdd(piAdd(one, piScale(a, big.NewRat(1, 2))), piScale(invA, big.NewRat(-1, 1)))
	c.Expression = "1+pi/4-2/pi"
	invA2 := piMonomial(-2, big.NewRat(4, 1))
	invA3 := piMonomial(-3, big.NewRat(8, 1))
	w := piAdd(piMul(piAdd(c, piConstant(big.NewRat(-2, 1))), invA2), piScale(invA3, big.NewRat(2, 1)))
	w.Expression = "1/pi-4/pi^2+8/pi^3"
	a.Expression = "pi/2"
	return OneRadiusTangencyCandidate{
		Radius: ExactRational{1, 1}, ContactFrequency: a, ExactCExpression: c, ExactWExpression: w,
		ContactMultiplicity: 2, Transform: OneRadiusFourierDensity,
		Derivation: "F(pi/2)=0 gives c=1+pi/4-2/pi; substituting in F'(pi/2)=0 gives w=1/pi-4/pi^2+8/pi^3",
	}
}

type rationalInterval struct{ lo, hi *big.Rat }

func ri(lo, hi *big.Rat) rationalInterval {
	return rationalInterval{new(big.Rat).Set(lo), new(big.Rat).Set(hi)}
}
func riPoint(q *big.Rat) rationalInterval { return ri(q, q) }
func riAdd(a, b rationalInterval) rationalInterval {
	return ri(new(big.Rat).Add(a.lo, b.lo), new(big.Rat).Add(a.hi, b.hi))
}
func riNeg(a rationalInterval) rationalInterval {
	return ri(new(big.Rat).Neg(a.hi), new(big.Rat).Neg(a.lo))
}
func riSub(a, b rationalInterval) rationalInterval { return riAdd(a, riNeg(b)) }
func riMul(a, b rationalInterval) rationalInterval {
	values := []*big.Rat{
		new(big.Rat).Mul(a.lo, b.lo), new(big.Rat).Mul(a.lo, b.hi),
		new(big.Rat).Mul(a.hi, b.lo), new(big.Rat).Mul(a.hi, b.hi),
	}
	lo, hi := values[0], values[0]
	for _, value := range values[1:] {
		if value.Cmp(lo) < 0 {
			lo = value
		}
		if value.Cmp(hi) > 0 {
			hi = value
		}
	}
	return ri(lo, hi)
}
func riReciprocal(a rationalInterval) rationalInterval {
	if a.lo.Sign() <= 0 && a.hi.Sign() >= 0 {
		panic("interval reciprocal crosses zero")
	}
	return ri(new(big.Rat).Inv(a.hi), new(big.Rat).Inv(a.lo))
}
func riPow(a rationalInterval, n int) rationalInterval {
	if n < 0 {
		return riPow(riReciprocal(a), -n)
	}
	out := riPoint(big.NewRat(1, 1))
	for k := 0; k < n; k++ {
		out = riMul(out, a)
	}
	return out
}

var m18Pi = rationalInterval{
	lo: big.NewRat(103993, 33102), // classical continued-fraction lower bound
	hi: big.NewRat(104348, 33215), // classical continued-fraction upper bound
}

func evalPi(e PiExpression) rationalInterval {
	out := riPoint(big.NewRat(0, 1))
	for _, term := range e.Terms {
		out = riAdd(out, riMul(riPoint(rat(term.Coefficient)), riPow(m18Pi, term.Power)))
	}
	return out
}

// VerifyPiExpressionRationalBounds exposes only the comparison needed by
// reports and artifacts; the large internal rational endpoints stay out of the
// int64-based ExactRational IR.
func VerifyPiExpressionRationalBounds(e PiExpression, lower, upper ExactRational) bool {
	if lower.Validate() != nil || upper.Validate() != nil {
		return false
	}
	b := evalPi(e)
	return b.lo.Cmp(rat(lower)) > 0 && b.hi.Cmp(rat(upper)) < 0
}

func (e PiExpression) Validate() error {
	if strings.TrimSpace(e.Expression) == "" || len(e.Terms) == 0 {
		return fmt.Errorf("empty exact pi expression")
	}
	last := -1000000
	for _, term := range e.Terms {
		if err := term.Coefficient.Validate(); err != nil {
			return err
		}
		if term.Power <= last || rat(term.Coefficient).Sign() == 0 {
			return fmt.Errorf("noncanonical pi Laurent expression")
		}
		last = term.Power
	}
	return nil
}

func (c OneRadiusTangencyCandidate) Validate() error {
	if c.Radius != (ExactRational{1, 1}) || c.ContactMultiplicity != 2 || c.Transform != OneRadiusFourierDensity {
		return fmt.Errorf("wrong M18 contact semantics")
	}
	for _, e := range []PiExpression{c.ContactFrequency, c.ExactCExpression, c.ExactWExpression} {
		if err := e.Validate(); err != nil {
			return err
		}
	}
	want := DeriveOneRadiusTangencyCandidate()
	if !samePiValue(c.ContactFrequency, want.ContactFrequency) || !samePiValue(c.ExactCExpression, want.ExactCExpression) || !samePiValue(c.ExactWExpression, want.ExactWExpression) {
		return fmt.Errorf("candidate does not solve the derived tangency system")
	}
	if evalPi(c.ExactCExpression).lo.Sign() <= 0 || evalPi(c.ExactWExpression).lo.Sign() < 0 {
		return fmt.Errorf("candidate parameter is not certified legal")
	}
	return nil
}

type OneRadiusExactTaylor struct {
	Constant  PiExpression `json:"constant_term"`
	Quadratic PiExpression `json:"t_squared_coefficient"`
	Quartic   PiExpression `json:"t_fourth_coefficient"`
	Sextic    PiExpression `json:"t_sixth_coefficient"`
}

func TangencyOriginTaylor(c OneRadiusTangencyCandidate) OneRadiusExactTaylor {
	one := piConstant(big.NewRat(1, 1))
	wMinusCPlusOne := piAdd(piAdd(one, c.ExactWExpression), piScale(c.ExactCExpression, big.NewRat(-1, 1)))
	a0 := piScale(wMinusCPlusOne, big.NewRat(2, 1))
	a0.Expression = "2*(1+w-c)"
	a2 := piAdd(piAdd(piScale(c.ExactCExpression, big.NewRat(1, 3)), piConstant(big.NewRat(-1, 4))), piScale(c.ExactWExpression, big.NewRat(-1, 1)))
	a2.Expression = "c/3-1/4-w"
	a4 := piAdd(piAdd(piConstant(big.NewRat(-1, 360)), piScale(piAdd(one, piScale(c.ExactCExpression, big.NewRat(-1, 1))), big.NewRat(1, 60))), piScale(c.ExactWExpression, big.NewRat(1, 12)))
	a4.Expression = "-1/360+(1-c)/60+w/12"
	a6 := piAdd(piAdd(piConstant(big.NewRat(1, 20160)), piScale(piAdd(c.ExactCExpression, piConstant(big.NewRat(-1, 1))), big.NewRat(1, 2520))), piScale(c.ExactWExpression, big.NewRat(-1, 360)))
	a6.Expression = "1/20160+(c-1)/2520-w/360"
	return OneRadiusExactTaylor{a0, a2, a4, a6}
}

func DeriveTangencySimpleUpper(c OneRadiusTangencyCandidate) PiExpression {
	j := piAdd(piConstant(big.NewRat(2, 1)), piScale(c.ExactCExpression, big.NewRat(-1, 1)))
	j.Expression = "1-pi/4+2/pi"
	return j
}

type ContactEquationCertificate struct {
	ValueEquation              string        `json:"value_equation"`
	DerivativeEquation         string        `json:"derivative_equation"`
	SecondDerivativeExpression PiExpression  `json:"second_derivative"`
	SecondDerivativeLower      ExactRational `json:"second_derivative_lower"`
	ThirdDerivativeAbsUpper    ExactRational `json:"third_derivative_abs_upper"`
	Verified                   bool          `json:"verified"`
}

func deriveContactSecond(c OneRadiusTangencyCandidate) PiExpression {
	// F''(a)=8/a^3-12/a^4+2(1-c)(-1/a+2/a^3), a=pi/2.
	first := piAdd(piMonomial(-3, big.NewRat(64, 1)), piMonomial(-4, big.NewRat(-192, 1)))
	secondFactor := piAdd(piMonomial(-1, big.NewRat(-2, 1)), piMonomial(-3, big.NewRat(16, 1)))
	second := piScale(piMul(piAdd(piConstant(big.NewRat(1, 1)), piScale(c.ExactCExpression, big.NewRat(-1, 1))), secondFactor), big.NewRat(2, 1))
	out := piAdd(first, second)
	out.Expression = "8/a^3-12/a^4+2*(1-c)*(-1/a+2/a^3), a=pi/2"
	return out
}

func VerifyTangencyAlgebra(c OneRadiusTangencyCandidate) (ContactEquationCertificate, error) {
	if err := c.Validate(); err != nil {
		return ContactEquationCertificate{}, err
	}
	one := piConstant(big.NewRat(1, 1))
	a := c.ContactFrequency
	invA := piMonomial(-1, big.NewRat(2, 1))
	invA2 := piMonomial(-2, big.NewRat(4, 1))
	invA3 := piMonomial(-3, big.NewRat(8, 1))
	value := piAdd(piAdd(one, piScale(invA2, big.NewRat(-2, 1))), piScale(piMul(piAdd(one, piScale(c.ExactCExpression, big.NewRat(-1, 1))), invA), big.NewRat(2, 1)))
	derivative := piAdd(piAdd(piScale(invA2, big.NewRat(-2, 1)), piScale(invA3, big.NewRat(4, 1))), piScale(piMul(piAdd(c.ExactCExpression, piConstant(big.NewRat(-1, 1))), invA2), big.NewRat(2, 1)))
	derivative = piAdd(derivative, piScale(c.ExactWExpression, big.NewRat(-2, 1)))
	if len(value.Terms) != 0 || len(derivative.Terms) != 0 {
		return ContactEquationCertificate{}, fmt.Errorf("exact contact equations do not vanish")
	}
	second := deriveContactSecond(c)
	if evalPi(second).lo.Cmp(big.NewRat(1, 8)) <= 0 {
		return ContactEquationCertificate{}, fmt.Errorf("quadratic contact lower bound failed")
	}
	// From F'''=-2 integral_0^1 x^3(c-x)sin(tx)dx+2w sin(t).
	m := riAdd(riSub(riMul(evalPi(c.ExactCExpression), riPoint(big.NewRat(1, 2))), riPoint(big.NewRat(2, 5))), riMul(evalPi(c.ExactWExpression), riPoint(big.NewRat(2, 1))))
	if m.hi.Cmp(big.NewRat(11, 20)) >= 0 {
		return ContactEquationCertificate{}, fmt.Errorf("third derivative bound failed")
	}
	_ = a
	return ContactEquationCertificate{
		ValueEquation:              "F(pi/2)=1-8/pi^2+4*(1-c)/pi=0",
		DerivativeEquation:         "F'(pi/2)=-8/pi^2+32/pi^3+8*(c-1)/pi^2-2*w=0",
		SecondDerivativeExpression: second, SecondDerivativeLower: ExactRational{1, 8}, ThirdDerivativeAbsUpper: ExactRational{11, 20}, Verified: true,
	}, nil
}

type CompactCellKind string

const (
	StrictPositiveCell                  CompactCellKind = "strict_positive"
	NonnegativeCellWithCertifiedContact CompactCellKind = "nonnegative_with_certified_contact"
)

type OneRadiusTangencyWholeLineCertificate struct {
	ID                  string          `json:"id"`
	PiLower             ExactRational   `json:"pi_lower"`
	PiUpper             ExactRational   `json:"pi_upper"`
	LocalInterval       string          `json:"local_interval"`
	LocalCellKind       CompactCellKind `json:"local_cell_kind"`
	LocalConvexityLower ExactRational   `json:"local_convexity_lower"`
	CompactIntervals    []string        `json:"strict_compact_intervals"`
	CompactCellKind     CompactCellKind `json:"compact_cell_kind"`
	GridStep            ExactRational   `json:"grid_step"`
	TaylorDegree        int             `json:"taylor_degree"`
	LipschitzUpper      ExactRational   `json:"lipschitz_upper"`
	TailInterval        string          `json:"tail_interval"`
	TailLowerBound      ExactRational   `json:"tail_lower_bound"`
	Coverage            string          `json:"coverage"`
	WholeLine           bool            `json:"whole_line"`
}

func candidatePointLower(c OneRadiusTangencyCandidate, t *big.Rat, degree int) *big.Rat {
	if t.Sign() == 0 {
		return evalPi(TangencyOriginTaylor(c).Constant).lo
	}
	cosLo, cosHi := trigTaylorBounds(t, degree, true)
	sinLo, sinHi := trigTaylorBounds(t, degree-1, false)
	cosI, sinI := ri(cosLo, cosHi), ri(sinLo, sinHi)
	tI := riPoint(t)
	t2 := riPoint(new(big.Rat).Mul(t, t))
	value := riPoint(big.NewRat(1, 1))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(riSub(cosI, riPoint(big.NewRat(1, 1))), riReciprocal(t2))))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(riMul(riSub(riPoint(big.NewRat(1, 1)), evalPi(c.ExactCExpression)), sinI), riReciprocal(tI))))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(evalPi(c.ExactWExpression), cosI)))
	return value.lo
}

func VerifyTangencyWholeLineCertificate(c OneRadiusTangencyCandidate, p OneRadiusTangencyWholeLineCertificate) error {
	contact, err := VerifyTangencyAlgebra(c)
	if err != nil {
		return err
	}
	if p.ID == "" || p.PiLower != (ExactRational{103993, 33102}) || p.PiUpper != (ExactRational{104348, 33215}) || p.LocalInterval != "3/2<=t<=17/10" || p.LocalCellKind != NonnegativeCellWithCertifiedContact || p.CompactCellKind != StrictPositiveCell || p.TailInterval != "|t|>=4" || !p.WholeLine {
		return fmt.Errorf("incomplete M18 equality-aware certificate")
	}
	if len(p.CompactIntervals) != 2 || p.CompactIntervals[0] != "0<=t<=3/2" || p.CompactIntervals[1] != "17/10<=t<=4" || p.GridStep != (ExactRational{1, 2500}) || p.TaylorDegree != 40 || p.LipschitzUpper != (ExactRational{83, 100}) {
		return fmt.Errorf("wrong M18 compact partition")
	}
	// pi/2 lies in (157/100,1571/1000), hence every point of the local
	// interval is within 13/100.  Taylor's theorem for F'' gives strict convexity.
	if new(big.Rat).Quo(m18Pi.lo, big.NewRat(2, 1)).Cmp(big.NewRat(157, 100)) <= 0 || new(big.Rat).Quo(m18Pi.hi, big.NewRat(2, 1)).Cmp(big.NewRat(1571, 1000)) >= 0 {
		return fmt.Errorf("pi bounds do not locate the contact inside the local distance proof")
	}
	localLower := new(big.Rat).Sub(rat(contact.SecondDerivativeLower), new(big.Rat).Mul(rat(contact.ThirdDerivativeAbsUpper), big.NewRat(13, 100)))
	if localLower.Sign() <= 0 || rat(p.LocalConvexityLower).Cmp(localLower) != 0 {
		return fmt.Errorf("local convexity certificate failed")
	}
	cI, wI := evalPi(c.ExactCExpression), evalPi(c.ExactWExpression)
	if cI.lo.Cmp(big.NewRat(1, 1)) <= 0 || evalPi(TangencyOriginTaylor(c).Constant).lo.Cmp(big.NewRat(1, 25)) <= 0 {
		return fmt.Errorf("origin or c>1 margin failed")
	}
	lipschitz := riAdd(riSub(cI, riPoint(big.NewRat(2, 3))), riMul(wI, riPoint(big.NewRat(2, 1))))
	if lipschitz.hi.Cmp(rat(p.LipschitzUpper)) >= 0 {
		return fmt.Errorf("declared Lipschitz upper bound is too small")
	}
	loss := new(big.Rat).Mul(rat(p.LipschitzUpper), new(big.Rat).Quo(rat(p.GridStep), big.NewRat(2, 1)))
	for k := int64(0); k <= 3750; k++ { // [0,3/2]
		t := new(big.Rat).Mul(big.NewRat(k, 1), rat(p.GridStep))
		if candidatePointLower(c, t, p.TaylorDegree).Cmp(loss) <= 0 {
			return fmt.Errorf("left strict compact cell fails at %s", t.RatString())
		}
	}
	for k := int64(4250); k <= 10000; k++ { // [17/10,4]
		t := new(big.Rat).Mul(big.NewRat(k, 1), rat(p.GridStep))
		if candidatePointLower(c, t, p.TaylorDegree).Cmp(loss) <= 0 {
			return fmt.Errorf("right strict compact cell fails at %s", t.RatString())
		}
	}
	tail := riSub(riSub(riPoint(big.NewRat(3, 4)), riMul(wI, riPoint(big.NewRat(2, 1)))), riMul(riSub(cI, riPoint(big.NewRat(1, 1))), riPoint(big.NewRat(1, 2))))
	if tail.lo.Cmp(rat(p.TailLowerBound)) <= 0 || p.TailLowerBound != (ExactRational{1, 3}) {
		return fmt.Errorf("tail lower bound failed")
	}
	if p.Coverage != "evenness; strict Taylor/Lipschitz cells on [0,3/2] and [17/10,4]; certified quadratic contact on [3/2,17/10]; strict analytic tail on [4,infinity)" {
		return fmt.Errorf("whole-line coverage has a gap")
	}
	return nil
}

type ExactOneRadiusOptimum struct {
	Family          string       `json:"family"`
	ExactValue      PiExpression `json:"exact_value"`
	LowerRoute      string       `json:"lower_route"`
	UpperRoute      string       `json:"upper_route"`
	EqualityDerived bool         `json:"equality_derived"`
}

func ComposeExactOneRadiusOptimum(c OneRadiusTangencyCandidate, p OneRadiusTangencyWholeLineCertificate, upper OneRadiusFamilyCeiling) (ExactOneRadiusOptimum, error) {
	if err := VerifyTangencyWholeLineCertificate(c, p); err != nil {
		return ExactOneRadiusOptimum{}, err
	}
	if upper != DeriveOneRadiusFamilyCeiling() || upper.ExactUpperExpression != c.ExactCExpression.Expression {
		return ExactOneRadiusOptimum{}, fmt.Errorf("M17 upper theorem does not match M18 witness")
	}
	return ExactOneRadiusOptimum{Family: "UnsaturatedOneRadius", ExactValue: c.ExactCExpression, LowerRoute: "M18 globally nonnegative exact r=1 tangency witness", UpperRoute: "M17 one-radius frequency envelope", EqualityDerived: true}, nil
}

func RejectAboveExactOneRadiusCeiling(delta ExactRational) error {
	if err := delta.Validate(); err != nil {
		return err
	}
	if rat(delta).Sign() <= 0 {
		return fmt.Errorf("perturbation must be positive")
	}
	return fmt.Errorf("at r=1 and c=c_M18+%s, F(pi/2)=-4*(%s)/pi<0 independently of w", rat(delta).RatString(), rat(delta).RatString())
}
