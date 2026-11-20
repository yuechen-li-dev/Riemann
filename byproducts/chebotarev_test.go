package byproducts

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"
)

func b3Certificates(t *testing.T) []GroupCertificate {
	t.Helper()
	data, err := os.ReadFile("../artifacts/b3/certificates.json")
	if err != nil {
		t.Fatal(err)
	}
	var report struct {
		Revision string
		Groups   []GroupCertificate
	}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	if report.Revision != Revision || len(report.Groups) != 3 {
		t.Fatal("wrong certificate inventory")
	}
	return report.Groups
}

func TestB3FiniteCertificates(t *testing.T) {
	expected := map[string][]int{"C5": {1, 1, 1, 1, 1}, "D8": {1, 1, 2, 2, 2}, "S4": {1, 3, 6, 6, 8}}
	for _, cert := range b3Certificates(t) {
		if err := ReplayGroup(cert); err != nil {
			t.Fatalf("%s: %v", cert.Name, err)
		}
		var sizes []int
		for _, class := range cert.Classes {
			sizes = append(sizes, len(class))
		}
		slices.Sort(sizes)
		if !slices.Equal(sizes, expected[cert.Name]) || len(cert.Subsets) != 32 {
			t.Fatalf("%s: wrong class sizes or subset count %v", cert.Name, sizes)
		}
		delete(expected, cert.Name)
	}
	if len(expected) != 0 {
		t.Fatal("missing pilot groups")
	}
}

func TestB3RejectCorruption(t *testing.T) {
	cases := []struct {
		name   string
		modify func(*GroupCertificate)
		want   string
	}{
		{"outside product", func(c *GroupCertificate) { c.Product[1][1] = 100 }, "outside group"},
		{"inverse", func(c *GroupCertificate) { c.Inverse[1] = 0 }, "identity or inverse"},
		{"associativity", func(c *GroupCertificate) { c.Product[1][2] = 1 }, "associativity"},
		{"overlap", func(c *GroupCertificate) { c.Classes[0] = append(c.Classes[0], c.Classes[1][0]) }, "disjoint partition"},
		{"mixed classes", func(c *GroupCertificate) { c.Classes[1][0], c.Classes[2][0] = c.Classes[2][0], c.Classes[1][0] }, "not conjugacy stable"},
		{"false cardinality", func(c *GroupCertificate) { c.Subsets[1].Cardinality = 2 }, "density mismatch"},
		{"false density", func(c *GroupCertificate) { c.Subsets[1].Density = "1/2" }, "density mismatch"},
		{"missing mask", func(c *GroupCertificate) { c.Subsets = c.Subsets[:31] }, "incomplete subset"},
		{"duplicate mask", func(c *GroupCertificate) { c.Subsets[1].Mask = 0 }, "duplicated class mask"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cert := b3Certificates(t)[2]
			tc.modify(&cert)
			if err := ReplayGroup(cert); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v, expected %s", err, tc.want)
			}
		})
	}
}
