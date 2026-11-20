package byproducts

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

// CyclicCertificate uses source order (vertex, increasing odd offset) for
// the primal and vertex order for the other vectors. Values are exact strings.
type CyclicCertificate struct {
	Base     int      `json:"l"`
	Length   int      `json:"d"`
	Capacity []string `json:"capacity"`
	Primal   []string `json:"primal"`
	Dual     []string `json:"dual"`
}

func (c CyclicCertificate) shape() (int, error) {
	// The replay resource bound covers every pilot instance.
	if c.Base < 3 || c.Length < 2 || c.Length > 64 || c.Length%2 != 0 {
		return 0, fmt.Errorf("invalid or out-of-pilot l=%d, d=%d", c.Base, c.Length)
	}
	return c.Length * (c.Length / 2), nil
}

func rationalVector(values []string, n int, name string) ([]*big.Rat, error) {
	if len(values) != n {
		return nil, fmt.Errorf("%s has %d entries, require %d", name, len(values), n)
	}
	result := make([]*big.Rat, n)
	for i, s := range values {
		q, ok := new(big.Rat).SetString(s)
		if !ok || q.Sign() < 0 {
			return nil, fmt.Errorf("%s[%d] is not a nonnegative rational", name, i)
		}
		result[i] = q
	}
	return result, nil
}

// ReplayCyclic verifies both feasible vectors and objective equality from
// scratch. Finite weak duality then proves optimality and attainment.
func ReplayCyclic(c CyclicCertificate) (*big.Rat, error) {
	edges, err := c.shape()
	if err != nil {
		return nil, err
	}
	b, err := rationalVector(c.Capacity, c.Length, "capacity")
	if err != nil {
		return nil, err
	}
	x, err := rationalVector(c.Primal, edges, "primal")
	if err != nil {
		return nil, err
	}
	sig, err := rationalVector(c.Dual, c.Length, "dual")
	if err != nil {
		return nil, err
	}

	d := c.Length
	one := big.NewRat(1, 1)
	load := make([]big.Rat, d)
	var primal, dual big.Rat
	for s := 0; s < d; s++ {
		if sig[s].Cmp(one) > 0 {
			return nil, fmt.Errorf("dual bound at vertex %d", s)
		}
		dual.Add(&dual, new(big.Rat).Mul(b[s], sig[s]))
		for i := 1; i < d; i += 2 {
			dest := (s + i) % d
			z := x[s*(d/2)+i/2]
			cost := new(big.Int).Exp(big.NewInt(int64(c.Base)), big.NewInt(int64(i)), nil)
			factor := new(big.Rat).SetInt(cost)
			load[s].Add(&load[s], new(big.Rat).Mul(factor, z))
			load[dest].Add(&load[dest], z)
			cover := new(big.Rat).Mul(factor, sig[s])
			if cover.Add(cover, sig[dest]).Cmp(one) < 0 {
				return nil, fmt.Errorf("dual inequality at vertex %d offset %d", s, i)
			}
			primal.Add(&primal, z)
		}
	}
	for s := range load {
		if load[s].Cmp(b[s]) > 0 {
			return nil, fmt.Errorf("primal capacity at vertex %d", s)
		}
	}
	if primal.Cmp(&dual) != 0 {
		return nil, fmt.Errorf("unequal objectives %s and %s", primal.RatString(), dual.RatString())
	}
	return &primal, nil
}

// ReadOctCyclic reads integer numerators and a common positive denominator.
func ReadOctCyclic(r io.Reader) ([]CyclicCertificate, error) {
	scanner := bufio.NewScanner(r)
	var result []CyclicCertificate
	for scanner.Scan() {
		f := strings.Fields(scanner.Text())
		if len(f) < 2 {
			return nil, fmt.Errorf("missing parameters at row %d", len(result)+1)
		}
		l, e1 := strconv.Atoi(f[0])
		d, e2 := strconv.Atoi(f[1])
		c := CyclicCertificate{Base: l, Length: d}
		edges, err := c.shape()
		if e1 != nil || e2 != nil || err != nil || len(f) != 3+2*d+edges {
			return nil, fmt.Errorf("bad dimensions at row %d", len(result)+1)
		}
		c.Capacity = f[2 : 2+d]
		den, ok := new(big.Int).SetString(f[2+d], 10)
		if !ok || den.Sign() <= 0 {
			return nil, fmt.Errorf("bad denominator at row %d", len(result)+1)
		}
		var values []string
		for _, s := range f[3+d:] {
			num, ok := new(big.Int).SetString(s, 10)
			if !ok {
				return nil, fmt.Errorf("bad numerator at row %d", len(result)+1)
			}
			values = append(values, new(big.Rat).SetFrac(num, den).RatString())
		}
		c.Primal, c.Dual = values[:edges], values[edges:]
		if _, err := ReplayCyclic(c); err != nil {
			return nil, fmt.Errorf("row %d: %w", len(result)+1, err)
		}
		result = append(result, c)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("empty Oct export")
	}
	return result, nil
}
