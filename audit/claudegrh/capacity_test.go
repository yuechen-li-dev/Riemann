package claudegrh

import (
	"math/big"
	"testing"
)

// All quantities in the final comparison are measured in 1/120 units.
// A successful replay has only the scope of these finite assembly inputs.
func TestA4SignedCapacity(t *testing.T) {
	shapes, cases := 0, 0
	for top := 0; top <= 5; top++ {
		for u := 0; u <= 5; u++ {
			for v := 0; v <= min(u, 5-u); v++ {
				shapes++
				terms, err := twoRowMinors(u, v)
				if err != nil || len(terms) != v+1 {
					t.Fatalf("shape (%d,%d): support %v, error %v", u, v, terms, err)
				}
				size := u + v
				gain := twoLetterGain5(0, int64(size))
				if gain.Cmp(big.NewRat(int64(5*size), 24)) != 0 {
					t.Fatalf("dual (1,0) is not optimal for D=(0,%d)", size)
				}
				for k := 0; k <= v; k++ {
					powers := minorPowers{u - k, k, k, v - k}
					want := new(big.Int).Binomial(int64(v), int64(k)).Int64()
					if k%2 == 1 {
						want = -want
					}
					if coeff := terms[powers]; int64(coeff) != want {
						t.Fatalf("coefficient %v: got %d want %d", powers, coeff, want)
					}
					rows := [2]int{powers[0] + powers[1], powers[2] + powers[3]}
					for _, labels := range signedRows5(top, size) {
						cases++
						pos, neg := 0, 0
						if rows[0]+rows[1] != size {
							t.Fatal("non-top total")
						}
						seen := map[int]bool{}
						for row, j := range labels {
							if j == 0 {
								continue
							}
							if seen[j] || (j != 1 && j != -1) || rows[row] > u {
								t.Fatalf("invalid label fibre: %v, masses %v", labels, rows)
							}
							seen[j] = true
							if j == 1 {
								pos = rows[row]
							} else {
								neg = rows[row]
							}
						}
						if seen[1] != (top >= 5) || seen[-1] != (size >= 5) {
							t.Fatalf("wrong signed reach at (%d,%d,%d): %v", top, u, v, labels)
						}
						// e = 5*epsilon(0); dTop(0) is infinite.
						e := 0
						if top == 5 && u > 0 {
							e = 1
						}
						topTerm, wcap := 0, 0
						if u == 0 {
							topTerm = 25 * top
						}
						if size == 5 {
							wcap = u
						}
						lower := 25*size + topTerm - 25*wcap
						realized := 25*(size+top) - 24*top*e - 120*pos*(1-e) - 24*neg
						if realized < lower {
							t.Fatalf("(%d,%d,%d) k=%d labels=%v: %d < %d", top, u, v, k, labels, realized, lower)
						}
					}
				}
			}
		}
	}
	if shapes != 72 || cases != 126 {
		t.Fatalf("shapes=%d allocations=%d", shapes, cases)
	}
	t.Logf("%s: %d partition pairs, %d monomial/label allocations", FiniteVerified, shapes, cases)
}

func TestA4PartitionBoundary(t *testing.T) {
	for _, shape := range [][2]int{{2, 3}, {4, 2}, {0, -1}, {6, 0}} {
		if _, err := twoRowMinors(shape[0], shape[1]); err == nil {
			t.Fatalf("accepted %v", shape)
		}
	}
}
