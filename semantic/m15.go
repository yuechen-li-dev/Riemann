package semantic

import (
	"fmt"
	"math/big"
)

// BoundaryAtomCompletion is the purpose-specific M15 family
// sigma_c = c 1_{|x|>1} dx + (c-1)(delta_-1+delta_1).
// Support and measure positivity are kept separate from Fourier positivity.
type BoundaryAtomCompletion struct {
	Constant           ExactRational `json:"constant"`
	ExteriorDensity    ExactRational `json:"exterior_density"`
	BoundaryAtomMass   ExactRational `json:"boundary_atom_mass"`
	DensitySupport     string        `json:"density_support"`
	AtomSupport        []int         `json:"atom_support"`
	SupportLegal       bool          `json:"support_legal"`
	MeasureNonnegative bool          `json:"measure_nonnegative"`
}

func (f BoundaryAtomCompletion) Validate() error {
	for _, value := range []ExactRational{f.Constant, f.ExteriorDensity, f.BoundaryAtomMass} {
		if err := value.Validate(); err != nil {
			return err
		}
	}
	if rat(f.ExteriorDensity).Cmp(rat(f.Constant)) != 0 {
		return fmt.Errorf("boundary-atom completion must cancel -c dx outside [-1,1]")
	}
	wantAtom := new(big.Rat).Sub(rat(f.Constant), big.NewRat(1, 1))
	if rat(f.BoundaryAtomMass).Cmp(wantAtom) != 0 {
		return fmt.Errorf("boundary atom mass must equal c-1")
	}
	if f.DensitySupport != "|x|>1" || len(f.AtomSupport) != 2 || f.AtomSupport[0] != -1 || f.AtomSupport[1] != 1 || !f.SupportLegal {
		return fmt.Errorf("sigma support must be exactly outside (-1,1), with legal boundary atoms")
	}
	if rat(f.ExteriorDensity).Sign() < 0 || rat(f.BoundaryAtomMass).Sign() < 0 || !f.MeasureNonnegative {
		return fmt.Errorf("sigma must be a nonnegative measure")
	}
	return nil
}

type WholeLinePDCertificate struct {
	ID                       string        `json:"id"`
	FourierVariable          string        `json:"fourier_variable"`
	FourierDensity           string        `json:"fourier_density"`
	ClearedDensity           string        `json:"cleared_density"`
	InnerInterval            string        `json:"inner_interval"`
	InnerSeries              string        `json:"inner_series"`
	FirstPositiveCoefficient ExactRational `json:"first_positive_coefficient"`
	TermRatioUpper           ExactRational `json:"term_ratio_upper"`
	OuterInterval            string        `json:"outer_interval"`
	OuterLowerBound          string        `json:"outer_lower_bound"`
	OuterAnchor              ExactRational `json:"outer_anchor"`
	OmittedDirections        string        `json:"omitted_direction_control"`
	WholeLine                bool          `json:"whole_line"`
}

// VerifyBoundaryAtomCertificate checks the exact constants on which the
// analytic proof rests. The proof object uses t=2*pi*xi and
// G(t)=4t^2*hat(P)(xi).
func VerifyBoundaryAtomCertificate(f BoundaryAtomCompletion, p WholeLinePDCertificate) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if rat(f.Constant).Cmp(big.NewRat(9, 8)) != 0 {
		return fmt.Errorf("M15 whole-line certificate is specialized to c=9/8")
	}
	const density = "hat(P)(xi)=1+2(cos(t)-1)/t^2+(cos(t)-sin(t)/t)/4, with value 0 at t=0"
	const cleared = "G(t)=4*t^2*hat(P)(xi)=t^2(4+cos(t))-t*sin(t)-8(1-cos(t))"
	const series = "G(t)=sum_{n>=3} (-1)^(n-1) 4(n-2)(n+1)t^(2n)/(2n)!"
	const outer = "G(t)>=3t^2-t-16>0 because |t|>=pi>3, q(3)=8, and q'(t)>0"
	if p.ID == "" || p.FourierVariable != "t=2*pi*xi" || p.FourierDensity != density || p.ClearedDensity != cleared || p.InnerSeries != series || p.OuterLowerBound != outer {
		return fmt.Errorf("incomplete Fourier representation")
	}
	local, err := BoundaryAtomLocalQuadraticNumerator(f.Constant)
	if err != nil || local.Sign() != 0 {
		return fmt.Errorf("endpoint does not cancel the local quadratic obstruction")
	}
	if rat(p.FirstPositiveCoefficient).Cmp(big.NewRat(1, 45)) != 0 {
		return fmt.Errorf("wrong leading coefficient in G(t)")
	}
	// On 0<=t<=pi, pi^2<10 and the largest consecutive-term ratio is
	// (5/2)*10/(8*7)=25/56<1, so the alternating tail decreases.
	if rat(p.TermRatioUpper).Cmp(big.NewRat(25, 56)) != 0 || rat(p.TermRatioUpper).Cmp(big.NewRat(1, 1)) >= 0 {
		return fmt.Errorf("inner alternating-series tail is not controlled")
	}
	// On t>=pi>3, G>=3t^2-t-16; its value at 3 is 8 and its
	// derivative is positive there and thereafter.
	if rat(p.OuterAnchor).Cmp(big.NewRat(8, 1)) != 0 {
		return fmt.Errorf("outer lower-bound anchor is wrong")
	}
	if p.OmittedDirections != "inner alternating series plus outer coercive bound covers every real t" {
		return fmt.Errorf("whole-line coverage text is inconsistent")
	}
	if p.InnerInterval != "0<=|t|<=pi" || p.OuterInterval != "|t|>=pi" || !p.WholeLine {
		return fmt.Errorf("certificate does not cover the whole line")
	}
	return nil
}

// BoundaryAtomLocalQuadraticNumerator is the exact numerator of the t^2
// coefficient after multiplying by 12. Positivity at the origin requires it
// to be nonnegative, hence c<=9/8 in this family.
func BoundaryAtomLocalQuadraticNumerator(c ExactRational) (*big.Rat, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return new(big.Rat).Sub(big.NewRat(9, 1), new(big.Rat).Mul(big.NewRat(8, 1), rat(c))), nil
}
