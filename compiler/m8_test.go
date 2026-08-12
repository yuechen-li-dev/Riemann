package compiler

import (
	"encoding/json"
	"math/cmplx"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func m8Basis(t *testing.T) semantic.OrderedBasis {
	t.Helper()
	m4, err := CompileM4()
	if err != nil {
		t.Fatal(err)
	}
	return m5WeilSpan(m4).Basis
}

func TestM8SingleSummandCoordinateOrientation(t *testing.T) {
	k, err := LowerZeroPointSummand(m8Basis(t), semantic.Point("rho"), 2)
	if err != nil {
		t.Fatal(err)
	}
	if got := k.Entries[1].Expression; !strings.Contains(got, "conjugate(M[f2](1-conjugate(rho)))*M[f3](rho)") {
		t.Fatalf("wrong row/column orientation: %s", got)
	}
	// Direct numerical check catches transpose or misplaced conjugation.
	v := [2]complex128{1 + 2i, -0.5 + 1i}
	w := [2]complex128{2 - 1i, 0.25 + 3i}
	c := [2]complex128{-1 + 0.5i, 2 + 1i}
	var matrix [2][2]complex128
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			matrix[i][j] = cmplx.Conj(w[i]) * v[j]
		}
	}
	direct := (c[0]*v[0] + c[1]*v[1]) * cmplx.Conj(c[0]*w[0]+c[1]*w[1])
	var coordinate complex128
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			coordinate += cmplx.Conj(c[i]) * matrix[i][j] * c[j]
		}
	}
	if cmplx.Abs(direct-coordinate) > 1e-12 {
		t.Fatalf("coordinate orientation mismatch: direct=%v matrix=%v", direct, coordinate)
	}
}

func TestM8CriticalContributionDerivedPSDWithRankCaveat(t *testing.T) {
	o, err := semantic.NewZeroOrbit("c", semantic.PointOnCriticalLine("rho"), semantic.CriticalLineOrbit, 1, []semantic.TheoremID{"symmetry"}, dlmfZerosReference)
	if err != nil {
		t.Fatal(err)
	}
	g, err := LowerZeroOrbitToMatrix(o, m8Basis(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(g.ReflectionPairs) != 2 || !g.Classification.PositiveSemidefinite || g.Classification.RankUpperBound != 2 {
		t.Fatalf("critical full orbit classification: %#v", g.Classification)
	}
	for _, p := range g.ReflectionPairs {
		if len(p.Points) != 1 || !p.Classification.PositiveSemidefinite || p.Classification.RankUpperBound != 1 || !p.Classification.RankOneIfNonzero {
			t.Fatalf("critical fixed block: %#v", p)
		}
	}
}

func TestM8OffCriticalPairConditionalIndefiniteness(t *testing.T) {
	o, err := semantic.NewZeroOrbit("o", semantic.Point("rho"), semantic.OffCriticalOrbit, 4, []semantic.TheoremID{"symmetry"}, dlmfZerosReference)
	if err != nil {
		t.Fatal(err)
	}
	g, err := LowerZeroOrbitToMatrix(o, m8Basis(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(g.ReflectionPairs) != 2 || !g.Classification.Hermitian || g.Classification.Indefinite {
		t.Fatalf("off orbit classification overstated: %#v", g)
	}
	for _, p := range g.ReflectionPairs {
		if len(p.Points) != 2 || p.Classification.RankUpperBound != 2 || p.Classification.Indefinite || p.Classification.IndefiniteCondition == "" || p.Classification.DegenerateCondition == "" {
			t.Fatalf("off pair classification: %#v", p.Classification)
		}
		for _, k := range p.PointContributions {
			if k.Multiplicity != 4 {
				t.Fatal("multiplicity discarded")
			}
		}
	}
}

func TestM8SyntheticIndependentAndDependentDiagnostics(t *testing.T) {
	d := m8ToyDiagnostics()
	if !(d[0].Determinant > -1e-12 && d[1].Determinant < 0 && d[1].Eigenvalues[0] < 0 && d[1].Eigenvalues[1] > 0) {
		t.Fatalf("unexpected critical/independent diagnostics: %#v", d)
	}
	if cmplx.Abs(complex(d[2].Determinant, 0)) > 1e-12 || d[2].Eigenvalues[0] < -1e-12 {
		t.Fatalf("dependent case did not degenerate: %#v", d[2])
	}
}

func TestM8AggregateDualIdentityAndNoReverseInference(t *testing.T) {
	r, err := testM8()
	if err != nil {
		t.Fatal(err)
	}
	if r.ZeroSide.SplitIdentity != "G=P+Q" || r.Dual.SemanticMatrixID != r.M7.Evaluation.Matrix.ID || r.Dual.IdentityTheorem != WeilExplicitFormulaTheoremID || r.Dual.NumericalIdentification {
		t.Fatalf("bad dual representation: %#v %#v", r.ZeroSide, r.Dual)
	}
	if r.FinitePSDReverseProof.Accepted || len(r.FinitePSDReverseProof.Diagnostics) == 0 {
		t.Fatalf("finite PSD reverse inference was not rejected: %#v", r.FinitePSDReverseProof)
	}
}

func TestM8ReportsDeterministicAndTyped(t *testing.T) {
	r, err := testM8()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M8JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M8JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatal("nondeterministic JSON")
	}
	var doc map[string]any
	if err := json.Unmarshal(a, &doc); err != nil {
		t.Fatal(err)
	}
	if doc["schema"] != "riemann.semantic-graph.m8" {
		t.Fatalf("schema=%v", doc["schema"])
	}
	h := M8HumanReport(r)
	for _, want := range []string{"K(p)_ij = conjugate", "G=P+Q", "same semantic matrix", "RH unresolved"} {
		if !strings.Contains(h, want) {
			t.Fatalf("report missing %q", want)
		}
	}
}
