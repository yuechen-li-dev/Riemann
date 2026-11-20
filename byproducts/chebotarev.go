package byproducts

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

type ClassSubset struct {
	Mask        int    `json:"class_mask"`
	Cardinality int    `json:"cardinality"`
	Density     string `json:"expected_density"`
}

type GroupCertificate struct {
	Name    string        `json:"group"`
	Product [][]int       `json:"product"`
	Inverse []int         `json:"inverse"`
	Classes [][]int       `json:"classes"`
	Subsets []ClassSubset `json:"stable_subsets"`
}

// ReplayGroup verifies the finite group and all its conjugacy-stable subsets.
// The density entries are predicted weights, with no analytic assertion.
func ReplayGroup(c GroupCertificate) error {
	n := len(c.Product)
	if n < 1 || n > 64 || len(c.Inverse) != n || len(c.Classes) == 0 || len(c.Classes) > 12 {
		return fmt.Errorf("invalid group dimensions")
	}
	for a, row := range c.Product {
		if len(row) != n || c.Inverse[a] < 0 || c.Inverse[a] >= n {
			return fmt.Errorf("invalid row or inverse at %d", a)
		}
		for _, p := range row {
			if p < 0 || p >= n {
				return fmt.Errorf("product outside group")
			}
		}
	}
	for a := range n {
		inv := c.Inverse[a]
		if c.Product[0][a] != a || c.Product[a][0] != a ||
			c.Product[a][inv] != 0 || c.Product[inv][a] != 0 {
			return fmt.Errorf("invalid identity or inverse at %d", a)
		}
		for b := range n {
			for z := range n {
				if c.Product[c.Product[a][b]][z] != c.Product[a][c.Product[b][z]] {
					return fmt.Errorf("associativity at (%d,%d,%d)", a, b, z)
				}
			}
		}
	}

	owner := make([]int, n)
	for a := range owner {
		owner[a] = -1
	}
	for i, members := range c.Classes {
		if len(members) == 0 {
			return fmt.Errorf("empty class %d", i)
		}
		for _, a := range members {
			if a < 0 || a >= n || owner[a] != -1 {
				return fmt.Errorf("classes are not a disjoint partition")
			}
			owner[a] = i
		}
	}
	for a, i := range owner {
		if i < 0 {
			return fmt.Errorf("class partition omits %d", a)
		}
		centralizer := 0
		for b := range n {
			if c.Product[a][b] == c.Product[b][a] {
				centralizer++
			}
			conjugate := c.Product[b][c.Product[a][c.Inverse[b]]]
			if owner[conjugate] != i {
				return fmt.Errorf("class %d is not conjugacy stable", i)
			}
		}
		// Stability gives a union of orbits. Orbit-stabilizer detects a
		// class that improperly merges two distinct orbits.
		if len(c.Classes[i])*centralizer != n {
			return fmt.Errorf("wrong class size at %d", a)
		}
	}
	choices := 1 << len(c.Classes)
	if len(c.Subsets) != choices {
		return fmt.Errorf("incomplete subset inventory")
	}
	seen := make([]bool, choices)
	for _, s := range c.Subsets {
		if s.Mask < 0 || s.Mask >= choices || seen[s.Mask] {
			return fmt.Errorf("invalid or duplicated class mask %d", s.Mask)
		}
		seen[s.Mask] = true
		count := 0
		for _, i := range owner {
			if s.Mask&(1<<i) != 0 {
				count++
			}
		}
		var sum big.Rat
		for i, members := range c.Classes {
			if s.Mask&(1<<i) != 0 {
				sum.Add(&sum, big.NewRat(int64(len(members)), int64(n)))
			}
		}
		got, ok := new(big.Rat).SetString(s.Density)
		if !ok || count != s.Cardinality || got.Cmp(big.NewRat(int64(count), int64(n))) != 0 || got.Cmp(&sum) != 0 {
			return fmt.Errorf("cardinality or density mismatch for mask %d", s.Mask)
		}
	}
	return nil
}

func ReadOctGroups(r io.Reader) ([]GroupCertificate, error) {
	scan := bufio.NewScanner(r)
	next := func(want int) ([]int, error) {
		if !scan.Scan() {
			return nil, fmt.Errorf("missing Oct table row: %v", scan.Err())
		}
		fields := strings.Fields(scan.Text())
		if want >= 0 && len(fields) != want {
			return nil, fmt.Errorf("wrong Oct row shape")
		}
		row := make([]int, len(fields))
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, err
			}
			row[i] = v
		}
		return row, nil
	}
	var result []GroupCertificate
	for scan.Scan() {
		c := GroupCertificate{Name: scan.Text()}
		dims, err := next(3)
		if err != nil {
			return nil, err
		}
		n, k, choices := dims[0], dims[1], dims[2]
		if n < 1 || n > 64 || k < 1 || k > 12 || choices != 1<<k {
			return nil, fmt.Errorf("invalid Oct dimensions")
		}
		if c.Inverse, err = next(n); err != nil {
			return nil, err
		}
		for range n {
			row, err := next(n)
			if err != nil {
				return nil, err
			}
			c.Product = append(c.Product, row)
		}
		for range k {
			row, err := next(-1)
			if err != nil {
				return nil, err
			}
			c.Classes = append(c.Classes, row)
		}
		for range choices {
			row, err := next(3)
			if err != nil {
				return nil, err
			}
			if row[2] != n {
				return nil, fmt.Errorf("incorrect density denominator")
			}
			c.Subsets = append(c.Subsets, ClassSubset{row[0], row[1], big.NewRat(int64(row[1]), int64(n)).RatString()})
		}
		if err := ReplayGroup(c); err != nil {
			return nil, fmt.Errorf("%s: %w", c.Name, err)
		}
		result = append(result, c)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("missing Oct groups")
	}
	return result, nil
}
