package claudegrh

import (
	"encoding/json"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/yuechen-li-dev/Riemann/compiler"
	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestGCDMinimumAndRationalValuations(t *testing.T) {
	for _, p := range []int64{2, 3, 5, 7} {
		for a := int64(-8); a <= 8; a++ {
			for b := int64(-8); b <= 8; b++ {
				in := Instance{[]string{big.NewInt(a).String(), big.NewInt(b).String()}, []Step{{Prime: uint32(p), Sigma: "0", Tau: 0}}}
				r, err := Replay(in)
				if err != nil || r.Outcome != FiniteVerified {
					t.Fatalf("%v: %v %v", in, r, err)
				}
				o := r.Steps[0]
				g := new(big.Int).GCD(nil, nil, big.NewInt(a), big.NewInt(b))
				if o.Content != g.String() || o.Minimum != integerOrder(g, big.NewInt(p)) {
					t.Fatalf("p=%d row=(%d,%d): %+v gcd=%s", p, a, b, o, g)
				}
			}
		}
	}
	p := big.NewInt(3)
	x := new(big.Int).Exp(p, big.NewInt(300), nil)
	if v := integerOrder(x, p); v != (Order{Value: 300}) {
		t.Fatal(v)
	}
	if v := valuation(new(big.Rat).SetFrac(big.NewInt(-2), x), p); v != (Order{Value: -300}) {
		t.Fatal(v)
	}
}

func TestAveragingCertificateAndSearch(t *testing.T) {
	data, err := os.ReadFile("testdata/averaging.json")
	if err != nil {
		t.Fatal(err)
	}
	var in Instance
	if err := json.Unmarshal(data, &in); err != nil {
		t.Fatal(err)
	}
	found, err := FirstAveragingFailure()
	if err != nil || !reflect.DeepEqual(found, in) {
		t.Fatalf("search: %v, %v", found, err)
	}
	r, err := Replay(in)
	if err != nil || r.Outcome != ReducedCounterexample {
		t.Fatalf("%+v %v", r, err)
	}
	o := r.Steps[0]
	if o.Minimum != (Order{Value: 0}) || o.Content != "1" || o.Sum != (Order{Value: 2}) || o.Average != (Order{Value: 1}) ||
		o.Pointwise || !o.Rounding || o.Integral || o.WitnessIndex != 0 || o.WitnessValue != "1/2" {
		t.Fatal(o)
	}
}

func TestShiftsAndDegeneracy(t *testing.T) {
	s := Step{Prime: 2, Sigma: "1", Tau: 1}
	for _, tt := range []struct {
		name    string
		in      Instance
		want    Outcome
		content string
	}{
		{"one", Instance{[]string{"2", "6"}, []Step{s}}, FiniteVerified, "2"},
		{"two", Instance{[]string{"4", "12"}, []Step{s, s}}, FiniteVerified, "2"},
		{"spent twice", Instance{[]string{"2", "6"}, []Step{s, s}}, ReducedCounterexample, "1"},
		{"two primes", Instance{[]string{"6", "18"}, []Step{s, {Prime: 3, Sigma: "1", Tau: 1}}}, FiniteVerified, "3"},
		{"unit prime", Instance{[]string{"1", "-3"}, []Step{{Prime: 5, Sigma: "0", Tau: 0}}}, FiniteVerified, "1"},
		{"ceiling attained", Instance{[]string{"4", "-12"}, []Step{{Prime: 2, Sigma: "3/2", Tau: 2}}}, FiniteVerified, "4"},
		{"rounding premise fails", Instance{[]string{"4"}, []Step{{Prime: 2, Sigma: "1", Tau: 2}}}, Inconclusive, "4"},
		{"valuation premise fails", Instance{[]string{"1"}, []Step{s}}, Inconclusive, "1"},
		{"negative shift", Instance{[]string{"-7"}, []Step{{Prime: 3, Sigma: "-3/2", Tau: -2}}}, FiniteVerified, "7"},
		{"zero and nonzero", Instance{[]string{"0", "14"}, []Step{s}}, FiniteVerified, "14"},
		{"zero functional", Instance{[]string{"0", "0"}, []Step{{Prime: 2, Sigma: "100", Tau: 100}}}, FiniteVerified, "0"},
		{"dimension zero", Instance{[]string{}, []Step{s}}, FiniteVerified, "0"},
		{"exact cancellation", Instance{[]string{"1", "-1"}, []Step{s}}, ReducedCounterexample, "1"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Replay(tt.in)
			if err != nil || r.Outcome != tt.want || len(r.Steps) != len(tt.in.Steps) {
				t.Fatalf("%+v %v", r, err)
			}
			o := r.Steps[len(r.Steps)-1]
			if o.Content != tt.content || o.HasAverage != (len(tt.in.Coefficients) > 0) || o.Minimum.Infinite != (o.Content == "0") {
				t.Fatal(o)
			}
		})
	}
}

