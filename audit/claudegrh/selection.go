package claudegrh

import (
	"fmt"
	"math/big"
)

// Each candidate is one chain for the full prime list. Its stage scalars
// multiply in Q. Bounds[i][j] is the proposed Fn for prime i and stage j.
type ChainCandidate struct {
	Name    string
	Scalars []string
	Bounds  [][]string
}

type SelectionBox struct {
	Primes     []uint32
	Column     string
	Low        []string
	Loss       []int64 // total dG over stages at each place
	Candidates []ChainCandidate
}

type LocalCheck struct {
	Candidate          string
	Prime              uint32
	Comparison, System bool
	FirstFailure       string
}

type SelectionResult struct {
	Outcome                                 Outcome
	Checks                                  []LocalCheck
	SeparateAtEveryPlace, JointAtEveryPlace bool
	Common                                  []string
	PlaceChoices                            []string // first joint choice for each prime, empty if none
}

// ReplaySelection tests a rank-one, one-complementary-direction model of
// CmpV and SysTw. None of its negative outcomes can refute the Lean theorem.
func ReplaySelection(box SelectionBox) (SelectionResult, error) {
	n := len(box.Primes)
	if n == 0 || len(box.Low) != n || len(box.Loss) != n {
		return SelectionResult{}, fmt.Errorf("incomplete prime data")
	}
	column, ok := new(big.Rat).SetString(box.Column)
	if !ok {
		return SelectionResult{}, fmt.Errorf("invalid actual column")
	}
	lows := make([]*big.Rat, n)
	seen := map[uint32]bool{}
	for i, p := range box.Primes {
		lows[i], ok = new(big.Rat).SetString(box.Low[i])
		if !ok || !prime(p) || seen[p] || box.Loss[i] < 0 {
			return SelectionResult{}, fmt.Errorf("invalid place %d", i)
		}
		seen[p] = true
	}
	cmp, sys, joint := make([]bool, n), make([]bool, n), make([]bool, n)
	result := SelectionResult{Outcome: Inconclusive, SeparateAtEveryPlace: true, JointAtEveryPlace: true, PlaceChoices: make([]string, n)}
	names := map[string]bool{}
	for _, c := range box.Candidates {
		if c.Name == "" || names[c.Name] || len(c.Bounds) != n || len(c.Scalars) == 0 || len(c.Scalars) != len(box.Candidates[0].Scalars) {
			return SelectionResult{}, fmt.Errorf("invalid candidate %q", c.Name)
		}
		names[c.Name] = true
		product := big.NewRat(1, 1)
		factors := make([]*big.Rat, len(c.Scalars))
		for j, s := range c.Scalars {
			factors[j], ok = new(big.Rat).SetString(s)
			if !ok {
				return SelectionResult{}, fmt.Errorf("%s: invalid scalar %d", c.Name, j)
			}
			product.Mul(product, factors[j])
		}
		common := true
		for i, p := range box.Primes {
			if len(c.Bounds[i]) != len(factors) {
				return SelectionResult{}, fmt.Errorf("%s: incomplete stage bounds", c.Name)
			}
			primeInt := new(big.Int).SetUint64(uint64(p))
			a, b := valuation(product, primeInt), valuation(column, primeInt)
			check := LocalCheck{Candidate: c.Name, Prime: p, Comparison: b.Infinite || (!a.Infinite && b.Value >= a.Value), System: true}
			if !check.Comparison {
				check.FirstFailure = "model CmpV: column order below chain order"
			}
			total := new(big.Rat).Neg(new(big.Rat).SetInt64(box.Loss[i]))
			for j, s := range c.Bounds[i] {
				f, valid := new(big.Rat).SetString(s)
				if !valid {
					return SelectionResult{}, fmt.Errorf("%s: invalid Fn", c.Name)
				}
				total.Add(total, f)
				if !valuation(factors[j], primeInt).atLeast(f) {
					check.System = false
					if check.FirstFailure == "" {
						check.FirstFailure = fmt.Sprintf("model LowAt: stage %d", j+1)
					}
				}
			}
			if total.Cmp(lows[i]) < 0 {
				check.System = false
				if check.FirstFailure == "" {
					check.FirstFailure = "model SysTw: sum Fn minus losses below low"
				}
			}
			both := check.Comparison && check.System
			if both && result.PlaceChoices[i] == "" {
				result.PlaceChoices[i] = c.Name
			}
			cmp[i] = cmp[i] || check.Comparison
			sys[i] = sys[i] || check.System
			joint[i] = joint[i] || both
			common = common && both
			result.Checks = append(result.Checks, check)
		}
		if common {
			result.Common = append(result.Common, c.Name)
		}
	}
	for i := range box.Primes {
		result.SeparateAtEveryPlace = result.SeparateAtEveryPlace && cmp[i] && sys[i]
		result.JointAtEveryPlace = result.JointAtEveryPlace && joint[i]
	}
	if len(result.Common) > 0 {
		result.Outcome = FiniteVerified
	} else if result.SeparateAtEveryPlace {
		result.Outcome = ReducedCounterexample
	}
	return result, nil
}
