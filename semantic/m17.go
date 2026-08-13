package semantic

import (
	"fmt"
	"math/big"
)

const OneRadiusFourierDensity = "hat(P)(t/(2*pi))=1+2(cos(t)-1)/t^2+2(1-c)sin(t)/t+2*w*cos(r*t)"

// UnsaturatedOneRadiusCompletion is exactly the M17 family
// sigma=c 1_{|x|>1}dx+w(delta_-r+delta_r). Origin saturation is deliberately
// absent: it is represented by SaturatedOneRadiusCompletion below.
type UnsaturatedOneRadiusCompletion struct {
	Constant           ExactRational   `json:"c"`
	Radius             ExactRational   `json:"r"`
	Weight             ExactRational   `json:"w"`
	ExteriorDensity    ExactRational   `json:"exterior_density"`
	DensitySupport     string          `json:"density_support"`
	SupportLegal       bool            `json:"support_legal"`
	MeasureNonnegative bool            `json:"measure_nonnegative"`
	ExactTransform     string          `json:"exact_transform"`
	Taylor             OneRadiusTaylor `json:"taylor_coefficients"`
	TheoremProvenance  Reference       `json:"theorem_provenance"`
}

type OneRadiusTaylor struct {
	Constant  ExactRational `json:"constant_term"`
	Quadratic ExactRational `json:"t_squared_coefficient"`
	Quartic   ExactRational `json:"t_fourth_coefficient"`
	Sextic    ExactRational `json:"t_sixth_coefficient"`
}

func deriveOneRadiusTaylor(c, r, w ExactRational) OneRadiusTaylor {
	cr, rr, wr := rat(c), rat(r), rat(w)
	r2 := new(big.Rat).Mul(rr, rr)
	r4 := new(big.Rat).Mul(r2, r2)
	r6 := new(big.Rat).Mul(r4, r2)
	a0 := new(big.Rat).Mul(big.NewRat(2, 1), new(big.Rat).Sub(new(big.Rat).Add(big.NewRat(1, 1), wr), cr))
	a2 := new(big.Rat).Sub(new(big.Rat).Sub(new(big.Rat).Quo(cr, big.NewRat(3, 1)), big.NewRat(1, 4)), new(big.Rat).Mul(wr, r2))
	a4 := new(big.Rat).Add(new(big.Rat).Add(big.NewRat(-1, 360), new(big.Rat).Quo(new(big.Rat).Sub(big.NewRat(1, 1), cr), big.NewRat(60, 1))), new(big.Rat).Quo(new(big.Rat).Mul(wr, r4), big.NewRat(12, 1)))
	a6 := new(big.Rat).Sub(new(big.Rat).Add(big.NewRat(1, 20160), new(big.Rat).Quo(new(big.Rat).Sub(cr, big.NewRat(1, 1)), big.NewRat(2520, 1))), new(big.Rat).Quo(new(big.Rat).Mul(wr, r6), big.NewRat(360, 1)))
	v0, _ := exactRat(a0)
	v2, _ := exactRat(a2)
	v4, _ := exactRat(a4)
	v6, _ := exactRat(a6)
	return OneRadiusTaylor{v0, v2, v4, v6}
}

func NewUnsaturatedOneRadiusCompletion(c, r, w ExactRational, provenance Reference) UnsaturatedOneRadiusCompletion {
	return UnsaturatedOneRadiusCompletion{
		Constant: c, Radius: r, Weight: w, ExteriorDensity: c,
		DensitySupport: "|x|>1", SupportLegal: true, MeasureNonnegative: true,
		ExactTransform: OneRadiusFourierDensity, Taylor: deriveOneRadiusTaylor(c, r, w),
		TheoremProvenance: provenance,
	}
}

func (f UnsaturatedOneRadiusCompletion) Validate() error {
	for _, q := range []ExactRational{f.Constant, f.Radius, f.Weight, f.ExteriorDensity} {
		if err := q.Validate(); err != nil {
			return err
		}
	}
	if rat(f.Constant).Sign() <= 0 || rat(f.ExteriorDensity).Cmp(rat(f.Constant)) != 0 {
		return fmt.Errorf("exterior density must equal positive c")
	}
	if rat(f.Radius).Cmp(big.NewRat(1, 1)) < 0 {
		return fmt.Errorf("atom radius lies inside (-1,1)")
	}
	if rat(f.Weight).Sign() < 0 {
		return fmt.Errorf("atom weight must be nonnegative")
	}
	if f.DensitySupport != "|x|>1" || !f.SupportLegal || !f.MeasureNonnegative {
		return fmt.Errorf("illegal nonnegative exterior measure semantics")
	}
	if f.ExactTransform != OneRadiusFourierDensity {
		return fmt.Errorf("wrong one-radius Fourier normalization")
	}
	if f.Taylor != deriveOneRadiusTaylor(f.Constant, f.Radius, f.Weight) {
		return fmt.Errorf("stored Taylor coefficients do not match the exact family")
	}
	if f.TheoremProvenance.Citation == "" {
		return fmt.Errorf("missing theorem provenance")
	}
	return nil
}