func TestInvalidCertificates(t *testing.T) {
	for _, in := range []Instance{
		{}, {nil, []Step{{Prime: 2, Sigma: "0", Tau: 0}}}, {[]string{"0"}, nil}, {[]string{"1/2"}, []Step{{Prime: 2, Sigma: "1", Tau: 1}}},
		{[]string{"1"}, []Step{{Prime: 4, Sigma: "1", Tau: 1}}}, {[]string{"1"}, []Step{{Prime: 1, Sigma: "1", Tau: 1}}},
		{[]string{"1"}, []Step{{Prime: 0, Sigma: "1", Tau: 1}}}, {[]string{"1"}, []Step{{Prime: 2, Sigma: "1/0", Tau: 1}}},
		{[]string{"1"}, []Step{{Prime: 2, Sigma: "0", Tau: 4097}}},
		{[]string{"1"}, []Step{{Prime: 3, Sigma: "0", RowMap: "2/3"}}},
	} {
		if _, err := Replay(in); err == nil {
			t.Fatalf("accepted %+v", in)
		}
	}
}

func TestTwistAndRowImageAtSelectedPrime(t *testing.T) {
	for _, tt := range []struct {
		coefficient, row, twist, content string
		unit, localized, integral        bool
		outcome                          Outcome
	}{
		{"3", "-5", "2", "30", true, true, true, FiniteVerified},
		{"1", "3", "1", "3", true, false, false, Inconclusive},
		{"1", "1", "3", "3", false, true, false, Inconclusive},
		{"1", "1", "0", "0", false, true, false, Inconclusive},
	} {
		in := Instance{[]string{tt.coefficient}, []Step{{Prime: 3, Sigma: "1", Tau: 1, RowMap: tt.row, UnitTwist: tt.twist}}}
		r, err := Replay(in)
		if err != nil || r.Outcome != tt.outcome {
			t.Fatalf("%+v: %+v %v", in, r, err)
		}
		o := r.Steps[0]
		if !o.Pointwise || !o.Rounding || o.ScalarUnit != tt.unit || o.RowLocalization != tt.localized || o.Integral != tt.integral || o.Content != tt.content {
			t.Fatal(o)
		}
	}
}

func TestGraphEvidenceAndSourceUse(t *testing.T) {
	g, err := ProofGraph()
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []statement{local, dual, cost, window, estimate, region, rh, lifts} {
		if ok, _ := g.Certify(s.id()); ok {
			t.Fatalf("certified open claim %v", s)
		}
	}
	for _, s := range []statement{finite, average, conditional} {
		if a := g.AttemptProof(s.id(), dual.id()); a.Accepted {
			t.Fatalf("accepted %+v", a)
		}
	}
	for _, r := range route {
		c, _ := g.Claim(r.output.id())
		if c.Provenance.Source.URI != source+r.path || !reflect.DeepEqual(c.Proposition, r.output.proposition()) ||
			c.Provenance.Theorem != semantic.TheoremID(r.name) || c.Evidence[0].Kind != semantic.DerivedEvidence {
			t.Fatal(c)
		}
		var parents []semantic.ClaimID
		for _, s := range r.inputs {
			parents = append(parents, s.id())
		}
		if !reflect.DeepEqual(parents, c.Provenance.Parents) {
			t.Fatal(c.Provenance)
		}
	}
	c, _ := g.Claim(conditional.id())
	if len(c.Assumptions) != 1 || c.Evidence[0].Kind != semantic.UnverifiedConjectureEvidence {
		t.Fatal(c)
	}
	for _, s := range []statement{finite, average} {
		c, _ := g.Claim(s.id())
		if c.Evidence[0].Kind != semantic.CertifiedComputationEvidence || c.Proposition.(semantic.ExternalStatement).Instance == "" {
			t.Fatal(c)
		}
	}
	ancestry, err := g.Lineage(rh.id())
	if err != nil {
		t.Fatal(err)
	}
	seen := map[semantic.ClaimID]bool{}
	for _, id := range ancestry {
		seen[id] = true
	}
	if !seen[dual.id()] || !seen[periods.id()] || !seen[cost.id()] || !seen[estimate.id()] || seen[average.id()] || seen[finite.id()] || seen[lifts.id()] {
		t.Fatal(ancestry)
	}
}

func TestRejectWrongQuantifiersAndMissingPremises(t *testing.T) {
	r := route[1] // the local-to-global rule
	for _, args := range [][]statement{{local}, {local, average}, {finite, shiftData}, {average, shiftData}, {conditional, shiftData}} {
		g := compiler.NewGraph()
		if err := r.apply(g, args...); err == nil || len(g.Claims()) != 0 {
			t.Fatalf("bound invalid signature %v", args)
		}
	}
	if err := r.apply(compiler.NewGraph(), local, shiftData); err == nil {
		t.Fatal("accepted absent graph premises")
	}
	g, err := ProofGraph()
	if err != nil {
		t.Fatal(err)
	}
	tr := compiler.Transformation{ID: "forged", Pass: "test", From: finite.id(), To: local.id(), Relation: compiler.Implies, Theorem: "forged", Trusted: true}
	if err := g.AddTransformation(tr); err == nil {
		t.Fatal("core graph allowed finite to universal")
	}
	tr.From, tr.To, tr.Relation = local.id(), finite.id(), compiler.Equivalent
	if err := g.AddTransformation(tr); err == nil {
		t.Fatal("core graph allowed universal generalization by reversing an equivalence")
	}
	tr.To, tr.Relation = local.id(), compiler.Implies
	tr.From = conditional.id()
	if err := g.AddTransformation(tr); err == nil {
		t.Fatal("core graph dropped an assumption")
	}
}
