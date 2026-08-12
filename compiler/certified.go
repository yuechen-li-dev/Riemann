package compiler

import (
	"fmt"
	"math"
	"math/big"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

// m7Interval is a deliberately small, positive-transcendental interval kernel.
// Every operation fixes the result precision and rounds toward the relevant
// infinity. It is not intended to become a general interval-analysis package.
type m7Interval struct{ lo, hi *big.Float }

type m7Context struct {
	prec  uint
	terms int
}

func newM7Context(prec uint) m7Context { return m7Context{prec: prec, terms: 112} }

func (c m7Context) f(mode big.RoundingMode) *big.Float {
	return new(big.Float).SetPrec(c.prec).SetMode(mode)
}
func (c m7Context) integer(n int64) m7Interval {
	return m7Interval{c.f(big.ToNegativeInf).SetInt64(n), c.f(big.ToPositiveInf).SetInt64(n)}
}
func (c m7Context) rational(n, d int64) m7Interval {
	r := new(big.Rat).SetFrac(big.NewInt(n), big.NewInt(d))
	return m7Interval{c.f(big.ToNegativeInf).SetRat(r), c.f(big.ToPositiveInf).SetRat(r)}
}
func (c m7Context) decimal(lo, hi string) (m7Interval, error) {
	l, _, err := big.ParseFloat(lo, 10, c.prec, big.ToNegativeInf)
	if err != nil {
		return m7Interval{}, err
	}
	h, _, err := big.ParseFloat(hi, 10, c.prec, big.ToPositiveInf)
	if err != nil {
		return m7Interval{}, err
	}
	return m7Interval{l, h}, nil
}
func (c m7Context) add(a, b m7Interval) m7Interval {
	return m7Interval{c.f(big.ToNegativeInf).Add(a.lo, b.lo), c.f(big.ToPositiveInf).Add(a.hi, b.hi)}
}
func (c m7Context) neg(a m7Interval) m7Interval {
	return m7Interval{c.f(big.ToNegativeInf).Neg(a.hi), c.f(big.ToPositiveInf).Neg(a.lo)}
}
func (c m7Context) sub(a, b m7Interval) m7Interval { return c.add(a, c.neg(b)) }
func minFloat(xs []*big.Float) *big.Float {
	m := xs[0]
	for _, x := range xs[1:] {
		if x.Cmp(m) < 0 {
			m = x
		}
	}
	return m
}
func maxFloat(xs []*big.Float) *big.Float {
	m := xs[0]
	for _, x := range xs[1:] {
		if x.Cmp(m) > 0 {
			m = x
		}
	}
	return m
}
func (c m7Context) mul(a, b m7Interval) m7Interval {
	ld := []*big.Float{c.f(big.ToNegativeInf).Mul(a.lo, b.lo), c.f(big.ToNegativeInf).Mul(a.lo, b.hi), c.f(big.ToNegativeInf).Mul(a.hi, b.lo), c.f(big.ToNegativeInf).Mul(a.hi, b.hi)}
	hu := []*big.Float{c.f(big.ToPositiveInf).Mul(a.lo, b.lo), c.f(big.ToPositiveInf).Mul(a.lo, b.hi), c.f(big.ToPositiveInf).Mul(a.hi, b.lo), c.f(big.ToPositiveInf).Mul(a.hi, b.hi)}
	return m7Interval{c.f(big.ToNegativeInf).Set(minFloat(ld)), c.f(big.ToPositiveInf).Set(maxFloat(hu))}
}
func (c m7Context) reciprocal(a m7Interval) (m7Interval, error) {
	if a.lo.Sign() <= 0 && a.hi.Sign() >= 0 {
		return m7Interval{}, fmt.Errorf("interval reciprocal crosses zero")
	}
	one := c.integer(1)
	return m7Interval{c.f(big.ToNegativeInf).Quo(one.lo, a.hi), c.f(big.ToPositiveInf).Quo(one.hi, a.lo)}, nil
}
func (c m7Context) div(a, b m7Interval) (m7Interval, error) {
	r, e := c.reciprocal(b)
	if e != nil {
		return m7Interval{}, e
	}
	return c.mul(a, r), nil
}
func (c m7Context) hull(a, b m7Interval) m7Interval {
	return m7Interval{c.f(big.ToNegativeInf).Set(minFloat([]*big.Float{a.lo, b.lo})), c.f(big.ToPositiveInf).Set(maxFloat([]*big.Float{a.hi, b.hi}))}
}

func (c m7Context) expPoint(x *big.Float, upper bool) *big.Float {
	mode := big.ToNegativeInf
	if upper {
		mode = big.ToPositiveInf
	}
	one := c.f(mode).SetInt64(1)
	sum := c.f(mode).Set(one)
	term := c.f(mode).Set(one)
	for n := 1; n <= c.terms; n++ {
		term = c.f(mode).Quo(c.f(mode).Mul(term, x), c.f(mode).SetInt64(int64(n)))
		sum = c.f(mode).Add(sum, term)
	}
	if !upper {
		return sum
	}
	next := c.f(big.ToPositiveInf).Quo(c.f(big.ToPositiveInf).Mul(term, x), c.f(big.ToPositiveInf).SetInt64(int64(c.terms+1)))
	ratio := c.f(big.ToPositiveInf).Quo(x, c.f(big.ToPositiveInf).SetInt64(int64(c.terms+2)))
	den := c.f(big.ToNegativeInf).Sub(c.f(big.ToNegativeInf).SetInt64(1), ratio)
	rem := c.f(big.ToPositiveInf).Quo(next, den)
	return c.f(big.ToPositiveInf).Add(sum, rem)
}
func (c m7Context) exp(a m7Interval) (m7Interval, error) {
	if a.lo.Sign() < 0 {
		return m7Interval{}, fmt.Errorf("m7 exp expects nonnegative interval")
	}
	return m7Interval{c.expPoint(a.lo, false), c.expPoint(a.hi, true)}, nil
}

// phi(x)=expm1(x)/x, with phi(0)=1. Its positive Taylor coefficients
// make endpoint Taylor bounds monotone and cancellation-free.
func (c m7Context) phiPoint(x *big.Float, upper bool) *big.Float {
	mode := big.ToNegativeInf
	if upper {
		mode = big.ToPositiveInf
	}
	one := c.f(mode).SetInt64(1)
	sum := c.f(mode).Set(one)
	term := c.f(mode).Set(one)
	for k := 1; k <= c.terms; k++ {
		term = c.f(mode).Quo(c.f(mode).Mul(term, x), c.f(mode).SetInt64(int64(k+1)))
		sum = c.f(mode).Add(sum, term)
	}
	if !upper {
		return sum
	}
	next := c.f(big.ToPositiveInf).Quo(c.f(big.ToPositiveInf).Mul(term, x), c.f(big.ToPositiveInf).SetInt64(int64(c.terms+2)))
	ratio := c.f(big.ToPositiveInf).Quo(x, c.f(big.ToPositiveInf).SetInt64(int64(c.terms+3)))
	den := c.f(big.ToNegativeInf).Sub(c.f(big.ToNegativeInf).SetInt64(1), ratio)
	return c.f(big.ToPositiveInf).Add(sum, c.f(big.ToPositiveInf).Quo(next, den))
}
func (c m7Context) phi(a m7Interval) (m7Interval, error) {
	if a.lo.Sign() < 0 {
		return m7Interval{}, fmt.Errorf("phi expects nonnegative interval")
	}
	return m7Interval{c.phiPoint(a.lo, false), c.phiPoint(a.hi, true)}, nil
}

func (c m7Context) logPointGE1(x *big.Float, upper bool) *big.Float {
	mode := big.ToNegativeInf
	if upper {
		mode = big.ToPositiveInf
	}
	one := c.f(mode).SetInt64(1)
	// z=(x-1)/(x+1): denominator rounding is opposite numerator/result
	// rounding because division is antitone in its positive denominator.
	num := c.f(mode).Sub(x, one)
	denMode := big.ToPositiveInf
	if upper {
		denMode = big.ToNegativeInf
	}
	den := c.f(denMode).Add(x, c.f(denMode).SetInt64(1))
	z := c.f(mode).Quo(num, den)
	z2 := c.f(mode).Mul(z, z)
	power := c.f(mode).Set(z)
	sum := c.f(mode).SetInt64(0)
	for k := 0; k <= c.terms; k++ {
		term := c.f(mode).Quo(power, c.f(mode).SetInt64(int64(2*k+1)))
		sum = c.f(mode).Add(sum, term)
		power = c.f(mode).Mul(power, z2)
	}
	sum = c.f(mode).Mul(c.f(mode).SetInt64(2), sum)
	if !upper {
		return sum
	}
	den1 := c.f(big.ToNegativeInf).SetInt64(int64(2*c.terms + 3))
	den2 := c.f(big.ToNegativeInf).Sub(c.f(big.ToNegativeInf).SetInt64(1), z2)
	rem := c.f(big.ToPositiveInf).Quo(c.f(big.ToPositiveInf).Mul(c.f(big.ToPositiveInf).SetInt64(2), power), c.f(big.ToNegativeInf).Mul(den1, den2))
	return c.f(big.ToPositiveInf).Add(sum, rem)
}
func (c m7Context) logPoint(x *big.Float, upper bool) (*big.Float, error) {
	if x.Sign() <= 0 {
		return nil, fmt.Errorf("log domain")
	}
	one := c.f(big.ToNearestEven).SetInt64(1)
	if x.Cmp(one) >= 0 {
		return c.logPointGE1(x, upper), nil
	}
	mode := big.ToPositiveInf
	if upper {
		mode = big.ToNegativeInf
	}
	inv := c.f(mode).Quo(c.f(mode).SetInt64(1), x)
	y := c.logPointGE1(inv, !upper)
	outMode := big.ToPositiveInf
	if !upper {
		outMode = big.ToNegativeInf
	}
	return c.f(outMode).Neg(y), nil
}
func (c m7Context) log(a m7Interval) (m7Interval, error) {
	l, e := c.logPoint(a.lo, false)
	if e != nil {
		return m7Interval{}, e
	}
	h, e := c.logPoint(a.hi, true)
	if e != nil {
		return m7Interval{}, e
	}
	return m7Interval{l, h}, nil
}

func (c m7Context) logInt(n int) (m7Interval, error) { return c.logRational(int64(n), 1) }

// Range reduction keeps the atanh-series argument at most 1/3, giving a
// uniformly tiny, explicitly bounded remainder for finite-sum logarithms.
func (c m7Context) logRational(n, d int64) (m7Interval, error) {
	if n <= 0 || d <= 0 {
		return m7Interval{}, fmt.Errorf("log rational domain")
	}
	k := int64(0)
	for n >= 2*d {
		d *= 2
		k++
	}
	for n < d {
		n *= 2
		k--
	}
	m, err := c.log(c.rational(n, d))
	if err != nil {
		return m7Interval{}, err
	}
	if k == 0 {
		return m, nil
	}
	l2, err := c.log(c.rational(2, 1))
	if err != nil {
		return m7Interval{}, err
	}
	return c.add(m, c.mul(c.integer(k), l2)), nil
}

func (c m7Context) semanticInterval(x m7Interval) semantic.ComplexInterval {
	lo, _ := x.lo.Float64()
	hi, _ := x.hi.Float64()
	lf := new(big.Float).SetFloat64(lo)
	if lf.Cmp(x.lo) > 0 {
		lo = math.Nextafter(lo, math.Inf(-1))
	}
	hf := new(big.Float).SetFloat64(hi)
	if hf.Cmp(x.hi) < 0 {
		hi = math.Nextafter(hi, math.Inf(1))
	}
	return semantic.ComplexInterval{RealLower: lo, RealUpper: hi, RealLowerExact: x.lo.Text('g', 60), RealUpperExact: x.hi.Text('g', 60), Representation: fmt.Sprintf("math/big.Float directed binary, %d bits", c.prec)}
}

type m7ArchParts struct {
	finite, tail, total m7Interval
	panels              int
	breakpoints         []string
}

func (c m7Context) rRatio(u m7Interval) (m7Interval, error) {
	p15 := c.mul(c.rational(3, 2), u)
	p2 := c.mul(c.integer(2), u)
	a, e := c.phi(p15)
	if e != nil {
		return m7Interval{}, e
	}
	b, e := c.phi(p2)
	if e != nil {
		return m7Interval{}, e
	}
	q, e := c.div(a, b)
	if e != nil {
		return m7Interval{}, e
	}
	return c.mul(c.rational(3, 4), q), nil
}
func (c m7Context) sRatio(u m7Interval) (m7Interval, error) {
	p15 := c.mul(c.rational(3, 2), u)
	p2 := c.mul(c.integer(2), u)
	e15, e := c.exp(p15)
	if e != nil {
		return m7Interval{}, e
	}
	ph, e := c.phi(p2)
	if e != nil {
		return m7Interval{}, e
	}
	q, e := c.div(e15, ph)
	if e != nil {
		return m7Interval{}, e
	}
	return c.mul(c.rational(1, 2), q), nil
}
func (c m7Context) plateauRange(u, h0 m7Interval) (m7Interval, error) {
	r, e := c.rRatio(u)
	if e != nil {
		return m7Interval{}, e
	}
	return c.mul(c.mul(c.integer(2), h0), r), nil
}
func (c m7Context) diagonalSlopeRange(u, h0 m7Interval) (m7Interval, error) {
	p, e := c.plateauRange(u, h0)
	if e != nil {
		return m7Interval{}, e
	}
	s, e := c.sRatio(u)
	if e != nil {
		return m7Interval{}, e
	}
	return c.sub(p, c.mul(c.integer(2), s)), nil
}
func (c m7Context) offdiagSlopeRange(u, h0, d0 m7Interval) (m7Interval, error) {
	p, e := c.plateauRange(u, h0)
	if e != nil {
		return m7Interval{}, e
	}
	d := c.sub(u, d0)
	p15 := c.mul(c.rational(3, 2), u)
	ee, e := c.exp(p15)
	if e != nil {
		return m7Interval{}, e
	}
	p2 := c.mul(c.integer(2), u)
	ph, e := c.phi(p2)
	if e != nil {
		return m7Interval{}, e
	}
	den := c.mul(c.mul(c.integer(2), u), ph)
	q, e := c.div(c.mul(c.integer(2), c.mul(d, ee)), den)
	if e != nil {
		return m7Interval{}, e
	}
	return c.sub(p, q), nil
}
func (c m7Context) postRange(u, h0 m7Interval) (m7Interval, error) {
	p2 := c.mul(c.integer(2), u)
	ph, e := c.phi(p2)
	if e != nil {
		return m7Interval{}, e
	}
	den := c.mul(p2, ph)
	q, e := c.div(c.mul(c.integer(-2), h0), den)
	return q, e
}

func (c m7Context) integrateBoxes(a, b *big.Float, panels int, rangeFn func(m7Interval) (m7Interval, error)) (m7Interval, error) {
	if a.Cmp(b) >= 0 {
		return c.integer(0), nil
	}
	aExact := m7Interval{c.f(big.ToNegativeInf).Set(a), c.f(big.ToPositiveInf).Set(a)}
	bExact := m7Interval{c.f(big.ToNegativeInf).Set(b), c.f(big.ToPositiveInf).Set(b)}
	segment := c.sub(bExact, aExact)
	point := func(i int) m7Interval {
		return c.add(aExact, c.mul(segment, c.rational(int64(i), int64(panels))))
	}
	sum := c.integer(0)
	for i := 0; i < panels; i++ {
		left, right := point(i), point(i+1)
		cell := m7Interval{left.lo, right.hi}
		rng, e := rangeFn(cell)
		if e != nil {
			return m7Interval{}, e
		}
		// The exact partition points lie in left and right. Therefore the
		// exact cell width lies in this directed interval, while cell encloses
		// every point of the mathematical cell.
		width := c.sub(right, left)
		sum = c.add(sum, c.mul(rng, width))
	}
	return sum, nil
}

func (c m7Context) certifiedArchimedean(p logBoxPair, panels int) (m7ArchParts, error) {
	lnMin, e := c.logInt(minInt(p.q, p.r))
	if e != nil {
		return m7ArchParts{}, e
	}
	h0 := c.mul(c.integer(4), lnMin)
	u, e := c.logInt(p.q * p.r)
	if e != nil {
		return m7ArchParts{}, e
	}
	u = c.mul(c.integer(2), u)
	d0 := c.integer(0)
	diagonal := p.q == p.r
	if !diagonal {
		ratio := maxInt(p.q, p.r)
		minq := minInt(p.q, p.r)
		lr, e := c.logRational(int64(ratio), int64(minq))
		if e != nil {
			return m7ArchParts{}, e
		}
		d0 = c.mul(c.integer(2), lr)
	}
	b := c.integer(5)
	if u.hi.Cmp(b.lo) >= 0 {
		return m7ArchParts{}, fmt.Errorf("tail start does not exceed support")
	}
	finite := c.integer(0)
	parts := 0
	addSeg := func(a, z *big.Float, n int, fn func(m7Interval) (m7Interval, error)) error {
		x, e := c.integrateBoxes(a, z, n, fn)
		if e == nil {
			finite = c.add(finite, x)
			parts += n
		}
		return e
	}
	zero := c.integer(0)
	if diagonal {
		if e = addSeg(zero.lo, u.lo, panels, func(x m7Interval) (m7Interval, error) { return c.diagonalSlopeRange(x, h0) }); e != nil {
			return m7ArchParts{}, e
		}
		if e = addSeg(u.lo, u.hi, 1, func(x m7Interval) (m7Interval, error) {
			a, _ := c.diagonalSlopeRange(x, h0)
			z, _ := c.postRange(x, h0)
			return c.hull(a, z), nil
		}); e != nil {
			return m7ArchParts{}, e
		}
	} else {
		if e = addSeg(zero.lo, d0.lo, panels/2, func(x m7Interval) (m7Interval, error) { return c.plateauRange(x, h0) }); e != nil {
			return m7ArchParts{}, e
		}
		if e = addSeg(d0.lo, d0.hi, 1, func(x m7Interval) (m7Interval, error) {
			a, _ := c.plateauRange(x, h0)
			z, _ := c.offdiagSlopeRange(x, h0, d0)
			return c.hull(a, z), nil
		}); e != nil {
			return m7ArchParts{}, e
		}
		if e = addSeg(d0.hi, u.lo, panels, func(x m7Interval) (m7Interval, error) { return c.offdiagSlopeRange(x, h0, d0) }); e != nil {
			return m7ArchParts{}, e
		}
		if e = addSeg(u.lo, u.hi, 1, func(x m7Interval) (m7Interval, error) {
			a, _ := c.offdiagSlopeRange(x, h0, d0)
			z, _ := c.postRange(x, h0)
			return c.hull(a, z), nil
		}); e != nil {
			return m7ArchParts{}, e
		}
	}
	if e = addSeg(u.hi, b.lo, panels/4, func(x m7Interval) (m7Interval, error) { return c.postRange(x, h0) }); e != nil {
		return m7ArchParts{}, e
	}
	// Integral from B to infinity is h0*log(1-exp(-2B)).
	e10, e := c.exp(c.integer(10))
	if e != nil {
		return m7ArchParts{}, e
	}
	inv, e := c.reciprocal(e10)
	if e != nil {
		return m7ArchParts{}, e
	}
	oneMinus := c.sub(c.integer(1), inv)
	lg, e := c.log(oneMinus)
	if e != nil {
		return m7ArchParts{}, e
	}
	tail := c.mul(h0, lg)
	// DLMF 5.2.3 and 3.12.1 print these decimal prefixes. Truncation and
	// adding one unit in the final printed place give exact decimal brackets.
	gamma, e := c.decimal("0.57721566490153286060", "0.57721566490153286061")
	if e != nil {
		return m7ArchParts{}, e
	}
	pi, e := c.decimal("3.14159265358979323846", "3.14159265358979323847")
	if e != nil {
		return m7ArchParts{}, e
	}
	logPi, e := c.log(pi)
	if e != nil {
		return m7ArchParts{}, e
	}
	total := c.add(c.mul(c.add(gamma, logPi), h0), c.add(finite, tail))
	return m7ArchParts{finite: finite, tail: tail, total: total, panels: parts, breakpoints: []string{"0", d0.lo.Text('g', 30) + ".." + d0.hi.Text('g', 30), u.lo.Text('g', 30) + ".." + u.hi.Text('g', 30), "5", "infinity"}}, nil
}

func (c m7Context) certifiedPrime(p logBoxPair) (m7Interval, int, error) {
	limit := p.q * p.r * p.q * p.r
	sum := c.integer(0)
	terms := 0
	mn, mx := minInt(p.q, p.r), maxInt(p.q, p.r)
	for _, prime := range primesThrough(limit) {
		for power := prime; power <= limit; {
			if power < limit {
				lp, e := c.logInt(prime)
				if e != nil {
					return m7Interval{}, 0, e
				}
				var h m7Interval
				if power*mn*mn <= mx*mx {
					lm, e := c.logInt(mn)
					if e != nil {
						return m7Interval{}, 0, e
					}
					h = c.mul(c.integer(4), lm)
				} else {
					lr, e := c.logRational(int64(limit), int64(power))
					if e != nil {
						return m7Interval{}, 0, e
					}
					h = lr
				}
				lk, e := c.logInt(power)
				if e != nil {
					return m7Interval{}, 0, e
				}
				root, e := c.exp(c.mul(c.rational(1, 2), lk))
				if e != nil {
					return m7Interval{}, 0, e
				}
				inv, e := c.reciprocal(root)
				if e != nil {
					return m7Interval{}, 0, e
				}
				sum = c.add(sum, c.mul(c.integer(2), c.mul(lp, c.mul(inv, h))))
				terms++
			}
			if power > limit/prime {
				break
			}
			power *= prime
		}
	}
	return sum, terms, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