func OneRadiusOriginTaylor(f UnsaturatedOneRadiusCompletion) (OneRadiusTaylor, error) {
	if err := f.Validate(); err != nil {
		return OneRadiusTaylor{}, err
	}
	return f.Taylor, nil
}

func VerifyOneRadiusLocalLegality(f UnsaturatedOneRadiusCompletion) error {
	if err := f.Validate(); err != nil {
		return err
	}
	for i, coefficient := range []ExactRational{f.Taylor.Constant, f.Taylor.Quadratic, f.Taylor.Quartic, f.Taylor.Sextic} {
		if rat(coefficient).Sign() > 0 {
			return nil
		}
		if rat(coefficient).Sign() < 0 {
			return fmt.Errorf("negative first nonzero Taylor coefficient at index %d", i)
		}
	}
	return fmt.Errorf("Taylor coefficients through t^6 vanish; higher-order analysis is required")
}

// SaturatedOneRadiusCompletion is a checked strict subtype, not the identity
// of the one-radius family.
type SaturatedOneRadiusCompletion struct {
	Family          UnsaturatedOneRadiusCompletion `json:"unsaturated_family"`
	OriginSaturated bool                           `json:"a0_equals_zero"`
	Provenance      string                         `json:"subfamily_provenance"`
}

func AsSaturatedOneRadius(f UnsaturatedOneRadiusCompletion) (SaturatedOneRadiusCompletion, error) {
	s := SaturatedOneRadiusCompletion{Family: f, OriginSaturated: true, Provenance: "M15 origin-saturation condition a0=0"}
	if err := s.Validate(); err != nil {
		return SaturatedOneRadiusCompletion{}, err
	}
	return s, nil
}

func (s SaturatedOneRadiusCompletion) Validate() error {
	if err := s.Family.Validate(); err != nil {
		return err
	}
	if !s.OriginSaturated || rat(s.Family.Taylor.Constant).Sign() != 0 {
		return fmt.Errorf("one-radius family is not origin-saturated")
	}
	return nil
}

// SaturatedOneRadiusCeiling recovers M15: a0=0 gives w=c-1, and r>=1
// gives a2<=c/3-1/4-(c-1)=(9-8c)/12.
func SaturatedOneRadiusCeiling(s SaturatedOneRadiusCompletion) error {
	if err := s.Validate(); err != nil {
		return err
	}
	if rat(s.Family.Constant).Cmp(big.NewRat(9, 8)) > 0 {
		return fmt.Errorf("origin-saturated one-radius branch exceeds its 9/8 ceiling")
	}
	return nil
}

