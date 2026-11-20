package claudegrh

import (
	"fmt"
	"math/big"
)

// inverseJetPower computes [x^n] g(x)^r for the two-letter orbit at p=5.
// g inverts x+x^25/5+x^625/25+... . Lagrange inversion below only
// needs the first two logarithm terms when n-r < 624.
func inverseJetPower(n, r int64) (*big.Rat, error) {
	if n < 0 || r < 0 {
		return nil, fmt.Errorf("negative series degree")
	}
	if r == 0 && n == 0 {
		return big.NewRat(1, 1), nil
	}
	if r == 0 || n < r {
		return new(big.Rat), nil
	}
	if n-r >= 624 {
		return nil, fmt.Errorf("degree difference requires higher logarithm terms")
	}
	if (n-r)%24 != 0 {
		return new(big.Rat), nil
	}
	j := (n - r) / 24
	// r/n * binom(-n,j) / 5^j. Form the rising product in big.Int.
	c := big.NewInt(1)
	for i := int64(0); i < j; i++ {
		c.Mul(c, new(big.Int).Add(big.NewInt(n), big.NewInt(i)))
		c.Quo(c, big.NewInt(i+1))
	}
	if j%2 == 1 {
		c.Neg(c)
	}
	c.Mul(c, big.NewInt(r))
	den := new(big.Int).Exp(big.NewInt(5), big.NewInt(j), nil)
	return new(big.Rat).SetFrac(c, den.Mul(den, big.NewInt(n))), nil
}

// twoLetterGain5 solves the two capacity constraints 5x+y <= n,
// x+5y <= m for nonnegative flows. n,m are nonnegative diagram sizes.
func twoLetterGain5(n, m int64) *big.Rat {
	total := new(big.Int).Add(big.NewInt(n), big.NewInt(m))
	flow := new(big.Rat).SetFrac(total, big.NewInt(6))
	for _, size := range []int64{n, m} {
		if bound := new(big.Rat).SetInt64(size); bound.Cmp(flow) < 0 {
			flow = bound
		}
	}
	return new(big.Rat).Sub(new(big.Rat).Mul(big.NewRat(5, 24), new(big.Rat).SetInt(total)), flow)
}
