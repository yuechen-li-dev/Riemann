package byproducts

import (
	"fmt"
	"math/big"
)

// Check rational certificates for LemmaFa on Z/2. The two flow capacities
// are l*x+y <= a and x+l*y <= b. Equality with the dual cost proves optimality.
func checkTwoLetterDuality() (int, error) {
	count := 0
	for l := int64(3); l <= 7; l++ {
		for a := int64(0); a <= 8; a++ {
			for b := int64(0); b <= 8; b++ {
				var x, y, s, t *big.Rat
				switch {
				case l*a <= b:
					x, y, s, t = big.NewRat(0, 1), big.NewRat(a, 1), big.NewRat(1, 1), big.NewRat(0, 1)
				case l*b <= a:
					x, y, s, t = big.NewRat(b, 1), big.NewRat(0, 1), big.NewRat(0, 1), big.NewRat(1, 1)
				default:
					x, y = big.NewRat(l*a-b, l*l-1), big.NewRat(l*b-a, l*l-1)
					s, t = big.NewRat(1, l+1), big.NewRat(1, l+1)
				}
				sum := func(k int64, u, v *big.Rat) *big.Rat {
					return new(big.Rat).Add(new(big.Rat).Mul(big.NewRat(k, 1), u), v)
				}
				primal := sum(1, x, y)
				dual := sum(a, s, sum(b, t, new(big.Rat)))
				one := big.NewRat(1, 1)
				if x.Sign() < 0 || y.Sign() < 0 || sum(l, x, y).Cmp(big.NewRat(a, 1)) > 0 || sum(l, y, x).Cmp(big.NewRat(b, 1)) > 0 ||
					s.Sign() < 0 || t.Sign() < 0 || s.Cmp(one) > 0 || t.Cmp(one) > 0 || sum(l, s, t).Cmp(one) < 0 || sum(l, t, s).Cmp(one) < 0 ||
					primal.Cmp(dual) != 0 {
					return count, fmt.Errorf("invalid two-letter certificate at l=%d, a=%d, b=%d", l, a, b)
				}
				count++
			}
		}
	}
	return count, nil
}
