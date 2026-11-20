package claudegrh

import "fmt"

// Exponents at the non-top letter in z00, z01, z10, z11 order.
// The top-letter factor is always z00^top, so it is kept separately.
type minorPowers [4]int

// twoRowMinors expands z00^(u-v) det(z)^v over Z. The shape bound
// implies v <= 2, so every intermediate coefficient fits in int.
func twoRowMinors(u, v int) (map[minorPowers]int, error) {
	if v < 0 || v > u || u > 5 || v > 5-u {
		return nil, fmt.Errorf("partition (%d,%d) outside the size-five box", u, v)
	}
	terms := map[minorPowers]int{{u - v, 0, 0, 0}: 1}
	for column := 0; column < v; column++ {
		next := make(map[minorPowers]int)
		for powers, coeff := range terms {
			diagonal, crossed := powers, powers
			diagonal[0]++
			diagonal[3]++
			crossed[1]++
			crossed[2]++
			next[diagonal] += coeff
			next[crossed] -= coeff
		}
		terms = next
	}
	return terms, nil
}

// Each entry gives the signed label on row 0 or 1 of letter 1. Zero
// means unlabelled. Enumerate every ordering of the qualifying labels.
func signedRows5(top, size int) [][2]int {
	if top == 5 && size == 5 {
		return [][2]int{{1, -1}, {-1, 1}}
	}
	if top == 5 {
		return [][2]int{{1, 0}}
	}
	if size == 5 {
		return [][2]int{{-1, 0}}
	}
	return [][2]int{{0, 0}}
}