type OneRadiusWholeLineCertificate struct {
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

func oneRadiusCompactPointLower(f UnsaturatedOneRadiusCompletion, t *big.Rat, degree int) *big.Rat {
	if t.Sign() == 0 {
		return rat(f.Taylor.Constant)
	}
	cosLo, _ := trigTaylorBounds(t, degree, true)
	_, sinHi := trigTaylorBounds(t, degree-1, false)
	t2 := new(big.Rat).Mul(t, t)
	value := big.NewRat(1, 1)
	value.Add(value, new(big.Rat).Quo(new(big.Rat).Mul(big.NewRat(2, 1), new(big.Rat).Sub(cosLo, big.NewRat(1, 1))), t2))
	value.Add(value, new(big.Rat).Quo(new(big.Rat).Mul(new(big.Rat).Mul(big.NewRat(2, 1), new(big.Rat).Sub(big.NewRat(1, 1), rat(f.Constant))), sinHi), t))
	x := new(big.Rat).Mul(rat(f.Radius), t)
	atomLo, _ := trigTaylorBounds(x, degree, true)
	value.Add(value, new(big.Rat).Mul(new(big.Rat).Mul(big.NewRat(2, 1), rat(f.Weight)), atomLo))
	return value
}

func VerifyOneRadiusCertificate(f UnsaturatedOneRadiusCompletion, p OneRadiusWholeLineCertificate) error {
	if err := VerifyOneRadiusLocalLegality(f); err != nil {
		return err
	}
	const tail = "hat(P)>=1-2*w-4/t^2-2*(c-1)/|t|>0 for |t|>=4"
	if p.ID == "" || p.FourierVariable != "t=2*pi*xi" || p.FourierDensity != OneRadiusFourierDensity || p.CompactInterval != "0<=|t|<=4" || p.TailInterval != "|t|>=4" || p.TailLowerBound != tail || !p.WholeLine {
		return fmt.Errorf("incomplete M17 Fourier certificate")
	}
	if p.TaylorDegree != 40 || rat(f.Constant).Cmp(big.NewRat(1, 1)) < 0 || rat(p.GridStep).Sign() <= 0 {
		return fmt.Errorf("unexpected compact or tail proof parameters")
	}
	steps := new(big.Rat).Quo(big.NewRat(4, 1), rat(p.GridStep))
	if !steps.IsInt() || !steps.Num().IsInt64() || steps.Num().Int64() > 100000 {
		return fmt.Errorf("compact grid must exactly partition [0,4] with at most 100000 cells")
	}
	wantTail := new(big.Rat).Sub(big.NewRat(3, 4), new(big.Rat).Mul(big.NewRat(2, 1), rat(f.Weight)))
	wantTail.Sub(wantTail, new(big.Rat).Quo(new(big.Rat).Sub(rat(f.Constant), big.NewRat(1, 1)), big.NewRat(2, 1)))
	if wantTail.Sign() <= 0 || rat(p.TailAnchor).Cmp(wantTail) != 0 {
		return fmt.Errorf("tail anchor is not the conservative t=4 bound")
	}
	wantL := new(big.Rat).Sub(rat(f.Constant), big.NewRat(2, 3))
	wantL.Add(wantL, new(big.Rat).Mul(big.NewRat(2, 1), new(big.Rat).Mul(rat(f.Weight), rat(f.Radius))))
	if rat(p.LipschitzBound).Cmp(wantL) != 0 {
		return fmt.Errorf("wrong Lipschitz bound")
	}
	halfStep := new(big.Rat).Quo(rat(p.GridStep), big.NewRat(2, 1))
	loss := new(big.Rat).Mul(wantL, halfStep)
	n := steps.Num().Int64()
	for k := int64(0); k <= n; k++ {
		center := new(big.Rat).Mul(big.NewRat(k, 1), rat(p.GridStep))
		if oneRadiusCompactPointLower(f, center, p.TaylorDegree).Cmp(loss) < 0 {
			return fmt.Errorf("compact enclosure fails at grid center %d", k)
		}
	}
	if p.OmittedDirections != "exact Taylor enclosures plus a Lipschitz cover on |t|<=4 and a radius-uniform analytic bound on |t|>=4 cover every real t" {
		return fmt.Errorf("whole-line coverage text is inconsistent")
	}
	return nil
}

type OneRadiusFamilyCeiling struct {
	ExactUpperExpression    string        `json:"exact_upper_expression"`
	RationalUpper           ExactRational `json:"rational_upper"`
	ContactFrequency        string        `json:"contact_frequency"`
	ParameterDomain         string        `json:"parameter_domain_theorem"`
	ProofRoute              string        `json:"proof_route"`
	StrictForRadiusAboveOne bool          `json:"strict_for_r_above_one"`
}

// DeriveOneRadiusFamilyCeiling encodes the calculus theorem used by M17. At
// t=pi/(2r), the atom vanishes. The remaining bound is h(t), where
// h(t)=1+t/(2 sin t)-tan(t/2)/t. Its positive-coefficient expansion on
// (0,pi/2] proves h is strictly increasing. Hence h(t)<=h(pi/2).
func DeriveOneRadiusFamilyCeiling() OneRadiusFamilyCeiling {
	return OneRadiusFamilyCeiling{
		ExactUpperExpression: "1+pi/4-2/pi", RationalUpper: ExactRational{1149, 1000},
		ContactFrequency: "t=pi/(2*r)", ParameterDomain: "r>=1, w>=0; global positivity also forces w<=1/2 by the atom-negative tail subsequence",
		ProofRoute:              "at t=pi/(2*r), cos(r*t)=0 and sin(t)>0, so c<=h(t)=1+t/(2*sin(t))-tan(t/2)/t; the positive even-power series of h-1 makes h strictly increasing on (0,pi/2]; Archimedes bounds 333/106<pi<355/113 give h(pi/2)<1149/1000",
		StrictForRadiusAboveOne: true,
	}
}

// RejectAboveOneRadiusRationalCeiling is a theorem-backed fast rejection. It
// is deliberately one-sided: candidates below 1149/1000 still require the
// whole-line verifier.
func RejectAboveOneRadiusRationalCeiling(f UnsaturatedOneRadiusCompletion) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if rat(f.Constant).Cmp(big.NewRat(1149, 1000)) >= 0 {
		return fmt.Errorf("candidate violates C_1R<1149/1000")
	}
	return nil
}
