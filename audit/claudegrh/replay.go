// Package claudegrh audits the local-to-global division used in the pinned Lean source.
package claudegrh

import (
	"fmt"
	"math/big"
)

// Outcome describes the scope of a result, including failures of a weaker premise.
type Outcome string

const (
	ExactCounterexample   Outcome = "EXACT_COUNTEREXAMPLE"
	ReducedCounterexample Outcome = "REDUCED_MODEL_COUNTEREXAMPLE"
	FiniteVerified        Outcome = "FINITE_VERIFIED"
	Conditional           Outcome = "CONDITIONAL"
	Inconclusive          Outcome = "INCONCLUSIVE"
)

// Order is an additive valuation. Zero has infinite order.
type Order struct {
	Infinite bool
	Value    int64
}

func prime(p uint32) bool {
	if p < 2 {
		return false
	}
	for d := uint32(2); d <= p/d; d++ {
		if p%d == 0 {
			return false
		}
	}
	return true
}

func integerOrder(n, p *big.Int) Order {
	if n.Sign() == 0 {
		return Order{Infinite: true}
	}
	q := new(big.Int).Abs(n)
	var v int64
	var rem big.Int
	for {
		q.QuoRem(q, p, &rem)
		if rem.Sign() != 0 {
			return Order{Value: v}
		}
		v++
	}
}

func valuation(x *big.Rat, p *big.Int) Order {
	v := integerOrder(x.Num(), p)
	if !v.Infinite {
		v.Value -= integerOrder(x.Denom(), p).Value
	}
	return v
}

func (v Order) atLeast(s *big.Rat) bool {
	return v.Infinite || new(big.Rat).SetInt64(v.Value).Cmp(s) >= 0
}

// Instance uses K=Q, R=M(w)=Z, and initial N=Z^d.
// B(x,y)=(sum c_i*x_i)*y. Only the indicated prime is shifted at each step.
type Instance struct {
	Coefficients []string `json:"coefficients"`
	Steps        []Step   `json:"steps"`
}

type Step struct {
	Prime uint32 `json:"prime"`
	Sigma string `json:"sigma"` // exact rational lower bound
	Tau   int64  `json:"tau"`
	// Local integer multipliers at Prime. Omitted values and the multipliers
	// at every other prime are 1. The global pairing stays fixed.
	RowMap    string `json:"row_map,omitempty"`
	UnitTwist string `json:"unit_twist,omitempty"`
}

type Observation struct {
	Sum, Average, Minimum       Order
	Content                     string // nonnegative generator of the integer pairing ideal
	HasAverage                  bool
	Pointwise, Rounding         bool
	ScalarUnit, RowLocalization bool
	Integral                    bool
	// First nonintegral pairing on a shifted basis vector and y=1.
	WitnessIndex int
	WitnessValue string
}

type Result struct {
	Outcome Outcome
	Steps   []Observation
}

func localInteger(s string) (*big.Int, bool) {
	if s == "" {
		return big.NewInt(1), true
	}
	return new(big.Int).SetString(s, 10)
}

// Replay checks every generator for the local premise and evaluates all shifted
// basis pairings for the conclusion. Bilinearity covers arbitrary lattice points.
func Replay(in Instance) (Result, error) {
	if in.Coefficients == nil {
		return Result{}, fmt.Errorf("missing coefficient row (use [] for dimension zero)")
	}
	if len(in.Steps) == 0 {
		return Result{}, fmt.Errorf("missing shift")
	}
	coeff := make([]*big.Rat, len(in.Coefficients))
	for i, text := range in.Coefficients {
		n, ok := new(big.Int).SetString(text, 10)
		if !ok {
			return Result{}, fmt.Errorf("coefficient %d must be an integer", i)
		}
		coeff[i] = new(big.Rat).SetInt(n)
	}
	result := Result{Outcome: FiniteVerified}
	for i, step := range in.Steps {
		sigma, ok := new(big.Rat).SetString(step.Sigma)
		if !ok || !prime(step.Prime) {
			return Result{}, fmt.Errorf("step %d needs a prime and a rational lower bound", i)
		}
		// A bound on certificate expansion in this small replay kernel.
		if step.Tau < -4096 || step.Tau > 4096 {
			return Result{}, fmt.Errorf("step %d exceeds replay exponent limit", i)
		}
		p := new(big.Int).SetUint64(uint64(step.Prime))
		row, rowOK := localInteger(step.RowMap)
		twist, twistOK := localInteger(step.UnitTwist)
		if !rowOK || !twistOK {
			return Result{}, fmt.Errorf("step %d needs integer row and twist multipliers", i)
		}
		multiplier := new(big.Rat).SetInt(new(big.Int).Mul(row, twist))
		// For integer tau, tau <= ceil(sigma) iff tau-1 < sigma.
		obs := Observation{
			Minimum: Order{Infinite: true}, WitnessIndex: -1, Integral: true,
			Rounding:        new(big.Rat).SetInt64(step.Tau-1).Cmp(sigma) < 0,
			ScalarUnit:      integerOrder(twist, p) == (Order{Value: 0}),
			RowLocalization: integerOrder(row, p) == (Order{Value: 0}),
		}
		var sum big.Rat
		var content big.Int
		for _, x := range coeff {
			localPair := new(big.Rat).Mul(x, multiplier)
			sum.Add(&sum, localPair)
			content.GCD(nil, nil, &content, localPair.Num()) // previous legal steps are integral
			v := valuation(localPair, p)
			if !v.Infinite && (obs.Minimum.Infinite || v.Value < obs.Minimum.Value) {
				obs.Minimum = v
			}
		}
		obs.Pointwise = obs.Minimum.atLeast(sigma)
		obs.Content = content.String()
		obs.Sum = valuation(&sum, p)
		if len(coeff) > 0 {
			obs.HasAverage = true
			mean := new(big.Rat).Quo(&sum, new(big.Rat).SetInt64(int64(len(coeff))))
			obs.Average = valuation(mean, p)
		}
		exp := big.NewInt(step.Tau)
		scale := new(big.Rat).SetInt(new(big.Int).Exp(p, exp.Abs(exp), nil))
		if step.Tau > 0 {
			scale.Inv(scale)
		}
		for j, x := range coeff {
			x.Mul(x, scale)
			if !x.IsInt() && obs.Integral {
				obs.Integral, obs.WitnessIndex, obs.WitnessValue = false, j, x.RatString()
			}
		}
		result.Steps = append(result.Steps, obs)
		switch {
		case !obs.Rounding || !obs.ScalarUnit || !obs.RowLocalization:
			result.Outcome = Inconclusive
		case !obs.Pointwise:
			result.Outcome = Inconclusive
			if obs.HasAverage && obs.Average.atLeast(sigma) && !obs.Integral {
				result.Outcome = ReducedCounterexample
			}
		case !obs.Integral:
			result.Outcome = ExactCounterexample
		}
		if result.Outcome != FiniteVerified {
			break
		}
	}
	return result, nil
}

// FirstAveragingFailure enumerates positive binary rows in lexicographic order.
// It requires a finite nonzero average and at least one failed generator.
func FirstAveragingFailure() (Instance, error) {
	for a := int64(1); a <= 4; a++ {
		for b := a; b <= 4; b++ {
			in := Instance{[]string{fmt.Sprint(a), fmt.Sprint(b)}, []Step{{Prime: 2, Sigma: "1", Tau: 1}}}
			r, err := Replay(in)
			if err != nil {
				return Instance{}, err
			}
			if r.Outcome == ReducedCounterexample {
				return in, nil
			}
		}
	}
	return Instance{}, fmt.Errorf("no averaging failure in the search box")
}
