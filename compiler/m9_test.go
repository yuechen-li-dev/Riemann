package compiler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

type exactSym2 struct{ a, b, d int }

func (x exactSym2) add(y exactSym2) exactSym2 { return exactSym2{x.a + y.a, x.b + y.b, x.d + y.d} }
func (x exactSym2) rank() int {
	if x.a*x.d-x.b*x.b != 0 {
		return 2
	}
	if x.a == 0 && x.b == 0 && x.d == 0 {
		return 0
	}
	return 1
}
func (x exactSym2) inertia() (int, int) {
	det, tr := x.a*x.d-x.b*x.b, x.a+x.d
	if det < 0 {
		return 1, 1
	}
	if det > 0 {
		if tr > 0 {
			return 2, 0
		}
		return 0, 2
	}
	if tr > 0 {
		return 1, 0
	}
	if tr < 0 {
		return 0, 1
	}
	return 0, 0
}
func outer2(x, y int) exactSym2 { return exactSym2{x * x, x * y, y * y} }
func pair2(ax, ay, bx, by int) exactSym2 {
	return exactSym2{2 * ax * bx, ax*by + ay*bx, 2 * ay * by}
}

func TestM9RankSubadditivityCriticalDirectionCollapseAndMultiplicity(t *testing.T) {
	u := outer2(1, 2)
	p := u.add(u) // multiplicity/repeated location direction changes weight only
	if p.rank() != 1 || p.rank() > u.rank()+u.rank() {
		t.Fatalf("rank accounting failed: P=%+v rank=%d", p, p.rank())
	}
	bound, err := CriticalRankLowerBound(2, 1)
	if err != nil || bound != 1 {
		t.Fatalf("critical lower bound=%d err=%v", bound, err)
	}
}

func TestM9OffCriticalLocalAndAggregateInertiaBudgets(t *testing.T) {
	independent := pair2(1, 0, 0, 1)
	p, n := independent.inertia()
	if p != 1 || n != 1 || independent.rank() != 2 {
		t.Fatalf("independent reflection pair: rank=%d inertia=(%d,%d)", independent.rank(), p, n)
	}
	dependent := pair2(1, 0, 2, 0)
	p, n = dependent.inertia()
	if p > 1 || n > 1 || dependent.rank() > 2 || (p == 1 && n == 1) {
		t.Fatalf("dependent case overstated: rank=%d inertia=(%d,%d)", dependent.rank(), p, n)
	}
	q1 := exactSym2{1, 0, -1}
	q2 := exactSym2{-1, 0, 1}
	q := q1.add(q2)
	p1, n1 := q1.inertia()
	p2, n2 := q2.inertia()
	pa, na := q.inertia()
	if pa > p1+p2 || na > n1+n2 || pa != 0 || na != 0 {
		t.Fatalf("aggregate subadditivity/cancellation failed: local=(%d,%d)+(%d,%d), aggregate=(%d,%d)", p1, n1, p2, n2, pa, na)
	}
}

func TestM9FinitePositiveDefiniteMatrixDoesNotExcludeOffCriticalBlock(t *testing.T) {
	pMatrix := exactSym2{0, 0, 2}
	qMatrix := exactSym2{1, 0, -1}
	gMatrix := pMatrix.add(qMatrix)
	gp, gn := gMatrix.inertia()
	qp, qn := qMatrix.inertia()
	if gp != 2 || gn != 0 || qp != 1 || qn != 1 {
		t.Fatalf("reverse-rejection fixture changed: G=(%d,%d), Q=(%d,%d)", gp, gn, qp, qn)
	}
}

func TestM9FiniteCriticalRankTheoremOnExactAdversarialFixtures(t *testing.T) {
	vectors := [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {1, -1}}
	for _, u := range vectors {
		pMatrix := outer2(u[0], u[1])
		for _, a := range vectors {
			for _, b := range vectors {
				qMatrix := pair2(a[0], a[1], b[0], b[1])
				gMatrix := pMatrix.add(qMatrix)
				gp, _ := gMatrix.inertia()
				qp, _ := qMatrix.inertia()
				if gp > pMatrix.rank()+qp {
					t.Fatalf("counterexample to exact theorem: P=%+v Q=%+v G=%+v", pMatrix, qMatrix, gMatrix)
				}
				lower, err := CriticalRankLowerBound(gp, qp)
				if err != nil || pMatrix.rank() < lower {
					t.Fatalf("derived lower bound failed: rank(P)=%d lower=%d", pMatrix.rank(), lower)
				}
			}
		}
	}
}

func TestM9CompileFusionContractsAndM7ReverseRejection(t *testing.T) {
	r, err := CompileM9()
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Contracts) != 5 || len(r.Invariants) != 1 {
		t.Fatalf("unexpected spectral surface: contracts=%d invariants=%d", len(r.Contracts), len(r.Invariants))
	}
	claim := r.Invariants[0]
	if claim.Invariant != semantic.PositiveIndexInvariant || claim.Bound != 2 || !claim.ExactTheoremPremise() {
		t.Fatalf("M7 invariant was not theorem-backed: %#v", claim)
	}
	if r.Fusion.SemanticMatrixID != r.M8.Dual.SemanticMatrixID || len(r.Fusion.ZeroSideFacts) == 0 || len(r.Fusion.ExplicitFormulaFacts) == 0 {
		t.Fatalf("representation lineage lost: %#v", r.Fusion)
	}
	if r.FinitePSDReverse.Accepted || r.DerivedTheorem.AsymptoticConsequenceDerived || r.ClaudeComparison.Reproduced || r.ClaudeComparison.RankTraceUsed {
		t.Fatal("M9 crossed a soundness or scope boundary")
	}
}

func TestM9ReportsDeterministicTypedAndExplicitAboutObstruction(t *testing.T) {
	r, err := CompileM9()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M9JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M9JSONReport(r)
	if err != nil || string(a) != string(b) {
		t.Fatal("M9 JSON is nondeterministic")
	}
	var doc map[string]any
	if err := json.Unmarshal(a, &doc); err != nil || doc["schema"] != "riemann.semantic-graph.m9" {
		t.Fatalf("bad typed report: schema=%v err=%v", doc["schema"], err)
	}
	h := M9HumanReport(r)
	for _, want := range []string{"ZERO-SIDE SPECTRAL BUDGET", "rank(P) >= max(0, n_plus(G)-B_off)", "REPRESENTATION FUSION", "NEWLY DERIVED MATHEMATICAL RESULT", "~1/2 reproduced: false", "rank-trace stage begun: false", "RH unresolved"} {
		if !strings.Contains(h, want) {
			t.Fatalf("human report missing %q", want)
		}
	}
}

func TestM9NegativeBudgetsRejected(t *testing.T) {
	if _, err := CriticalRankLowerBound(-1, 0); err == nil {
		t.Fatal("negative exact invariant accepted")
	}
}
