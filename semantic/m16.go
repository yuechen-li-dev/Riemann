package semantic

import (
	"fmt"
	"math/big"
)

// SymmetricExteriorAtom is one nonnegative pair w(delta_-r+delta_r).
type SymmetricExteriorAtom struct {
	Radius ExactRational `json:"radius"`
	Weight ExactRational `json:"weight"`
}

// TwoRadiusCompletion is the exact M16 family
// sigma=c 1_{|x|>1}dx+sum_j w_j(delta_-r_j+delta_r_j).
type TwoRadiusCompletion struct {
	Constant           ExactRational           `json:"constant"`
	ExteriorDensity    ExactRational           `json:"exterior_density"`
	Atoms              []SymmetricExteriorAtom `json:"atoms"`
	DensitySupport     string                  `json:"density_support"`
	SupportLegal       bool                    `json:"support_legal"`
	MeasureNonnegative bool                    `json:"measure_nonnegative"`
}

func (f TwoRadiusCompletion) Canonical() TwoRadiusCompletion {
	out := f
	out.Atoms = append([]SymmetricExteriorAtom(nil), f.Atoms...)
	if len(out.Atoms) == 2 && rat(out.Atoms[0].Radius).Cmp(rat(out.Atoms[1].Radius)) > 0 {
		out.Atoms[0], out.Atoms[1] = out.Atoms[1], out.Atoms[0]
	}
	return out
}

func (f TwoRadiusCompletion) Validate() error {
	if err := f.Constant.Validate(); err != nil {
		return err
	}
	if err := f.ExteriorDensity.Validate(); err != nil {
		return err
	}
	if rat(f.Constant).Sign() <= 0 || rat(f.ExteriorDensity).Cmp(rat(f.Constant)) != 0 {
		return fmt.Errorf("exterior density must be the positive constant c that cancels -c dx")
	}
	if len(f.Atoms) != 2 {
		return fmt.Errorf("M16 requires exactly two typed atom radii")
	}
	for _, a := range f.Atoms {
		if err := a.Radius.Validate(); err != nil {
			return err
		}
		if err := a.Weight.Validate(); err != nil {
			return err
		}
		if rat(a.Radius).Cmp(big.NewRat(1, 1)) < 0 {
			return fmt.Errorf("atom radius lies inside (-1,1)")
		}
		if rat(a.Weight).Sign() < 0 {
			return fmt.Errorf("atom weight must be nonnegative")
		}
	}
	if rat(f.Atoms[0].Radius).Cmp(rat(f.Atoms[1].Radius)) > 0 {
		return fmt.Errorf("atom radii are not canonical")
	}
	if f.DensitySupport != "|x|>1" || !f.SupportLegal {
		return fmt.Errorf("exterior support is not certified")
	}
	if !f.MeasureNonnegative {
		return fmt.Errorf("sigma is not certified nonnegative")
	}
	return nil
}

// CollapsedAtoms returns the number of distinct radii and their total weights.
func (f TwoRadiusCompletion) CollapsedAtoms() []SymmetricExteriorAtom {
	c := f.Canonical()
	if len(c.Atoms) == 2 && rat(c.Atoms[0].Radius).Cmp(rat(c.Atoms[1].Radius)) == 0 {
		w, _ := exactRat(new(big.Rat).Add(rat(c.Atoms[0].Weight), rat(c.Atoms[1].Weight)))
		return []SymmetricExteriorAtom{{Radius: c.Atoms[0].Radius, Weight: w}}
	}
	return append([]SymmetricExteriorAtom(nil), c.Atoms...)
}

type TwoRadiusTaylor struct {
	Constant  ExactRational `json:"constant_term"`
	Quadratic ExactRational `json:"t_squared_coefficient"`
	Quartic   ExactRational `json:"t_fourth_coefficient"`
	Sextic    ExactRational `json:"t_sixth_coefficient"`
}

