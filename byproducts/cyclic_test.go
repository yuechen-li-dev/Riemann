package byproducts

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/compiler"
	"github.com/yuechen-li-dev/Riemann/semantic"
)

func balancedCertificate() CyclicCertificate {
	return CyclicCertificate{3, 2, []string{"1", "1"}, []string{"1/4", "1/4"}, []string{"1/4", "1/4"}}
}

func TestB2RejectCorruptedCertificates(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*CyclicCertificate)
		want   string
	}{
		{"negative primal", func(c *CyclicCertificate) { c.Primal[0] = "-1/4" }, "nonnegative"},
		{"primal capacity", func(c *CyclicCertificate) { c.Primal[0] = "1/3" }, "primal capacity"},
		{"dual inequality", func(c *CyclicCertificate) { c.Dual[0] = "1/5" }, "dual inequality"},
		{"dual box", func(c *CyclicCertificate) { c.Dual[0] = "2" }, "dual bound"},
		{"different objectives", func(c *CyclicCertificate) { c.Primal[0] = "1/5" }, "unequal objectives"},
		{"source parameters", func(c *CyclicCertificate) { c.Base = 2 }, "out-of-pilot"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := balancedCertificate()
			tc.change(&c)
			if _, err := ReplayCyclic(c); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v", err)
			}
		})
	}
}

func TestB2ArtifactReplay(t *testing.T) {
	f, err := os.Open("../artifacts/b2/certificates.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var report struct {
		Revision     string
		Certificates []CyclicCertificate
	}
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&report); err != nil {
		t.Fatal(err)
	}
	if report.Revision != Revision || len(report.Certificates) != 435 {
		t.Fatal("wrong pilot inventory")
	}
	if n, err := checkTwoLetterDuality(); n != 405 || err != nil {
		t.Fatalf("B1 replay: %d %v", n, err)
	}
	for i, c := range report.Certificates {
		value, err := ReplayCyclic(c)
		if err != nil {
			t.Fatalf("certificate %d: %v", i, err)
		}
		if i < 405 {
			want := fmt.Sprintf("%d %d", i%81/9, i%9)
			if c.Length != 2 || c.Base != 3+i/81 || strings.Join(c.Capacity, " ") != want {
				t.Fatalf("B1 case %d has different source parameters", i)
			}
		}
		if i == 10 && value.RatString() != "1/2" {
			t.Fatal("wrong balanced optimum")
		}
	}
	// A legal example that needs the offset-five column. With this column
	// removed the only two supported vertices have no forward short arc.
	c := report.Certificates[432]
	if c.Length != 6 || c.Primal[2] != "1" {
		t.Fatal("missing offset-five regression")
	}
	c.Primal[2] = "0"
	if _, err := ReplayCyclic(c); err == nil {
		t.Fatal("accepted missing long arc")
	}
}

func TestB2ProofBoundary(t *testing.T) {
	g, err := CyclicGraph(balancedCertificate())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := g.Claim("b2.optimality"); !ok {
		t.Fatal("missing finite conclusion")
	}
	if ok, why := g.Certify("b2.optimality"); !ok {
		t.Fatalf("finite conclusion did not certify: %v", why)
	}
	for _, id := range []string{"b2.universal", "b2.lean-proof", "b2.formal-status"} {
		if ok, _ := g.Certify(semantic.ClaimID(id)); ok {
			t.Fatalf("%s promoted", id)
		}
	}
	err = g.AddTransformation(compiler.Transformation{ID: "b2.invalid", Pass: "B2", From: "b2.optimality", To: "b2.universal", Relation: compiler.Implies})
	if err == nil || !strings.Contains(err.Error(), "generalizes a finite") {
		t.Fatalf("got %v", err)
	}
}
