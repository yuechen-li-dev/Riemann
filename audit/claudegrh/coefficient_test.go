package claudegrh

import (
	"math/big"
	"testing"
)

// One-row minors, both letters in matrix entry (0,0), TT empty, ep=sg=1.
// The audit derives the full admissible power series for each integral U.
// h fixes a=5^h exactly, avoiding floating-point comparisons of logarithms.
func TestA3MixedJets(t *testing.T) {
	cases := []struct {
		name                         string
		n, m, r, s, h, rank, next, k int64
		raw                          string
		order                        Order
	}{
		{"smallest positive rows", 1, 1, 1, 1, 1, 5, 6, 3, "1", Order{}},
		{"first mixed monomial", 5, 1, 0, 0, 1, 5, 6, 3, "1", Order{Value: 1}},
		{"first inverse denominator", 5, 25, 0, 0, 2, 7, 8, 3, "-1/5", Order{Value: 6}},
		{"vanishing jet", 6, 1, 0, 0, 2, 7, 8, 3, "0", Order{Infinite: true}},
		{"positive Cor2", 78125, 1, 78096, 0, 7, 18, 20, 18, "-78101/5", Order{Value: 19530}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := new(big.Int).Exp(big.NewInt(5), big.NewInt(tc.h), nil)
			if tc.n < 1 || tc.m < 1 || !a.IsInt64() || a.Int64() < max(tc.n, tc.m) || tc.h < 1 || tc.rank < 2*tc.h+3 || tc.next <= tc.rank || tc.k < 3 {
				t.Fatal("invalid diagram, rank or Hida setting")
			}
			power := func(n, r int64) *big.Rat {
				v, err := inverseJetPower(n, r)
				if err != nil {
					t.Fatal(err)
				}
				return v
			}
			// g(x)^r g(y)^s (1 + g(x)^5 g(y)).
			// Skip a zero factor before requesting an unnecessary high-degree coefficient.
			term := func(r, s int64) *big.Rat {
				y := power(tc.m, s)
				if y.Sign() == 0 {
					return y
				}
				return y.Mul(y, power(tc.n, r))
			}
			raw := new(big.Rat).Add(term(tc.r, tc.s), term(tc.r+5, tc.s+1))
			if raw.RatString() != tc.raw {
				t.Fatalf("raw coefficient %s", raw)
			}
			v := valuation(raw, big.NewInt(5))
			if !v.Infinite {
				for _, degree := range []int64{tc.n, tc.m} {
					for q := degree / 5; q > 0; q /= 5 {
						v.Value += q
					}
				}
			}
			if v != tc.order {
				t.Fatalf("jet order %+v", v)
			}
			gain := twoLetterGain5(tc.n, tc.m)
			thmA := new(big.Rat).Sub(gain, new(big.Rat).SetInt64(2*tc.rank*tc.next*(tc.h+1)))
			// E=0 and w=D/(k-1) meet both Cor2 shape equations.
			cor2 := new(big.Rat).Sub(thmA, big.NewRat(11*(tc.n+tc.m), 5*(tc.k-1)))
			if !v.atLeast(thmA) || !v.atLeast(cor2) || (tc.name == "positive Cor2" && cor2.Sign() <= 0) {
				t.Fatal("bound failed", thmA, cor2)
			}
			t.Logf("%s raw=%s jet=%+v ThmA=%s Cor2=%s", FiniteVerified, raw.RatString(), v, thmA.RatString(), cor2.RatString())
		})
	}
}

func TestA3InverseWindow(t *testing.T) {
	if _, err := inverseJetPower(625, 1); err == nil {
		t.Fatal("missing higher logarithm term")
	}
	// L(g)=x through degree 49 gives g_25=-1/5, g_49=1.
	g, err := inverseJetPower(49, 1)
	if err != nil || g.Cmp(big.NewRat(1, 1)) != 0 {
		t.Fatal(g, err)
	}
}

// One height-two minor and a one-box minor at the other letter.
// At this total degree the inverse logs are linear and the mixed
// exponential starts in degree six. U has the two indicated monomials.
func TestA3AlternatingMinor(t *testing.T) {
	for _, tc := range []struct {
		c    int64
		want Order
	}{{-4, Order{Value: 1}}, {1, Order{Infinite: true}}, {0, Order{}}} {
		input := map[[2]int]int64{{0, 1}: 1, {1, 0}: tc.c}
		jet := big.NewInt(0)
		for column := 0; column < 2; column++ {
			term := big.NewInt(input[[2]int{column, 1 - column}])
			if column == 1 {
				term.Neg(term)
			}
			jet.Add(jet, term)
		}
		order := integerOrder(jet, big.NewInt(5))
		thmA := new(big.Rat).Sub(twoLetterGain5(2, 1), big.NewRat(120, 1))
		if order != tc.want || !order.atLeast(thmA) {
			t.Fatal(tc, jet, order, thmA)
		}
	}
}

// Replay the one-letter split factors in the proof of LemmaS at the
// first two exceptional powers and at every intervening ordinary index.
func TestA3SplitOrders(t *testing.T) {
	for _, p := range []int64{3, 5, 7} {
		for n := int64(2); n <= p*p+1; n++ {
			var content big.Int
			for j := int64(1); j < n; j++ {
				content.GCD(nil, nil, &content, new(big.Int).Binomial(n, j))
			}
			m := n
			for m%p == 0 {
				m /= p
			}
			expected := int64(0)
			if m == 1 {
				expected = 1
			}
			if v := integerOrder(&content, big.NewInt(p)); v != (Order{Value: expected}) {
				t.Fatalf("p=%d n=%d gcd=%s", p, n, &content)
			}
		}
	}
}