func TwoRadiusOriginTaylor(f TwoRadiusCompletion) (TwoRadiusTaylor, error) {
	if err := f.Validate(); err != nil {
		return TwoRadiusTaylor{}, err
	}
	c := rat(f.Constant)
	wSum, wr2, wr4, wr6 := new(big.Rat), new(big.Rat), new(big.Rat), new(big.Rat)
	for _, a := range f.Atoms {
		r, w := rat(a.Radius), rat(a.Weight)
		r2 := new(big.Rat).Mul(r, r)
		r4 := new(big.Rat).Mul(r2, r2)
		r6 := new(big.Rat).Mul(r4, r2)
		wSum.Add(wSum, w)
		wr2.Add(wr2, new(big.Rat).Mul(w, r2))
		wr4.Add(wr4, new(big.Rat).Mul(w, r4))
		wr6.Add(wr6, new(big.Rat).Mul(w, r6))
	}
	a0 := new(big.Rat).Mul(big.NewRat(2, 1), new(big.Rat).Sub(new(big.Rat).Add(big.NewRat(1, 1), wSum), c))
	a2 := new(big.Rat).Sub(new(big.Rat).Sub(new(big.Rat).Quo(c, big.NewRat(3, 1)), big.NewRat(1, 4)), wr2)
	a4 := new(big.Rat).Add(new(big.Rat).Add(big.NewRat(-1, 360), new(big.Rat).Quo(new(big.Rat).Sub(big.NewRat(1, 1), c), big.NewRat(60, 1))), new(big.Rat).Quo(wr4, big.NewRat(12, 1)))
	a6 := new(big.Rat).Sub(new(big.Rat).Add(big.NewRat(1, 20160), new(big.Rat).Quo(new(big.Rat).Sub(c, big.NewRat(1, 1)), big.NewRat(2520, 1))), new(big.Rat).Quo(wr6, big.NewRat(360, 1)))
	v0, _ := exactRat(a0)
	v2, _ := exactRat(a2)
	v4, _ := exactRat(a4)
	v6, _ := exactRat(a6)
	return TwoRadiusTaylor{Constant: v0, Quadratic: v2, Quartic: v4, Sextic: v6}, nil
}

func VerifyTwoRadiusLocalLegality(f TwoRadiusCompletion) error {
	t, err := TwoRadiusOriginTaylor(f)
	if err != nil {
		return err
	}
	for i, coefficient := range []ExactRational{t.Constant, t.Quadratic, t.Quartic, t.Sextic} {
		s := rat(coefficient).Sign()
		if s > 0 {
			return nil
		}
		if s < 0 {
			return fmt.Errorf("negative first nonzero Taylor coefficient at index %d", i)
		}
	}
	return fmt.Errorf("Taylor coefficients through t^6 vanish; higher-order analysis is required")
}

// SaturatedTwoRadiusCeiling proves the origin-saturated branch has c<=9/8:
// sum w=c-1 and r>=1 imply sum(w r^2)>=c-1, so a2<=(9-8c)/12.
func SaturatedTwoRadiusCeiling(f TwoRadiusCompletion) error {
	if err := f.Validate(); err != nil {
		return err
	}
	w := new(big.Rat).Add(rat(f.Atoms[0].Weight), rat(f.Atoms[1].Weight))
	if w.Cmp(new(big.Rat).Sub(rat(f.Constant), big.NewRat(1, 1))) != 0 {
		return fmt.Errorf("branch is not origin-saturated")
	}
	if rat(f.Constant).Cmp(big.NewRat(9, 8)) > 0 {
		return fmt.Errorf("origin-saturated two-radius branch exceeds its 9/8 ceiling")
	}
	return nil
}

type TwoRadiusWholeLineCertificate struct {
	ID                string        `json:"id"`
	FourierVariable   string        `json:"fourier_variable"`
	FourierDensity    string        `json:"fourier_density"`
	CompactInterval   string        `json:"compact_interval"`
	GridStep          ExactRational `json:"grid_step"`
	TaylorDegree      int           `json:"taylor_degree"`
	LipschitzBound    ExactRational `json:"lipschitz_bound"`
	TailInterval      string        `json:"tail_interval"`
	TailLowerBound    string        `json:"tail_lower_bound"`
	TailAnchor        ExactRational `json:"tail_anchor"`
	OmittedDirections string        `json:"omitted_direction_control"`
	WholeLine         bool          `json:"whole_line"`
}

func factorial(n int) *big.Int {
	z := big.NewInt(1)
	for k := 2; k <= n; k++ {
		z.Mul(z, big.NewInt(int64(k)))
	}
	return z
}

func powRat(x *big.Rat, n int) *big.Rat {
	z := big.NewRat(1, 1)
	for k := 0; k < n; k++ {
		z.Mul(z, x)
	}
	return z
}

func trigTaylorBounds(x *big.Rat, degree int, cosine bool) (*big.Rat, *big.Rat) {
	s := new(big.Rat)
	start := 0
	if !cosine {
		start = 1
	}
	power := big.NewRat(1, 1)
	if start == 1 {
		power.Set(x)
	}
	x2 := new(big.Rat).Mul(x, x)
	for n := start; n <= degree; n += 2 {
		term := new(big.Rat).Quo(power, new(big.Rat).SetInt(factorial(n)))
		if ((n-start)/2)%2 == 0 {
			s.Add(s, term)
		} else {
			s.Sub(s, term)
		}
		power.Mul(power, x2)
	}
	remainder := new(big.Rat).Quo(powRat(x, degree+1), new(big.Rat).SetInt(factorial(degree+1)))
	return new(big.Rat).Sub(s, remainder), new(big.Rat).Add(s, remainder)
}

func compactPointLower(f TwoRadiusCompletion, t *big.Rat, degree int) *big.Rat {
	if t.Sign() == 0 {
		taylor, _ := TwoRadiusOriginTaylor(f)
		return rat(taylor.Constant)
	}
	cosLo, _ := trigTaylorBounds(t, degree, true)
	_, sinHi := trigTaylorBounds(t, degree-1, false)
	t2 := new(big.Rat).Mul(t, t)
	value := big.NewRat(1, 1)
	value.Add(value, new(big.Rat).Quo(new(big.Rat).Mul(big.NewRat(2, 1), new(big.Rat).Sub(cosLo, big.NewRat(1, 1))), t2))
	value.Add(value, new(big.Rat).Quo(new(big.Rat).Mul(new(big.Rat).Mul(big.NewRat(2, 1), new(big.Rat).Sub(big.NewRat(1, 1), rat(f.Constant))), sinHi), t))
	for _, a := range f.Atoms {
		x := new(big.Rat).Mul(rat(a.Radius), t)
		lo, _ := trigTaylorBounds(x, degree, true)
		value.Add(value, new(big.Rat).Mul(new(big.Rat).Mul(big.NewRat(2, 1), rat(a.Weight)), lo))
	}
	return value
}

// VerifyTwoRadiusCertificate performs the exact finite compact enclosure and
// analytic tail proof encoded by the M16 witness artifact.
func VerifyTwoRadiusCertificate(f TwoRadiusCompletion, p TwoRadiusWholeLineCertificate) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if err := VerifyTwoRadiusLocalLegality(f); err != nil {
		return err
	}
	const density = "hat(P)(t/(2*pi))=1+2(cos(t)-1)/t^2+2(1-c)sin(t)/t+2*w1*cos(r1*t)+2*w2*cos(r2*t)"
	const tail = "hat(P)>=1-2(w1+w2)-4/t^2-2(c-1)/|t|>0 for |t|>=4"
	if p.ID == "" || p.FourierVariable != "t=2*pi*xi" || p.FourierDensity != density || p.CompactInterval != "0<=|t|<=4" || p.TailInterval != "|t|>=4" || p.TailLowerBound != tail || !p.WholeLine {
		return fmt.Errorf("incomplete M16 Fourier certificate")
	}
	if p.GridStep != (ExactRational{1, 1000}) || p.TaylorDegree != 40 || rat(f.Constant).Cmp(big.NewRat(1, 1)) < 0 {
		return fmt.Errorf("unexpected compact or tail proof parameters")
	}
	wSum := new(big.Rat).Add(rat(f.Atoms[0].Weight), rat(f.Atoms[1].Weight))
	wantTail := new(big.Rat).Sub(big.NewRat(3, 4), new(big.Rat).Mul(big.NewRat(2, 1), wSum))
	wantTail.Sub(wantTail, new(big.Rat).Quo(new(big.Rat).Sub(rat(f.Constant), big.NewRat(1, 1)), big.NewRat(2, 1)))
	if wantTail.Sign() <= 0 || rat(p.TailAnchor).Cmp(wantTail) != 0 {
		return fmt.Errorf("tail anchor is not the conservative t=4 bound")
	}
	// |F'| <= 2 int_0^1 x(c-x)dx + 2 sum(wr)
	wantL := new(big.Rat).Sub(rat(f.Constant), big.NewRat(2, 3))
	for _, a := range f.Atoms {
		wantL.Add(wantL, new(big.Rat).Mul(big.NewRat(2, 1), new(big.Rat).Mul(rat(a.Weight), rat(a.Radius))))
	}
	if rat(p.LipschitzBound).Cmp(wantL) != 0 {
		return fmt.Errorf("wrong Lipschitz bound")
	}
	halfStep := big.NewRat(1, 2000)
	loss := new(big.Rat).Mul(wantL, halfStep)
	for k := int64(0); k <= 4000; k++ {
		center := big.NewRat(k, 1000)
		if compactPointLower(f, center, p.TaylorDegree).Cmp(loss) < 0 {
			return fmt.Errorf("compact enclosure fails at grid center %d/1000", k)
		}
	}
	if p.OmittedDirections != "exact Taylor enclosures plus a Lipschitz cover on |t|<=4 and an analytic bound on |t|>=4 cover every real t" {
		return fmt.Errorf("whole-line coverage text is inconsistent")
	}
	return nil
}
